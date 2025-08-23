package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/bloom"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/config"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/scanner"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	// Load addresses
	addressMap, err := storage.ReadAddresses(cfg.AddressesFilePath)
	if err != nil {
		log.Fatalf("Failed to load addresses: %v", err)
	}

	// Build bloom filter
	bloomFilter := bloom.New(cfg.BloomFilterSize, cfg.BloomFilterHash)
	for addr := range addressMap {
		bloomFilter.Add(addr)
	}

	// Init scanner
	w, err := scanner.New(ctx, cfg, bloomFilter, addressMap)
	if err != nil {
		log.Fatalf("Failed to init scanner: %v", err)
	}

	if err := w.Start(); err != nil {
		log.Fatalf("Failed to start scanner: %v", err)
	}
	defer w.Stop()

	log.Println("Ethereum scanner started")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
