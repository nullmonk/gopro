package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/micahjmartin/gopro"
)

func main() {
	stat, _ := os.Stdin.Stat()
	var buf []byte
	var err error
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		buf, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from stdin")
			os.Exit(1)
		}
	} else {
		if len(os.Args) < 2 {
			fmt.Fprintln(os.Stderr, "Pass protobuf blobs in through either stdin or a filename")
			fmt.Fprintln(os.Stderr, "USAGE: gopro <filename>")
			os.Exit(1)
		}
		buf, err = ioutil.ReadFile(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from %s\n", os.Args[2])
			os.Exit(1)
		}
	}

	msg, err := gopro.Decode(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	gopro.DumpMessage(msg)
}
