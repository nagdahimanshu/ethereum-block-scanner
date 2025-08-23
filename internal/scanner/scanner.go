package scanner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/bloom"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/config"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/storage"
)

type Scanner struct {
	client         *ethclient.Client
	bloomFilter    *bloom.AddressBloomFilter
	addressMap     map[string]string
	nodeURL        string
	headersChan    chan *types.Header
	ctx            context.Context
	checkpointFile string
	lastBlock      uint64
	//nolint:typecheck
	subscription ethereum.Subscription
}

func New(ctx context.Context, cfg *config.Config, bloomFilter *bloom.AddressBloomFilter, addressMap map[string]string) (*Scanner, error) {
	client, err := ethclient.Dial(cfg.EthereumNodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	lastBlock, err := storage.ReadLastProcessedBlock(cfg.CheckpointFile)
	if err != nil {
		log.Printf("Warning: Could not read checkpoint: %v", err)
	}

	return &Scanner{
		client:         client,
		bloomFilter:    bloomFilter,
		addressMap:     addressMap,
		nodeURL:        cfg.EthereumNodeURL,
		headersChan:    make(chan *types.Header),
		ctx:            ctx,
		checkpointFile: cfg.CheckpointFile,
		lastBlock:      lastBlock,
	}, nil
}

func (s *Scanner) tryReconnect() {
	s.Stop()

	for {
		client, err := ethclient.Dial(s.nodeURL)
		if err != nil {
			log.Printf("Failed to reconnect: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		s.client = client
		s.headersChan = make(chan *types.Header)

		if err := s.Start(); err != nil {
			log.Printf("Failed to restart watcher: %v. Retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Successfully reconnected to Ethereum node")
		break
	}
}
