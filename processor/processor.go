// Package processor contains the code to generate text from embedded sourcecode.
package processor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	// Indicates a file was processed, but no gocog markers were found in it
	NoCogCode = errors.New("NoCogCode")
)

// New creates a new Processor with the given options.
func New(file string, opt *Options) *Processor {
	if opt == nil {
		opt = &Options{}
	}

	var logger *log.Logger
	if opt.Quiet {
		logger = log.New(io.Discard, "", log.LstdFlags)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}
	return &Processor{
		file,
		filepath.Dir(file),
		filepath.Base(file),
		regexp.MustCompile(opt.GenStart),
		regexp.MustCompile(opt.GenEnd),
		regexp.MustCompile(opt.OutStart),
		regexp.MustCompile(opt.OutEnd),
		opt,
		logger,
	}
}

// Processor holds the data for generating code for a specific file.
type Processor struct {
	File           string
	FileDir        string
	FileName       string
	GenStartRegexp *regexp.Regexp
	GenEndRegexp   *regexp.Regexp
	OutStartRegexp *regexp.Regexp
	OutEndRegexp   *regexp.Regexp
	*Options
	*log.Logger
}

// tracef will only log if verbose output is enabled.
func (p *Processor) tracef(format string, v ...interface{}) {
	if p.Verbose {
		p.Printf(format, v...)
	}
}

// Run processes the input file with the options specified.
// This will read the file, rewriting to a temporary file
// then run any embedded code, using the given options.
// It cleans up any code files it writes, and only overwrites the
// original if generation was successful.
func (p *Processor) Run() error {
	p.tracef("Processing file '%s'", p.File)

	output, err := p.tryCog()
	p.tracef("Output file: '%s'", output)

	if err == NoCogCode {
		p.tracef("Removing output file: '%s'\n", output)
		if err := os.Remove(output); err != nil {
			p.Println(err)
		}
		p.Printf("No generator code found in file '%s'", p.File)
		return err
	}

	// this is the success case - got to the end of the file without any other errors
	if err == io.EOF {
		p.tracef("Removing original file: '%s'\n", p.File)
		if err := os.Remove(p.File); err != nil {
			p.Printf("Error removing original file '%s': %s", p.File, err)
			return err
		}
		p.tracef("Renaming output file '%s' to original filename '%s'", output, p.File)
		if err := os.Rename(output, p.File); err != nil {
			p.Printf("Error renaming cog file '%s': %s", output, err)
			return err
		}
		p.Printf("Successfully processed '%s'", p.File)
		return nil
	} else {
		p.Printf("Error processing cog file '%s': %s", p.File, err)
		if output != "" {
			p.tracef("Removing output file: '%s'\n", output)
			if err := os.Remove(output); err != nil {
				p.Println(err)
			}
		}
		return err
	}
}

// tryCog encapsulates opening the original file, and creating the temporary output file.
// If output is nil, no output file was created, otherwise output is a valid file on disk
// that needs to be cleaned up after this function exits.
func (p *Processor) tryCog() (output string, err error) {
	in, err := os.Open(p.File)
	if err != nil {
		return "", err
	}
	defer in.Close()

	r := bufio.NewReader(in)

	values := map[string]string{
		"TMP":  os.TempDir(),
		"DIR":  p.FileDir,
		"FILE": p.FileName,
	}
	output = os.Expand(p.OutFile, func(s string) string { return values[s] })
	p.tracef("Writing output to %s", output)
	out, err := createNew(output)
	if err != nil {
		return "", err
	}
	defer out.Close()

	return output, p.gen(r, out)
}

// gen enacapsulates the process of generating text from an input and writing to an output.
func (p *Processor) gen(r *bufio.Reader, w io.Writer) error {
	for counter := 1; ; counter++ {
		prefix, submatches, err := p.cogPlainText(r, w, counter == 1)
		if err != nil {
			return err
		}

		if len(submatches) == 0 {
			return fmt.Errorf("No generator code processor defined")
		}
		cmd := submatches[0]

		var output string
		if 1 < len(submatches) {
			lines := []string{submatches[1]}
			output, err = p.generate(lines, prefix, counter, cmd)
		} else {
			output, err = p.cogGeneratorCode(r, w, prefix, counter, cmd)
		}
		if err != nil {
			return err
		}

		if err := p.cogToEnd(r, w, output); err != nil {
			return err
		}
	}
}

// cogPlainText reads any plaintext up to and including the GenStart regexp.
// If this is the first time we've read the file and we reach the end before
// finding the GenStart regexp, we won't write anything to the output file.
// Otherwise we'll write this plaintext back out to the output file as-is.
// Any prefix before the GenStart regexp is returned so we can handle single line comment tags.
// The match for the first regexp group (if any) in GenStart is returned.
func (p *Processor) cogPlainText(
	r *bufio.Reader,
	w io.Writer,
	firstRun bool,
) (prefix string, submatches []string, err error) {
	p.tracef("cogging plaintext")
	lines, found, prefix, submatches, err := readUntil(r, p.GenStartRegexp)
	if err == io.EOF {
		if found {
			// found gocog statement, but nothing after it
			if 1 < len(submatches) && p.UseEOF {
				// The file ended with a line containing the --genstart regexp (and has
				// no final newline).  This is only valid if the --genstart regexp has a
				// second submatch (containing the generator code to run) and options
				// permit EOF to end the generator output.
				lines[len(lines)-1] += "\n"
				err = nil
			} else {
				return "", []string{}, io.ErrUnexpectedEOF
			}
		} else if firstRun {
			// default case - no cog code, don't bother to write out anything
			return "", []string{}, NoCogCode
		}
		// didn't find it, but this isn't the first time we've run
		// so no big deal, we just ran off the end of the file.
	}
	if err != nil && err != io.EOF {
		return "", []string{}, err
	}

	// we can just write out the non-cog code to the output file
	// this also writes out the cog start line (if any)
	for _, line := range lines {
		if _, err := w.Write([]byte(line)); err != nil {
			return "", []string{}, err
		}
	}
	p.tracef("Wrote %v lines to output file", len(lines))

	if !found {
		return "", []string{}, err
	}

	return prefix, submatches, err
}

// Reads generator code lines from the reader until reaching the GenEnd regexp
// If the OutStart regexp is defined then further reads lines until reaching it
// Writes out the generator code to a file with the given name
// Any lines that start with the prefix or indent of the first line
// will have the prefix (for single line comments) or indent (for multi-
// line comments) removed.
// counter is a unique integer for each generator block in the same file.
// cmd is the command (if any) explicitly specified in the start marker.
func (p *Processor) cogGeneratorCode(
	r *bufio.Reader,
	w io.Writer,
	prefix string,
	counter int,
	cmd string,
) (output string, err error) {
	p.tracef("cogging generator code")
	lines, _, _, _, err := readUntil(r, p.GenEndRegexp)
	if err == io.EOF {
		return "", io.ErrUnexpectedEOF
	}
	if err != nil {
		return "", err
	}
	generatorLines := len(lines) - 1

	// if we have an optional marker for the start of output, read more lines until
	// we find it
	if p.OutStart != "" {
		moreLines, _, _, _, err := readUntil(r, p.OutStartRegexp)
		if err == io.EOF {
			return "", io.ErrUnexpectedEOF
		}
		if err != nil {
			return "", err
		}
		lines = append(lines, moreLines...)
	}

	// we have to write this out both to the output file and to the code file that we'll be running
	for _, line := range lines {
		if _, err := w.Write([]byte(line)); err != nil {
			return "", err
		}
	}
	p.tracef("Wrote %v lines to output file", len(lines))

	if !p.Excise && len(lines) > 0 {
		return p.generate(lines[:generatorLines], prefix, counter, cmd)
	}

	return "", nil
}

// generate writes out the generator code to a file and runs it.
// If running the code doesn't return any errors, the output is written to the output file.
// The file with the generator code is always deleted at the end of this function.
func (p *Processor) generate(
	lines []string,
	prefix string,
	counter int,
	cmd string,
) (output string, err error) {
	p.tracef("generating runnable code")
	values := map[string]string{
		"TMP":  os.TempDir(),
		"DIR":  p.FileDir,
		"FILE": p.FileName,
		"CTR":  strconv.Itoa(counter),
	}
	gen := os.Expand(p.GenFile, func(s string) string { return values[s] })
	if !p.Retain {
		defer func() {
			p.tracef("Removing generator code file: '%s'\n", gen)
			os.Remove(gen)
		}()
	}

	// write all but the last line to the generator file
	p.tracef("Creating generator code file: '%s'\n", gen)
	if err := writeNewFile(gen, lines, prefix); err != nil {
		return "", err
	}

	b := bytes.Buffer{}
	if err := p.runFile(gen, &b, cmd); err != nil {
		return "", err
	}

	return b.String(), nil
}

// runFile executes the given file with the command line specified in the Processor's options.
// If the process exits without an error, the output is written to the writer.
func (p *Processor) runFile(f string, w io.Writer, cmd string) error {
	p.tracef("output file %v", f)
	if p.Verbose {
		contents, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		p.tracef("file contents:\n%s", contents)
	}

	// When we run the generator command the current directory should be the same directory
	// as the input file, so that the generator code in the input file can reference other
	// files relatively.  Use an absolute path path for the temporary file containing the
	// generator code since it may not be in the same directory as the input file.
	f, err := filepath.Abs(f)
	if err != nil {
		return err
	}

	if strings.Contains(cmd, "%s") {
		cmd = fmt.Sprintf(cmd, f)
	}
	args := make([]string, len(p.Args))
	for i, s := range p.Args {
		if strings.Contains(s, "%s") {
			args[i] = fmt.Sprintf(s, f)
		} else {
			args[i] = s
		}
	}

	if err := run(p.FileDir, cmd, args, w, p.Logger); err != nil {
		return fmt.Errorf("Error generating code from source: %s", err)
	}
	return nil
}

// cogToEnd reads the old generateed code, up until the end tag. All but the last line is discarded.
func (p *Processor) cogToEnd(r *bufio.Reader, w io.Writer, output string) error {
	p.tracef("cogging to end")
	// we'll drop all but the COG_END line, so no need to keep them in memory
	line, found, err := findLine(r, p.OutEndRegexp)
	if err == io.EOF && !found {
		if !p.UseEOF {
			return io.ErrUnexpectedEOF
		}
		p.tracef("No gocog end statement, treating EOF as end statement.")
	}
	if err != nil && err != io.EOF {
		return err
	}

	// indent each output line according to the indent of the line with the end mark
	trimmedLine := strings.TrimLeftFunc(line, unicode.IsSpace)
	indent := strings.Repeat(" ", len(line)-len(trimmedLine))
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := fmt.Fprintf(w, "%s%s\n", indent, line); err != nil {
			return err
		}
	}

	if found {
		if _, err := w.Write([]byte(line)); err != nil {
			return err
		}
		p.tracef("Wrote 1 line to output file")
	}
	return err
}
