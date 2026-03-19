package order

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "status 429",
			err:      fmt.Errorf("order service returned status 429"),
			expected: true,
		},
		{
			name:     "status 500",
			err:      fmt.Errorf("order service returned status 500"),
			expected: true,
		},
		{
			name:     "status 502",
			err:      fmt.Errorf("order service returned status 502"),
			expected: true,
		},
		{
			name:     "connection timeout",
			err:      fmt.Errorf("connection timeout"),
			expected: true,
		},
		{
			name:     "connection refused",
			err:      fmt.Errorf("connection refused"),
			expected: true,
		},
		{
			name:     "status 404",
			err:      fmt.Errorf("order service returned status 404"),
			expected: false,
		},
		{
			name:     "other error",
			err:      fmt.Errorf("some other error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, shouldRetry(tt.err))
		})
	}
}

func TestContains(t *testing.T) {
	assert.True(t, contains("hello world", "world"))
	assert.True(t, contains("status 429", "429"))
	assert.False(t, contains("hello", "world"))
	assert.False(t, contains("", "test"))
	assert.True(t, contains("status 502", "status 5"))
}
