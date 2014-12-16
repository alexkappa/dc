package flag

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/ogier/pflag"
)

// Args represents the available command line arguments available to the user.
type Flag struct {
	Template,
	Data,
	Static,
	Addr string
	Serve,
	ShowHelp,
	ShowVersion bool
	Usage func()
}

// Validate that at least a template and some data were specified via the CLI.
func (f Flag) Valid() bool {
	return f.Template != "" && f.Data != ""
}

// New parses the command line arguments and returns an Args object.
func New(args []string) (Flag, error) {
	var f Flag
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.BoolVarP(&f.ShowHelp, "help", "h", false, "Prints this message")
	fs.BoolVarP(&f.ShowVersion, "version", "V", false, "Prints this version and exits")
	fs.StringVarP(&f.Template, "template", "t", "", "Input template. Can be a URL, file, directory or raw string.")
	fs.StringVarP(&f.Data, "data", "d", "", "Data to populate the template with. Can be a URL, file or raw string.")
	fs.StringVarP(&f.Static, "static", "s", "", "Static file root.")
	fs.BoolVarP(&f.Serve, "serve", "S", false, "Serve the template using HTTP.")
	fs.StringVarP(&f.Addr, "address", "a", "localhost:8080", "Server address.")

	f.Usage = func() {
		fmt.Printf("Usage: %s [options]\n", filepath.Base(args[0]))
		fs.PrintDefaults()
		os.Exit(2)
	}

	return f, fs.Parse(args[1:])
}
