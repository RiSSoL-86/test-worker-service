package app_settings

import (
	"app/src/app_settings/brokers"

	"github.com/joho/godotenv"
)

type Config struct {
	Brokers  *brokers.Settings
	Postgres *PostgresSettings
	Grpc     *GrpcSettings
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load("src/.env")

	brokerSettings, err := brokers.NewSettings()
	if err != nil {
		return nil, err
	}

	postgresSettings, err := NewPostgresSettings()
	if err != nil {
		return nil, err
	}

	grpcSettings, err := NewGrpcSettings()
	if err != nil {
		return nil, err
	}

	return &Config{
		Brokers:  brokerSettings,
		Postgres: postgresSettings,
		Grpc:     grpcSettings,
	}, nil
}
