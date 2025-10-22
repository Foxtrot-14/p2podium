package dht

import (
	"encoding/binary"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func (d *DHT) GetPeers(initialNodes []Node) {
	in := make(chan Node, 100)
	out := make(chan Node, 100)
	visited := make(map[[20]byte]bool)
	var visMu sync.Mutex

	var wg sync.WaitGroup
	var counter int

	for _, n := range initialNodes {
		in <- n
	}

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				node, ok := <-in
				if !ok {
					return
				}

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
					if strings.Contains(err.Error(), "connection refused") ||
						strings.Contains(err.Error(), "no route to host") {
						continue
					}
				}

				if len(resp.R.Nodes) > 0 {
					raw := []byte(resp.R.Nodes)
					var nodes []Node
					for i := 0; i+26 <= len(raw); i += 26 {
						var nodeID [20]byte
						copy(nodeID[:], raw[i:i+20])
						ip := net.IP(raw[i+20 : i+24])
						port := binary.BigEndian.Uint16(raw[i+24 : i+26])
						nodes = append(nodes, Node{NodeID: nodeID, Address: ip, Port: port})
					}

					closestNodes := ReturnClosestN(nodes, d.InfoHash, 3)
					for _, c := range closestNodes {
						visMu.Lock()
						if !visited[c.NodeID] {
							visited[c.NodeID] = true
							counter++
							visMu.Unlock()
							out <- c
						} else {
							visMu.Unlock()
						}
					}

				} else if len(resp.R.Values) > 0 {
					for _, v := range resp.R.Values {
						raw := []byte(v)
						if len(raw)%6 != 0 {
							log.Printf("[WARN] Unexpected peer value length: %d", len(raw))
							continue
						}

						type peer struct {
							IP   net.IP
							port uint16
						}

						for i := 0; i+6 <= len(raw); i += 6 {
							ip := net.IP(raw[i : i+4])
							port := binary.BigEndian.Uint16(raw[i+4 : i+6])
							log.Printf("[PEER] %v", peer{IP: ip, port: port})
						}
					}
					continue
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
			if counter == lastCount {
				log.Println("[STOP] No new nodes discovered. Stopping DHT crawl.")
				close(in)
				close(out)
				return
			}
			lastCount = counter
		}
	}()

	wg.Wait()
}
