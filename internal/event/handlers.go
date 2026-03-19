package event

import (
	"context"
	"fmt"
	"service-courier/internal/model"
)

type OrderCreatedHandler struct {
	deliveryService DeliveryService
	orderGateway    OrderGateway
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, event *model.OrderEvent) error {
	if h.orderGateway != nil {
		externalOrder, err := h.orderGateway.GetOrderByID(ctx, event.OrderID)
		if err != nil {
			return fmt.Errorf("failed to verify order status: %w", err)
		}

		if externalOrder.Status != "created" {
			return fmt.Errorf("order status mismatch: expected 'created', got '%s'", externalOrder.Status)
		}
	}

	req := &model.AssignRequest{
		OrderID: event.OrderID,
	}

	_, err := h.deliveryService.AssignCourier(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to assign courier: %w", err)
	}

	return nil
}

type OrderCancelledHandler struct {
	deliveryService DeliveryService
	orderGateway    OrderGateway
}

func (h *OrderCancelledHandler) Handle(ctx context.Context, event *model.OrderEvent) error {
	if h.orderGateway != nil {
		externalOrder, err := h.orderGateway.GetOrderByID(ctx, event.OrderID)
		if err != nil {
			return fmt.Errorf("failed to verify order status: %w", err)
		}

		if externalOrder.Status != "cancelled" {
			return fmt.Errorf("order status mismatch: expected 'cancelled', got '%s'", externalOrder.Status)
		}
	}

	req := &model.UnassignRequest{
		OrderID: event.OrderID,
	}

	_, err := h.deliveryService.UnassignCourier(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to unassign courier: %w", err)
	}

	return nil
}

type OrderCompletedHandler struct {
	deliveryService DeliveryService
	orderGateway    OrderGateway
}

func (h *OrderCompletedHandler) Handle(ctx context.Context, event *model.OrderEvent) error {
	if h.orderGateway != nil {
		externalOrder, err := h.orderGateway.GetOrderByID(ctx, event.OrderID)
		if err != nil {
			return fmt.Errorf("failed to verify order status: %w", err)
		}

		if externalOrder.Status != "completed" {
			return fmt.Errorf("order status mismatch: expected 'completed', got '%s'", externalOrder.Status)
		}
	}

	return h.freeCourier(ctx, event.OrderID)
}

func (h *OrderCompletedHandler) freeCourier(ctx context.Context, orderID string) error {
	req := &model.UnassignRequest{
		OrderID: orderID,
	}

	_, err := h.deliveryService.UnassignCourier(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to free courier: %w", err)
	}

	return nil
}
