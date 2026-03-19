package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"service-courier/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type courierRepository struct {
	pool *pgxpool.Pool
}

func NewCourierRepository(pool *pgxpool.Pool) CourierRepository {
	return &courierRepository{pool: pool}
}

func (r *courierRepository) CreateCourier(ctx context.Context, courier *model.Courier) error {
	query := `
		INSERT INTO couriers (name, phone, status, transport_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		courier.Name,
		courier.Phone,
		courier.Status,
		courier.TransportType,
		time.Now(),
		time.Now(),
	).Scan(&courier.ID)

	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrDuplicatePhone
		}
		return fmt.Errorf("failed to create courier: %w", err)
	}

	return nil
}

func (r *courierRepository) GetCourier(ctx context.Context, id int64) (*model.Courier, error) {
	query := `
		SELECT id, name, phone, status, transport_type, created_at, updated_at
		FROM couriers
		WHERE id = $1
	`

	var courier model.Courier
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&courier.ID,
		&courier.Name,
		&courier.Phone,
		&courier.Status,
		&courier.TransportType,
		&courier.CreatedAt,
		&courier.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCourierNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get courier: %w", err)
	}

	return &courier, nil
}

func (r *courierRepository) GetAllCouriers(ctx context.Context) ([]*model.Courier, error) {
	query := `
		SELECT id, name, phone, status, transport_type, created_at, updated_at
		FROM couriers
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get couriers: %w", err)
	}
	defer rows.Close()

	var couriers []*model.Courier
	for rows.Next() {
		var courier model.Courier
		if err := rows.Scan(
			&courier.ID,
			&courier.Name,
			&courier.Phone,
			&courier.Status,
			&courier.TransportType,
			&courier.CreatedAt,
			&courier.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan courier: %w", err)
		}
		couriers = append(couriers, &courier)
	}

	return couriers, nil
}

func (r *courierRepository) UpdateCourier(ctx context.Context, courier *model.Courier) error {
	query := `
		UPDATE couriers 
		SET name = $1, phone = $2, status = $3, transport_type = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.pool.Exec(ctx, query,
		courier.Name,
		courier.Phone,
		courier.Status,
		courier.TransportType,
		time.Now(),
		courier.ID,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrDuplicatePhone
		}
		return fmt.Errorf("failed to update courier: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrCourierNotFound
	}

	return nil
}

func isDuplicateKeyError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}
