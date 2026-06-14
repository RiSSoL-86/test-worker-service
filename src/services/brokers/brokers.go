package brokers

import (
	brokersettings "app/src/app_settings/brokers"
	"app/src/core/brokers"
	"app/src/services/brokers/kafka"
)

func NewKafkaBroker(settings *brokersettings.KafkaSettings) brokers.Broker {
	return kafka.NewBroker(settings)
}
