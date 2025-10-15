package dht

import (
	"fmt"
)

func (d *DHT) GetPeerList() {

	d.JoinDHT()

	fmt.Printf("Contents of closest bucket: %v", d.Table.Buckets[2])

	//Check each health

	//Send request to retrive peers
	//Run until peers found then sleep for 10 mins
}
