package factory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeliveryTimeFactory_CalculateDeadline(t *testing.T) {
	factory := NewDeliveryTimeFactory()

	tests := []struct {
		name          string
		transportType string
		wantOffset    time.Duration
	}{
		{"on_foot", "on_foot", 30 * time.Minute},
		{"scooter", "scooter", 15 * time.Minute},
		{"car", "car", 5 * time.Minute},
		{"default for unknown", "unknown", 30 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			deadline := factory.CalculateDeadline(tt.transportType)

			expected := start.Add(tt.wantOffset)
			diff := deadline.Sub(expected).Abs()

			assert.True(t, diff < time.Second,
				"CalculateDeadline(%q) = %v, want ~%v (diff: %v)",
				tt.transportType, deadline, expected, diff)
		})
	}
}
