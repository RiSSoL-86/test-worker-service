package kafka

import (
	brokersettings "app/src/app_settings/brokers"
	"app/src/core/brokers"
	"context"
	"errors"
	"sync"

	kgo "github.com/segmentio/kafka-go"
)

var ErrConsumerClosed = errors.New("consumer is closed")

type Consumer struct {
	brokerAddresses []string
	mu              sync.Mutex
	readers         map[*kgo.Reader]struct{}
	closed          bool
}

func NewConsumer(settings *brokersettings.KafkaSettings) *Consumer {
	return &Consumer{
		brokerAddresses: settings.BrokerAddresses,
		readers:         make(map[*kgo.Reader]struct{}),
	}
}

func (c *Consumer) Consume(ctx context.Context, topic string, groupID string, handler brokers.MessageHandler) error {
	reader := kgo.NewReader(kgo.ReaderConfig{
		Brokers: c.brokerAddresses,
		Topic:   topic,
		GroupID: groupID,
	})

	if err := c.registerReader(reader); err != nil {
		_ = reader.Close()
		return err
	}
	defer c.unregisterReader(reader)

	for {
		kafkaMessage, err := reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		headers := make(map[string]string, len(kafkaMessage.Headers))
		for _, header := range kafkaMessage.Headers {
			headers[header.Key] = string(header.Value)
		}

		message := brokers.Message{
			Topic:   kafkaMessage.Topic,
			Key:     kafkaMessage.Key,
			Value:   kafkaMessage.Value,
			Headers: headers,
		}

		if err := handler(ctx, message); err != nil {
			return err
		}

		if err := reader.CommitMessages(ctx, kafkaMessage); err != nil {
			return err
		}
	}
}

func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.closed = true

	var err error
	for reader := range c.readers {
		err = errors.Join(err, reader.Close())
		delete(c.readers, reader)
	}

	return err
}

func (c *Consumer) registerReader(reader *kgo.Reader) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return ErrConsumerClosed
	}

	c.readers[reader] = struct{}{}
	return nil
}

func (c *Consumer) unregisterReader(reader *kgo.Reader) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.readers, reader)
	_ = reader.Close()
}
