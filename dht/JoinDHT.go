package dht

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"

	bencode "github.com/jackpal/bencode-go"
)

func (d *DHT) JoinDHT() {
	bootstrapNodes := []string{
		"dht.libtorrent.org:25401",
		"router.bittorrent.com:6881",
		"router.utorrent.com:6881",
		"dht.transmissionbt.com:6881",
	}

	for _, node := range bootstrapNodes {
		msg := DHTRequest{
			T: generateTransactionID(),
			Y: "q",
			Q: "find_node",
			A: map[string][]byte{
				"id":     d.NodeID[:],
				"target": d.NodeID[:],
			},
		}

		var buf bytes.Buffer
		if err := bencode.Marshal(&buf, msg); err != nil {
			log.Printf("[ERROR] bencode marshal: %v", err)
			continue
		}

		raddr, err := net.ResolveUDPAddr("udp", node)
		if err != nil {
			log.Printf("[ERROR] resolve %s: %v", node, err)
			continue
		}

		conn, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			log.Printf("[ERROR] dial %s: %v", node, err)
			continue
		}

		defer conn.Close()

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Printf("[ERROR] send ping to %s: %v", node, err)
			continue
		}

		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		resp := make([]byte, 2048)
		n, addr, err := conn.ReadFromUDP(resp)
		if err != nil {
			log.Printf("[WARN] no response from %s", node)
			continue
		}

		dhtRes, err := parseJoinDHTResponse(resp[:n])
		if err != nil {
			log.Printf("[ERROR] error while formatting response %s from %s", err.Error(), addr)
			continue
		}

		nodes := []byte(dhtRes.R.Nodes)
		for i := 0; i <= len(nodes); i += 26 {
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
			d.Table.TableLock.Lock()
			d.Table.Buckets[bucketNumber] = append(d.Table.Buckets[bucketNumber], node)
			d.Table.TableLock.Unlock()
		}
	}
}
