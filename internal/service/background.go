package service

import (
	"context"
	"log/slog"
	"service-courier/internal/repository"
	"time"
)

type BackgroundService struct {
	deliveryRepo  repository.DeliveryRepository
	checkInterval time.Duration
}

func NewBackgroundService(deliveryRepo repository.DeliveryRepository, checkInterval time.Duration) *BackgroundService {
	return &BackgroundService{
		deliveryRepo:  deliveryRepo,
		checkInterval: checkInterval,
	}
}

func (s *BackgroundService) StartOverdueCheck(ctx context.Context) {
	ticker := time.NewTicker(s.checkInterval)
	defer ticker.Stop()

	slog.Info("Background overdue check started", "interval", s.checkInterval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Background overdue check stopped")
			return
		case <-ticker.C:
			s.checkOverdueDeliveries(ctx)
		}
	}
}

func (s *BackgroundService) checkOverdueDeliveries(ctx context.Context) {
	overdueDeliveries, err := s.deliveryRepo.GetOverdueDeliveries(ctx)
	if err != nil {
		slog.Error("Error getting overdue deliveries", "error", err)
		return
	}

	if len(overdueDeliveries) == 0 {
		return
	}

	courierIDs := make([]int64, 0, len(overdueDeliveries))

	for _, delivery := range overdueDeliveries {
		courierIDs = append(courierIDs, delivery.CourierID)
	}

	if err := s.deliveryRepo.FreeCouriers(ctx, courierIDs); err != nil {
		slog.Error("Error freeing couriers", "error", err)
		return
	}

	slog.Info("Freed couriers from overdue deliveries", "count", len(courierIDs))
}
