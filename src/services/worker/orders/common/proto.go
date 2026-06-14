package common

import (
	ordermodel "app/src/core/models/orders"
	orderspb "app/src/core/proto/orders"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToOrderProto(order ordermodel.Order) *orderspb.Order {
	return &orderspb.Order{
		Id:            order.ID,
		Title:         order.Title,
		Description:   order.Description,
		CustomerName:  order.CustomerName,
		CustomerPhone: order.CustomerPhone,
		Address:       order.Address,
		Priority:      order.Priority,
		Status:        order.Status,
		CreatedAt:     timestamppb.New(order.CreatedAt),
		UpdatedAt:     timestamppb.New(order.UpdatedAt),
	}
}
