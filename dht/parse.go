package dht

import (
	"bytes"
	"github.com/jackpal/bencode-go"
)

func parseJoinDHTResponse(res []byte) (DHTMessage, error) {
	var msg DHTMessage

	reader := bytes.NewReader(res)
	if err := bencode.Unmarshal(reader, &msg); err != nil {
		return DHTMessage{}, err
	}

	return msg, nil
}

