package process

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func run(name string, stdout io.Writer, errLog *log.Logger) error {
	errOut := bytes.Buffer{}
	cmd := exec.Command("go", "run", name)
	cmd.Stdout = stdout
	cmd.Stderr = &errOut

	err := cmd.Run()
	if errOut.Len() > 0 {
		errLog.Printf("%s", errOut.String())
	}
	return err
}

func writeNewFile(name string, lines []string) error {
	out, err := createNew(name)
	if err != nil {
		return err
	}
	for _, line := range lines {
		if _, err := out.Write([]byte(line)); err != nil {
			if err2 := out.Close(); err2 != nil {
				return fmt.Errorf("Error writing to and closing newfile %s: %s%s", name, err, err2)
			}
			return fmt.Errorf("Error writing to newfile %s: %s", name, err)
		}
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("Error closing newfile %s: %s", name, err)
	}
	return nil
}

func readUntil(r *bufio.Reader, marker string) ([]string, error) {
	lines := make([]string, 0, 50)
	var err error
	for err == nil {
		var line string
		line, err = r.ReadString('\n')
		if line != "" {
			lines = append(lines, line)
		}
		if strings.Contains(line, marker) {
			return lines, err
		}
	}
	return lines, err
}

func findLine(r *bufio.Reader, marker string) (string, error) {
	var err error
	for err == nil {
		line, err := r.ReadString('\n')
		if err == nil || err == io.EOF {
			if strings.Contains(line, marker) {
				return line, nil
			}
		}
	}
	return "", err
}

func createNew(filename string) (*os.File, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if os.IsExist(err) {
		return f, fmt.Errorf("File '%s' already exists.", filename)
	}
	return f, err
}