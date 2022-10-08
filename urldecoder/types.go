package urldecoder

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/micahjmartin/gopro/decoder"
)

/*
func UrlType(t string) string {
	switch t {
	case "B":
		return decoder.TYPE_BYTES
	case "b":
		return decoder.TYPE_BOOL
	case "d":
		return decoder.TYPE_DOUBLE
	case "e":
		ENUM
	case "f":
		return decoder.TYPE_FLOAT

	case "g":
		return decoder.TYPE_SFIXED32
	case "h":
		return decoder.TYPE_SFIXED64
	case "i":
		return decoder.TYPE_INT32
	case "j":
		return decoder.TYPE_INT64
	case "m":
		return decoder.TYPE_MESSAGE
	case "n":
		return decoder.TYPE_SINT32
	case "o":
		return decoder.TYPE_SINT64
	case "s":
		return decoder.TYPE_STRING
	case "u":
		return decoder.TYPE_UINT32
	case "v":
		return decoder.TYPE_UINT64
	case "x":
		return decoder.TYPE_FIXED32
	case "y":
		return decoder.TYPE_FIXED64
	case "z":
		return decoder.SPECIAL_TYPE_BASE64
	}
	return ""
}
*/

func setVarint(value string, i *decoder.Item, signed bool) error {
	if signed {
		//sint64
		i.WireType = 0
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		i.Raw = make([]byte, 8)
		v = binary.PutVarint(i.Raw, int64(v))
		i.Raw = i.Raw[:v]
		return nil
	}
	i.WireType = 0
	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	i.Raw = make([]byte, 8)
	v = binary.PutUvarint(i.Raw, uint64(v))
	i.Raw = i.Raw[:v]
	return nil
}

func ToItem(fn int, typ, value string) (decoder.Item, error) {
	i := decoder.Item{
		FieldNumber: fn,
	}
	switch typ {
	// Non-zigzag varints
	case "b":
		i.SetType("bool")
		setVarint(value, &i, false)
	case "i":
		i.SetType("int32")
		setVarint(value, &i, false)
	case "j":
		i.SetType("int64")
		setVarint(value, &i, false)
	case "u":
		i.SetType("uint32")
		setVarint(value, &i, false)
	case "v":
		i.SetType("uint64")
		setVarint(value, &i, false)
	case "e":
		i.SetType("enum")
		setVarint(value, &i, false)
	case "n":
		i.SetType("sint32")
		setVarint(value, &i, true)
	case "o":
		i.SetType("sint64")
		setVarint(value, &i, true)
	case "d":
		// Double
		i.SetType("double")
		d, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return i, err
		}
		i.WireType = 1
		i.Raw = make([]byte, 8)
		binary.LittleEndian.PutUint64(i.Raw, uint64(d))

	case "f":
		// Float
		i.SetType("float")
		d, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return i, err
		}
		i.WireType = 5
		i.Raw = make([]byte, 4)
		binary.LittleEndian.PutUint32(i.Raw, uint32(d))
	case "x":
		// fixed32
		i.SetType("fixed32")
		d, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return i, err
		}
		i.WireType = 5
		i.Raw = make([]byte, 4)
		binary.LittleEndian.PutUint32(i.Raw, uint32(d))
	case "y":
		// fixed64
		i.SetType("fixed64")
		d, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return i, err
		}
		i.WireType = 1
		i.Raw = make([]byte, 8)
		binary.LittleEndian.PutUint64(i.Raw, uint64(d))
	case "z":
		// base64
		i.WireType = 2
		var err error
		i.Raw, err = base64.URLEncoding.DecodeString(value)
		if err != nil {
			return i, fmt.Errorf("error decoding base64: %s", err)
		}
	case "s":
		i.WireType = 2
		i.SetType("string")
		i.Raw = []byte(value)
	default:
		return i, fmt.Errorf("unknown type '%s'", typ)
	}

	if len(i.Raw) == 0 && len(value) != 0 {
		return i, fmt.Errorf("%d%s%s: no value parsed, but value is specified", fn, typ, value)
	}
	return i, nil
}
