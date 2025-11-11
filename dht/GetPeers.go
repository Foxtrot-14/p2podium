package dht

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func (d *DHT) GetPeers(initialNodes []Node) {
	const workerIdleTimeout = 6 * time.Second

	in := make(chan Node, 1024)
	out := make(chan Node, 1024)
	defer close(out)

	var (
		visited     = make(map[[20]byte]bool)
		peerVisited = make(map[string]bool)
		visMu       sync.Mutex
		wg          sync.WaitGroup
	)

	go func() {
		for _, n := range initialNodes {
			in <- n
		}
	}()

	go func() {
		for n := range out {
			select {
			case in <- n:
			default:
				// in is full, drop this node (or you can block if you prefer)
			}
		}
	}()

	worker := func() {
		defer wg.Done()
		for {
			select {
			case node := <-in:
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
					continue
				}

				if len(resp.R.Nodes) == 0 && len(resp.R.Values) == 0 {
					continue
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

					closest := ReturnClosestN(nodes, d.InfoHash, 15)
					for _, c := range closest {
						visMu.Lock()
						if !visited[c.NodeID] {
							visited[c.NodeID] = true
							visMu.Unlock()
							select {
							case out <- c:
							default:
							}
							go d.HandleNewNodes(c)
						} else {
							visMu.Unlock()
						}
					}
				}

				if len(resp.R.Values) > 0 {
					for _, v := range resp.R.Values {
						raw := []byte(v)
						if len(raw)%6 != 0 {
							continue
						}
						for i := 0; i+6 <= len(raw); i += 6 {
							ip := net.IP(raw[i : i+4])
							port := binary.BigEndian.Uint16(raw[i+4 : i+6])
							ipPort := fmt.Sprintf("%s:%d", ip, port)

							visMu.Lock()
							if !peerVisited[ipPort] {
								peerVisited[ipPort] = true
								visMu.Unlock()
								select {
								case d.PeerChan <- Peer{IP: ip, Port: port}:
								default:
								}
							} else {
								visMu.Unlock()
							}
						}
					}
				}

			case <-time.After(workerIdleTimeout):
				return
			}
		}
	}

	wg.Add(60)
	for range 60 {
		go worker()
	}

	wg.Wait()

	close(in)

	log.Printf("[DHT] GetPeers finished: discovered %d nodes, %d peers", len(visited), len(peerVisited))
}
