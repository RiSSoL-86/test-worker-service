package kafka

import (
	brokersettings "app/src/app_settings/brokers"
	"app/src/core/brokers"
	"context"

	kgo "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kgo.Writer
}

func NewProducer(settings *brokersettings.KafkaSettings) *Producer {
	return &Producer{
		writer: &kgo.Writer{
			Addr:         kgo.TCP(settings.BrokerAddresses...),
			Balancer:     &kgo.LeastBytes{},
			RequiredAcks: kgo.RequireAll,
			Async:        false,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, message brokers.Message) error {
	headers := make([]kgo.Header, 0, len(message.Headers))
	for key, value := range message.Headers {
		headers = append(headers, kgo.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	return p.writer.WriteMessages(ctx, kgo.Message{
		Topic:   message.Topic,
		Key:     message.Key,
		Value:   message.Value,
		Headers: headers,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
