package state

import (
	"github.com/Foxtrot-14/p2podium/dht"
)

func LoadDHT() (*dht.DHT, error) {
	path, err := getStateFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("state file does not exist")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sDHT SerializableDHT
	if err := json.Unmarshal(data, &sDHT); err != nil {
		return nil, err
	}

	dht := &dht.DHT{
		NodeID:   sDHT.NodeID,
		Table:    &RoutingTable{Buckets: make(map[int]*KBucket)},
		InfoHash: sDHT.InfoHash,
		Peers:    sDHT.Peers,
		Done:     make(chan struct{}),
		Errc:     make(chan error),
	}

	for k, v := range sDHT.Table.Buckets {
		nodes := make([]Node, len(v.Nodes))
		for i, sn := range v.Nodes {
			nodes[i] = Node{
				NodeID:   sn.NodeID,
				Address:  sn.Address,
				LastSeen: time.Unix(sn.LastSeen, 0),
			}
		}
		dht.Table.Buckets[k] = &KBucket{Nodes: nodes}
	}

	return dht, nil
}
