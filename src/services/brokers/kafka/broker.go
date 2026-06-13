package kafka

import (
	brokersettings "app/src/app_settings/brokers"
	"app/src/core/brokers"
	"errors"
)

type Broker struct {
	producer *Producer
	consumer *Consumer
}

func NewBroker(settings *brokersettings.KafkaSettings) *Broker {
	return &Broker{
		producer: NewProducer(settings),
		consumer: NewConsumer(settings),
	}
}

func (b *Broker) Producer() brokers.Producer {
	return b.producer
}

func (b *Broker) Consumer() brokers.Consumer {
	return b.consumer
}

func (b *Broker) Close() error {
	return errors.Join(
		b.producer.Close(),
		b.consumer.Close(),
	)
}
