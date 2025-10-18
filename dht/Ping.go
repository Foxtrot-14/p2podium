package dht

import (
	"bytes"
	"log"
	"net"
	"time"

	"github.com/jackpal/bencode-go"
)

func Ping(node Node, myID [20]byte) bool {
	msg := DHTRequest{
		T: generateTransactionID(),
		Y: "q",
		Q: "ping",
		A: map[string][]byte{
			"id": myID[:],
		},
	}

	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, msg); err != nil {
		log.Printf("[ERROR] bencode marshal: %v", err)
		return false
	}

	udpAddr := &net.UDPAddr{
		IP:   node.Address,
		Port: int(node.Port),
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Printf("[ERROR] dial %s: %v", udpAddr, err)
		return false
	}

	defer conn.Close()

	if _, err = conn.Write(buf.Bytes()); err != nil {
		log.Printf("[ERROR] send ping to %s: %v", udpAddr, err)
		return false
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	resp := make([]byte, 2048)
	n, addr, err := conn.ReadFromUDP(resp)
	if err != nil {
		log.Printf("[WARN] no response from %s: %v", udpAddr, err)
		return false
	}

	dhtRes, err := parseJoinDHTResponse(resp[:n])
	if err != nil {
		log.Printf("[ERROR] parse response from %s: %v", addr, err)
		return false
	}

	if dhtRes.R.ID == string(node.NodeID[:]) {
		return true
	} else {
		return false
	}
}
