/* @GENERATE go@
package main
import(
  "fmt"
  "os"
  "os/exec"
)
func main() {
  fmt.Println("")
  fmt.Print("/", "*", "\n")
  fmt.Println("Package main creates an executable that will generate text from inline sourcecode.\n")
  cmd := exec.Command("gocog")
  cmd.Stdout = os.Stdout
  cmd.Run()
  fmt.Print("*","/", "\n")
  fmt.Println("package main")
}
@OUTPUT@ */

/*
Package main creates an executable that will generate text from inline sourcecode.

Usage:

	gocog [OPTIONS] [INFILE | @FILELIST] ...

	Runs gocog over each infile.
	Strings prepended with @ are assumed to be files continaing newline delimited lists of gocog command lines.
	Command line options are passed to each command line in the file list, but options on the file list line
	will override command line options. You may have filelists specified inside filelist files.

Application Options:

	-z, --eof       The end marker can be assumed at eof.
	-v, --verbose   enables verbose output
	-q, --quiet     turns off all output
	-S, --serial    Write to the specified cog files serially
	-a, --args=     Comma separated arguments to cmd, %s for the code file
	                (default: [%s])
	-g, --genstart= Regexp that starts gocog statements (default:
	                \[\[\[generate\s+([^]]+)\]\]\])
	-G, --genend=   Regexp that ends gocog statements (default:
	                \[\[\[output\]\]\])
	-o, --outstart= Optional regexp that starts gocog output
	-O, --outend=   Regexp that ends gocog output (default: \[\[\[end\]\]\])
	-f, --genfile=  Filename template for temp generator code files (default:
	                $DIR/cog_${FILE}_cog_${CTR}_.txt)
	-F, --outfile=  Filename template for temp output files (default:
	                $DIR/${FILE}_cog)
	-x, --excise    Excise all the generated output without running the
	                generators.
	-r, --retain    Don't delete temporary files containing generator code.
	-V, --version   Display the version of gocog

Help Options:

	-h, --help      Show this help message
*/
package main
