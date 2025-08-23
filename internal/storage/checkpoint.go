package storage

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadLastProcessedBlock(filename string) (uint64, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return 0, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		blockStr := strings.TrimSpace(scanner.Text())
		if blockNum, err := strconv.ParseUint(blockStr, 10, 64); err == nil {
			return blockNum, nil
		}
	}
	return 0, nil
}

func WriteLastProcessedBlock(filename string, blockNumber uint64) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(strconv.FormatUint(blockNumber, 10))
	return err
}
