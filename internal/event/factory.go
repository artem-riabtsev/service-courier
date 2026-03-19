package event

import (
	"context"
	"fmt"
	"service-courier/internal/model"
)

type DeliveryService interface {
	AssignCourier(ctx context.Context, req *model.AssignRequest) (*model.AssignResponse, error)
	UnassignCourier(ctx context.Context, req *model.UnassignRequest) (*model.UnassignResponse, error)
}

type EventHandler interface {
	Handle(ctx context.Context, event *model.OrderEvent) error
}

type EventHandlerFactory struct {
	deliveryService DeliveryService
	orderGateway    OrderGateway
}

func NewEventHandlerFactory(deliveryService DeliveryService, orderGateway OrderGateway) *EventHandlerFactory {
	return &EventHandlerFactory{
		deliveryService: deliveryService,
		orderGateway:    orderGateway,
	}
}

func (f *EventHandlerFactory) GetHandler(status string) (EventHandler, error) {
	switch status {
	case "created", "confirmed":
		return &OrderCreatedHandler{
			deliveryService: f.deliveryService,
			orderGateway:    f.orderGateway,
		}, nil
	case "cancelled":
		return &OrderCancelledHandler{
			deliveryService: f.deliveryService,
			orderGateway:    f.orderGateway,
		}, nil
	case "completed":
		return &OrderCompletedHandler{
			deliveryService: f.deliveryService,
			orderGateway:    f.orderGateway,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported order status: %s", status)
	}
}
