package dht

import (
	"crypto/rand"
)

func GenerateNodeID() ([20]byte, error) {
	var nodeID [20]byte
	_, err := rand.Read(nodeID[:])
	if err != nil {
		return nodeID, err
	}
	return nodeID, nil
}
