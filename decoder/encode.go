package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
)

func (i Item) Encode(w io.Writer) (bytes int, err error) {
	header := uint8(i.FieldNumber)<<3 | uint8(i.WireType)
	bytes, err = w.Write([]byte{header})
	if err != nil {
		return
	}
	b := 0
	fmt.Fprintf(os.Stderr, "%d %d (%x) %v (%x)\n", i.FieldNumber, i.WireType, header, i.Raw, i.Raw)

	if i.WireType == 2 {
		varint := make([]byte, 8)
		i := binary.PutUvarint(varint, uint64(len(i.Raw)))
		varint = varint[:i]
		b, err = w.Write(varint)
		if err != nil {
			return
		}
		bytes += b
	}
	b, err = w.Write(i.Raw)
	bytes += b
	return
}

func (m Message) Encode(fn int, w io.Writer) (count int, err error) {
	// Sort the keys
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	buf := new(bytes.Buffer)
	var length, l int
	for _, k := range keys {
		for _, v := range m[k] {
			if sm, ok := v.(Message); ok {
				l, err = sm.Encode(k, buf)
				if err != nil {
					return
				}
				length += l
			} else if sm, ok := v.(Item); ok {
				l, err = sm.Encode(buf)
				if err != nil {
					return
				}
				length += l
			}
		}
	}

	// Write the header if the FieldNumber is not 0
	if fn > 0 {
		header := uint8(fn<<3 | 2)
		l, err = w.Write([]byte{header})
		if err != nil {
			return
		}
		count += l
		varint := make([]byte, 8)
		i := binary.PutUvarint(varint, uint64(length))
		varint = varint[:i]
		l, err = w.Write(varint)
		if err != nil {
			return
		}
		count += l
	}

	l, err = w.Write(buf.Bytes())
	count += l
	return
}
