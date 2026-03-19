package database

import (
	"context"
	"testing"

	"service-courier/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestNew_InvalidConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg := config.DatabaseConfig{
		Host:     "invalid-host-that-definitely-does-not-exist-12345",
		Port:     9999,
		User:     "invalid",
		Password: "invalid",
		DBName:   "invalid",
	}

	pool, err := New(ctx, cfg)
	assert.Error(t, err)
	assert.Nil(t, pool)

	assert.Contains(t, err.Error(), "unable to")
}

func TestNew_ValidConfigButUnreachable(t *testing.T) {
	t.Parallel()

	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "test",
	}

	ctx := context.Background()

	_, err := New(ctx, cfg)

	assert.Error(t, err)
}
