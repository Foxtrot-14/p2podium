package scraper

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	"github.com/jackpal/bencode-go"
)

func createExtendedHandshake() ([]byte, error) {
	payload := map[string]any{
		"m": map[string]int{
			"ut_metadata": 1,
		},
	}

	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, payload); err != nil {
		return nil, err
	}
	payloadBytes := buf.Bytes()

	msg := new(bytes.Buffer)
	msg.WriteByte(20)
	msg.WriteByte(0)
	msg.Write(payloadBytes)

	fullMsg := new(bytes.Buffer)
	binary.Write(fullMsg, binary.BigEndian, int32(msg.Len()))
	fullMsg.Write(msg.Bytes())

	return fullMsg.Bytes(), nil
}

func (s *Scraper) RequestMetaData(conn net.Conn) {
	req, _ := createExtendedHandshake()
	conn.Write(req)

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	//Not reponsponding to Extended Handshake
}
