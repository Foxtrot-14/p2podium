package scraper

func (s *Scraper) PieceDownloader(peer string) {
	for {
		s.TableLock.Lock()
		if len(s.PendingPieces) == 0 {
			s.TableLock.Unlock()
			return
		}
		index := s.PendingPieces[0]
		s.PendingPieces = s.PendingPieces[1:]
		s.TableLock.Unlock()

		data := downloadPieceFromPeer(peer, index)
		s.PieceChan <- Piece{Index: index, Data: data}

		s.TableLock.Lock()
		s.DownloadedPieces = append(s.DownloadedPieces, index)
		s.TableLock.Unlock()
	}
}
