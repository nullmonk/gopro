package decoder

import "fmt"

type ProtobufDecodeError struct {
	s string // Error message
	i int    // Index of error
}

func (p *ProtobufDecodeError) Error() string {
	if p.i != -255 {
		return fmt.Sprintf("Protobuf Decode Error at Index %d: %s", p.i, p.s)
	}
	return fmt.Sprintf("Protobuf Decode Error: %s", p.s)
}

func IsProtobufError(e error) *ProtobufDecodeError {
	if p, ok := e.(*ProtobufDecodeError); ok {
		return p
	}
	return nil
}
