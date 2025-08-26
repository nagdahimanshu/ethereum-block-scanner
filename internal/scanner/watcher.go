package scanner

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/storage"
)

const (
	Confirmations = 12  // Number of confirmations before processing a block to handle reorgs
	NumWorkers    = 4   // Number of worker for transaction processing
	JobQueueSize  = 100 // Bounded channel size for jobs
)

func (s *Scanner) Start() error {
	sub, err := s.client.SubscribeNewHead(s.ctx, s.headersChan)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %v", err)
	}
	s.subscription = sub

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

// processHeaders processes incoming Ethereum headers with safe confirmations
func (s *Scanner) processHeaders() {
	// Queue to hold incoming blocks until they have enough confirmations
	blockQueue := make([]*types.Header, 0)

	for {
		select {
		case <-s.ctx.Done():
			return
		case err := <-s.subscription.Err():
			s.logger.Infof("Subscription error: %v. Reconnecting...", err)
			time.Sleep(5 * time.Second)
			s.tryReconnect()
		case header := <-s.headersChan:
			blockQueue = append(blockQueue, header)
			s.logger.Infof("Adding block to the queue")

			// Process blocks that have enough confirmations
			for len(blockQueue) > Confirmations {
				toProcess := blockQueue[0]
				blockQueue = blockQueue[1:]

				blockNumber := toProcess.Number.Uint64()
				if blockNumber > s.lastBlock {
					s.processNewBlock(blockNumber)
					s.lastBlock = blockNumber
				}
			}
		}
	}
}
