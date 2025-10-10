package tracker

func mergePeers(oldPeers, newPeers []string) []string {
	peerSet := make(map[string]struct{}, len(oldPeers))
	for _, p := range oldPeers {
		peerSet[p] = struct{}{}
	}

	for _, p := range newPeers {
		if _, exists := peerSet[p]; !exists {
			oldPeers = append(oldPeers, p)
			peerSet[p] = struct{}{}
		}
	}

	return oldPeers
}
