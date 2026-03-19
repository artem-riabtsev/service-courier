package service

import (
	"context"
	"service-courier/internal/factory"
	"service-courier/internal/model"
	"service-courier/internal/repository"
	"time"
)

type deliveryService struct {
	deliveryRepo repository.DeliveryRepository
	courierRepo  repository.CourierRepository
	factory      *factory.DeliveryTimeFactory
}

func NewDeliveryService(
	deliveryRepo repository.DeliveryRepository,
	courierRepo repository.CourierRepository,
	factory *factory.DeliveryTimeFactory,
) DeliveryService {
	return &deliveryService{
		deliveryRepo: deliveryRepo,
		courierRepo:  courierRepo,
		factory:      factory,
	}
}

func (s *deliveryService) AssignCourier(ctx context.Context, req *model.AssignRequest) (*model.AssignResponse, error) {
	if req.OrderID == "" {
		return nil, ErrInvalidInput
	}

	courier, err := s.deliveryRepo.GetAvailableCourierWithMinLoad(ctx)
	if err != nil {
		if err == repository.ErrNoAvailableCouriers {
			return nil, ErrNoAvailableCouriers
		}
		return nil, err
	}

	deadline := s.factory.CalculateDeadline(courier.TransportType)

	delivery := &model.Delivery{
		CourierID:  courier.ID,
		OrderID:    req.OrderID,
		AssignedAt: time.Now(),
		Deadline:   deadline,
	}

	courier.Status = "busy"
	if err := s.courierRepo.UpdateCourier(ctx, courier); err != nil {
		return nil, err
	}

	if err := s.deliveryRepo.CreateDelivery(ctx, delivery); err != nil {
		courier.Status = "available"
		_ = s.courierRepo.UpdateCourier(ctx, courier)
		return nil, err
	}

	return &model.AssignResponse{
		CourierID:        courier.ID,
		OrderID:          req.OrderID,
		TransportType:    courier.TransportType,
		DeliveryDeadline: deadline,
	}, nil
}

func (s *deliveryService) UnassignCourier(ctx context.Context, req *model.UnassignRequest) (*model.UnassignResponse, error) {
	if req.OrderID == "" {
		return nil, ErrInvalidInput
	}

	delivery, err := s.deliveryRepo.GetDeliveryByOrderID(ctx, req.OrderID)
	if err != nil {
		if err == repository.ErrDeliveryNotFound {
			return nil, ErrDeliveryNotFound
		}
		return nil, err
	}

	courier, err := s.courierRepo.GetCourier(ctx, delivery.CourierID)
	if err != nil {
		return nil, err
	}

	courier.Status = "available"
	if err := s.courierRepo.UpdateCourier(ctx, courier); err != nil {
		return nil, err
	}

	if err := s.deliveryRepo.DeleteDelivery(ctx, req.OrderID); err != nil {
		courier.Status = "busy"
		_ = s.courierRepo.UpdateCourier(ctx, courier)
		return nil, err
	}

	return &model.UnassignResponse{
		OrderID:   req.OrderID,
		Status:    "unassigned",
		CourierID: delivery.CourierID,
	}, nil
}
