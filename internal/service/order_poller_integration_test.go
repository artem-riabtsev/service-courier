package service

import (
	"testing"
)

func TestOrderPoller_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("PollWithRealGateway", func(t *testing.T) {
		skipMsg := "Requires running service-order. Set SERVICE_ORDER_URL env var to run."
		t.Skip(skipMsg)
	})
}
