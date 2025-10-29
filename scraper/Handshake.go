package scraper

import (
	"io"
	"net"

	"github.com/Foxtrot-14/p2podium/dht"
)

func Handshake(peer dht.Peer, infohash [20]byte, peerID [20]byte) []byte {
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
	handshakeMsg := Handshake(peer, s.InfoHash, s.PeerID)

	if _, err := conn.Write(handshakeMsg); err != nil {
		return false
	}

	resp := make([]byte, 68)
	if _, err := io.ReadFull(conn, resp); err != nil {
		return false
	}

	if string(resp[28:48]) != string(s.InfoHash[:]) {
		return false
	}

	return true
}
