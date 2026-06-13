package app_settings

import (
	"errors"
	"os"
	"strings"
)

type PostgresSettings struct {
	DSN string
}

func NewPostgresSettings() (*PostgresSettings, error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	return &PostgresSettings{DSN: dsn}, nil
}
