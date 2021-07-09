package detector_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/micahjmartin/gopro/decoder"
	"github.com/micahjmartin/gopro/detector"
)

func TestReadStrings(t *testing.T) {
	c := make(chan detector.StringIndex)
	f, err := ioutil.ReadFile("test/blob.bin")
	if err != nil {
		t.Error(err)
	}
	go detector.FindStrings(c, bytes.NewReader(f))

	for idx := range c {
		if detector.IsPbString(f, idx) > 0 {
			fmt.Printf("[%d/%d] %s\n", idx.Index, len(idx.String), idx.String)
		} else {
			fmt.Println(idx.String)
		}
	}
}

func testVarint(t *testing.T, i uint64, b []byte) {
	result, n := decoder.ReadVarintReverse(b, len(b)-1)
	if n < 1 {
		t.Error("Could not read varint")
		return
	}
	if result != i {
		t.Error(fmt.Sprintf("Expected %d, read %d", i, result))
		return
	}
	fmt.Printf("%d == %d\n", result, i)
}
func TestReadVarintReverse(t *testing.T) {
	testVarint(t, 19, []byte{19})
	testVarint(t, 0x7F, []byte{0x7F})
	testVarint(t, 0x80, []byte{0x80, 0x01})
	testVarint(t, 123444, []byte{0x01, 0xb4, 0xc4, 0x07})
}
