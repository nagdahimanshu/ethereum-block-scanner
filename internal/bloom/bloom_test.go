package bloom_test

import (
	"testing"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/bloom"
)

func TestBloomFilter_Test(t *testing.T) {
	filter := bloom.New(500000, 5)

	address := "0xtest_address_1"
	if filter.Test(address) {
		t.Errorf("Expected address %s to not exist yet", address)
	}

	filter.Add(address)
	if !filter.Test(address) {
		t.Errorf("Expected address %s to exist after adding", address)
	}
}

func TestBloomFilter_BatchTest(t *testing.T) {
	filter := bloom.New(500000, 5)

	addresses := []string{
		"0xtest_address_1",
		"0xtest_address_2",
		"0xtest_address_3",
	}

	// Add only first two addresses
	filter.Add(addresses[0])
	filter.Add(addresses[1])
	filter.Add(addresses[2])

	matches := filter.BatchTest(addresses)

	if len(matches) != len(addresses) {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}

	expected := map[string]bool{
		addresses[0]: true,
		addresses[1]: true,
		addresses[2]: true,
	}

	for _, addr := range matches {
		if !expected[addr] {
			t.Errorf("Unexpected match: %s", addr)
		}
	}
}
