package scanner

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/metrics"
	"go.uber.org/multierr"
)

// TxEvent represents a normalized blockchain transaction event
type TxEvent struct {
	UserID      string `json:"userId"`
	From        string `json:"from"`
	To          string `json:"to"`
	AmountWei   string `json:"amountWei"`
	AmountEth   string `json:"amountEth"`
	Hash        string `json:"hash"`
	BlockNumber uint64 `json:"blockNumber"`
	Timestamp   string `json:"timestamp"`
}

// ValidateEvent checks transaction event
func ValidateEvent(e TxEvent) error {
	var validatorErr error

	if e.UserID == "" {
		validatorErr = multierr.Append(validatorErr, fmt.Errorf("userId is required"))
	}
	if e.From == "" {
		validatorErr = multierr.Append(validatorErr, fmt.Errorf("from is required"))
	}
	if e.To == "" {
		validatorErr = multierr.Append(validatorErr, fmt.Errorf("to is required"))
	}
	if e.AmountWei == "" {
		validatorErr = multierr.Append(validatorErr, fmt.Errorf("amountWei is required"))
	}
	if e.AmountEth == "" {
		validatorErr = multierr.Append(validatorErr, fmt.Errorf("amountEth is required"))
	}
	if e.Hash == "" {
		validatorErr = multierr.Append(validatorErr, fmt.Errorf("hash is required"))
	}
	if e.BlockNumber == 0 {
		validatorErr = multierr.Append(validatorErr, fmt.Errorf("blockNumber is required"))
	}
	if e.Timestamp == "" {
		validatorErr = multierr.Append(validatorErr, fmt.Errorf("timestamp is required"))
	}

	return validatorErr
}

// processNewBlock processes all transactions in a block using worker pool
func (s *Scanner) processNewBlock(blockNumber uint64) {
	metrics.CurrentBlock.Set(float64(blockNumber))
	metrics.BlocksProcessed.Inc()

	block, err := s.client.BlockByNumber(s.ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		s.logger.Infof("Error getting block %d: %v", blockNumber, err)
		return
	}

	s.logger.Infow("Processing block",
		"block", blockNumber,
		"hash", block.Hash().Hex(),
	)

	metrics.TransactionsProcessed.Add(float64(len(block.Transactions())))

	var addressesToCheck []string
	for _, tx := range block.Transactions() {
		from, err := GetSenderAddress(tx)
		if err == nil {
			addressesToCheck = append(addressesToCheck, from)
		}
		if to := tx.To(); to != nil {
			addressesToCheck = append(addressesToCheck, strings.ToLower(to.Hex()))
		}
	}

	potentialMatches := s.bloomFilter.BatchTest(addressesToCheck)

	if len(potentialMatches) > 0 {
		s.processTransactions(block, potentialMatches)
	} else {
		s.logger.Infof("No transactions detected from the list of addresses at block: %d", blockNumber)
	}
}

// processTransactions uses a bounded worker pool to process transactions
func (s *Scanner) processTransactions(block *types.Block, addresses []string) {
	jobs := make(chan *types.Transaction, JobQueueSize)
	done := make(chan bool)

	// Start workers
	for i := 0; i < NumWorkers; i++ {
		go func() {
			for tx := range jobs {
				s.ProcessTransaction(tx, block, addresses)
			}
			done <- true
		}()
	}

	// Add jobs into the queue
	for _, tx := range block.Transactions() {
		jobs <- tx
	}
	close(jobs)

	// Wait for all workers to finish
	for i := 0; i < NumWorkers; i++ {
		<-done
	}
}

// ProcessTransaction processes a single transaction
// Checks if sender or receiver is in monitored addresses
// Logs and publishes events if a match is found
func (s *Scanner) ProcessTransaction(tx *types.Transaction, block *types.Block, addresses []string) {
	addressSet := make(map[string]bool)
	for _, addr := range addresses {
		addressSet[addr] = true
	}

	from, err := GetSenderAddress(tx)
	if err != nil {
		return
	}
	to := tx.To()
	if to == nil {
		return
	}

	fromStr := strings.ToLower(from)
	toStr := strings.ToLower(to.Hex())

	// Check if sender or receiver is in our address map
	if userID, ok := s.addressMap[fromStr]; ok && addressSet[fromStr] {
		s.logTransaction(userID, fromStr, toStr, tx, block)
		s.publishTransaction(userID, fromStr, toStr, tx, block)
	} else if userID, ok := s.addressMap[toStr]; ok && addressSet[toStr] {
		s.logTransaction(userID, fromStr, toStr, tx, block)
		s.publishTransaction(userID, fromStr, toStr, tx, block)
	}
}

// publishTransaction constructs a TxEvent, validates it and publishes to Kafka
func (s *Scanner) publishTransaction(userID, from, to string, tx *types.Transaction, block *types.Block) {
	event := TxEvent{
		UserID:      userID,
		From:        from,
		To:          to,
		AmountWei:   tx.Value().String(),
		AmountEth:   WeiToEther(tx.Value()),
		Hash:        tx.Hash().Hex(),
		BlockNumber: block.Number().Uint64(),
		Timestamp:   time.Unix(int64(block.Time()), 0).Format(time.RFC3339),
	}

	if err := ValidateEvent(event); err != nil {
		s.logger.Errorf("Invalid transaction event: ", err)
	}

	s.producer.PublishEvent(event)
	metrics.KafkaEventsPublished.Inc()
}

func (s *Scanner) logTransaction(userID, from, to string, tx *types.Transaction, block *types.Block) {
	event := TxEvent{
		UserID:      userID,
		From:        from,
		To:          to,
		AmountWei:   tx.Value().String(),
		AmountEth:   WeiToEther(tx.Value()),
		Hash:        tx.Hash().Hex(),
		BlockNumber: block.Number().Uint64(),
		Timestamp:   time.Unix(int64(block.Time()), 0).Format(time.RFC3339),
	}

	s.logger.Infof("Transaction detected: %+v", event)
}

// WeiToEther converts wei amount to Ether
func WeiToEther(wei *big.Int) string {
	eth := new(big.Float).SetInt(wei)
	eth.Quo(eth, big.NewFloat(1e18))
	return eth.Text('f', 8)
}
