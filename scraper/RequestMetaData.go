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

func (s *Scraper) RequestMetaData(conn net.Conn) {
	time.Sleep(100 * time.Millisecond)

	req, err := createExtendedHandshake()
	if err != nil {
		log.Printf("[ERROR] while creating extended handshake")
		return
	}

	_, err = conn.Write(req)
	if err != nil {
		log.Printf("[ERROR] while sending request")
		return
	}
	log.Println("[DEBUG] Sent extended handshake, waiting for response...")

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("[ERROR] failed to read raw response: %v", err)
		return
	}

	log.Printf("[DEBUG] Received %d bytes", n)
	log.Printf("[DEBUG] Raw response bytes: %x", buf[:n])
	log.Printf("[DEBUG] Raw response string: %q", string(buf[:n]))
	// resp, err := readExtendedHandshake(conn)
	// if err != nil {
	// 	log.Printf("[ERROR] failed to read extended handshake response: %v", err)
	// 	return
	// }
	//
	// log.Printf("[INFO] Peer supports ut_metadata: %d", resp.M["ut_metadata"])
	// log.Printf("[INFO] Metadata size: %d bytes", resp.MetadataSize)
	// log.Printf("[INFO] Client version: %s", resp.Version)
}
