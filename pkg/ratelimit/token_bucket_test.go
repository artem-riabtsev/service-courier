package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_Simple(t *testing.T) {
	tb := NewTokenBucket(5, 1)

	for i := 0; i < 5; i++ {
		assert.True(t, tb.Allow())
	}

	assert.False(t, tb.Allow())

	time.Sleep(1 * time.Second)

	assert.True(t, tb.Allow())
	assert.False(t, tb.Allow())
}
