package scraper

func (s *Scraper) StartScraper() {
	for peer := range s.PeerChan {
		//TODO if torrent empty get meta data
		// if len(s.Torrent) == 0 {
		s.GetMetaData(peer)
		// }

		//TODO check if torrent already in use
		if _, ok := s.ActivePeers[peer.IP.String()]; ok {
			continue
		}

		//TODO if torrent not empty spawn a go routine
		s.ActivePeers[peer.IP.String()] += 1
	}
}
