package tracker

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

func (t *Tracker) sendConnectionRequest() {
	if len(t.Trackers) == 0 {
		t.Errc <- fmt.Errorf("no trackers available")
		return
	}

	idx := t.TrackerIdx
	tracker := t.Trackers[idx]

	if strings.HasPrefix(tracker, "wss") {
		t.TrackerIdx = (idx + 1) % len(t.Trackers)
		t.Errc <- fmt.Errorf("tracker %v uses WSS, skipping", tracker)
		return
	}

	host, err := url.Parse(tracker)
	if err != nil {
		t.Errc <- err
		return
	}

	hostPort := host.Host
	if hostPort == "" {
		t.Errc <- fmt.Errorf("invalid tracker URL: %v", tracker)
		return
	}

	addr, err := net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		t.Errc <- fmt.Errorf("resolve error %v: %w", hostPort, err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		t.Errc <- err
		return
	}
	defer conn.Close()

	payload := connectionRequestPayload(t.TransactionID)
	_, err = conn.Write(payload)
	if err != nil {
		t.Errc <- err
		return
	}

	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			log.Printf("[WARN] Tracker %v read timeout: %v", tracker, err)
			t.TrackerIdx = (idx + 1) % len(t.Trackers)
			return
		}
		if opErr, ok := err.(*net.OpError); ok {
			if strings.Contains(opErr.Err.Error(), "connection refused") {
				log.Printf("[WARN] Tracker %v connection refused", tracker)
				t.TrackerIdx = (idx + 1) % len(t.Trackers)
				return
			}
		}
		t.Errc <- err
		return
	}

	connectionID, err := parseConnectionResponse(buf[:n], t.TransactionID)
	if err != nil {
		t.Errc <- err
		return
	}

	t.ConnectionID = connectionID
	t.ConnectionIDTime = time.Now().Add(60 * time.Second)
	t.TrackerIdx = (idx + 1) % len(t.Trackers)
	return
}
