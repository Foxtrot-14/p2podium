package dht

type DHTMessage struct {
	T string         `bencode:"t"`
	Y string         `bencode:"y"`
	Q string         `bencode:"q,omitempty"`
	A map[string]any `bencode:"a,omitempty"`
	R map[string]any `bencode:"r,omitempty"`
	E []any          `bencode:"e,omitempty"`
}
