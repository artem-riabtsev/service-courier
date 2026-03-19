package model

import (
	"time"
)

type Courier struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	Status        string    `json:"status"`
	TransportType string    `json:"transport_type"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

type CreateCourierRequest struct {
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Status        string `json:"status"`
	TransportType string `json:"transport_type"`
}

type UpdateCourierRequest struct {
	ID            int64   `json:"id"`
	Name          *string `json:"name,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Status        *string `json:"status,omitempty"`
	TransportType *string `json:"transport_type,omitempty"`
}

type Delivery struct {
	ID         int64     `json:"id"`
	CourierID  int64     `json:"courier_id"`
	OrderID    string    `json:"order_id"`
	AssignedAt time.Time `json:"assigned_at"`
	Deadline   time.Time `json:"deadline"`
}

type AssignRequest struct {
	OrderID string `json:"order_id"`
}

type AssignResponse struct {
	CourierID        int64     `json:"courier_id"`
	OrderID          string    `json:"order_id"`
	TransportType    string    `json:"transport_type"`
	DeliveryDeadline time.Time `json:"delivery_deadline"`
}

type UnassignRequest struct {
	OrderID string `json:"order_id"`
}

type UnassignResponse struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	CourierID int64  `json:"courier_id"`
}

type CourierStats struct {
	CourierID        int64 `json:"courier_id"`
	ActiveDeliveries int   `json:"active_deliveries"`
}
