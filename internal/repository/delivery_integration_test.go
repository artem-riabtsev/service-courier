package repository

import (
	"context"
	"testing"
	"time"

	"service-courier/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeliveryRepository_CreateDelivery_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	courierRepo := NewCourierRepository(pool)
	deliveryRepo := NewDeliveryRepository(pool)
	ctx := context.Background()

	courier := &model.Courier{
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}
	err := courierRepo.CreateCourier(ctx, courier)
	require.NoError(t, err)

	delivery := &model.Delivery{
		CourierID:  courier.ID,
		OrderID:    "test-order-123",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(30 * time.Minute),
	}

	err = deliveryRepo.CreateDelivery(ctx, delivery)
	require.NoError(t, err)
	assert.Greater(t, delivery.ID, int64(0), "Delivery ID should be set")

	fetched, err := deliveryRepo.GetDeliveryByOrderID(ctx, "test-order-123")
	require.NoError(t, err)
	assert.Equal(t, delivery.ID, fetched.ID)
	assert.Equal(t, courier.ID, fetched.CourierID)
	assert.Equal(t, "test-order-123", fetched.OrderID)
}

func TestDeliveryRepository_GetDeliveryByOrderID_NotFound_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	repo := NewDeliveryRepository(pool)
	ctx := context.Background()

	delivery, err := repo.GetDeliveryByOrderID(ctx, "non-existent-order")
	assert.Error(t, err)
	assert.Equal(t, ErrDeliveryNotFound, err)
	assert.Nil(t, delivery)
}

func TestDeliveryRepository_DeleteDelivery_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	courierRepo := NewCourierRepository(pool)
	deliveryRepo := NewDeliveryRepository(pool)
	ctx := context.Background()

	courier := &model.Courier{
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}
	err := courierRepo.CreateCourier(ctx, courier)
	require.NoError(t, err)

	delivery := &model.Delivery{
		CourierID:  courier.ID,
		OrderID:    "test-order-123",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(30 * time.Minute),
	}
	err = deliveryRepo.CreateDelivery(ctx, delivery)
	require.NoError(t, err)

	err = deliveryRepo.DeleteDelivery(ctx, "test-order-123")
	require.NoError(t, err)

	deleted, err := deliveryRepo.GetDeliveryByOrderID(ctx, "test-order-123")
	assert.Error(t, err)
	assert.Equal(t, ErrDeliveryNotFound, err)
	assert.Nil(t, deleted)
}

func TestDeliveryRepository_GetAvailableCourier_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	courierRepo := NewCourierRepository(pool)
	deliveryRepo := NewDeliveryRepository(pool)
	ctx := context.Background()

	availableCourier := &model.Courier{
		Name:          "Available Courier",
		Phone:         "111",
		Status:        "available",
		TransportType: "car",
	}
	err := courierRepo.CreateCourier(ctx, availableCourier)
	require.NoError(t, err)

	busyCourier := &model.Courier{
		Name:          "Busy Courier",
		Phone:         "222",
		Status:        "busy",
		TransportType: "scooter",
	}
	err = courierRepo.CreateCourier(ctx, busyCourier)
	require.NoError(t, err)

	courier, err := deliveryRepo.GetAvailableCourier(ctx)
	require.NoError(t, err)
	assert.Equal(t, availableCourier.Name, courier.Name)
	assert.Equal(t, "available", courier.Status)
}

func TestDeliveryRepository_GetAvailableCourier_NoAvailable_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	courierRepo := NewCourierRepository(pool)
	deliveryRepo := NewDeliveryRepository(pool)
	ctx := context.Background()

	busyCourier := &model.Courier{
		Name:          "Busy Courier",
		Phone:         "111",
		Status:        "busy",
		TransportType: "car",
	}
	err := courierRepo.CreateCourier(ctx, busyCourier)
	require.NoError(t, err)

	courier, err := deliveryRepo.GetAvailableCourier(ctx)
	assert.Error(t, err)
	assert.Equal(t, ErrNoAvailableCouriers, err)
	assert.Nil(t, courier)
}

func TestDeliveryRepository_GetOverdueDeliveries_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	courierRepo := NewCourierRepository(pool)
	deliveryRepo := NewDeliveryRepository(pool)
	ctx := context.Background()

	courier := &model.Courier{
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}
	err := courierRepo.CreateCourier(ctx, courier)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '1 hour')
	`, courier.ID, "overdue-order")
	require.NoError(t, err)

	notOverdueDelivery := &model.Delivery{
		CourierID:  courier.ID,
		OrderID:    "not-overdue-order",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(2 * time.Hour),
	}
	err = deliveryRepo.CreateDelivery(ctx, notOverdueDelivery)
	require.NoError(t, err)

	overdueDeliveries, err := deliveryRepo.GetOverdueDeliveries(ctx)
	require.NoError(t, err)

	require.Len(t, overdueDeliveries, 1, "Should have exactly 1 overdue delivery")
	assert.Equal(t, "overdue-order", overdueDeliveries[0].OrderID)
}
func TestDeliveryRepository_FreeCouriers_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	courierRepo := NewCourierRepository(pool)
	deliveryRepo := NewDeliveryRepository(pool)
	ctx := context.Background()

	courier := &model.Courier{
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "busy",
		TransportType: "car",
	}
	err := courierRepo.CreateCourier(ctx, courier)
	require.NoError(t, err)

	err = deliveryRepo.FreeCouriers(ctx, []int64{courier.ID})
	require.NoError(t, err)

	freedCourier, err := courierRepo.GetCourier(ctx, courier.ID)
	require.NoError(t, err)
	assert.Equal(t, "available", freedCourier.Status)
}

func TestDeliveryRepository_GetAvailableCourierWithMinLoad_Integration(t *testing.T) {
	t.Parallel()

	pool := setupTestDatabase(t)
	courierRepo := NewCourierRepository(pool)
	deliveryRepo := NewDeliveryRepository(pool)
	ctx := context.Background()

	courier1 := &model.Courier{
		Name:          "Courier 1",
		Phone:         "111",
		Status:        "available",
		TransportType: "car",
	}
	err := courierRepo.CreateCourier(ctx, courier1)
	require.NoError(t, err)

	courier2 := &model.Courier{
		Name:          "Courier 2",
		Phone:         "222",
		Status:        "available",
		TransportType: "scooter",
	}
	err = courierRepo.CreateCourier(ctx, courier2)
	require.NoError(t, err)

	delivery := &model.Delivery{
		CourierID:  courier1.ID,
		OrderID:    "order-for-courier1",
		AssignedAt: time.Now(),
		Deadline:   time.Now().Add(1 * time.Hour),
	}
	err = deliveryRepo.CreateDelivery(ctx, delivery)
	require.NoError(t, err)

	availableCourier, err := deliveryRepo.GetAvailableCourierWithMinLoad(ctx)
	require.NoError(t, err)
	assert.Equal(t, courier2.ID, availableCourier.ID)
	assert.Equal(t, "Courier 2", availableCourier.Name)
}
