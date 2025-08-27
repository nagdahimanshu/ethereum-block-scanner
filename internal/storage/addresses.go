package storage

import (
	"encoding/csv"
	"os"
	"strings"
)

// ReadAddresses reads a CSV file containing user IDs and Ethereum addresses
func ReadAddresses(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	addresses := make(map[string]string)
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	startIndex := 0
	if len(records) > 0 {
		firstRowLower := strings.ToLower(strings.Join(records[0], ","))
		if strings.Contains(firstRowLower, "address") && strings.Contains(firstRowLower, "userid") {
			startIndex = 1
		}
	}

	for i := startIndex; i < len(records); i++ {
		if len(records[i]) >= 2 {
			userId := strings.TrimSpace(records[i][0])
			address := strings.ToLower(strings.TrimSpace(records[i][1]))
			if address != "" {
				addresses[address] = userId
			}
		}
	}

	return addresses, nil
}
