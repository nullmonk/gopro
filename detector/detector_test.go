package detector_test

import (
	"fmt"
	"testing"

	"github.com/micahjmartin/gopro/detector"
)

var BYTES = []byte{0xfd, 0x2f, 0x05, 0x6a, 0x9f, 0x50, 0x40, 0x0a, 0x09, 0x4a, 0x65, 0x72, 0x69, 0x20, 0x44, 0x69, 0x61, 0x7a, 0x15, 0x00, 0x00, 0x92, 0xc2, 0x1d}

func TestIsPbString(t *testing.T) {
	idx := detector.NewStringIndex(9)
	idx.String = "Jeri Diaz"
	i := detector.IsPbString(BYTES, *idx)
	if i < 1 {
		t.Error("Did not detect valid pb string")
	}
	fmt.Println(i)
	//fmt.Println(string(BYTES[9:]))
}
