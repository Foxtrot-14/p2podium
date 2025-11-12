package scraper

func (s *Scraper) StartScraper() {
	
	// go StartDownloader()
		
	for peer := range s.PeerChan {
		//Check Peer Health
		go s.Handshake(peer)
	}
}
