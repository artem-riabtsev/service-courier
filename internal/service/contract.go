package service

import (
	"context"
	"service-courier/internal/model"
)

type CourierService interface {
	CreateCourier(ctx context.Context, req *model.CreateCourierRequest) (*model.Courier, error)
	GetCourier(ctx context.Context, id int64) (*model.Courier, error)
	GetAllCouriers(ctx context.Context) ([]*model.Courier, error)
	UpdateCourier(ctx context.Context, req *model.UpdateCourierRequest) (*model.Courier, error)
}

type DeliveryService interface {
	AssignCourier(ctx context.Context, req *model.AssignRequest) (*model.AssignResponse, error)
	UnassignCourier(ctx context.Context, req *model.UnassignRequest) (*model.UnassignResponse, error)
}
