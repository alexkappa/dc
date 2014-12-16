package main

import (
	"fmt"
	"os"

	"github.com/alexkappa/dc/console"
	"github.com/alexkappa/dc/flag"
	"github.com/alexkappa/dc/server"
)

type args struct {
	template,
	data,
	static,
	addr string
	serve bool
}

func main() {
	f, err := flag.New(os.Args)
	if f.ShowHelp || err != nil {
		f.Usage()
		os.Exit(0)
	}

	if !f.Valid() {
		fmt.Fprintf(os.Stderr, "Usage %s:\n", os.Args[0])
		f.Usage()
		os.Exit(1)
	}

	if f.Serve {
		s := server.New(f)
		go s.Open()
		if err = s.Listen(); err != nil {
			fmt.Fprintf(os.Stderr, "Server failed to start. %s\n", err)
			os.Exit(1)
		}
	} else {
		c := console.New(f)
		if err := c.Render(); err != nil {
			fmt.Fprintf(os.Stderr, "Template faled to render. %s\n", err)
			os.Exit(1)
		}
	}
}
