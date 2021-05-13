package kafka

import (
	"context"

	"github.com/Mario-Jimenez/pricescraper/subscriber"
	"github.com/juju/errors"
	"github.com/segmentio/kafka-go"
)

// Consumer receives messages from the broker
type Consumer struct {
	consumer *kafka.Reader
}

// NewConsumer creates a consumer that receives messages from kafka
func NewConsumer(topic, groupID string, brokers []string) *Consumer {
	newConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:       brokers,
		GroupID:       groupID,
		Topic:         topic,
		QueueCapacity: 1,
		StartOffset:   kafka.FirstOffset,
	})

	return &Consumer{
		consumer: newConsumer,
	}
}

// Fetch message from the broker
func (c *Consumer) Fetch(ctx context.Context) (*subscriber.Message, error) {
	m, err := c.consumer.FetchMessage(ctx)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &subscriber.Message{
		Message:   m.Value,
		Topic:     m.Topic,
		Partition: m.Partition,
		Offset:    m.Offset,
	}, nil
}

// Commit message to the broker
// if you use fetch, you have to commit the message's offset to the broker when finished
func (c *Consumer) Commit(ctx context.Context, message *subscriber.Message) error {
	m := kafka.Message{
		Topic:     message.Topic,
		Partition: message.Partition,
		Offset:    message.Offset,
	}

	return c.consumer.CommitMessages(ctx, m)
}

// Close consumer
func (c *Consumer) Close() error {
	return c.consumer.Close()
}
