package dht

import (
	"crypto/rand"
	"encoding/binary"
)

func GenerateID() ([20]byte, error) {
	var nodeID [20]byte
	_, err := rand.Read(nodeID[:])
	if err != nil {
		return nodeID, err
	}
	return nodeID, nil
}

func generateTransactionID() string {
	id := make([]byte, 2)
	binary.BigEndian.PutUint16(id, randomUint16())
	return string(id)
}

func randomUint16() uint16 {
	var b [2]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return binary.BigEndian.Uint16(b[:])
}
