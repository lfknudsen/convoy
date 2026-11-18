package main

import (
	"fmt"
	"os"
	"path/filepath"

	. "convoy/pkg"

	flag "github.com/spf13/pflag"
)

var Port = flag.IntP("port", "p", 8080, "Port to listen on.")
var Help = flag.BoolP("help", "h", false, "Print this message and exit.")

func printHelp() {
	fmt.Printf("CONVOY client\n\nUsage:\n\t%s [OPTIONS] <strings...>\n\nOptions:\n",
		filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printHelp
	flag.Parse()
	if *Help {
		printHelp()
		os.Exit(0)
	}
	if flag.NArg() < 1 {
		printHelp()
		os.Exit(1)
	}

	client := NewConvoyClient(*Port)
	if client == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Could not create Convoy client.")
		os.Exit(1)
	}
	for _, arg := range flag.Args() {
		client.Push(arg)
	}
}
