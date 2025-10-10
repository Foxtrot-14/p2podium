package tracker

import (
	"encoding/binary"
	"fmt"
	"net"
)

func parseConnectionResponse(res []byte, sentTransactionID uint32) (uint64, error) {
	if len(res) < 16 {
		return 0, fmt.Errorf("[ERROR] invalid response length: %d", len(res))
	}

	action := binary.BigEndian.Uint32(res[0:4])
	transactionID := binary.BigEndian.Uint32(res[4:8])
	connectionID := binary.BigEndian.Uint64(res[8:16])

	if action != 0 {
		return 0, fmt.Errorf("[ERROR] invalid action code: %d", action)
	}

	if transactionID != sentTransactionID {
		return 0, fmt.Errorf("[ERROR] transaction ID mismatch: got %d, expected %d", transactionID, sentTransactionID)
	}

	return connectionID, nil
}

func parseAnnounceResponse(res []byte, sentTransactionID uint32) ([]string, uint32, uint32, uint32, error) {
	if len(res) < 20 {
		return nil, 0, 0, 0, fmt.Errorf("[ERROR] invalid response length: %d", len(res))
	}

	action := binary.BigEndian.Uint32(res[0:4])
	transactionID := binary.BigEndian.Uint32(res[4:8])

	if action != 1 {
		return nil, 0, 0, 0, fmt.Errorf("[ERROR] invalid action code: %d", action)
	}

	if transactionID != sentTransactionID {
		return nil, 0, 0, 0, fmt.Errorf("[ERROR] transaction ID mismatch: got %d, expected %d", transactionID, sentTransactionID)
	}

	interval := binary.BigEndian.Uint32(res[8:12])
	leechers := binary.BigEndian.Uint32(res[12:16])
	seeders := binary.BigEndian.Uint32(res[16:20])

	peers := []string{}
	for i := 20; i+6 <= len(res); i += 6 {
		ip := net.IPv4(res[i], res[i+1], res[i+2], res[i+3]).String()
		port := binary.BigEndian.Uint16(res[i+4 : i+6])
		peers = append(peers, fmt.Sprintf("%s:%d", ip, port))
	}

	return peers, interval, leechers, seeders, nil
}
