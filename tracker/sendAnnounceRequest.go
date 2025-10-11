package tracker

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

func (t *Tracker) sendAnnounceRequest() {

	if time.Now().After(t.ConnectionIDTime) {
		return
	}

	rawURL := t.Trackers[(t.TrackerIdx-1+len(t.Trackers))%len(t.Trackers)]
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Errc <- fmt.Errorf("invalid tracker URL %v: %w", rawURL, err)
		return
	}

	hostPort := u.Host
	if !strings.Contains(hostPort, ":") {
		hostPort += ":6969"
	}

	addr, err := net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		t.Errc <- fmt.Errorf("resolve error %v: %w", hostPort, err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)

	defer conn.Close()

	data := announceRequestPayload(t.ConnectionID, t.TransactionID, t.InfoHash, t.PeerID)
	_, err = conn.Write(data)
	if err != nil {
		t.Errc <- err
		return
	}

	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			log.Printf("[WARN] Tracker %v read timeout: %v", addr, err)
			return
		}
		if opErr, ok := err.(*net.OpError); ok {
			if strings.Contains(opErr.Err.Error(), "connection refused") {
				log.Printf("[WARN] Tracker %v connection refused", addr)
				return
			}
		}
		t.Errc <- err
		return
	}

	if n == 0 {
		return
	}
	buf = buf[:n]

	peers, interval, _, _, err := parseAnnounceResponse(buf, t.TransactionID)
	if err != nil {
		t.Errc <- err
		return
	}

	t.PeerLock.Lock()
	t.Peers = mergePeers(t.Peers, peers)

	if time.Duration(interval)*time.Second > t.PeersTime {
		t.PeersTime = time.Duration(interval) * time.Second
	}

	t.PeerLock.Unlock()

	return
}
