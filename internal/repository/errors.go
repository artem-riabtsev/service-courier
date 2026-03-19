package repository

import "errors"

var (
	ErrCourierNotFound     = errors.New("courier not found")
	ErrDuplicatePhone      = errors.New("courier with this phone already exists")
	ErrDeliveryNotFound    = errors.New("delivery not found")
	ErrNoAvailableCouriers = errors.New("no available couriers")
)
