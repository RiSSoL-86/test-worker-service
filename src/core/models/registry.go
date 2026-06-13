package models

import (
	"app/src/core/models/orders"
)

func GetModels() []any {
	return []any{
		&orders.Order{},
		// Add new models here
	}
}
