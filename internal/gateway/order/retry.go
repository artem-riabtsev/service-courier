package order

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	gatewayRetriesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gateway_retries_total",
			Help: "Total number of gateway retry attempts",
		},
	)
)

func (g *Gateway) getOrderByIDWithRetry(ctx context.Context, orderID string) (*ExternalOrder, error) {
	maxAttempts := 3
	baseDelay := 100 * time.Millisecond

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		order, err := g.getOrderByIDDirect(ctx, orderID)
		if err == nil {
			return order, nil
		}

		lastErr = err

		if !shouldRetry(err) || attempt == maxAttempts {
			break
		}

		delay := baseDelay * time.Duration(1<<uint(attempt-1))
		if delay > 2*time.Second {
			delay = 2 * time.Second
		}

		gatewayRetriesTotal.Inc()

		select {
		case <-time.After(delay):
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	return contains(errStr, "status 429") ||
		contains(errStr, "status 5") ||
		contains(errStr, "timeout") ||
		contains(errStr, "connection refused")
}

func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
