package scraper

import (
	"sync"

	"github.com/Foxtrot-14/p2podium/dht"
)

type TorrentScraper interface {
	StartScraper()
	PieceDownloader()
	WriteFile()
	GetMetaData()
}

type Piece struct {
	Index int
	Data  []byte
}

type Torrent struct {
}

type Scraper struct {
	PeerList         []dht.Peer
	ActivePeers      []dht.Peer
	InfoHash         [20]byte
	Torrent          Torrent
	PendingPieces    []int
	DownloadedPieces []int
	PieceChan        chan Piece
	TableLock        sync.Mutex
}
