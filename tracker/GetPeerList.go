package tracker

import (
	"log"
	"time"
)

func (t *Tracker) GetPeerList() {
	for {
		t.TransactionID = generateTransactionID()
		t.sendConnectionRequest()

		t.sendAnnounceRequest()

		if len(t.Peers) >= 1 && t.ConnectionID != 0 {
			log.Printf("[INFO] Received Peers %v from Tracker %v", t.Peers, t.Trackers[(t.TrackerIdx-1+len(t.Trackers))%len(t.Trackers)])
			log.Printf("[INFO] Sleeping for %v Seconds", t.PeersTime)
			time.Sleep(t.PeersTime)
		} else {
			time.Sleep(5 * time.Second)
			continue
		}
	}
}
