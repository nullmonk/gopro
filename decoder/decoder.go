package decoder

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"
	"unicode"
)

var INDENT = "  "

type Message map[int][]interface{}

func (m Message) Add(item Item) {
	var i interface{}
	if item.WireType == 2 {
		i = getType2(item)
	} else {
		i = item.String()
	}

	// Check if there are multiple items
	if _, ok := m[item.FieldNumber]; !ok {
		m[item.FieldNumber] = make([]interface{}, 0, 1)
	}
	m[item.FieldNumber] = append(m[item.FieldNumber], i)
}

// A generic protobuf item
type Item struct {
	WireType    int
	FieldNumber int
	Raw         []byte
}

func (i *Item) Dump() {
	fmt.Printf("%d (%d) = %s\n", i.FieldNumber, i.WireType, i.String())
}

func (i *Item) String() string {
	switch i.WireType {
	case 0:
		v, _ := binary.Uvarint(i.Raw)
		return fmt.Sprintf("%d", v)
	case 1:
		return fmt.Sprintf("0x%s", hex.EncodeToString(i.Raw))
	case 2:
		if IsString(i.Raw) {
			return string(i.Raw)
		} else {
			// Do something different for bytes?
			return fmt.Sprintf("0x%s", hex.EncodeToString(i.Raw))
		}
	case 5:
		return fmt.Sprintf("0x%s", hex.EncodeToString(i.Raw))
	}
	return ""
}

// read a protobuf item from the buffer
func ReadNextItem(b *Buffer) (item Item, err error) {
	fieldNumber, wireType, err := b.ReadKey()
	if fieldNumber == 0 || wireType > 5 {
		return Item{}, b.Error(fmt.Errorf("Invalid Field Number (%d) or Bad Wiretype (%d)", fieldNumber, wireType))
	}
	if err != nil {
		return Item{}, err
	}

	item = Item{
		WireType:    wireType,
		FieldNumber: int(fieldNumber),
	}

	switch wireType {
	case 0:
		item.Raw, err = b.ReadVarintBytes()
	case 1:
		item.Raw = make([]byte, 8)
		_, err = b.Read(item.Raw)
	case 2:
		item.Raw, err = b.ReadLenDelim()
	case 3:
		fallthrough
	case 4:
		return Item{}, b.Error(fmt.Errorf("Groups detected. Parser does not handle groups"))
	case 5:
		item.Raw = make([]byte, 4)
		_, err = b.Read(item.Raw)
	}
	if err != nil {
		return Item{}, b.Error(err)
	}
	return item, nil
}

// Read an entire message from a buffer
func ReadMessage(b *Buffer) (message Message, err error) {
	message = make(Message)
	for !b.Empty() {
		item, err := ReadNextItem(b)
		if err != nil && IsProtobufError(err) != nil {
			return nil, err
		}

		message.Add(item)
	}
	return message, nil
}

// Wire type 2 can be either Bytes, Str, or a submessage
func getType2(i Item) interface{} {
	if i.WireType != 2 {
		return nil
	}
	msg, err := ReadMessage(NewBuffer(i.Raw))
	if err == nil {
		return msg
	}
	// String
	if IsString(i.Raw) {
		return string(i.Raw)
	}
	// Bytes
	return i.Raw
}

// Dump a message to stdout
func DumpMessage(m Message, indent string) {
	// Sort the keys
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		for _, v := range m[k] {
			if sm, ok := v.(Message); ok {
				fmt.Printf("%s%d:\n", indent, k)
				DumpMessage(sm, indent+INDENT)
			} else {
				fmt.Printf("%s%d: %v\n", indent, k, v)
			}
		}
	}
}

// Return true if the given slice is a string
func IsString(b []byte) bool {
	for _, l := range b {
		if !unicode.IsPrint(rune(l)) {
			return false
		}
	}
	return true
}

// Decode the given bytes to a message
func Decode(b []byte) (Message, error) {
	return ReadMessage(NewBuffer(b))
}

/* Reading a varint in reverse is actually quite difficult, luckily, we know the byte before is */

/* Read a Varint in reverse starting at index i. N is the number of bytes read */
func ReadVarintReverse(buf []byte, idx int) (uint64, int) {
	if buf[idx] >= 0x80 {
		return 0, 0
	}

	varint := make([]byte, 0, 2)
	for idx >= 0 {
		// Read a byte, add it to the result
		byt := buf[idx] // Always read the first byte
		if byt < 0x80 && len(varint) != 0 {
			break
		}
		varint = append([]byte{byt}, varint...)
		idx--
	}

	return binary.Uvarint(varint)
}

/* Read a Varint in reverse starting at index i until wiretype X is read. N is the number of bytes read */
func ReadVarintReverseWiretype(buf []byte, idx, wiretype int) (uint64, int) {
	if buf[idx] >= 0x80 {
		return 0, 0
	}

	varint := make([]byte, 0, 2)
	for idx >= 0 {
		// Read a byte, add it to the result
		byt := buf[idx] // Always read the first byte

		if byt < 0x80 && len(varint) != 0 {
			break
		}
		varint = append([]byte{byt}, varint...)
		idx--
		// we have the wiretype that we want, stop
		if int(byt&0x7) == wiretype {
			break
		}
	}

	return binary.Uvarint(varint)
}
