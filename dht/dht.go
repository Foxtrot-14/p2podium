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
	HandleNewNodes(nodes []Node)
	AnnouncePresence()
	Listen()
}

type Node struct {
	NodeID   [20]byte
	Address  net.IP
	Port     uint16
	LastSeen time.Time
}

type RoutingTable struct {
	Buckets   map[int][]Node
	Peers     map[[20]byte][]string
	TableLock *sync.Mutex
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
