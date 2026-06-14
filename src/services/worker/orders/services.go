package orders

import (
	ordermodel "app/src/core/models/orders"
	"app/src/services/worker/orders/models"
	"context"
	"errors"

	"github.com/google/uuid"
)

const defaultOrderStatus = "new"

var ErrInvalidOrderID = errors.New("invalid order id")

type OrderService struct {
	repository *OrderRepository
}

func NewOrderService(repository *OrderRepository) *OrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) Get(ctx context.Context, orderID string) (ordermodel.Order, error) {
	if _, err := uuid.Parse(orderID); err != nil {
		return ordermodel.Order{}, ErrInvalidOrderID
	}

	return s.repository.Get(ctx, orderID)
}

func (s *OrderService) Create(ctx context.Context, orderID string, payload models.CreateOrderPayload) error {
	if _, err := uuid.Parse(orderID); err != nil {
		return ErrInvalidOrderID
	}

	order := ordermodel.Order{
		ID:            orderID,
		Title:         payload.Title,
		Description:   payload.Description,
		CustomerName:  payload.CustomerName,
		CustomerPhone: payload.CustomerPhone,
		Address:       payload.Address,
		Priority:      payload.Priority,
		Status:        defaultOrderStatus,
	}

	return s.repository.Create(ctx, order)
}

func (s *OrderService) Update(ctx context.Context, orderID string, payload models.UpdateOrderPayload) error {
	if _, err := uuid.Parse(orderID); err != nil {
		return ErrInvalidOrderID
	}

	return s.repository.Update(ctx, orderID, payload)
}

func (s *OrderService) Delete(ctx context.Context, orderID string) error {
	if _, err := uuid.Parse(orderID); err != nil {
		return ErrInvalidOrderID
	}

	return s.repository.Delete(ctx, orderID)
}
