package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "service-courier/proto/order"
)

type Gateway struct {
	client pb.OrdersServiceClient
	conn   *grpc.ClientConn
}

func NewGateway(addr string) (*Gateway, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}

	client := pb.NewOrdersServiceClient(conn)

	return &Gateway{
		client: client,
		conn:   conn,
	}, nil
}

func (g *Gateway) Close() error {
	if g.conn != nil {
		return g.conn.Close()
	}
	return nil
}

func (g *Gateway) GetOrders(ctx context.Context, from time.Time) ([]*ExternalOrder, error) {
	req := &pb.GetOrdersRequest{
		From: timestamppb.New(from),
	}

	resp, err := g.client.GetOrders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	var orders []*ExternalOrder
	for _, pbOrder := range resp.Orders {
		order := &ExternalOrder{
			ID: pbOrder.Id,
		}

		if pbOrder.CreatedAt != nil {
			order.CreatedAt = pbOrder.CreatedAt.AsTime()
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (g *Gateway) GetOrderByID(ctx context.Context, orderID string) (*ExternalOrder, error) {
	return g.getOrderByIDWithRetry(ctx, orderID)
}

func (g *Gateway) getOrderByIDDirect(ctx context.Context, orderID string) (*ExternalOrder, error) {
	url := fmt.Sprintf("http://service-order:8080/public/api/v1/order/%s", orderID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("order service returned status %d", resp.StatusCode)
	}

	var orderData struct {
		OrderID   string    `json:"id"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&orderData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ExternalOrder{
		ID:        orderData.OrderID,
		Status:    orderData.Status,
		CreatedAt: orderData.CreatedAt,
	}, nil
}
