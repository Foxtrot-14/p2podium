package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Foxtrot-14/p2podium/dht"
	"github.com/Foxtrot-14/p2podium/magnet"
	_ "github.com/Foxtrot-14/p2podium/tracker"
)

var (
	// Build information derived from ldflags -X
	buildRevision string
	buildVersion  string
	buildTime     string

	// Command-line flags
	versionFlag bool
	helpFlag    bool

	// Program control channels
	Done = make(chan struct{})
	Errc = make(chan error)

	// Shared peer list and mutex
	PeerList []string
	PeerLock sync.Mutex
)

func init() {
	flag.BoolVar(&versionFlag, "version", false, "Show current version and exit")
	flag.BoolVar(&helpFlag, "help", false, "Show usage information and exit")
}

func setBuildVariables() {
	if buildRevision == "" {
		buildRevision = "dev"
	}
	if buildVersion == "" {
		buildVersion = "dev"
	}
	if buildTime == "" {
		buildTime = time.Now().UTC().Format(time.RFC3339)
	}
}

func parseFlags() {
	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		fmt.Printf("%s %s %s\n", buildRevision, buildVersion, buildTime)
		os.Exit(0)
	}
}

func handleInterrupts() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	sig := <-interrupt
	log.Printf("[INFO] Interrupt signal caught: %v", sig)

	// Trigger graceful shutdown for all goroutines
	close(Done)
}

func main() {
	setBuildVariables()
	parseFlags()

	go handleInterrupts()

	// TODO: Obtain magnet link from user input
	// Example magnet link for testing
	sample := "magnet:?xt=urn:btih:dd8255ecdc7ca55fb0bbf81323d87062db1f6d1c&dn=Big+Buck+Bunny&tr=udp%3A%2F%2Fexplodie.org%3A6969&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.empire-js.us%3A1337&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=wss%3A%2F%2Ftracker.btorrent.xyz&tr=wss%3A%2F%2Ftracker.fastcast.nz&tr=wss%3A%2F%2Ftracker.openwebtorrent.com&ws=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2F&xs=https%3A%2F%2Fwebtorrent.io%2Ftorrents%2Fbig-buck-bunny.torrent"
	//sample := "magnet:?xt=urn:btih:EA71F8020C11A3677C7C4C17611622C029FF5CB8&dn=Jaws%20(1975)%201080p%20BrRip%20x264%20-%20YIFY&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Ftracker.bittor.pw%3A1337%2Fannounce&tr=udp%3A%2F%2Fpublic.popcorn-tracker.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969&tr=udp%3A%2F%2Fopen.demonii.com%3A1337%2Fannounce"
	m, err := magnet.ParseMagnet(sample)
	if err != nil {
		log.Fatalf("[ERROR] Failed to parse magnet link: %s", err)
	}

	log.Printf("[INFO] Total %d trackers found", len(m.Trackers))
	log.Printf("[INFO] InfoHash: %s", m.InfoHash)

	// Note: If no peers are retrieved from the DHT, fallback to querying the UDP trackers

	nodeID, err := dht.GenerateNodeID()
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
	var infoHash [20]byte
	copy(infoHash[:], []byte(m.InfoHash))

	d := &dht.DHT{
		NodeID:   nodeID,
		InfoHash: infoHash,
		Peers:    PeerList,
		PeerLock: &PeerLock,
		Done:     Done,
		Errc:     Errc,
	}

	d.Table = &dht.RoutingTable{
		Buckets:   make(map[int][]dht.Node),
		TableLock: new(sync.Mutex),
	}

	go d.GetPeerList()

	// t := &tracker.Tracker{
	// 	PeerID:   tracker.GeneratePeerID(),
	// 	InfoHash: m.InfoHash,
	// 	Trackers: m.Trackers,
	// 	Peers:    PeerList,
	// 	PeerLock: &PeerLock,
	// 	Done:     Done,
	// 	Errc:     Errc,
	// }

	// Start the tracker peer updater in a separate goroutine
	// go t.GetPeerList()

	// Wait for either an interrupt signal or a fatal error
	select {
	case <-Done:
		log.Println("[INFO] Shutting down gracefully...")
	case err := <-Errc:
		log.Printf("[ERROR] %v", err)
		close(Done)
	}

	// Allow a short delay for goroutines to finish cleanup
	time.Sleep(1 * time.Second)
	log.Println("[INFO] Exiting.")
}
