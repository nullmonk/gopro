package detector

import (
	"io"
	"strings"
	"unicode"
)

var MIN_STR_LEN = 5

type StringIndex struct {
	Index  int
	String string

	strb strings.Builder
}

func NewStringIndex(i int) *StringIndex {
	return &StringIndex{
		Index: i,
		strb:  strings.Builder{},
	}
}

func (s *StringIndex) Finalize() string {
	s.String = s.strb.String()
	return s.String
}

func (s *StringIndex) IsValid() bool {
	return s.strb.Len() >= MIN_STR_LEN
}

func (s *StringIndex) Append(b byte) {
	s.strb.WriteByte(b)
}

// FindStrings detects strings in the given reader and passes tehm into the channel
func FindStrings(c chan StringIndex, r io.ByteReader) error {
	i := 0
	var idx *StringIndex

	for {
		b, err := r.ReadByte()
		if err != nil {
			close(c)
			return err
		}
		// If its a printable character, add to the string length
		if unicode.IsPrint(rune(b)) {
			if idx == nil {
				// Create a new string right here
				idx = NewStringIndex(i)
			}
			idx.Append(b)
		} else if idx != nil {
			// We have reached the end of a string, if its long enough, pass it thorugh the channel
			if idx.IsValid() {
				idx.Finalize()
				c <- *idx
			}
			idx = nil
		}
		i++
	}
}
