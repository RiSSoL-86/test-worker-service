package orders

import "gorm.io/gorm"

type Dependencies struct {
	Handler *OrderHandler
}

func NewDependencies(db *gorm.DB) *Dependencies {
	repository := NewOrderRepository(db)
	service := NewOrderService(repository)

	return &Dependencies{
		Handler: NewOrderHandler(service),
	}
}
