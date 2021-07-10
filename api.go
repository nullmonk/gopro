package gopro

import (
	"github.com/micahjmartin/gopro/decoder"
)

// DumpMessage dumps a Message to stdout
func DumpMessage(m decoder.Message) {
	decoder.DumpMessage(m, "")
}

// Decode the given bytes to a message
func Decode(b []byte) (decoder.Message, error) {
	return decoder.Decode(b)
}
