package repository

import (
	"context"
	"service-courier/internal/model"
)

type CourierRepository interface {
	CreateCourier(ctx context.Context, courier *model.Courier) error
	GetCourier(ctx context.Context, id int64) (*model.Courier, error)
	GetAllCouriers(ctx context.Context) ([]*model.Courier, error)
	UpdateCourier(ctx context.Context, courier *model.Courier) error
}

type DeliveryRepository interface {
	CreateDelivery(ctx context.Context, delivery *model.Delivery) error
	GetDeliveryByOrderID(ctx context.Context, orderID string) (*model.Delivery, error)
	DeleteDelivery(ctx context.Context, orderID string) error
	GetAvailableCourier(ctx context.Context) (*model.Courier, error)
	GetOverdueDeliveries(ctx context.Context) ([]*model.Delivery, error)
	FreeCouriers(ctx context.Context, courierIDs []int64) error
	GetAvailableCourierWithMinLoad(ctx context.Context) (*model.Courier, error)
}
