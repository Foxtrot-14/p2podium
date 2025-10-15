package dht

import (
	"bytes"
	"github.com/jackpal/bencode-go"
)

func parseJoinDHTResponse(res []byte) (DHTResponse, error) {
	var msg DHTResponse

	reader := bytes.NewReader(res)
	if err := bencode.Unmarshal(reader, &msg); err != nil {
		return DHTResponse{}, err
	}

	return msg, nil
}
