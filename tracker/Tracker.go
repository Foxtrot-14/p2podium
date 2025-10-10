package tracker

import (
	"sync"
	"time"
)

type PeerTracker interface {
	GetPeerList() //Gateway into the package
	sendConnectionRequest() (connectionID uint64, err error)
	sendAnnounceRequest()
}

type Tracker struct {
	PeerID           [20]byte
	InfoHash         string
	Trackers         []string
	TransactionID    uint32
	ConnectionID     uint64
	TrackerIdx       int
	ConnectionIDTime time.Time
	Peers            []string
	PeersTime        time.Duration
	PeerLock         *sync.Mutex
	Done             chan struct{}
	Errc             chan error
}

type ConnectRequest struct {
	ProtocolID    uint64
	Action        uint32
	TransactionID uint32
}

type AnnounceRequest struct {
	ConnectionID  uint64
	Action        uint32
	TransactionID uint32
	InfoHash      [20]byte
	PeerID        [20]byte
	Downloaded    uint64
	Left          uint64
	Uploaded      uint64
	Event         uint32
	IPAddress     uint32
	Key           uint32
	NumWant       int32
	Port          uint16
}
