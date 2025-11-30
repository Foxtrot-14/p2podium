# P2Podium

P2Podium is a Native BitTorrent client. This remains the most difficult project I have attempted to build.
I may be employable, but whether or not I am a self-reliant engineer still remains to be seen.  
It is completely dependent on whether P2Podium comes to fruition. Until then, I will be iterating over it
and building it regularly.

## Tech Stack
- Go for backend  
- Wails (Go) for UI

## To-Do
- [x] Get Magnet Link
- [x] Parse Magnet Link
- [x] Obtain Trackers
- [x] Configure and store a Transaction ID
- [x] Obtain a Connection ID and store the same
- [x] Obtain a Peer List
- [ ] Handle Token storage
- [x] Announce presence to the DHT
- [x] Get list of peers
- [ ] Start listening to incoming DHT requests
- [x] Make initial handshake
- [ ] Obtain file info and write a skeleton of the same locally
  - [ ] Implement uTP
  - [ ] Get metadata in chunks
  - [ ] Figure out Merkle Trees or `.torrent` structure and obtain file structure
  - [ ] Write the same file structure to disk
  - [ ] Implement Rarest Piece First algorithm
- [ ] Spawn a goroutine for each peer and start downloading blocks
  - [ ] Figure out NAT traversal
  - [ ] Determine what messages to send to peers
- [ ] Push downloaded blocks into a channel
- [ ] Implement network request pipelining
- [ ] Create a separate goroutine to write downloaded pieces
