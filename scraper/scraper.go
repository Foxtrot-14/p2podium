package scraper

import (
	"sync"
	"sync/atomic"

	"github.com/Foxtrot-14/p2podium/dht"
)

type TorrentScraper interface {
	StartScraper() //Gateway
	PieceDownloader()
	WriteFile()
	HandShake()
	GetMetaData()
}

type Piece struct {
	Index int
	Data  []byte
}

// 16384 ~ 16KB
type File struct {
	Name string
	Size int64
}

type Directory struct {
	Name        string
	Files       []File
	Directories []*Directory
}

type Torrent struct {
	Name string
	Root *Directory
}

type Scraper struct {
	PeerID           [20]byte
	PeerChan         chan dht.Peer
	ActivePeers      sync.Map
	InfoHash         [20]byte
	Torrent          Torrent
	metaRequested    atomic.Bool
	PendingPieces    []int
	DownloadedPieces []int
	PieceChan        chan Piece
	TableLock        sync.Mutex
}

type ExtendedHandshake struct {
	M            map[string]int `bencode:"m"`
	MetadataSize int            `bencode:"metadata_size"`
	Version      string         `bencode:"v"`
	ReqQ         int            `bencode:"reqq"`
}
