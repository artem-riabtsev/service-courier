package service

import "errors"

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrCourierNotFound     = errors.New("courier not found")
	ErrDuplicatePhone      = errors.New("courier with this phone already exists")
	ErrDeliveryNotFound    = errors.New("delivery not found")
	ErrNoAvailableCouriers = errors.New("no available couriers")
)
