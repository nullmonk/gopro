package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Custom bytes.Reader which has useful functions for protobuf
type Buffer struct {
	r   *bytes.Reader
	idx int
}

func NewBuffer(b []byte) *Buffer {
	return &Buffer{
		r: bytes.NewReader(b),
	}
}

func (b *Buffer) Index() int {
	return b.idx
}

func (b *Buffer) Empty() bool {
	return int64(b.idx) >= b.r.Size()
}

/* ReadVarintBytes reads a single varint from the buffer and returns the bytes.
   we could use the binary.ReadUvarint method, but then we would lose track of the index
*/
func (b *Buffer) ReadVarintBytes() ([]byte, error) {
	buf := make([]byte, 0, 4)
	for {
		bt, err := b.r.ReadByte()
		b.idx++
		if err != nil {
			return buf, err
		}
		buf = append(buf, bt)
		if bt < 0x80 {
			break
		}
	}
	return buf, nil
}

// Read a varint from the buffer
func (b *Buffer) ReadVarint() (uint64, error) {
	buf, err := b.ReadVarintBytes()
	if err != nil {
		return 0, nil
	}
	result, _ := binary.Uvarint(buf)
	return result, nil
}

// Read a protobuf key from the buffer
func (b *Buffer) ReadKey() (fieldNumber uint64, wireType int, err error) {
	v, err := b.ReadVarint()
	if err != nil {
		return 0, 0, err
	}
	fieldNumber = uint64(v >> 3)
	wireType = int(v & 0x7)
	return
}

// Read a length delimited item from the buffer
func (b *Buffer) ReadLenDelim() ([]byte, error) {
	ln, err := b.ReadVarint()
	if err != nil {
		return nil, err
	}
	if ln == 0 {
		return nil, fmt.Errorf("String should not be raw")
	}
	res := make([]byte, ln)
	n, err := b.r.Read(res)
	b.idx += n
	if err != nil {
		if err.Error() == "EOF" {
			return nil, fmt.Errorf("EOF: Reading len delim")

		}
		return res, err
	}
	if uint64(n) != ln {
		return nil, fmt.Errorf("EOF: Reading len delim")
	}
	return res, nil
}

func (b *Buffer) Read(buf []byte) (int, error) {
	b.idx += len(buf)
	return b.r.Read(buf)
}

// Return an error which marks the proper index of the error
func (b *Buffer) Error(i interface{}) *ProtobufDecodeError {
	if p, ok := i.(*ProtobufDecodeError); ok {
		p.i = b.Index()
		return p
	}
	if e, ok := i.(error); ok {
		return &ProtobufDecodeError{
			i: b.Index(),
			s: e.Error(),
		}
	}
	return &ProtobufDecodeError{
		i: b.Index(),
		s: fmt.Sprint(i),
	}
}
