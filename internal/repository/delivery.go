package repository

import (
	"context"
	"fmt"
	"service-courier/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type deliveryRepository struct {
	pool *pgxpool.Pool
}

func NewDeliveryRepository(pool *pgxpool.Pool) DeliveryRepository {
	return &deliveryRepository{pool: pool}
}

func (r *deliveryRepository) CreateDelivery(ctx context.Context, delivery *model.Delivery) error {
	query := `
		INSERT INTO delivery (courier_id, order_id, assigned_at, deadline)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		delivery.CourierID,
		delivery.OrderID,
		delivery.AssignedAt,
		delivery.Deadline,
	).Scan(&delivery.ID)

	if err != nil {
		return fmt.Errorf("failed to create delivery: %w", err)
	}

	return nil
}

func (r *deliveryRepository) GetDeliveryByOrderID(ctx context.Context, orderID string) (*model.Delivery, error) {
	query := `
		SELECT id, courier_id, order_id, assigned_at, deadline
		FROM delivery
		WHERE order_id = $1
	`

	var delivery model.Delivery
	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&delivery.ID,
		&delivery.CourierID,
		&delivery.OrderID,
		&delivery.AssignedAt,
		&delivery.Deadline,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrDeliveryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	return &delivery, nil
}

func (r *deliveryRepository) DeleteDelivery(ctx context.Context, orderID string) error {
	query := `DELETE FROM delivery WHERE order_id = $1`

	result, err := r.pool.Exec(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete delivery: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrDeliveryNotFound
	}

	return nil
}

func (r *deliveryRepository) GetAvailableCourier(ctx context.Context) (*model.Courier, error) {
	query := `
		SELECT id, name, phone, status, transport_type
		FROM couriers
		WHERE status = 'available'
		LIMIT 1
	`

	var courier model.Courier
	err := r.pool.QueryRow(ctx, query).Scan(
		&courier.ID,
		&courier.Name,
		&courier.Phone,
		&courier.Status,
		&courier.TransportType,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNoAvailableCouriers
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get available courier: %w", err)
	}

	return &courier, nil
}

func (r *deliveryRepository) GetOverdueDeliveries(ctx context.Context) ([]*model.Delivery, error) {
	query := `
		SELECT id, courier_id, order_id, assigned_at, deadline
		FROM delivery 
		WHERE deadline < NOW()
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue deliveries: %w", err)
	}
	defer rows.Close()

	var deliveries []*model.Delivery
	for rows.Next() {
		var delivery model.Delivery
		if err := rows.Scan(
			&delivery.ID,
			&delivery.CourierID,
			&delivery.OrderID,
			&delivery.AssignedAt,
			&delivery.Deadline,
		); err != nil {
			return nil, fmt.Errorf("failed to scan delivery: %w", err)
		}
		deliveries = append(deliveries, &delivery)
	}

	return deliveries, nil
}

func (r *deliveryRepository) FreeCouriers(ctx context.Context, courierIDs []int64) error {
	if len(courierIDs) == 0 {
		return nil
	}

	query := `
        UPDATE couriers c
        SET status = 'available', updated_at = NOW()
        WHERE c.id = ANY($1) 
          AND c.status = 'busy'
          AND NOT EXISTS (
            SELECT 1 FROM delivery d
            WHERE d.courier_id = c.id
              AND d.deadline > NOW()  -- Есть другие НЕ просроченные доставки
          )
    `

	_, err := r.pool.Exec(ctx, query, courierIDs)
	if err != nil {
		return fmt.Errorf("failed to free couriers: %w", err)
	}

	return nil
}

func (r *deliveryRepository) GetAvailableCourierWithMinLoad(ctx context.Context) (*model.Courier, error) {
	query := `
		SELECT c.id, c.name, c.phone, c.status, c.transport_type,
		       COUNT(d.id) as active_deliveries
		FROM couriers c
		LEFT JOIN delivery d ON c.id = d.courier_id AND d.deadline > NOW()
		WHERE c.status = 'available'
		GROUP BY c.id
		ORDER BY active_deliveries ASC
		LIMIT 1
	`

	var courier model.Courier
	var activeDeliveries int

	err := r.pool.QueryRow(ctx, query).Scan(
		&courier.ID,
		&courier.Name,
		&courier.Phone,
		&courier.Status,
		&courier.TransportType,
		&activeDeliveries,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrNoAvailableCouriers
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get available courier: %w", err)
	}

	return &courier, nil
}
