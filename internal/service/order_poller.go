package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"service-courier/internal/gateway/order"
	"service-courier/internal/model"
)

type OrderPoller struct {
	gateway    order.OrderGateway
	delivery   DeliveryService
	interval   time.Duration
	lastCursor time.Time
	mu         sync.RWMutex
}

func NewOrderPoller(gateway order.OrderGateway, delivery DeliveryService) *OrderPoller {
	initialCursor := time.Now()

	return &OrderPoller{
		gateway:    gateway,
		delivery:   delivery,
		interval:   5 * time.Second,
		lastCursor: initialCursor,
	}
}

func (p *OrderPoller) Start(ctx context.Context) {
	slog.Info("Order poller started", "interval", p.interval)

	go p.poll(ctx)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Order poller stopped")
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *OrderPoller) poll(ctx context.Context) {
	now := time.Now()
	from := now.Add(-5 * time.Second)

	p.mu.Lock()
	p.lastCursor = now
	p.mu.Unlock()

	slog.Debug("Polling orders", "from", from, "to", now)

	orders, err := p.gateway.GetOrders(ctx, from)
	if err != nil {
		slog.Error("Failed to get orders from gateway", "error", err)
		return
	}

	if len(orders) == 0 {
		slog.Debug("No new orders")
		return
	}

	slog.Info("Got new orders", "count", len(orders))

	successCount := 0
	failCount := 0

	for _, order := range orders {
		req := &model.AssignRequest{
			OrderID: order.ID,
		}

		_, err := p.delivery.AssignCourier(ctx, req)
		if err != nil {
			slog.Error("Failed to assign courier", "order_id", order.ID, "error", err)
			failCount++
			continue
		}

		slog.Info("Successfully assigned courier", "order_id", order.ID, "created_at", order.CreatedAt)
		successCount++
	}

	if successCount > 0 || failCount > 0 {
		slog.Info("Assignment summary", "succeeded", successCount, "failed", failCount)
	}
}

func (p *OrderPoller) SetInterval(interval time.Duration) {
	p.interval = interval
}

func (p *OrderPoller) GetLastCursor() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastCursor
}
