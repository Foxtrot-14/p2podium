package dht

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func (d *DHT) GetPeers(initialNodes []Node) {
	in := make(chan Node, 200)
	out := make(chan Node, 200)
	var wg sync.WaitGroup

	seen := make(map[[20]byte]bool)
	var mu sync.Mutex
	var counter int

	for _, n := range initialNodes {
		in <- n
	}

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for node := range in {
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
					if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
						continue
					} else if err == io.EOF {
						continue
					}
					continue
				}

				nodes := []byte(resp.R.Nodes)
				if len(nodes) > 0 {
					for i := 0; i+26 <= len(nodes); i += 26 {
						var nodeID [20]byte
						copy(nodeID[:], nodes[i:i+20])

						mu.Lock()
						if seen[nodeID] {
							mu.Unlock()
							continue
						}
						seen[nodeID] = true
						counter++
						log.Printf("[Counter] %d", counter)
						mu.Unlock()

						ip := net.IP(nodes[i+20 : i+24])
						port := binary.BigEndian.Uint16(nodes[i+24 : i+26])
						newNode := Node{NodeID: nodeID, Address: ip, Port: port}

						out <- newNode
					}
				} else if len(resp.R.Values) > 0 {
					for _, v := range resp.R.Values {
						log.Printf("[PEER] %v", v)
					}
				}
			}
		}()
	}

	go func() {
		for newNode := range out {
			in <- newNode
		}
	}()

	go func() {
		lastCount := 0
		for {
			time.Sleep(10 * time.Second)
			mu.Lock()
			if counter == lastCount {
				close(in)
				close(out)
				mu.Unlock()
				return
			}
			lastCount = counter
			mu.Unlock()
		}
	}()

	wg.Wait()
}

