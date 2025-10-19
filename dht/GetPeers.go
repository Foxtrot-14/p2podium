package dht

import (
	"log"
)

func (d *DHT) GetPeers(nodes []Node) {
	buckets := make(map[int]int)
	for _, node := range nodes {
		msg := DHTRequest{
			T: generateTransactionID(),
			Y: "q",
			Q: "get_peers",
			A: map[string][]byte{
				"id":        d.NodeID[:],
				"info_hash": d.InfoHash[:],
			},
		}

		resp, err := SendRequest(msg, node)
		if err != nil {
			log.Printf("%s", err)
			return
		}

		nodes := []byte(resp.R.Nodes)
		for i := 0; i < len(nodes); i += 26 {
			var nodeID [20]byte
			copy(nodeID[:], nodes[i:i+20])
			bucketNumber, err := GetBucketNumber(nodeID[:], d.NodeID[:])
			if err != nil {
				log.Printf("[ERROR] error while calculating bucket number %s", err.Error())
				continue
			}

			// ip := net.IP(nodes[i+20 : i+24])
			// port := binary.BigEndian.Uint16(nodes[i+24 : i+26])
			// node := Node{
			// 	NodeID:  nodeID,
			// 	Address: ip,
			// 	Port:    port,
			// }
			buckets[bucketNumber] += 1
			//Handle recursivity
			//Store new nodes in this scope also pass to HandleNewNodes
		}
	}
	log.Printf("[INFO] %v", buckets)
}
