package worker

import (
	"app/src/services/worker/orders"

	"gorm.io/gorm"
)

type Dependencies struct {
	Orders *orders.Dependencies
}

func NewDependencies(db *gorm.DB) *Dependencies {
	return &Dependencies{
		Orders: orders.NewDependencies(db),
	}
}
