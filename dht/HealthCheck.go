package dht

import (
	"sync"
	"time"
)

func (d *DHT) HealthCheck() {
	for bucketID, nodes := range d.Table.Buckets {
		var wg sync.WaitGroup
		activeNodesChan := make(chan Node, len(nodes))

		for i := range nodes {
			node := nodes[i]
			wg.Add(1)
			go func(n Node) {
				defer wg.Done()
				if Ping(n, d.NodeID) {
					n.LastSeen = time.Now()
					activeNodesChan <- n
				}
			}(node)
		}

		wg.Wait()
		close(activeNodesChan)

		activeNodes := make([]Node, 0, len(nodes))
		for n := range activeNodesChan {
			activeNodes = append(activeNodes, n)
		}

		d.Table.TableLock.Lock()
		d.Table.Buckets[bucketID] = activeNodes
		d.Table.TableLock.Unlock()
	}
}
