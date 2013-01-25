gocog - generate code for any language, with any language
=====

gocog is a command line executable that processes in-line code in a file and outputs the results into the same file.

Binaries for popular OSes are available [on the wiki](https://github.com/natefinch/gocog/wiki)

Read the [Godoc documentation](http://godoc.org/github.com/natefinch/gocog)

Design of gocog is heavily based on [cog.py](http://nedbatchelder.com/code/cog/).  Many thanks to Ned Batchelder for a really great design.
<!-- {{{gocog
package main
import(
  "bytes"
  "fmt"
  "os/exec"
)
func main() {
  b := &bytes.Buffer{}
  cmd := exec.Command("gocog")
  cmd.Stdout = b
  cmd.Run()
  for {
    line, err := b.ReadString(byte('\n'))
    if len(line) > 0 {
      fmt.Print("\t", line)
    }
    if err != nil {
      break
    }
  }
}
gocog}}} -->
	Usage:
	  gocog [OPTIONS] [INFILE | @FILELIST] ...
	
	  Runs gocog over each infile. 
	  Strings prepended with @ are assumed to be files continaing newline delimited lists of files to be processed.
	
	Help Options:
	  -h, --help         Show this help message
	
	Application Options:
	  -z, --eof          The end marker can be assumed at eof.
	  -v, --verbose      enables verbose output
	  -q, --quiet        turns off all output
	  -S, --serial       Write to the specified cog files serially
	  -c, --cmd          The command used to run the generator code (go)
	  -a, --args         Comma separated arguments to cmd, %s for the code file
	                     ([run, %s])
	  -e, --ext          Extension to append to the generator filename (.go)
	  -M, --startmark    String that starts gocog statements ([[[)
	  -E, --endmark      String that ends gocog statements (]]])
	  -x, --excise       Excise all the generated output without running the
	                     generators.
	  -V, --version      Display the version of gocog
<!-- {{{end}}} -->

How it works
------

Code is embedded in comments in the given files, delimited thusly:

    [[[gocog
      <generator code that will be run to generate output>
    gocog]]]
    [[[end]]]

Anything written to standard out from the generator code will be injected between gocog]]] and [[[end]]]

The generator code embedded in the file is written out to a temporary file on disk by gocog named filename_cog.ext (where filename is the original filename, and ext is the appropriate extension for the generator language. This file is then run using the specified command line tool.  Standard output generated by the generator code is piped to a new file named filename_cog, along with the original text. If generation is successful for all gocog blocks in a file, this output file is then used to replace the original file.

If at any time there is an error while running gocog over a file, the original file is not replaced. Errors from the generator code will be piped to gocog's stderr.

By default, each file is processed in parallel, to speed the processing of large numbers of files.

The gocog marker tags can be preceded by any text (such as comment tags to prevent your compiler/interpreter from barfing on them).

Any non-whitespace text that precedes the gocog start mark will be treated as a single line comment tag and will be removed in the generator code that is written out - for example:

	# [[[gocog
	# do something here
	#     and some indent
	# gocog]]]
	# [[[end]]]

output code:

	do something here
	    and some indent

You can rerun gocog over the same file multiple times. Previously generated text will be discarded and replaced by the newly generated text.

You can have multiple blocks of gocog generator code inside the same file.

Current Limitations
----------

* All marker tags must be on different lines

Todo
----
Gocog is a work in progress. Here's some stuff I'll be adding soon

* Support for single line gocog statements e.g. [[[gocog your_code_here gocog]]]
* Anything commented out in [options.go](https://github.com/natefinch/gocog/blob/master/processor/options.go)
* Better support for correct indentation
* Pre and post-run commands
* Support for running across an entire directory / tree
* Support for adding command line flags to files in file lists
* Support for standardized header and footer text for extracted generator code (to remove boilerplate)
* Support for running different generator blocks in the same file in parallel (currently they're run serially)

Examples
------
Gocog uses gocog! Check out [README.md](https://raw.github.com/natefinch/gocog/master/README.md), [main.go](https://github.com/natefinch/gocog/blob/master/main.go), [gocog.go](https://github.com/natefinch/gocog/blob/master/gocog.go) and [doc.go](https://github.com/natefinch/gocog/blob/master/doc.go) (doc.go uses single line comments as an example of how that works).
The command line I use for gocog's own use is in [update.sh](https://github.com/natefinch/gocog/blob/master/update.sh)

(I use different start and end markers so gocog won't get tripped up by my documentation that uses the same markers)

Now for a toy example:
Using generator code written in Go to write out properties for a C# class

    using System;
    
    namespace foo 
    {
      public class Foo
      {
        /* [[[gocog
        package main
        import "fmt"
        func main() {
          for _, s := range []string{ "Bar", "Baz", "Bat", "Stuff" } {
            fmt.Printf("\t\tpublic String %s { get; set; }\n", s)
          }
        }
        gocog]]]  */
        // [[[end]]]
      }
    }

Output:

    using System;
    
    namespace foo 
    {
      public class Foo
      {
        /* [[[gocog
        package main
        import "fmt"
        func main() {
          for _, s := range []string{ "Bar", "Baz", "Bat", "Stuff" } {
            fmt.Printf("\t\tpublic String %s { get; set; }\n", s)
          }
        }
        gocog]]]  */
        public String Bar { get; set; }
        public String Baz { get; set; }
        public String Bat { get; set; }
        public String Stuff { get; set; }
        // [[[end]]]
      }
    }
    
Things to note:
The generator code and gocog markers are all hidden from the original file's compiler by comments, so the file is always valid.

The generator code stays in the file even after running through gocog. This keeps the generator code and the target close together so there's no need to worry about one getting lost. It also makes it a lot more clear where and how the output will be used in the original file.
