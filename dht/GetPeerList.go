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

	log.Printf("[INFO] Closest Bucket to Torrent: %d\n", bucketNumber)
	//GetPeer
	d.GetPeers(d.Table.Buckets[bucketNumber])
	for key, bucket := range d.Table.Buckets {
		count := len(bucket)
		log.Printf("[INFO] Bucket Number: %d has %d nodes\n", key, count)
	}
	// TODO:
	//Annouce presence
	// - Sleep for 10 mins before next refresh
}
