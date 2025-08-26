package scanner

import (
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func (s *Scanner) processNewBlock(blockNumber uint64) {
	block, err := s.client.BlockByNumber(s.ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		s.logger.Infof("Error getting block %d: %v", blockNumber, err)
		return
	}

	s.logger.Infow("Processing block",
		"block", blockNumber,
		"hash", block.Hash().Hex(),
	)

	var addressesToCheck []string
	for _, tx := range block.Transactions() {
		from, err := s.getSenderAddress(tx)
		if err == nil {
			addressesToCheck = append(addressesToCheck, from)
		}
		if to := tx.To(); to != nil {
			addressesToCheck = append(addressesToCheck, strings.ToLower(to.Hex()))
		}
	}

	potentialMatches := s.bloomFilter.BatchTest(addressesToCheck)
	s.processTransactions(block, potentialMatches)

	if len(potentialMatches) > 0 {
		s.processTransactions(block, potentialMatches)
	} else {
		s.logger.Infof("No transactions detected from the list of addresses at block: %d", blockNumber)
	}
}

func (s *Scanner) processTransactions(block *types.Block, addresses []string) {
	addressSet := make(map[string]bool)
	for _, addr := range addresses {
		addressSet[addr] = true
	}

	for _, tx := range block.Transactions() {
		from, err := s.getSenderAddress(tx)
		if err != nil {
			continue
		}
		to := tx.To()
		if to == nil {
			continue
		}

		fromStr := strings.ToLower(from)
		toStr := strings.ToLower(to.Hex())

		if userID, ok := s.addressMap[fromStr]; ok && addressSet[fromStr] {
			s.logTransaction(userID, fromStr, toStr, tx, block)
		} else if userID, ok := s.addressMap[toStr]; ok && addressSet[toStr] {
			s.logTransaction(userID, fromStr, toStr, tx, block)
			s.publishTransaction(userID, fromStr, toStr, tx, block)
		}
	}
}

func (s *Scanner) publishTransaction(userID, from, to string, tx *types.Transaction, block *types.Block) {
	event := map[string]interface{}{
		"userId":      userID,
		"from":        from,
		"to":          to,
		"amountWei":   tx.Value().String(),
		"amountEth":   weiToEther(tx.Value()),
		"hash":        tx.Hash().Hex(),
		"blockNumber": block.Number().Uint64(),
		"timestamp":   time.Unix(int64(block.Time()), 0).Format(time.RFC3339),
	}

	s.producer.PublishEvent(event)
}

func (s *Scanner) logTransaction(userID, from, to string, tx *types.Transaction, block *types.Block) {
	info := map[string]interface{}{
		"userId":      userID,
		"from":        from,
		"to":          to,
		"amountWei":   tx.Value().String(),
		"amountEth":   weiToEther(tx.Value()),
		"hash":        tx.Hash().Hex(),
		"blockNumber": block.Number().Uint64(),
		"timestamp":   time.Unix(int64(block.Time()), 0).Format(time.RFC3339),
	}
	s.logger.Infof("Transaction detected: %+v", info)
}

func weiToEther(wei *big.Int) string {
	eth := new(big.Float).SetInt(wei)
	eth.Quo(eth, big.NewFloat(1e18))
	return eth.Text('f', 8)
}
