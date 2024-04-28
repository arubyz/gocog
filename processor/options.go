package processor

type Options struct {
	UseEOF   bool     `short:"z" long:"eof" description:"The end marker can be assumed at eof."`
	Verbose  bool     `short:"v" long:"verbose" description:"enables verbose output"`
	Quiet    bool     `short:"q" long:"quiet" description:"turns off all output"`
	Serial   bool     `short:"S" long:"serial" description:"Write to the specified cog files serially"`
	Args     []string `short:"a" long:"args" description:"Comma separated arguments to cmd, %s for the code file"`
	GenStart string   `short:"g" long:"genstart" description:"Regexp that starts gocog statements"`
	GenEnd   string   `short:"G" long:"genend" description:"Regexp that ends gocog statements"`
	OutStart string   `short:"o" long:"outstart" description:"Optional regexp that starts gocog output"`
	OutEnd   string   `short:"O" long:"outend" description:"Regexp that ends gocog output"`
	GenFile  string   `short:"f" long:"genfile" description:"Filename template for temp generator code files"`
	OutFile  string   `short:"F" long:"outfile" description:"Filename template for temp output files"`
	Excise   bool     `short:"x" long:"excise" description:"Excise all the generated output without running the generators."`
	Retain   bool     `short:"r" long:"retain" description:"Don't delete temporary files containing generator code."`
	Version  bool     `short:"V" long:"version" description:"Display the version of gocog"`
	//	Checksum bool              `short:"c" description:"Checksum the output to protect it against accidental change."`
	//	Delete   bool              `short:"d" description:"Delete the generator code from the output file."`
	//	Define   map[string]string `short:"D" description:"Define a global string available to your generator code."`
	//	Include  string            `short:"I" description:"Add PATH to the list of directories for data files and modules."`
	//	Output   string            `short:"o" description:"Write the output to OUTNAME."`
	//	Suffix   string            `short:"s" description:"Suffix all generated output lines with STRING."`
	//	Unix     bool              `short:"U" description:"Write the output with Unix newlines (only LF line-endings)."`
	//	WriteCmd string            `short:"w" description:"Use CMD if the output file needs to be made writable. A %s in the CMD will be filled with the filename."`
}
