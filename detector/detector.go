package detector

import (
	"bytes"
	"fmt"

	"github.com/micahjmartin/gopro/decoder"
)

func IsPbString(buf []byte, s StringIndex) int {
	length, n := decoder.ReadVarintReverse(buf, s.Index-1)
	// Couldnt decode varint
	if n < 1 {
		return -1
	}
	if int(length) != len(s.String) {
		return -1
	}

	idx := s.Index - n
	// Read the varint before that, this one will contain the field number and the wire_type
	pbKey, n := decoder.ReadVarintReverseWiretype(buf, idx-1, 2)
	idx -= n
	fieldNumber := uint64(pbKey >> 3)
	wireType := pbKey & 0x7
	if wireType != 2 {
		return -1
	}
	if fieldNumber > 0x80 || fieldNumber == 0 {
		return -1
	}
	return idx
}

type protoLocation struct {
	idx       int
	itemCount int
}

func DetectProtobuf(buf []byte) error {
	c := make(chan StringIndex)

	go FindStrings(c, bytes.NewReader(buf))

	for idx := range c {
		if pbidx := IsPbString(buf, idx); pbidx > 0 {
			// Try and read more protobuf items after the string
			b := decoder.NewBuffer(buf[pbidx:])
			msg := make(decoder.Message)
			var err error = nil
			var i decoder.Item
			for err == nil {
				i, err = decoder.ReadNextItem(b)
				if err != nil {
					break
				}
				msg.Add(i)
			}
			if len(msg) > 2 {
				fmt.Printf("Protobuf blobs detected in string at index: %d\n", idx.Index)
				decoder.DumpMessage(msg, "  ")
			}
		}
	}
	return nil
}
