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

// A generic protobuf item
type Item struct {
	WireType    int
	FieldNumber int
	Raw         []byte
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
func readNextItem(b *Buffer) (item Item, err error) {
	fieldNumber, wireType, err := b.ReadKey()
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
func readMessage(b *Buffer) (message Message, err error) {
	message = make(Message)
	for !b.Empty() {
		item, err := readNextItem(b)
		if err != nil && IsProtobufError(err) != nil {
			return nil, err
		}

		var i interface{}
		if item.WireType == 2 {
			i = getType2(item)
		} else {
			i = item.String()
		}

		// Check if there are multiple items
		if _, ok := message[item.FieldNumber]; !ok {
			message[item.FieldNumber] = make([]interface{}, 0)
		}
		message[item.FieldNumber] = append(message[item.FieldNumber], i)

	}
	return message, nil
}

// Wire type 2 can be either Bytes, Str, or a submessage
func getType2(i Item) interface{} {
	if i.WireType != 2 {
		return nil
	}
	msg, err := readMessage(NewBuffer(i.Raw))
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
	return readMessage(NewBuffer(b))
}
