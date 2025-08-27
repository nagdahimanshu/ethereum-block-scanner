package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/bloom"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/config"
	kafka "github.com/nagdahimanshu/ethereum-block-scanner/internal/events"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/logger"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/storage"
)

// Scanner is the main struct for Ethereum block scanning
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
	logger       logger.Logger
	producer     *kafka.KafkaProducer
}

// New initializes a new Scanner instance
func New(ctx context.Context, cfg *config.Config, logger logger.Logger, bloomFilter *bloom.AddressBloomFilter, addressMap map[string]string, producer *kafka.KafkaProducer) (*Scanner, error) {
	client, err := ethclient.Dial(cfg.EthereumNodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	lastBlock, err := storage.ReadLastProcessedBlock(cfg.CheckpointFile)
	if err != nil {
		logger.Infof("Could not read checkpoint: %v", err)
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
		logger:         logger,
		producer:       producer,
	}, nil
}

// tryReconnect attempts to reconnect to the Ethereum node on connection failure
func (s *Scanner) tryReconnect() {
	s.Stop()

	for {
		client, err := ethclient.Dial(s.nodeURL)
		if err != nil {
			s.logger.Warnw("Failed to reconnect",
				"error", err,
				"retry_in", "5s",
				"node", s.nodeURL,
			)
			time.Sleep(5 * time.Second)
			continue
		}

		s.client = client
		s.headersChan = make(chan *types.Header)

		if err := s.Start(); err != nil {
			s.logger.Infof("Failed to restart watcher: %v. Retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		s.logger.Infof("Successfully reconnected to Ethereum node")
		break
	}
}
