package dht

import (
	"encoding/binary"
	"log"
	"net"
	"sync"
)

func (d *DHT) JoinDHT() {
	bootstrapNodes := []string{
		"dht.libtorrent.org:25401",
		"router.bittorrent.com:6881",
		"router.utorrent.com:6881",
		"dht.transmissionbt.com:6881",
	}

	var wg sync.WaitGroup

	for _, addr := range bootstrapNodes {
		msg := DHTRequest{
			T: generateTransactionID(),
			Y: "q",
			Q: "find_node",
			A: map[string][]byte{
				"id":     d.NodeID[:],
				"target": d.NodeID[:],
			},
		}

		raddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			log.Printf("[ERROR] resolve %s: %v", addr, err)
			continue
		}

		wg.Add(1)
		go func(raddr *net.UDPAddr) {
			defer wg.Done()

			resp, err := SendRequest(msg, Node{
				Address: raddr.IP,
				Port:    uint16(raddr.Port),
			})
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					return
				}
				log.Printf("[ERROR] while sending request to %s: %v", raddr, err)
				return
			}

			nodes := []byte(resp.R.Nodes)
			for i := 0; i+26 <= len(nodes); i += 26 {
				var nodeID [20]byte
				copy(nodeID[:], nodes[i:i+20])

				bucketNumber, err := GetBucketNumber(nodeID[:], d.NodeID[:])
				if err != nil {
					log.Printf("[ERROR] calculating bucket number: %v", err)
					continue
				}

				ip := net.IP(nodes[i+20 : i+24])
				port := binary.BigEndian.Uint16(nodes[i+24 : i+26])
				node := Node{
					NodeID:  nodeID,
					Address: ip,
					Port:    port,
				}

				d.Table.BucketLock[bucketNumber].Lock()
				d.Table.Buckets[bucketNumber] = append(d.Table.Buckets[bucketNumber], node)
				d.Table.BucketLock[bucketNumber].Unlock()

			}
		}(raddr)
	}
	wg.Wait()
}
