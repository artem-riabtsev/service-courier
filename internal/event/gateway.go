package event

import (
	"context"
	"service-courier/internal/gateway/order"
)

type OrderGateway interface {
	GetOrderByID(ctx context.Context, orderID string) (*order.ExternalOrder, error)
}
