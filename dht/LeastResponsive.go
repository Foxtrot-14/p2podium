package dht

func LeastResponsive(nodes []Node) int {
	oldest := nodes[0].LastSeen
	var index int
	
	for idx, node := range nodes {
		if node.LastSeen.Before(oldest) {
			oldest = node.LastSeen
			index = idx
		}
		continue
	}

	return index
}
