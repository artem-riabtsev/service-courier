package service

import (
	"context"
	"errors"
	"service-courier/internal/factory"
	"service-courier/internal/mocks"
	"service-courier/internal/model"
	"service-courier/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeliveryService_AssignCourier_Success(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.AssignRequest{OrderID: "test-order-123"}

	courier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "123",
		Status:        "available",
		TransportType: "car",
	}

	mockDeliveryRepo.On("GetAvailableCourierWithMinLoad", ctx).Return(courier, nil)
	mockCourierRepo.On("UpdateCourier", ctx, mock.MatchedBy(func(c *model.Courier) bool {
		return c.ID == 1 && c.Status == "busy"
	})).Return(nil)

	mockDeliveryRepo.On("CreateDelivery", ctx, mock.MatchedBy(func(d *model.Delivery) bool {
		return d.CourierID == 1 &&
			d.OrderID == "test-order-123" &&
			!d.Deadline.IsZero()
	})).Return(nil)

	resp, err := service.AssignCourier(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.CourierID)
	assert.Equal(t, "test-order-123", resp.OrderID)
	assert.Equal(t, "car", resp.TransportType)
	assert.False(t, resp.DeliveryDeadline.IsZero())

	mockDeliveryRepo.AssertExpectations(t)
	mockCourierRepo.AssertExpectations(t)
}

func TestDeliveryService_AssignCourier_InvalidInput(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.AssignRequest{OrderID: ""}
	resp, err := service.AssignCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
	assert.Nil(t, resp)

	mockDeliveryRepo.AssertNotCalled(t, "GetAvailableCourierWithMinLoad")
	mockCourierRepo.AssertNotCalled(t, "UpdateCourier")
}

func TestDeliveryService_AssignCourier_NoAvailableCouriers(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.AssignRequest{OrderID: "test-order-123"}

	mockDeliveryRepo.On("GetAvailableCourierWithMinLoad", ctx).
		Return(nil, repository.ErrNoAvailableCouriers)

	resp, err := service.AssignCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrNoAvailableCouriers, err)
	assert.Nil(t, resp)

	mockDeliveryRepo.AssertExpectations(t)
	mockCourierRepo.AssertNotCalled(t, "UpdateCourier")
}

func TestDeliveryService_AssignCourier_CourierUpdateFails(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.AssignRequest{OrderID: "test-order-123"}

	courier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "123",
		Status:        "available",
		TransportType: "car",
	}

	mockDeliveryRepo.On("GetAvailableCourierWithMinLoad", ctx).Return(courier, nil)
	mockCourierRepo.On("UpdateCourier", ctx, mock.AnythingOfType("*model.Courier")).
		Return(errors.New("update failed"))

	resp, err := service.AssignCourier(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update failed")
	assert.Nil(t, resp)

	mockDeliveryRepo.AssertExpectations(t)
	mockCourierRepo.AssertExpectations(t)
	mockDeliveryRepo.AssertNotCalled(t, "CreateDelivery")
}

func TestDeliveryService_AssignCourier_DeliveryCreateFails(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.AssignRequest{OrderID: "test-order-123"}

	courier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "123",
		Status:        "available",
		TransportType: "car",
	}

	mockDeliveryRepo.On("GetAvailableCourierWithMinLoad", ctx).Return(courier, nil)
	mockCourierRepo.On("UpdateCourier", ctx, mock.AnythingOfType("*model.Courier")).Return(nil)
	mockDeliveryRepo.On("CreateDelivery", ctx, mock.AnythingOfType("*model.Delivery")).
		Return(errors.New("create failed"))

	mockCourierRepo.On("UpdateCourier", ctx, mock.MatchedBy(func(c *model.Courier) bool {
		return c.ID == 1 && c.Status == "available"
	})).Return(nil)

	resp, err := service.AssignCourier(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create failed")
	assert.Nil(t, resp)

	mockDeliveryRepo.AssertExpectations(t)
	mockCourierRepo.AssertNumberOfCalls(t, "UpdateCourier", 2)
}

func TestDeliveryService_UnassignCourier_Success(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.UnassignRequest{OrderID: "test-order-123"}

	delivery := &model.Delivery{
		ID:         1,
		CourierID:  1,
		OrderID:    "test-order-123",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(30 * time.Minute),
	}
	courier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "123",
		Status:        "busy",
		TransportType: "car",
	}

	mockDeliveryRepo.On("GetDeliveryByOrderID", ctx, "test-order-123").Return(delivery, nil)
	mockCourierRepo.On("GetCourier", ctx, int64(1)).Return(courier, nil)
	mockCourierRepo.On("UpdateCourier", ctx, mock.MatchedBy(func(c *model.Courier) bool {
		return c.ID == 1 && c.Status == "available"
	})).Return(nil)
	mockDeliveryRepo.On("DeleteDelivery", ctx, "test-order-123").Return(nil)

	resp, err := service.UnassignCourier(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-order-123", resp.OrderID)
	assert.Equal(t, "unassigned", resp.Status)
	assert.Equal(t, int64(1), resp.CourierID)

	mockDeliveryRepo.AssertExpectations(t)
	mockCourierRepo.AssertExpectations(t)
}

func TestDeliveryService_UnassignCourier_InvalidInput(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.UnassignRequest{OrderID: ""}
	resp, err := service.UnassignCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
	assert.Nil(t, resp)

	mockDeliveryRepo.AssertNotCalled(t, "GetDeliveryByOrderID")
	mockCourierRepo.AssertNotCalled(t, "GetCourier")
}

func TestDeliveryService_UnassignCourier_DeliveryNotFound(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.UnassignRequest{OrderID: "non-existent-order"}

	mockDeliveryRepo.On("GetDeliveryByOrderID", ctx, "non-existent-order").
		Return(nil, repository.ErrDeliveryNotFound)

	resp, err := service.UnassignCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrDeliveryNotFound, err)
	assert.Nil(t, resp)

	mockDeliveryRepo.AssertExpectations(t)
	mockCourierRepo.AssertNotCalled(t, "GetCourier")
}

func TestDeliveryService_UnassignCourier_CourierNotFound(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.UnassignRequest{OrderID: "test-order-123"}

	delivery := &model.Delivery{
		ID:         1,
		CourierID:  999,
		OrderID:    "test-order-123",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(30 * time.Minute),
	}

	mockDeliveryRepo.On("GetDeliveryByOrderID", ctx, "test-order-123").Return(delivery, nil)
	mockCourierRepo.On("GetCourier", ctx, int64(999)).
		Return(nil, repository.ErrCourierNotFound)

	resp, err := service.UnassignCourier(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "courier not found")
	assert.Nil(t, resp)

	mockDeliveryRepo.AssertExpectations(t)
	mockCourierRepo.AssertExpectations(t)
	mockCourierRepo.AssertNotCalled(t, "UpdateCourier")
	mockDeliveryRepo.AssertNotCalled(t, "DeleteDelivery")
}

func TestDeliveryService_UnassignCourier_CourierUpdateFails(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.UnassignRequest{OrderID: "test-order-123"}

	delivery := &model.Delivery{
		ID:         1,
		CourierID:  1,
		OrderID:    "test-order-123",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(30 * time.Minute),
	}
	courier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "123",
		Status:        "busy",
		TransportType: "car",
	}

	mockDeliveryRepo.On("GetDeliveryByOrderID", ctx, "test-order-123").Return(delivery, nil)
	mockCourierRepo.On("GetCourier", ctx, int64(1)).Return(courier, nil)
	mockCourierRepo.On("UpdateCourier", ctx, mock.AnythingOfType("*model.Courier")).
		Return(errors.New("update failed"))

	resp, err := service.UnassignCourier(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update failed")
	assert.Nil(t, resp)

	mockDeliveryRepo.AssertExpectations(t)
	mockCourierRepo.AssertExpectations(t)
	mockDeliveryRepo.AssertNotCalled(t, "DeleteDelivery")
}

func TestDeliveryService_UnassignCourier_DeliveryDeleteFails(t *testing.T) {
	mockDeliveryRepo := &mocks.DeliveryRepository{}
	mockCourierRepo := &mocks.CourierRepository{}
	realFactory := factory.NewDeliveryTimeFactory()

	service := NewDeliveryService(mockDeliveryRepo, mockCourierRepo, realFactory)
	ctx := context.Background()

	req := &model.UnassignRequest{OrderID: "test-order-123"}

	delivery := &model.Delivery{
		ID:         1,
		CourierID:  1,
		OrderID:    "test-order-123",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(30 * time.Minute),
	}
	courier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "123",
		Status:        "busy",
		TransportType: "car",
	}

	mockDeliveryRepo.On("GetDeliveryByOrderID", ctx, "test-order-123").Return(delivery, nil)
	mockCourierRepo.On("GetCourier", ctx, int64(1)).Return(courier, nil)
	mockCourierRepo.On("UpdateCourier", ctx, mock.AnythingOfType("*model.Courier")).Return(nil)
	mockDeliveryRepo.On("DeleteDelivery", ctx, "test-order-123").
		Return(errors.New("delete failed"))

	mockCourierRepo.On("UpdateCourier", ctx, mock.MatchedBy(func(c *model.Courier) bool {
		return c.ID == 1 && c.Status == "busy"
	})).Return(nil)

	resp, err := service.UnassignCourier(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete failed")
	assert.Nil(t, resp)

	mockDeliveryRepo.AssertExpectations(t)
	mockCourierRepo.AssertNumberOfCalls(t, "UpdateCourier", 2)
}
