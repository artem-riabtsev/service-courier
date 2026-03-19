package order

import (
	"context"
	"time"
)

type OrderGateway interface {
	GetOrders(ctx context.Context, from time.Time) ([]*ExternalOrder, error)
	GetOrderByID(ctx context.Context, orderID string) (*ExternalOrder, error)
}

type ExternalOrder struct {
	ID        string
	Status    string `json:"status"`
	CreatedAt time.Time
}
