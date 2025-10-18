package dht

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jackpal/bencode-go"
)

func (d *DHT) GetPeers(nodes []Node) {
	//Recursively query to peers
	for _, node := range nodes {
		msg := DHTRequest{
			T: generateTransactionID(),
			Y: "q",
			Q: "get_peers",
			A: map[string][]byte{
				"id":        d.NodeID[:],
				"info_hash": d.InfoHash[:],
			},
		}

		var buf bytes.Buffer
		if err := bencode.Marshal(&buf, msg); err != nil {
			log.Printf("[ERROR] bencode marshal: %v", err)
		}

		udpAddr := &net.UDPAddr{
			IP:   node.Address,
			Port: int(node.Port),
		}
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			log.Printf("[ERROR] dial %s: %v", udpAddr, err)
		}

		defer conn.Close()

		if _, err = conn.Write(buf.Bytes()); err != nil {
			log.Printf("[ERROR] send ping to %s: %v", udpAddr, err)
		}

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		resp := make([]byte, 2048)
		n, addr, err := conn.ReadFromUDP(resp)
		if err != nil {
			log.Printf("[WARN] no response from %s: %v", udpAddr, err)
		}

		dhtRes, err := parseJoinDHTResponse(resp[:n])
		if err != nil {
			log.Printf("[ERROR] parse response from %s: %v", addr, err)
		}

		fmt.Printf("This is the response:%v", dhtRes.R.Nodes)

		nodes := []byte(dhtRes.R.Nodes)
		for i := 0; i < len(nodes); i += 26 {
			var nodeID [20]byte
			copy(nodeID[:], nodes[i:i+20])
			bucketNumber, err := GetBucketNumber(nodeID[:], d.NodeID[:])
			if err != nil {
				log.Printf("[ERROR] error while calculating bucket number %s", err.Error())
				continue
			}
			ip := net.IP(nodes[i+20 : i+24])
			port := binary.BigEndian.Uint16(nodes[i+24 : i+26])
			node := Node{
				NodeID:  nodeID,
				Address: ip,
				Port:    port,
			}
			log.Printf("[INFO] Found Node:%v belongs to node:%v", node, bucketNumber)
		}
	}
}
