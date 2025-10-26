package scraper

import (
	"sync"
)

type TorrentScraper interface {
	StartScraper()
	PieceDownloader()
	WriteFile()
}

type Piece struct {
	Index int
	Data  []byte
}

type Scraper struct {
	PeerList         []string
	InfoHash         [20]byte
	PendingPieces    []int
	DownloadedPieces []int
	PieceChan        chan Piece
	TableLock        sync.Mutex
}

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

