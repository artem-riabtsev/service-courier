package model

import (
	"encoding/json"
	"errors"
	"time"
)

type OrderEvent struct {
	OrderID   string    `json:"order_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrInvalidInput = errors.New("invalid input")
)

func (e *OrderEvent) Validate() error {
	if e.OrderID == "" {
		return ErrInvalidInput
	}
	if e.Status == "" {
		return ErrInvalidInput
	}
	return nil
}

func (e *OrderEvent) UnmarshalJSON(data []byte) error {
	type Alias OrderEvent
	aux := &struct {
		CreatedAt string `json:"created_at"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.CreatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, aux.CreatedAt)
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02T15:04:05Z", aux.CreatedAt)
			if err != nil {
				return err
			}
		}
		e.CreatedAt = parsedTime
	}

	return nil
}
