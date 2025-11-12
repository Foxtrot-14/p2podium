package scraper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
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

func readExtendedHandshake(conn net.Conn) (*ExtendedHandshake, error) {
	var length int32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to read length prefix: %w", err)
	}

	if length <= 0 || length > 64*1024 {
		return nil, fmt.Errorf("invalid message length: %d", length)
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, fmt.Errorf("failed to read message body: %w", err)
	}

	if len(buf) < 2 {
		return nil, fmt.Errorf("message too short")
	}
	if buf[0] != 20 {
		return nil, fmt.Errorf("expected message ID 20 (extended), got %d", buf[0])
	}
	if buf[1] != 0 {
		return nil, fmt.Errorf("expected extended handshake (ext_id=0), got %d", buf[1])
	}
	payload := buf[2:]

	var handshake ExtendedHandshake
	if err := bencode.Unmarshal(bytes.NewReader(payload), &handshake); err != nil {
		return nil, fmt.Errorf("bencode decode error: %w", err)
	}

	return &handshake, nil
}

func readNextMessage(conn net.Conn) ([]byte, error) {
	var length int32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}

	if length <= 0 || length > 64*1024 {
		return nil, fmt.Errorf("invalid message length: %d", length)
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, fmt.Errorf("failed to read message body: %w", err)
	}

	return buf, nil
}

func (s *Scraper) RequestMetaData(conn net.Conn) {
	req, _ := createExtendedHandshake()
	conn.Write(req)

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	for {
		msg, err := readNextMessage(conn)
		if err != nil {
			s.metaRequested.CompareAndSwap(true, false)
			log.Printf("[ERROR] from %s: %v", conn.RemoteAddr(), err)
			return
		}

		if len(msg) < 2 {
			continue
		}

		id := msg[0]
		if id != 20 {
			log.Printf("[DEBUG] Skipping message id=%d, len=%d", id, len(msg))
			s.metaRequested.CompareAndSwap(true, false)
			return
		}

		extID := msg[1]
		if extID == 0 {
			var handshake ExtendedHandshake
			if err := bencode.Unmarshal(bytes.NewReader(msg[2:]), &handshake); err != nil {
				log.Printf("[ERROR] bencode decode failed: %v", err)
				continue
			}

			if _, ok := handshake.M["ut_pex"]; ok {
				s.metaRequested.CompareAndSwap(true, false)
				return
			}
			log.Printf("[INFO] Got Extended Handshake: %+v", handshake)
		} else {
			log.Printf("[DEBUG] Got extended msg %d, skipping", extID)
		}
	}
	//Not reponsponding to Extended Handshake
}
