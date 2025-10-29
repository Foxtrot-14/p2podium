package scraper

import (
	"log"

	utp "github.com/anacrolix/go-libutp"
)

func (s *Scraper) StartScraper() {

	sock, err := utp.NewSocket("udp", ":6881")
	if err != nil {
		log.Printf("[ERROR] Creating uTP socket: %v", err)
		return
	}

	defer sock.Close()

	for peer := range s.PeerChan {
		//TODO if torrent empty get meta data
		// if len(s.Torrent) == 0 {
		go s.GetMetaData(peer, sock)
		// }

		//TODO check if torrent already in use
		if _, ok := s.ActivePeers[peer.IP.String()]; ok {
			continue
		}

		//TODO if torrent not empty spawn a go routine
		s.ActivePeers[peer.IP.String()] += 1
	}
}
