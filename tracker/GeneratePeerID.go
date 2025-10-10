package tracker

import (
	"crypto/rand"
)

func GeneratePeerID() [20]byte {
	var peerID [20]byte
	prefix := "-GO0001-"
	copy(peerID[:], []byte(prefix))

	randomBytes := make([]byte, 20-len(prefix))
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	for i, b := range randomBytes {
		peerID[len(prefix)+i] = '0' + (b % 10)
	}

	return peerID
}
