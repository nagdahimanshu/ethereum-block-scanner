package scanner_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/scanner"
	"github.com/stretchr/testify/assert"
)

func TestScanner_GetSenderAddress(t *testing.T) {
	privKey, _ := crypto.GenerateKey()
	fromAddr := crypto.PubkeyToAddress(privKey.PublicKey).Hex()
	toAddr := crypto.PubkeyToAddress(privKey.PublicKey)

	// Create Legacy transaction
	legacyTx := types.NewTransaction(1, toAddr, big.NewInt(1000), 21000, big.NewInt(1), nil)
	signedLegacy, _ := types.SignTx(legacyTx, types.HomesteadSigner{}, privKey)

	sender, err := scanner.GetSenderAddress(signedLegacy)
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress(sender), common.HexToAddress(fromAddr))

	// Create EIP-155 transaction (protected)
	chainID := big.NewInt(1)
	eip155Tx := types.NewTransaction(1, toAddr, big.NewInt(2000), 21000, big.NewInt(1), nil)
	signedEIP155, _ := types.SignTx(eip155Tx, types.NewEIP155Signer(chainID), privKey)

	sender, err = scanner.GetSenderAddress(signedEIP155)
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress(sender), common.HexToAddress(fromAddr))
}
