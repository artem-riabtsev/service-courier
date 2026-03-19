package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"service-courier/internal/model"
)

type EventProcessor struct {
	factory *EventHandlerFactory
}

func NewEventProcessor(factory *EventHandlerFactory) *EventProcessor {
	return &EventProcessor{
		factory: factory,
	}
}

func (p *EventProcessor) ProcessMessage(ctx context.Context, message []byte) error {
	var event model.OrderEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	slog.Info("Processing event", "order_id", event.OrderID, "status", event.Status)

	handler, err := p.factory.GetHandler(event.Status)
	if err != nil {
		return fmt.Errorf("failed to get handler: %w", err)
	}

	if err := handler.Handle(ctx, &event); err != nil {
		return fmt.Errorf("failed to handle event: %w", err)
	}

	slog.Info("Successfully processed event", "order_id", event.OrderID)
	return nil
}
