package dht

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	myIDHex := "9f5a12cd8b0e9c7af014bc5e27d33a98f4d120ab"
	AHex := "1f5a12cd8b0e9c7af014bc5e27d33a98f4d120ab"

	myID, err := hex.DecodeString(myIDHex)
	if err != nil {
		panic(err)
	}
	A, err := hex.DecodeString(AHex)
	if err != nil {
		panic(err)
	}

	var final [20]byte
	for i := 0; i < 20; i++ {
		final[i] = myID[i] ^ A[i]
	}

	fmt.Printf("Distance: %b\n", final)
	distance := new(big.Int).SetBytes(final[:])

	msb := distance.BitLen()
	fmt.Println("MSB index:", msb)
	fmt.Println("Bucket Number:", 160-msb)
}
