package dht

import (
	"sync"
	"time"
)

type PeerRetriever interface {
	GetPeerList() //Gateway into the package
	AnnouncePresence()
	Repopulate() //Should re-populate the DHT
}

type Node struct {
	NodeID   [20]byte
	Address  string
	LastSeen time.Time
}

type KBucket struct {
	Nodes []Node
	mu    sync.Mutex
}

type RoutingTable struct {
	Buckets map[int]*KBucket
	Peers   map[string][]string
	mu      sync.Mutex
}

type DHT struct {
	NodeID   [20]byte
	Table    *RoutingTable
	InfoHash [20]byte
	Peers    []string
	PeerLock *sync.Mutex
	Done     chan struct{}
	Errc     chan error
}
