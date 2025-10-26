package dht

func (d *DHT) HandleNewNodes() {
	for node := range d.NodeC {
		bucketNumber, err := GetBucketNumber(node.NodeID[:], d.InfoHash[:])
		if err != nil {
			continue
		}

		d.Table.TableLock.Lock()

		bucket := d.Table.Buckets[bucketNumber]
		if len(bucket) < 8 {
			d.Table.Buckets[bucketNumber] = append(bucket, node)
		} else {
			index := LeastResponsive(bucket)
			if !Ping(bucket[index], d.InfoHash) {
				bucket[index] = node
			}
			d.Table.Buckets[bucketNumber] = bucket
		}

		d.Table.TableLock.Unlock()
	}
}
