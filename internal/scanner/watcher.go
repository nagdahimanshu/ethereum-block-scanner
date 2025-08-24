package scanner

import (
	"fmt"
	"time"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/storage"
)

func (s *Scanner) Start() error {
	sub, err := s.client.SubscribeNewHead(s.ctx, s.headersChan)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err)
	}
	s.subscription = sub

	s.logger.Infof("Started Ethereum block scanner...")
	go s.processHeaders()

	return nil
}

func (s *Scanner) Stop() {
	if s.subscription != nil {
		s.subscription.Unsubscribe()
	}

	if err := storage.WriteLastProcessedBlock(s.checkpointFile, s.lastBlock); err != nil {
		s.logger.Infof("Failed to save checkpoint: %v", err)
	}

	s.logger.Infof("Stopped Ethereum block scanner")
}

func (s *Scanner) processHeaders() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case err := <-s.subscription.Err():
			s.logger.Infof("Subscription error: %v. Reconnecting...", err)
			time.Sleep(5 * time.Second)
			s.tryReconnect()
		case header := <-s.headersChan:
			blockNumber := header.Number.Uint64()
			s.logger.Infof("New block arrived with number: %d", blockNumber)
			if blockNumber > s.lastBlock {
				s.processNewBlock(blockNumber)
				s.lastBlock = blockNumber
			}
		}
	}
}
