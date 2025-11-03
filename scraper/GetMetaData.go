package scraper

import (
	"log"
	"net"
	"time"
	"github.com/Foxtrot-14/p2podium/dht"
)

func (s *Scraper) GetMetaData(peer dht.Peer) {
	
	remoteAddr := &net.UDPAddr{
		IP:   peer.IP,
		Port: int(peer.Port),
	}

	log.Printf("[INFO] Connecting via TCP to %s:%d", peer.IP, peer.Port)

	conn, err := net.DialTimeout("tcp", remoteAddr.String(), 5*time.Second)
	if err != nil {
		log.Printf("[ERROR] uTP dial failed: %v", err)
		return
	}

	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	if !s.SendHandshake(conn, peer) {
		log.Printf("[WARN] Handshake failed with %s", peer.IP)
		return
	}

	log.Printf("[INFO] Handshake successful with %s", peer.IP)

	// Next steps: request metadata, parse response, etc.
	// s.RequestMetaData(conn)
}
