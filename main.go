package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Foxtrot-14/p2podium/dht"
	"github.com/Foxtrot-14/p2podium/magnet"
	"github.com/Foxtrot-14/p2podium/scraper"
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
	PeerChan = make(chan dht.Peer, 100)
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
	// TODO: Handle Seeder mode
	// TODO: Add Context to Handle Failures
	// Example magnet link for testing
	sample := "magnet:?xt=urn:btih:40DB23EFBAC9031E88A1198142B58A7DBEAACB80&dn=One.Battle.After.Another.2025.1080p.TELESYNC.x264-SyncUP&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Ftracker.bittor.pw%3A1337%2Fannounce&tr=udp%3A%2F%2Fpublic.popcorn-tracker.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969&tr=udp%3A%2F%2Fopen.demonii.com%3A1337%2Fannounce"
	m, err := magnet.ParseMagnet(sample)
	if err != nil {
		log.Fatalf("[ERROR] Failed to parse magnet link: %s", err)
	}

	log.Printf("[INFO] Display Name: %s", m.DisplayName)
	log.Printf("[INFO] InfoHash: %X", m.InfoHash)

	nodeID, err := dht.GenerateID()
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	log.Printf("[INFO] NodeID: %x", nodeID)

	d := &dht.DHT{
		NodeID:   nodeID,
		InfoHash: m.InfoHash,
		PeerChan: PeerChan,
		Done:     Done,
		Errc:     Errc,
	}

	d.Table = &dht.RoutingTable{
		Peers: make(map[[20]byte][]string),
	}

	go d.GetPeerList()

	peerId, err := dht.GenerateID()
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	sc := &scraper.Scraper{
		ActivePeers: make(map[string]int),
		PeerChan:    PeerChan,
		PeerID:      peerId,
		InfoHash:    m.InfoHash,
	}

	go sc.StartScraper()

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
