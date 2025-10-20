package dht

func Ping(node Node, myID [20]byte) bool {
	msg := DHTRequest{
		T: generateTransactionID(),
		Y: "q",
		Q: "ping",
		A: map[string][]byte{
			"id": myID[:],
		},
	}

	resp, err := SendRequest(msg, node)
	if err != nil {
		return false
	}

	if resp.R.ID == string(node.NodeID[:]) {
		return true
	} else {
		return false
	}
}
