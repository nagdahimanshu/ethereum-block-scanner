package storage_test

import (
	"os"
	"testing"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestScanner_ReadAddresses(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]string
	}{
		{
			name: "with header and valid addresses",
			content: `userId,address
user1,0xABCDEF1234567890
user2,0x1234567890ABCDEF
user4,0xabcdefabcdefabcd
`,
			expected: map[string]string{
				"0xabcdef1234567890": "user1",
				"0x1234567890abcdef": "user2",
				"0xabcdefabcdefabcd": "user4",
			},
		},
		{
			name: "no header",
			content: `user1,0xAAAABBBBCCCCDDDD
user2,0x1111222233334444
`,
			expected: map[string]string{
				"0xaaaabbbbccccdddd": "user1",
				"0x1111222233334444": "user2",
			},
		},
		{
			name:     "empty file",
			content:  ``,
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary CSV file
			tmpFile, err := os.CreateTemp("", "addresses.csv")
			assert.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(tt.content)
			assert.NoError(t, err)
			tmpFile.Close()

			addressMap, err := storage.ReadAddresses(tmpFile.Name())
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, addressMap)
		})
	}
}
