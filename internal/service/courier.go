package service

import (
	"context"
	"service-courier/internal/model"
	"service-courier/internal/repository"
)

type courierService struct {
	repo repository.CourierRepository
}

func NewCourierService(repo repository.CourierRepository) CourierService {
	return &courierService{repo: repo}
}

func (s *courierService) CreateCourier(ctx context.Context, req *model.CreateCourierRequest) (*model.Courier, error) {
	// Валидация бизнес-правил
	if req.Name == "" || req.Phone == "" || req.Status == "" {
		return nil, ErrInvalidInput
	}

	if !isValidStatus(req.Status) {
		return nil, ErrInvalidStatus
	}

	if req.TransportType == "" {
		req.TransportType = "on_foot"
	}
	if !isValidTransportType(req.TransportType) {
		return nil, ErrInvalidInput
	}

	courier := &model.Courier{
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
	}

	if err := s.repo.CreateCourier(ctx, courier); err != nil {
		if err == repository.ErrDuplicatePhone {
			return nil, ErrDuplicatePhone
		}
		return nil, err
	}

	return courier, nil
}

func (s *courierService) GetCourier(ctx context.Context, id int64) (*model.Courier, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}

	courier, err := s.repo.GetCourier(ctx, id)
	if err != nil {
		if err == repository.ErrCourierNotFound {
			return nil, ErrCourierNotFound
		}
		return nil, err
	}

	return courier, nil
}

func (s *courierService) GetAllCouriers(ctx context.Context) ([]*model.Courier, error) {
	return s.repo.GetAllCouriers(ctx)
}

func (s *courierService) UpdateCourier(ctx context.Context, req *model.UpdateCourierRequest) (*model.Courier, error) {
	if req.ID <= 0 {
		return nil, ErrInvalidInput
	}

	existing, err := s.repo.GetCourier(ctx, req.ID)
	if err != nil {
		if err == repository.ErrCourierNotFound {
			return nil, ErrCourierNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.Status != nil {
		if !isValidStatus(*req.Status) {
			return nil, ErrInvalidStatus
		}
		existing.Status = *req.Status
	}

	if req.TransportType != nil {
		if !isValidTransportType(*req.TransportType) {
			return nil, ErrInvalidInput
		}
		existing.TransportType = *req.TransportType
	}

	if err := s.repo.UpdateCourier(ctx, existing); err != nil {
		if err == repository.ErrDuplicatePhone {
			return nil, ErrDuplicatePhone
		}
		return nil, err
	}

	return existing, nil
}

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"available": true,
		"busy":      true,
		"paused":    true,
	}
	return validStatuses[status]
}

func isValidTransportType(transportType string) bool {
	validTransportTypes := map[string]bool{
		"on_foot": true,
		"scooter": true,
		"car":     true,
	}
	return validTransportTypes[transportType]
}
