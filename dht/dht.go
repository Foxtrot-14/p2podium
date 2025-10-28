package dht

import (
	"net"
	"sync"
	"time"
)

type PeerRetriever interface {
	GetPeerList() //Gateway into the package
	JoinDHT()
	HealthCheck()
	GetPeers()
	HandleNewNodes()
	AnnouncePresence()
	Listen()
}

type Node struct {
	NodeID   [20]byte
	Address  net.IP
	Port     uint16
	LastSeen time.Time
}

type Peer struct {
	IP   net.IP
	Port uint16
}

type RoutingTable struct {
	Buckets    [160][]Node
	BucketLock [160]sync.Mutex
	Peers      map[[20]byte][]string
}

type DHT struct {
	NodeID   [20]byte
	Table    *RoutingTable
	InfoHash [20]byte
	PeerChan chan Peer
	Done     chan struct{}
	Errc     chan error
}
