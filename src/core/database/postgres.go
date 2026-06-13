package database

import (
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectGORM(ctx context.Context, dsn string) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return conn, nil
}
