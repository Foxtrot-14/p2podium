package scraper

import (
	"log"
	"net"
	"time"

	"github.com/Foxtrot-14/p2podium/dht"
)

func (s *Scraper) GetMetaData(peer dht.Peer) {
	addr := net.TCPAddr{
		IP:   peer.IP,
		Port: int(peer.Port),
	}

	conn, err := net.DialTimeout("tcp", addr.String(), 5*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	if !s.SendHandshake(conn, peer) {
		return
	}

	log.Printf("[INFO] Handshake successful with %s", peer.IP)

	//Send meta data request
	//Parse and write to s.torrent
	//Remove from ActivePeers
}
