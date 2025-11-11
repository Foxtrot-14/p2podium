package scraper

func (s *Scraper) StartScraper() {

	for peer := range s.PeerChan {
		//Check Peer Health
		go s.Handshake(peer)
	}
}
