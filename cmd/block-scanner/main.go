package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/bloom"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/config"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/events"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/logger"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/scanner"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/server"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	logger, err := logger.NewDefaultProductionLogger(cfg.LogLevel)
	if err != nil {
		logger.Fatalf("Failed to create logger: %v", err)
	}

	// Initialize HTTP server
	httpServer := server.NewServer(ctx, logger, cfg.Port)
	go httpServer.Start(ctx)

	// Load addresses
	addressMap, err := storage.ReadAddresses(cfg.AddressesFilePath)
	if err != nil {
		logger.Fatalf("Failed to load addresses: %v", err)
	}

	// Build bloom filter
	bloomFilter := bloom.New(cfg.BloomFilterSize, cfg.BloomFilterHash)
	for addr := range addressMap {
		bloomFilter.Add(addr)
	}

	// Kafka setup
	brokers := cfg.KafkaBrokers
	topic := cfg.KafkaTopic
	producer := events.NewProducer(ctx, logger, brokers, topic)
	defer producer.Close()

	// Init scanner
	watcher, err := scanner.New(ctx, cfg, logger, bloomFilter, addressMap, producer)
	if err != nil {
		logger.Fatalf("Failed to init scanner: %v", err)
	}

	if err := watcher.Start(); err != nil {
		logger.Fatalf("Failed to start scanner: %v", err)
	}
	defer watcher.Stop()

	logger.Infow("Scanner started",
		"node", cfg.EthereumNodeURL,
		"bloom_size", cfg.BloomFilterSize,
	)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Infof("Shutting down...")
}
