package gopro

import (
	"bytes"
	"encoding/binary"
	"io"
)

var WIRE_TYPES = map[int]string{
	0: "Varint",
	1: "64bit",
	2: "Len-delim",
	3: "Group Start",
	4: "Group End",
	5: "32bit",
}

/*
message Data {
    string id = 1;
    float lat = 2;
    float lon = 3;
    CryptoKey key = 4;
    repeated string ssids = 5;
}

0 	Varint 	int32, int64, uint32, uint64, sint32, sint64, bool, enum
1 	64-bit 	fixed64, sfixed64, double
2 	Length-delimited 	string, bytes, embedded messages, packed repeated fields
5 	32-bit 	fixed32, sfixed32, float
*/

type CryptoKey struct{}

type TestMessage struct {
	Id    string    `FieldNumber:"1" Name:"id"`
	Lat   float64   `FieldNumber:"2" Name:"lat"`
	Lon   float64   `FieldNumber:"3" Name:"lon"`
	Key   CryptoKey `FieldNumber:"4" Name:"key"`
	SSIDs []string  `FieldNumber:"5" Name:"ssids"`
	Test  int
}

func SerializeVarint(w bytes.Buffer, key []byte, value int) {
	w.WriteByte(0x08) // Write the key
	buf := make([]byte, 10)
	n := binary.PutUvarint(buf, 1)
	w.Write(buf[:n])
}

func (t *TestMessage) ProtobufSize() int {
	return 0
}

func SerializeMessage(m interface{}) ([]byte, error) {
	//result := make([]byte, 0, 100)
	return nil, nil
}

// WriteVarint writes a uint64 to a ByteWriter
func WriteVarint(w io.ByteWriter, x uint64) int {
	i := 0
	for x >= 0x80 {
		w.WriteByte(byte(x) | 0x80)
		x >>= 7
		i++
	}
	w.WriteByte(byte(x))
	return i + 1
}

func SerializeLengthDelimited(fieldNumber int, buf []byte) ([]byte, error) {
	varint := uint64(fieldNumber << 3 & 0x2)
	result := &bytes.Buffer{}
	WriteVarint(result, varint)
	WriteVarint(result, uint64(len(buf)))
	result.Write(buf)
	return result.Bytes(), nil
}
