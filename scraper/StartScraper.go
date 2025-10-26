package scraper

func (s *Scraper) StartScraper() {
	for _, peer := range s.PeerList {
		go s.PieceDownloader(peer)
	}
	go s.WriteFile()
}
