package brokers

import (
	"errors"
	"os"
	"strings"
)

type KafkaSettings struct {
	BrokerAddresses []string
}

func NewKafkaSettings() (*KafkaSettings, error) {
	brokerAddresses, err := parseRequiredList("KAFKA_BROKER_ADDRESSES")
	if err != nil {
		return nil, err
	}

	return &KafkaSettings{
		BrokerAddresses: brokerAddresses,
	}, nil
}

func parseRequiredList(name string) ([]string, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return nil, errors.New(name + " is not set")
	}

	values := strings.Split(raw, ",")
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}

	if len(result) == 0 {
		return nil, errors.New(name + " is empty")
	}

	return result, nil
}
