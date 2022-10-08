package urldecoder

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/micahjmartin/gopro/decoder"
)

var re = regexp.MustCompile(`(\d+)([A-Za-z])(.*)`)

// Decode URL protobuf format and reencode it to pb

var Delimiter = "!"

func Decode(pb string) (decoder.Message, error) {
	items := strings.Split(pb, Delimiter)
	return decodeMessage(items[1:])
}

func decodeMessage(items []string) (decoder.Message, error) {
	rootMessage := decoder.Message{}
	i := 0
	for i < len(items) {
		matches := re.FindStringSubmatch(items[i])
		fn, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid field number: %s", matches[1])
		}

		if matches[2] == "m" {
			i++
			continue
			// Parse submessage
			msgLen, err := strconv.Atoi(matches[3])
			if err != nil {
				return nil, fmt.Errorf("invalid message length: %s", matches[3])
			}
			i++
			msg, err := decodeMessage(items[i : i+msgLen])
			if err != nil {
				return nil, err
			}
			rootMessage.AddSubmessage(fn, msg)
			i += msgLen
			continue
		}
		item, err := ToItem(fn, matches[2], matches[3])
		if err != nil {
			return nil, err
		}
		fmt.Fprint(os.Stderr, "!"+matches[0])
		rootMessage.Add(item)
		i += 1
	}
	fmt.Fprintln(os.Stderr, "")
	return rootMessage, nil
}
