package repository

import (
	"context"
	"testing"
	"time"

	"service-courier/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")

	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx), "Failed to terminate container")
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get connection string")

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err, "Failed to create connection pool")

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS couriers (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			phone VARCHAR(20) UNIQUE NOT NULL,
			status VARCHAR(50) NOT NULL,
			transport_type VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS delivery (
			id BIGSERIAL PRIMARY KEY,
			courier_id BIGINT NOT NULL,
			order_id VARCHAR(255) UNIQUE NOT NULL,
			assigned_at TIMESTAMP NOT NULL,
			deadline TIMESTAMP NOT NULL
		);
	`)
	require.NoError(t, err, "Failed to create tables")

	return pool
}

func TestCourierRepository_CreateCourier_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	repo := NewCourierRepository(pool)
	ctx := context.Background()

	courier := &model.Courier{
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}

	err := repo.CreateCourier(ctx, courier)
	require.NoError(t, err)
	assert.Greater(t, courier.ID, int64(0), "Courier ID should be set")

	fetched, err := repo.GetCourier(ctx, courier.ID)
	require.NoError(t, err)
	assert.Equal(t, courier.ID, fetched.ID)
	assert.Equal(t, "Test Courier", fetched.Name)
	assert.Equal(t, "1234567890", fetched.Phone)
	assert.Equal(t, "available", fetched.Status)
	assert.Equal(t, "car", fetched.TransportType)
}

func TestCourierRepository_CreateCourier_DuplicatePhone_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	repo := NewCourierRepository(pool)
	ctx := context.Background()

	courier1 := &model.Courier{
		Name:          "Courier 1",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}

	err := repo.CreateCourier(ctx, courier1)
	require.NoError(t, err)

	courier2 := &model.Courier{
		Name:          "Courier 2",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "scooter",
	}

	err = repo.CreateCourier(ctx, courier2)
	assert.Error(t, err)
	assert.Equal(t, ErrDuplicatePhone, err)
}

func TestCourierRepository_GetCourier_NotFound_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	repo := NewCourierRepository(pool)
	ctx := context.Background()

	courier, err := repo.GetCourier(ctx, 999)
	assert.Error(t, err)
	assert.Equal(t, ErrCourierNotFound, err)
	assert.Nil(t, courier)
}

func TestCourierRepository_GetAllCouriers_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	repo := NewCourierRepository(pool)
	ctx := context.Background()

	couriersToCreate := []*model.Courier{
		{Name: "Courier 1", Phone: "111", Status: "available", TransportType: "car"},
		{Name: "Courier 2", Phone: "222", Status: "busy", TransportType: "scooter"},
		{Name: "Courier 3", Phone: "333", Status: "available", TransportType: "on_foot"},
	}

	for _, courier := range couriersToCreate {
		err := repo.CreateCourier(ctx, courier)
		require.NoError(t, err)
	}

	couriers, err := repo.GetAllCouriers(ctx)
	require.NoError(t, err)
	assert.Len(t, couriers, 3)

	for i, courier := range couriers {
		assert.Greater(t, courier.ID, int64(0))
		assert.Equal(t, couriersToCreate[i].Name, courier.Name)
		assert.Equal(t, couriersToCreate[i].Phone, courier.Phone)
		assert.Equal(t, couriersToCreate[i].Status, courier.Status)
		assert.Equal(t, couriersToCreate[i].TransportType, courier.TransportType)
	}
}

func TestCourierRepository_UpdateCourier_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	repo := NewCourierRepository(pool)
	ctx := context.Background()

	courier := &model.Courier{
		Name:          "Old Name",
		Phone:         "111",
		Status:        "available",
		TransportType: "car",
	}
	err := repo.CreateCourier(ctx, courier)
	require.NoError(t, err)

	courier.Name = "New Name"
	courier.Phone = "222"
	courier.Status = "busy"
	courier.TransportType = "scooter"

	err = repo.UpdateCourier(ctx, courier)
	require.NoError(t, err)

	updated, err := repo.GetCourier(ctx, courier.ID)
	require.NoError(t, err)

	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "222", updated.Phone)
	assert.Equal(t, "busy", updated.Status)
	assert.Equal(t, "scooter", updated.TransportType)
}

func TestCourierRepository_UpdateCourier_NotFound_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	repo := NewCourierRepository(pool)
	ctx := context.Background()

	nonExistentCourier := &model.Courier{
		ID:            999,
		Name:          "Non-existent",
		Phone:         "999",
		Status:        "available",
		TransportType: "car",
	}

	err := repo.UpdateCourier(ctx, nonExistentCourier)
	assert.Error(t, err)
	assert.Equal(t, ErrCourierNotFound, err)
}
