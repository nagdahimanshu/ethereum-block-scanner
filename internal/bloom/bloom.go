package bloom

import bloom "github.com/bits-and-blooms/bloom/v3"

type AddressBloomFilter struct {
	filter *bloom.BloomFilter
}

func New(size uint, hash uint) *AddressBloomFilter {
	return &AddressBloomFilter{
		filter: bloom.NewWithEstimates(size, 0.0001), // 0.01% false positive rate
	}
}

func (b *AddressBloomFilter) Add(address string) {
	b.filter.Add([]byte(address))
}

func (b *AddressBloomFilter) Test(address string) bool {
	return b.filter.Test([]byte(address))
}

func (b *AddressBloomFilter) BatchTest(addresses []string) []string {
	var matches []string
	for _, addr := range addresses {
		if b.Test(addr) {
			matches = append(matches, addr)
		}
	}
	return matches
}
