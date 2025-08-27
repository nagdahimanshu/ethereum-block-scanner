package bloom

import bloom "github.com/bits-and-blooms/bloom/v3"

// AddressBloomFilter wraps a bloom.BloomFilter for Ethereum addresses
type AddressBloomFilter struct {
	filter *bloom.BloomFilter
}

// New creates a new instance of AddressBloomFilter
func New(size uint, hash uint) *AddressBloomFilter {
	return &AddressBloomFilter{
		filter: bloom.NewWithEstimates(size, 0.0001), // 0.01% false positive rate
	}
}

// Add inserts a single address into the bloom filter
func (b *AddressBloomFilter) Add(address string) {
	b.filter.Add([]byte(address))
}

// Test checks if a given address might be in the filter
func (b *AddressBloomFilter) Test(address string) bool {
	return b.filter.Test([]byte(address))
}

// BatchTest checks multiple addresses at once
func (b *AddressBloomFilter) BatchTest(addresses []string) []string {
	var matches []string
	for _, addr := range addresses {
		if b.Test(addr) {
			matches = append(matches, addr)
		}
	}
	return matches
}
