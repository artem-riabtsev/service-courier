package service

import (
	"context"
	"service-courier/internal/mocks"
	"service-courier/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackgroundService_CheckOverdueDeliveries(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	service := NewBackgroundService(mockDeliveryRepo, 1*time.Second)
	ctx := context.Background()

	overdueDeliveries := []*model.Delivery{
		{
			ID:         1,
			CourierID:  1,
			OrderID:    "order-1",
			AssignedAt: time.Now().Add(-1 * time.Hour),
			Deadline:   time.Now().Add(-30 * time.Minute),
		},
		{
			ID:         2,
			CourierID:  2,
			OrderID:    "order-2",
			AssignedAt: time.Now().Add(-45 * time.Minute),
			Deadline:   time.Now().Add(-15 * time.Minute),
		},
	}

	mockDeliveryRepo.On("GetOverdueDeliveries", ctx).Return(overdueDeliveries, nil)
	mockDeliveryRepo.On("FreeCouriers", ctx, []int64{1, 2}).Return(nil)

	service.checkOverdueDeliveries(ctx)

	mockDeliveryRepo.AssertExpectations(t)
}

func TestBackgroundService_CheckOverdueDeliveries_NoOverdue(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	service := NewBackgroundService(mockDeliveryRepo, 1*time.Second)
	ctx := context.Background()

	mockDeliveryRepo.On("GetOverdueDeliveries", ctx).Return([]*model.Delivery{}, nil)

	service.checkOverdueDeliveries(ctx)

	mockDeliveryRepo.AssertExpectations(t)
	mockDeliveryRepo.AssertNotCalled(t, "FreeCouriers")
}

func TestBackgroundService_CheckOverdueDeliveries_Error(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	service := NewBackgroundService(mockDeliveryRepo, 1*time.Second)
	ctx := context.Background()

	mockDeliveryRepo.On("GetOverdueDeliveries", ctx).
		Return(nil, assert.AnError)

	service.checkOverdueDeliveries(ctx)

	mockDeliveryRepo.AssertExpectations(t)
	mockDeliveryRepo.AssertNotCalled(t, "FreeCouriers")
}
