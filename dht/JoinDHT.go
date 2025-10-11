package dht

import (
	"bytes"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"log"
	"net"
	"time"
)

func (d *DHT) JoinDHT() {
	fmt.Printf("[INFO] Joining DHT network...\n")
	fmt.Printf("[INFO] NodeID: %x\n", d.NodeID)
	fmt.Printf("[INFO] InfoHash: %x\n", d.InfoHash)

	bootstrapNodes := []string{
		"dht.libtorrent.org:25401",
		"router.bittorrent.com:6881",
		"router.utorrent.com:6881",
		"dht.transmissionbt.com:6881",
	}

	for _, node := range bootstrapNodes {
		msg := DHTMessage{
			T: generateTransactionID(),
			Y: "q",
			Q: "ping",
			A: map[string]any{
				"id": string(d.NodeID[:]),
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

		log.Printf("[INFO] Sent ping to %s", node)

		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		resp := make([]byte, 2048)
		n, addr, err := conn.ReadFromUDP(resp)
		if err != nil {
			log.Printf("[WARN] no response from %s", node)
			continue
		}

		dhtRes, err := parseJoinDHTResponse(resp[:n])
		if err != nil {
			log.Printf("[ERROR] error while formatting reponse %s", err.Error())
			continue
		}

		log.Printf("[RECV] recevied message: %v\n from %v", dhtRes, addr)
	}

	fmt.Printf("[INFO] DHT bootstrap complete.\n")
}
