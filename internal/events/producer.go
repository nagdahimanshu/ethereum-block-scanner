package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nagdahimanshu/ethereum-block-scanner/internal/logger"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
	logger logger.Logger
	ctx    context.Context
}

func NewProducer(ctx context.Context, logger logger.Logger, brokers []string, topic string) *KafkaProducer {
	if err := createTopicIfNotExists(brokers[0], topic, logger); err != nil {
		logger.Errorf("Failed to create topic %s: %v", topic, err)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Async:    false,
	})

	return &KafkaProducer{
		writer: writer,
		topic:  topic,
		logger: logger,
		ctx:    ctx,
	}
}

func (p *KafkaProducer) PublishEvent(event interface{}) {
	data, err := json.Marshal(event)
	if err != nil {
		p.logger.Errorf("Failed to parse event", err)
		return
	}

	if err := p.writer.WriteMessages(p.ctx,
		kafka.Message{
			Key:   []byte(time.Now().Format(time.RFC3339Nano)),
			Value: data,
		},
	); err != nil {
		p.logger.Errorf("Failed to publish event", err)
	} else {
		p.logger.Infof("Published event to topic %s", p.topic)
	}
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

// createTopicIfNotExists creates the topic if it doesn't exist
func createTopicIfNotExists(broker, topic string, log logger.Logger) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	// check if topic already exists
	partitions, err := conn.ReadPartitions()
	if err == nil {
		for _, p := range partitions {
			if p.Topic == topic {
				log.Infof("Topic %s already exists", topic)
				return nil
			}
		}
	}

	topicConfig := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	if err := conn.CreateTopics(topicConfig...); err != nil {
		return err
	}

	log.Infof("Created topic %s", topic)
	return nil
}
