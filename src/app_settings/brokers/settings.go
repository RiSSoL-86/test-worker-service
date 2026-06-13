package brokers

type Settings struct {
	KafkaSettings *KafkaSettings
}

func NewSettings() (*Settings, error) {
	kafkaSettings, err := NewKafkaSettings()
	if err != nil {
		return nil, err
	}

	return &Settings{
		KafkaSettings: kafkaSettings,
	}, nil
}
