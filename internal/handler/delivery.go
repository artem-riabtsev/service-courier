package handler

import (
	"encoding/json"
	"net/http"
	"service-courier/internal/model"
	"service-courier/internal/service"
)

func (h *Handler) assignCourier(w http.ResponseWriter, r *http.Request) {
	var req model.AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	response, err := h.deliveryService.AssignCourier(r.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			writeError(w, http.StatusBadRequest, ErrInvalidInput)
		case service.ErrNoAvailableCouriers:
			writeError(w, http.StatusConflict, "No available couriers")
		default:
			writeError(w, http.StatusInternalServerError, ErrInternalServer)
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) unassignCourier(w http.ResponseWriter, r *http.Request) {
	var req model.UnassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	response, err := h.deliveryService.UnassignCourier(r.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			writeError(w, http.StatusBadRequest, ErrInvalidInput)
		case service.ErrDeliveryNotFound:
			writeError(w, http.StatusNotFound, "Delivery not found")
		default:
			writeError(w, http.StatusInternalServerError, ErrInternalServer)
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}
