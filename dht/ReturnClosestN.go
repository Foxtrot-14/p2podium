package dht

import (
	"math/big"
	"sort"
)

func xorDistance(a, b [20]byte) *big.Int {
	var result [20]byte
	for i := range 20 {
		result[i] = a[i] ^ b[i]
	}
	return new(big.Int).SetBytes(result[:])
}

func ReturnClosestN(nodes []Node, target [20]byte, N int) []Node {
	if len(nodes) == 0 {
		return nil
	}

	type nodeDist struct {
		node Node
		dist *big.Int
	}
	distances := make([]nodeDist, 0, len(nodes))
	for _, n := range nodes {
		dist := xorDistance(n.NodeID, target)
		distances = append(distances, nodeDist{n, dist})
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i].dist.Cmp(distances[j].dist) < 0
	})

	limit := N
	if len(distances) < N {
		limit = len(distances)
	}

	result := make([]Node, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, distances[i].node)
	}

	return result
}
