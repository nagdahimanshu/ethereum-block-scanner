package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	BlocksProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "block_scanner_blocks_processed_total",
		Help: "Total number of blocks processed",
	})

	TransactionsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "block_scanner_transactions_processed_total",
		Help: "Total number of transactions processed",
	})

	KafkaEventsPublished = promauto.NewCounter(prometheus.CounterOpts{
		Name: "block_scanner_kafka_events_published_total",
		Help: "Total number of events published to Kafka",
	})

	Reconnections = promauto.NewCounter(prometheus.CounterOpts{
		Name: "block_scanner_reconnections_total",
		Help: "Total number of reconnections to Ethereum node",
	})

	CurrentBlock = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "block_scanner_current_block",
		Help: "Current block number being processed",
	})
)
