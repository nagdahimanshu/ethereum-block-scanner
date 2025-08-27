package scanner_test

import (
	"math/big"
	"testing"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/scanner"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestScanner_WeiToEther(t *testing.T) {
	tests := []struct {
		name     string
		wei      string
		expected string
	}{
		{"zero wei", "0", "0.00000000"},
		{"one wei", "1", "0.00000000"},
		{"one ether", "1000000000000000000", "1.00000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weiBigInt, ok := new(big.Int).SetString(tt.wei, 10)
			if !ok {
				t.Fatalf("invalid test input: %s", tt.wei)
			}
			got := scanner.WeiToEther(weiBigInt)
			if got != tt.expected {
				t.Errorf("weiToEther(%s) = %s; want %s", tt.wei, got, tt.expected)
			}
		})
	}
}

func TestScanner_ValidateEvent(t *testing.T) {
	tests := []struct {
		name      string
		event     scanner.TxEvent
		wantError bool
		errorMsgs []string
	}{
		{
			name:      "all fields valid",
			event:     scanner.TxEvent{"user1", "fromAddr", "toAddr", "1000", "0.00000100", "0xhash", 1, "2025-08-27T12:00:00Z"},
			wantError: false,
		},
		{
			name:      "missing UserID",
			event:     scanner.TxEvent{"", "fromAddr", "toAddr", "1000", "0.00000100", "0xhash", 1, "2025-08-27T12:00:00Z"},
			wantError: true,
			errorMsgs: []string{"userId is required"},
		},
		{
			name: "missing multiple fields",
			event: scanner.TxEvent{
				UserID:      "",
				From:        "",
				To:          "",
				AmountWei:   "",
				AmountEth:   "",
				Hash:        "",
				BlockNumber: 0,
				Timestamp:   "",
			},
			wantError: true,
			errorMsgs: []string{
				"userId is required",
				"from is required",
				"to is required",
				"amountWei is required",
				"amountEth is required",
				"hash is required",
				"blockNumber is required",
				"timestamp is required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := scanner.ValidateEvent(tt.event)
			if tt.wantError {
				assert.Error(t, err, "expected an error")
				errMessages := []string{}
				for _, e := range multierr.Errors(err) {
					errMessages = append(errMessages, e.Error())
				}
				assert.ElementsMatch(t, tt.errorMsgs, errMessages, "error messages should match")
			} else {
				assert.NoError(t, err, "expected no error")
			}
		})
	}

}
