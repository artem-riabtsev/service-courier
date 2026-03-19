package service

import (
	"context"
	"errors"
	"service-courier/internal/mocks"
	"service-courier/internal/model"
	"service-courier/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCourierService_CreateCourier_Success(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	req := &model.CreateCourierRequest{
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}

	mockRepo.On("CreateCourier", mock.Anything, mock.MatchedBy(func(c *model.Courier) bool {
		return c.Name == "Test Courier" &&
			c.Phone == "1234567890" &&
			c.Status == "available" &&
			c.TransportType == "car"
	})).Run(func(args mock.Arguments) {
		courier := args.Get(1).(*model.Courier)
		courier.ID = 1
	}).Return(nil)

	result, err := service.CreateCourier(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Test Courier", result.Name)
	assert.Equal(t, "1234567890", result.Phone)
	assert.Equal(t, "available", result.Status)
	assert.Equal(t, "car", result.TransportType)

	mockRepo.AssertExpectations(t)
}

func TestCourierService_CreateCourier_InvalidStatus(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	req := &model.CreateCourierRequest{
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "invalid_status",
		TransportType: "car",
	}

	result, err := service.CreateCourier(ctx, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidStatus))
	assert.Nil(t, result)

	mockRepo.AssertNotCalled(t, "CreateCourier")
}

func TestCourierService_GetCourier_Success(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	expectedCourier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}

	mockRepo.On("GetCourier", mock.Anything, int64(1)).
		Return(expectedCourier, nil)

	result, err := service.GetCourier(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedCourier, result)
	mockRepo.AssertExpectations(t)
}

func TestCourierService_GetCourier_NotFound(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetCourier", mock.Anything, int64(999)).
		Return(nil, errors.New("courier not found"))

	result, err := service.GetCourier(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCourierService_CreateCourier_InvalidTransportType(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	req := &model.CreateCourierRequest{
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "invalid_transport",
	}

	result, err := service.CreateCourier(ctx, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
	assert.Nil(t, result)

	mockRepo.AssertNotCalled(t, "CreateCourier")
}

func TestCourierService_CreateCourier_DuplicatePhone(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	req := &model.CreateCourierRequest{
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}

	mockRepo.On("CreateCourier", mock.Anything, mock.AnythingOfType("*model.Courier")).
		Return(repository.ErrDuplicatePhone)

	result, err := service.CreateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrDuplicatePhone, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCourierService_GetAllCouriers_Success(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	expectedCouriers := []*model.Courier{
		{ID: 1, Name: "Courier 1", Phone: "111", Status: "available", TransportType: "car"},
		{ID: 2, Name: "Courier 2", Phone: "222", Status: "busy", TransportType: "on_foot"},
	}

	mockRepo.On("GetAllCouriers", mock.Anything).
		Return(expectedCouriers, nil)

	result, err := service.GetAllCouriers(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedCouriers, result)
	mockRepo.AssertExpectations(t)
}

func TestCourierService_UpdateCourier_Success(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	existingCourier := &model.Courier{
		ID:            1,
		Name:          "Old Name",
		Phone:         "111",
		Status:        "available",
		TransportType: "on_foot",
	}

	newName := "New Name"
	newPhone := "222"
	newStatus := "busy"
	newTransport := "car"

	req := &model.UpdateCourierRequest{
		ID:            1,
		Name:          &newName,
		Phone:         &newPhone,
		Status:        &newStatus,
		TransportType: &newTransport,
	}

	mockRepo.On("GetCourier", mock.Anything, int64(1)).
		Return(existingCourier, nil)
	mockRepo.On("UpdateCourier", mock.Anything, mock.MatchedBy(func(c *model.Courier) bool {
		return c.Name == "New Name" &&
			c.Phone == "222" &&
			c.Status == "busy" &&
			c.TransportType == "car"
	})).Return(nil)

	result, err := service.UpdateCourier(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, "222", result.Phone)
	assert.Equal(t, "busy", result.Status)
	assert.Equal(t, "car", result.TransportType)

	mockRepo.AssertExpectations(t)
}

func TestCourierService_UpdateCourier_NotFound(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	newName := "New Name"
	req := &model.UpdateCourierRequest{
		ID:   999,
		Name: &newName,
	}

	mockRepo.On("GetCourier", mock.Anything, int64(999)).
		Return(nil, repository.ErrCourierNotFound)

	result, err := service.UpdateCourier(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, ErrCourierNotFound, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCourierService_UpdateCourier_InvalidStatus(t *testing.T) {
	mockRepo := &mocks.CourierRepository{}
	service := NewCourierService(mockRepo)
	ctx := context.Background()

	existingCourier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "111",
		Status:        "available",
		TransportType: "car",
	}

	invalidStatus := "invalid_status"
	req := &model.UpdateCourierRequest{
		ID:     1,
		Status: &invalidStatus,
	}

	mockRepo.On("GetCourier", mock.Anything, int64(1)).
		Return(existingCourier, nil)

	result, err := service.UpdateCourier(ctx, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidStatus))
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "UpdateCourier")
}
