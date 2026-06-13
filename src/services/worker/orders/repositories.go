package orders

import (
	ordermodel "app/src/core/models/orders"
	"app/src/services/worker/orders/models"
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrOrderNotFound = errors.New("order not found")

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Get(ctx context.Context, id string) (ordermodel.Order, error) {
	var order ordermodel.Order
	err := r.db.WithContext(ctx).First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ordermodel.Order{}, ErrOrderNotFound
		}
		return ordermodel.Order{}, err
	}

	return order, nil
}

func (r *OrderRepository) Create(ctx context.Context, order ordermodel.Order) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&order).
		Error
}

func (r *OrderRepository) Update(ctx context.Context, id string, payload models.UpdateOrderPayload) error {
	updates := make(map[string]any)
	if payload.Title != nil {
		updates["title"] = *payload.Title
	}
	if payload.Description != nil {
		updates["description"] = *payload.Description
	}
	if payload.CustomerName != nil {
		updates["customer_name"] = *payload.CustomerName
	}
	if payload.CustomerPhone != nil {
		updates["customer_phone"] = *payload.CustomerPhone
	}
	if payload.Address != nil {
		updates["address"] = *payload.Address
	}
	if payload.Priority != nil {
		updates["priority"] = *payload.Priority
	}
	if payload.Status != nil {
		updates["status"] = *payload.Status
	}

	tx := r.db.WithContext(ctx).Model(&ordermodel.Order{}).Where("id = ?", id)

	if len(updates) == 0 {
		var count int64
		if err := tx.Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return ErrOrderNotFound
		}
		return nil
	}

	result := tx.Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrOrderNotFound
	}

	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&ordermodel.Order{}).
		Error
}
