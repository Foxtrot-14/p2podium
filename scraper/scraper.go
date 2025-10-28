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

//16384 ~ 16KB

type Torrent struct {

}

type Scraper struct {
	PeerID			 [20]byte
	PeerChan         chan dht.Peer
	ActivePeers      map[string]int
	InfoHash         [20]byte
	Torrent          Torrent
	PendingPieces    []int
	DownloadedPieces []int
	PieceChan        chan Piece
	TableLock        sync.Mutex
}
