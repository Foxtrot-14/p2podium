package scraper

import (
	"bytes"
	"io"
	"net"
	"time"

	"github.com/Foxtrot-14/p2podium/dht"
)

func Handshake(infohash [20]byte, peerID [20]byte) []byte {
	pstr := "BitTorrent protocol"
	buf := make([]byte, len(pstr)+49)
	buf[0] = byte(len(pstr))
	curr := 1
	curr += copy(buf[curr:], pstr)
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], infohash[:])
	curr += copy(buf[curr:], peerID[:])
	return buf
}

func (s *Scraper) SendHandshake(conn net.Conn, peer dht.Peer) bool {
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	handshakeMsg := Handshake(s.InfoHash, s.PeerID)
	if _, err := conn.Write(handshakeMsg); err != nil {
		return false
	}

	pstrlen := make([]byte, 1)
	if _, err := io.ReadFull(conn, pstrlen); err != nil {
		return false
	}

	total := int(pstrlen[0]) + 49 - 1
	resp := make([]byte, total)
	resp[0] = pstrlen[0]

	if _, err := io.ReadFull(conn, resp[1:]); err != nil {
		return false
	}

	if !bytes.Equal(resp[28:48], s.InfoHash[:]) {
		return false
	}

	return true
}
