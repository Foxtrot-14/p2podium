package dht

import (
	"sync"
	"time"
)

func (d *DHT) HealthCheck() {
	const maxConcurrentPings = 50

	var globalWG sync.WaitGroup
	sem := make(chan struct{}, maxConcurrentPings)

	for bucketID, nodes := range d.Table.Buckets {
		globalWG.Add(1)
		go func() {
			defer globalWG.Done()

			activeNodes := make([]Node, 0, len(nodes))
			var mu sync.Mutex
			var wg sync.WaitGroup

			for _, node := range nodes {
				sem <- struct{}{}
				wg.Add(1)

				go func(n Node) {
					defer wg.Done()
					defer func() { <-sem }()

					if Ping(n, d.NodeID) {
						n.LastSeen = time.Now()
						mu.Lock()
						activeNodes = append(activeNodes, n)
						mu.Unlock()
					}
				}(node)
			}

			wg.Wait()
			d.Table.BucketLock[bucketID].Lock()
			d.Table.Buckets[bucketID] = activeNodes
			d.Table.BucketLock[bucketID].Unlock()
		}()
	}

	globalWG.Wait()
}
