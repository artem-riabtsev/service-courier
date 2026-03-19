package handler

import (
	"net/http"

	"service-courier/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	courierService  service.CourierService
	deliveryService service.DeliveryService
}

func NewHandler(courierService service.CourierService, deliveryService service.DeliveryService) http.Handler {
	r := chi.NewRouter()
	SetupRoutes(r, courierService, deliveryService)
	return r
}

func SetupRoutes(r chi.Router, courierService service.CourierService, deliveryService service.DeliveryService) {
	h := &Handler{
		courierService:  courierService,
		deliveryService: deliveryService,
	}

	r.Get("/ping", HandlePing)
	r.Head("/healthcheck", HandleHealthCheck)

	r.Get("/couriers", h.getCouriers)
	r.Post("/courier", h.createCourier)
	r.Get("/courier/{id}", h.getCourier)
	r.Put("/courier", h.updateCourier)

	r.Post("/delivery/assign", h.assignCourier)
	r.Post("/delivery/unassign", h.unassignCourier)
}
