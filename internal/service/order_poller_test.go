package service

import (
	"context"
	"testing"
	"time"

	"service-courier/internal/gateway/order"
	"service-courier/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderGateway struct {
	mock.Mock
}

func (m *MockOrderGateway) GetOrders(ctx context.Context, from time.Time) ([]*order.ExternalOrder, error) {
	args := m.Called(ctx, from)
	return args.Get(0).([]*order.ExternalOrder), args.Error(1)
}

func (m *MockOrderGateway) GetOrderByID(ctx context.Context, orderID string) (*order.ExternalOrder, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.ExternalOrder), args.Error(1)
}

func (m *MockOrderGateway) Close() error {
	args := m.Called()
	return args.Error(0)
}

type MockDeliveryService struct {
	mock.Mock
}

func (m *MockDeliveryService) AssignCourier(ctx context.Context, req *model.AssignRequest) (*model.AssignResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AssignResponse), args.Error(1)
}

func (m *MockDeliveryService) UnassignCourier(ctx context.Context, req *model.UnassignRequest) (*model.UnassignResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UnassignResponse), args.Error(1)
}

func TestOrderPoller_NewOrderPoller(t *testing.T) {
	mockGateway := &MockOrderGateway{}
	mockDelivery := &MockDeliveryService{}

	poller := NewOrderPoller(mockGateway, mockDelivery)

	assert.NotNil(t, poller)
	assert.Equal(t, 5*time.Second, poller.interval)
}

func TestOrderPoller_Poll_WithOrders(t *testing.T) {
	ctx := context.Background()
	mockGateway := &MockOrderGateway{}
	mockDelivery := &MockDeliveryService{}

	now := time.Now()
	testOrders := []*order.ExternalOrder{
		{ID: "order-1", CreatedAt: now.Add(-2 * time.Second)},
		{ID: "order-2", CreatedAt: now.Add(-1 * time.Second)},
	}

	mockGateway.On("GetOrders", mock.Anything, mock.Anything).
		Return(testOrders, nil)

	mockDelivery.On("AssignCourier", mock.Anything, mock.MatchedBy(func(req *model.AssignRequest) bool {
		return req.OrderID == "order-1" || req.OrderID == "order-2"
	})).Return(&model.AssignResponse{}, nil).Times(2)

	poller := NewOrderPoller(mockGateway, mockDelivery)
	poller.poll(ctx)

	mockGateway.AssertExpectations(t)
	mockDelivery.AssertExpectations(t)
}

func TestOrderPoller_Poll_NoOrders(t *testing.T) {
	ctx := context.Background()
	mockGateway := &MockOrderGateway{}
	mockDelivery := &MockDeliveryService{}

	mockGateway.On("GetOrders", mock.Anything, mock.Anything).
		Return([]*order.ExternalOrder{}, nil)

	poller := NewOrderPoller(mockGateway, mockDelivery)
	poller.poll(ctx)

	mockGateway.AssertExpectations(t)
	mockDelivery.AssertNotCalled(t, "AssignCourier")
}

func TestOrderPoller_Poll_GatewayError(t *testing.T) {
	ctx := context.Background()
	mockGateway := &MockOrderGateway{}
	mockDelivery := &MockDeliveryService{}

	mockGateway.On("GetOrders", mock.Anything, mock.Anything).
		Return([]*order.ExternalOrder{}, assert.AnError)

	poller := NewOrderPoller(mockGateway, mockDelivery)
	poller.poll(ctx)

	mockGateway.AssertExpectations(t)
	mockDelivery.AssertNotCalled(t, "AssignCourier")
}

func TestOrderPoller_Poll_AssignCourierError(t *testing.T) {
	ctx := context.Background()
	mockGateway := &MockOrderGateway{}
	mockDelivery := &MockDeliveryService{}

	testOrders := []*order.ExternalOrder{
		{ID: "order-1", CreatedAt: time.Now()},
	}

	mockGateway.On("GetOrders", mock.Anything, mock.Anything).
		Return(testOrders, nil)

	mockDelivery.On("AssignCourier", mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	poller := NewOrderPoller(mockGateway, mockDelivery)
	poller.poll(ctx)

	mockGateway.AssertExpectations(t)
	mockDelivery.AssertExpectations(t)
}

func TestOrderPoller_StartAndStop(t *testing.T) {
	mockGateway := &MockOrderGateway{}
	mockDelivery := &MockDeliveryService{}

	poller := NewOrderPoller(mockGateway, mockDelivery)
	poller.SetInterval(100 * time.Millisecond)

	mockGateway.On("GetOrders", mock.Anything, mock.Anything).
		Return([]*order.ExternalOrder{}, nil).
		Maybe()

	ctx, cancel := context.WithCancel(context.Background())

	go poller.Start(ctx)

	time.Sleep(350 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)

	assert.Greater(t, len(mockGateway.Calls), 0, "GetOrders should have been called")
}

func TestOrderPoller_CursorLogic(t *testing.T) {
	mockGateway := &MockOrderGateway{}
	mockDelivery := &MockDeliveryService{}

	poller := NewOrderPoller(mockGateway, mockDelivery)

	initialCursor := poller.GetLastCursor()
	assert.True(t, time.Since(initialCursor) < time.Second)
}
