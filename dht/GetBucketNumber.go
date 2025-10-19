package dht

import (
	"errors"
	"math/big"
)

func GetBucketNumber(target, source []byte) (int, error) {
	if len(target) != 20 || len(source) != 20 {
		return -1, errors.New("both target and source IDs must be 20 bytes (160 bits)")
	}

	var xorResult [20]byte
	for i := range 20 {
		xorResult[i] = target[i] ^ source[i]
	}

	distance := new(big.Int).SetBytes(xorResult[:])
	msb := distance.BitLen()

	if msb == 0 {
		return 0, nil
	}

	bucket := 160 - msb
	return bucket, nil
}
