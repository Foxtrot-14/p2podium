package scraper

func (s *Scraper) WriteFile() {
	for piece := range s.PieceChan {
		if verifyPiece(piece) {
			writePieceToDisk(piece)
		}
	}
}

func downloadPieceFromPeer(peer string, index int) []byte {
	return []byte{}
}

func verifyPiece(p Piece) bool {
	return true
}

func writePieceToDisk(p Piece) {}
