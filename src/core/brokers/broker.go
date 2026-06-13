package brokers

import "context"

type Broker interface {
	Producer() Producer
	Consumer() Consumer
	Close() error
}

type Producer interface {
	Publish(ctx context.Context, message Message) error
}

type Consumer interface {
	Consume(ctx context.Context, topic string, groupID string, handler MessageHandler) error
}

type MessageHandler func(ctx context.Context, message Message) error
