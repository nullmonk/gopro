package gopro

import (
	"io"

	"github.com/micahjmartin/gopro/decoder"
	"github.com/micahjmartin/gopro/urldecoder"
)

// DumpMessage dumps a Message to stdout
func DumpMessage(m decoder.Message) {
	decoder.DumpMessage(m, "")
}

// Decode the given bytes to a message
func Decode(b []byte) (decoder.Message, error) {
	return decoder.Decode(b)
}

func DecodeUrl(pb string) (decoder.Message, error) {
	return urldecoder.Decode(pb)
}

func Encode(m decoder.Message, w io.Writer) error {
	_, err := m.Encode(0, w)
	return err
}
