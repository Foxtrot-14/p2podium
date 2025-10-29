package scraper

import (
	"log"
	"net"
	"time"
	"github.com/anacrolix/go-libutp"
	"github.com/Foxtrot-14/p2podium/dht"
)

func (s *Scraper) GetMetaData(peer dht.Peer, sock *utp.Socket) {
	
	remoteAddr := &net.UDPAddr{
		IP:   peer.IP,
		Port: int(peer.Port),
	}

	log.Printf("[INFO] Connecting via uTP to %s:%d", peer.IP, peer.Port)

	conn, err := sock.Dial(remoteAddr.String())
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
