package dht

import "log"

func (d *DHT) GetPeerList() {
	//Join DHT
	d.JoinDHT()

	//Ping All Nodes
	d.HealthCheck()

	bucketNumber, err := GetBucketNumber(d.NodeID[:], d.InfoHash[:])
	if err != nil {
		log.Printf("[ERROR] while calculating distance")
	}

	for bucketNumber >= 0 {
		nodes, ok := d.Table.Buckets[bucketNumber]
		if ok && len(nodes) > 0 {
			break
		} else {
			bucketNumber--
		}

		if bucketNumber < 0 {
			log.Printf("[WARN] no suitable bucket found, routing table might be empty")
			break
		}
	}

	go d.HandleNewNodes()

	//GetPeer
	d.GetPeers(d.Table.Buckets[bucketNumber])

	// TODO:
	//Annouce presence
	// - Sleep for 10 mins before next refresh
}
