package dht

type DHTRequest struct {
	T string            `bencode:"t"`
	Y string            `bencode:"y"`
	Q string            `bencode:"q,omitempty"`
	A map[string][]byte `bencode:"a"`
}

type DHTResponse struct {
	IP string       `bencode:"ip,omitempty"`
	T  string       `bencode:"t"`
	Y  string       `bencode:"y"`
	V  string       `bencode:"v,omitempty"`
	R  DHTResponseR `bencode:"r"`
}

type DHTResponseR struct {
	ID     string   `bencode:"id"`
	Nodes  string   `bencode:"nodes,omitempty"`
	Values []string `bencode:"values,omitempty"`
	Token  string   `bencode:"token,omitempty"`
}
