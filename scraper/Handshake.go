package scraper

import (
	"github.com/Foxtrot-14/p2podium/dht"
	"log"
	"net"
	"time"
)

func (s *Scraper) Handshake(peer dht.Peer) {

	remoteAddr := &net.UDPAddr{
		IP:   peer.IP,
		Port: int(peer.Port),
	}

	conn, err := net.DialTimeout("tcp", remoteAddr.String(), 5*time.Second)
	if err != nil {
		return
	}

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	if !s.SendHandshake(conn, peer) {
		log.Printf("[WARN] Handshake failed with %s", peer.IP)
		conn.Close()
		return
	}

	log.Printf("[INFO] Handshake successful with %s", peer.IP)

	//Check if metadata is already fetched
	if s.Torrent.Root == nil && s.metaRequested.CompareAndSwap(false, true) {
		go s.RequestMetaData(conn)
	}

	// Next steps: add to ActivePeers
	if old, ok := s.ActivePeers.Load(remoteAddr.String()); ok {
		if oldConn, ok := old.(net.Conn); ok {
			oldConn.Close()
		}
	}

	s.ActivePeers.Store(remoteAddr.String(), conn)
}
