package orders

import (
	"app/src/core/brokers"
	orderspb "app/src/core/proto/orders"
	ordercommon "app/src/services/worker/orders/common"
	"app/src/services/worker/orders/models"
	"context"
	"encoding/json"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	orderspb.UnimplementedOrdersServiceServer
	service *OrderService
}

func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Get(ctx context.Context, req *orderspb.GetRequest) (*orderspb.GetResponse, error) {
	order, err := h.service.Get(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidOrderID):
			return nil, status.Error(codes.InvalidArgument, "invalid order id")
		case errors.Is(err, ErrOrderNotFound):
			return nil, status.Error(codes.NotFound, "order not found")
		default:
			log.Printf("orders grpc: get order failed (id=%s): %v", req.GetId(), err)
			return nil, status.Error(codes.Internal, "failed to get order")
		}
	}

	return &orderspb.GetResponse{Order: ordercommon.ToOrderProto(order)}, nil
}

func (h *OrderHandler) Create(ctx context.Context, notification brokers.Notification) error {
	var payload models.CreateOrderPayload
	if err := json.Unmarshal(notification.Payload, &payload); err != nil {
		log.Printf("orders consumer: skipping create with bad payload (id=%s): %v", notification.EntityID, err)
		return nil
	}

	err := h.service.Create(ctx, notification.EntityID, payload)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrInvalidOrderID) {
		log.Printf("orders consumer: skipping create with invalid order id (id=%s): %v", notification.EntityID, err)
		return nil
	}

	if errors.Is(err, ErrOrderNotFound) {
		log.Printf("orders consumer: skipping create for missing order (id=%s): %v", notification.EntityID, err)
		return nil
	}

	log.Printf("orders consumer: create failed (id=%s), will retry: %v", notification.EntityID, err)
	return err
}

func (h *OrderHandler) Update(ctx context.Context, notification brokers.Notification) error {
	var payload models.UpdateOrderPayload
	if err := json.Unmarshal(notification.Payload, &payload); err != nil {
		log.Printf("orders consumer: skipping update with bad payload (id=%s): %v", notification.EntityID, err)
		return nil
	}

	err := h.service.Update(ctx, notification.EntityID, payload)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrInvalidOrderID) {
		log.Printf("orders consumer: skipping update with invalid order id (id=%s): %v", notification.EntityID, err)
		return nil
	}

	if errors.Is(err, ErrOrderNotFound) {
		log.Printf("orders consumer: skipping update for missing order (id=%s): %v", notification.EntityID, err)
		return nil
	}

	log.Printf("orders consumer: update failed (id=%s), will retry: %v", notification.EntityID, err)
	return err
}

func (h *OrderHandler) Delete(ctx context.Context, notification brokers.Notification) error {
	err := h.service.Delete(ctx, notification.EntityID)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrInvalidOrderID) {
		log.Printf("orders consumer: skipping delete with invalid order id (id=%s): %v", notification.EntityID, err)
		return nil
	}

	if errors.Is(err, ErrOrderNotFound) {
		log.Printf("orders consumer: skipping delete for missing order (id=%s): %v", notification.EntityID, err)
		return nil
	}

	log.Printf("orders consumer: delete failed (id=%s), will retry: %v", notification.EntityID, err)
	return err
}
