package factory

import (
	"time"
)

type DeliveryTimeFactory struct{}

func NewDeliveryTimeFactory() *DeliveryTimeFactory {
	return &DeliveryTimeFactory{}
}

func (f *DeliveryTimeFactory) CalculateDeadline(transportType string) time.Time {
	now := time.Now()

	switch transportType {
	case "on_foot":
		return now.Add(30 * time.Minute)
	case "scooter":
		return now.Add(15 * time.Minute)
	case "car":
		return now.Add(5 * time.Minute)
	default:
		return now.Add(30 * time.Minute) // default
	}
}
