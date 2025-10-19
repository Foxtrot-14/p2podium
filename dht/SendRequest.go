package dht

import (
	"bytes"
	"net"
	"time"

	"github.com/jackpal/bencode-go"
)

func SendRequest(req DHTRequest, node Node) (DHTResponse, error) {
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, req); err != nil {
		return DHTResponse{}, err
	}

	udpAddr := &net.UDPAddr{
		IP:   node.Address,
		Port: int(node.Port),
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return DHTResponse{}, err
	}

	defer conn.Close()

	if _, err = conn.Write(buf.Bytes()); err != nil {
		return DHTResponse{}, err
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	resp := make([]byte, 2048)
	n, _, err := conn.ReadFromUDP(resp)
	if err != nil {
		return DHTResponse{}, err
	}

	dhtRes, err := parseJoinDHTResponse(resp[:n])
	if err != nil {
		return DHTResponse{}, err
	}
	return dhtRes, nil
}
