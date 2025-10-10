package tracker

import (
	"bytes"
	"encoding/binary"
)

func connectionRequestPayload(transactionID uint32) []byte {
	req := ConnectRequest{
		ProtocolID:    0x41727101980,
		Action:        0,
		TransactionID: transactionID,
	}

	var buf bytes.Buffer
	
	binary.Write(&buf, binary.BigEndian, req.ProtocolID)
	binary.Write(&buf, binary.BigEndian, req.Action)
	binary.Write(&buf, binary.BigEndian, req.TransactionID)
	
	return buf.Bytes()
}

func announceRequestPayload(connectionID uint64, transactionID uint32, infoHashStr string, peerID [20]byte) []byte {
	var infoHash [20]byte
	copy(infoHash[:], []byte(infoHashStr))

	req := AnnounceRequest{
		ConnectionID:  connectionID,
		Action:        1,
		TransactionID: transactionID,
		InfoHash:      infoHash,
		PeerID:        peerID,
		Downloaded:    0,
		Left:          0,
		Uploaded:      0,
		Event:         0,
		IPAddress:     0,
		Key:           0,
		NumWant:       50,
		Port:          6881,
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, req.ConnectionID)
	binary.Write(&buf, binary.BigEndian, req.Action)
	binary.Write(&buf, binary.BigEndian, req.TransactionID)
	binary.Write(&buf, binary.BigEndian, req.InfoHash)
	binary.Write(&buf, binary.BigEndian, req.PeerID)
	binary.Write(&buf, binary.BigEndian, req.Downloaded)
	binary.Write(&buf, binary.BigEndian, req.Left)
	binary.Write(&buf, binary.BigEndian, req.Uploaded)
	binary.Write(&buf, binary.BigEndian, req.Event)
	binary.Write(&buf, binary.BigEndian, req.IPAddress)
	binary.Write(&buf, binary.BigEndian, req.Key)
	binary.Write(&buf, binary.BigEndian, req.NumWant)
	binary.Write(&buf, binary.BigEndian, req.Port)

	return buf.Bytes()
}

