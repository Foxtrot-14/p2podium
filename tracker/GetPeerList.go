package tracker

import (
	"time"
	"log"
)

func (t *Tracker) GetPeerList() {
	for {
		t.TransactionID = generateTransactionID()
		t.sendConnectionRequest()

		t.sendAnnounceRequest()

		if len(t.Peers) > 3 && t.ConnectionID != 0 {
			log.Printf("[DEBUG] Received Peers %v from Tracker %v", t.Peers, t.Trackers[t.TrackerIdx])
			log.Printf("[DEBUG] Sleeping for %v Seconds", t.PeersTime)
			time.Sleep(t.PeersTime)
		} else {
			time.Sleep(5*time.Second)
			continue
		}
	}
}

