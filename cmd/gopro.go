package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/micahjmartin/gopro"
	"github.com/micahjmartin/gopro/decoder"
	"github.com/micahjmartin/gopro/detector"
	"github.com/micahjmartin/gopro/urldecoder"
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
		buf, err = ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from %s\n", os.Args[1])
			os.Exit(1)
		}
	}

	// Handle URLPB decoding
	var msg decoder.Message
	if buf[0] == urldecoder.Delimiter[0] {
		msg, err = gopro.DecodeUrl(string(buf))
	} else {
		msg, err = gopro.Decode(buf)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s. Falling back to protobuf detection.\n", err)
		detector.DetectProtobuf(buf)
		os.Exit(1)
	}

	dump := false
	for _, f := range os.Args {
		if f == "-dump" {
			dump = true
		}
	}
	if dump {
		gopro.Encode(msg, os.Stdout)
	} else {
		gopro.DumpMessage(msg)
	}
}
