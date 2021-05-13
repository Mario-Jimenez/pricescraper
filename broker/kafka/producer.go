package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Producer sends messages to the broker
type Producer struct {
	producer *kafka.Writer
}

// NewProducer creates a producer to send messages to kafka
func NewProducer(topic string, brokers []string) *Producer {
	newProducer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.RoundRobin{},
		RequiredAcks: kafka.RequireOne,
		BatchSize:    1,
	}

	return &Producer{
		producer: newProducer,
	}
}

// Publish message to kafka
func (p *Producer) Publish(ctx context.Context, msg []byte) error {
	return p.producer.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
}

// Close producer
func (p *Producer) Close() error {
	return p.producer.Close()
}
