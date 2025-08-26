package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	EthereumNodeURL   string
	AddressesFilePath string
	BloomFilterSize   uint
	BloomFilterHash   uint
	BatchSize         int
	CheckpointFile    string
	LogLevel          string
	KafkaBrokers      []string
	KafkaTopic        string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment/default values")
	}

	return &Config{
		EthereumNodeURL:   getEnv("ETH_NODE_URL", "wss://ethereum-rpc.com"),
		AddressesFilePath: getEnv("ADDRESSES_FILE", "addresses.csv"),
		BloomFilterSize:   getEnvAsUint("BLOOM_FILTER_SIZE", 10000000),
		BloomFilterHash:   getEnvAsUint("BLOOM_FILTER_HASH", 7),
		BatchSize:         getEnvAsInt("BATCH_SIZE", 1000),
		CheckpointFile:    getEnv("CHECKPOINT_FILE", "checkpoint.txt"),
		KafkaBrokers:      getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9093"}, ","),
		KafkaTopic:        getEnv("KAFKA_TOPIC", "ethereum-tx-events"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsUint(key string, defaultValue uint) uint {
	if value := os.Getenv(key); value != "" {
		if uintValue, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint(uintValue)
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultVal []string, sep string) []string {
	if val := os.Getenv(key); val != "" {
		parts := strings.Split(val, sep)
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return defaultVal
}
