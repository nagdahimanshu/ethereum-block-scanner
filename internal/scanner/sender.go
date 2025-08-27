package scanner

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

// GetSenderAddress extracts the sender address from a transaction.
func GetSenderAddress(tx *types.Transaction) (string, error) {
	if tx.ChainId().BitLen() > 0 {
		if sender, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx); err == nil {
			return strings.ToLower(sender.Hex()), nil
		}
	}
	if tx.Protected() {
		signer := types.NewEIP155Signer(tx.ChainId())
		if sender, err := types.Sender(signer, tx); err == nil {
			return strings.ToLower(sender.Hex()), nil
		}
	}
	signer := types.HomesteadSigner{}
	if sender, err := types.Sender(signer, tx); err == nil {
		return strings.ToLower(sender.Hex()), nil
	}
	return "", fmt.Errorf("cannot extract sender")
}
