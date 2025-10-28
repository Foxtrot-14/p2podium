package dht

import "log"

func (d *DHT) HandleNewNodes(node Node) {
	bucketNumber, err := GetBucketNumber(node.NodeID[:], d.NodeID[:])
	if err != nil {
		log.Printf("[ERROR] fetching bucket number: %v", err)
		return
	}

	lock := &d.Table.BucketLock[bucketNumber]
	lock.Lock()
	defer lock.Unlock()

	bucket := d.Table.Buckets[bucketNumber]
	if len(bucket) < 8 {
		d.Table.Buckets[bucketNumber] = append(bucket, node)
		return
	}

	index := LeastResponsive(bucket)
	if !Ping(bucket[index], d.NodeID) {
		bucket[index] = node
	}
	d.Table.Buckets[bucketNumber] = bucket
}

