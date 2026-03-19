package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "service-courier/proto/order"
)

type MockOrdersServiceClient struct {
	pb.OrdersServiceClient
	GetOrdersFunc func(ctx context.Context, req *pb.GetOrdersRequest, opts ...grpc.CallOption) (*pb.GetOrdersResponse, error)
}

func (m *MockOrdersServiceClient) GetOrders(ctx context.Context, req *pb.GetOrdersRequest, opts ...grpc.CallOption) (*pb.GetOrdersResponse, error) {
	return m.GetOrdersFunc(ctx, req, opts...)
}

func TestGateway_GetOrders_Success(t *testing.T) {
	mockClient := &MockOrdersServiceClient{
		GetOrdersFunc: func(ctx context.Context, req *pb.GetOrdersRequest, opts ...grpc.CallOption) (*pb.GetOrdersResponse, error) {
			assert.NotNil(t, req.From)

			return &pb.GetOrdersResponse{
				Orders: []*pb.Order{
					{
						Id:        "test-order-1",
						CreatedAt: timestamppb.New(time.Now()),
					},
				},
			}, nil
		},
	}

	gateway := &Gateway{
		client: mockClient,
	}

	ctx := context.Background()
	from := time.Now().Add(-5 * time.Minute)

	orders, err := gateway.GetOrders(ctx, from)

	assert.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, "test-order-1", orders[0].ID)
	assert.False(t, orders[0].CreatedAt.IsZero())
}

func TestGateway_GetOrders_EmptyResponse(t *testing.T) {
	mockClient := &MockOrdersServiceClient{
		GetOrdersFunc: func(ctx context.Context, req *pb.GetOrdersRequest, opts ...grpc.CallOption) (*pb.GetOrdersResponse, error) {
			return &pb.GetOrdersResponse{
				Orders: []*pb.Order{},
			}, nil
		},
	}

	gateway := &Gateway{client: mockClient}
	ctx := context.Background()

	orders, err := gateway.GetOrders(ctx, time.Now())

	assert.NoError(t, err)
	assert.Empty(t, orders)
}
