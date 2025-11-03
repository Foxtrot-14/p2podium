package scraper

import (
	"io"
	"log"
	"net"

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
	handshakeMsg := Handshake(s.InfoHash, s.PeerID)

	if _, err := conn.Write(handshakeMsg); err != nil {
		return false
	}

	resp := make([]byte, 68)
	if _, err := io.ReadFull(conn, resp); err != nil {
		return false
	}

	log.Printf("[INFO] Response from handshake: %q", resp[:])

	return string(resp[28:48]) == string(s.InfoHash[:])
}
