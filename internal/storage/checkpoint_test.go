package storage_test

import (
	"os"
	"testing"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestScanner_LastProcessedBlock(t *testing.T) {
	tests := []struct {
		name        string
		writeBlock  uint64
		expectRead  uint64
		expectError bool
	}{
		{
			name:       "write and read block",
			writeBlock: 100,
			expectRead: 100,
		},
		{
			name:       "write and read zero block",
			writeBlock: 0,
			expectRead: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test.txt")
			assert.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			// Write block to file
			err = storage.WriteLastProcessedBlock(tmpFile.Name(), tt.writeBlock)
			assert.NoError(t, err)

			// Read block
			block, err := storage.ReadLastProcessedBlock(tmpFile.Name())
			assert.NoError(t, err)
			assert.Equal(t, tt.expectRead, block)
		})
	}
}
