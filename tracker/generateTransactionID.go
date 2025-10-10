package tracker

import (
	"math/rand"
)

func generateTransactionID() uint32 {
	return rand.Uint32()
}
