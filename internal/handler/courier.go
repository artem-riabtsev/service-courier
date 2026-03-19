package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"service-courier/internal/model"
	"service-courier/internal/service"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) getCourier(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidCourierID)
		return
	}

	courier, err := h.courierService.GetCourier(r.Context(), id)
	if err != nil {
		switch err {
		case service.ErrInvalidInput:
			writeError(w, http.StatusBadRequest, ErrInvalidInput)
		case service.ErrCourierNotFound:
			writeError(w, http.StatusNotFound, ErrCourierNotFound)
		default:
			writeError(w, http.StatusInternalServerError, ErrInternalServer)
		}
		return
	}

	writeJSON(w, http.StatusOK, courier)
}

func (h *Handler) getCouriers(w http.ResponseWriter, r *http.Request) {
	couriers, err := h.courierService.GetAllCouriers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	writeJSON(w, http.StatusOK, couriers)
}

func (h *Handler) createCourier(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	courier, err := h.courierService.CreateCourier(r.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrInvalidInput, service.ErrInvalidStatus:
			writeError(w, http.StatusBadRequest, err.Error())
		case service.ErrDuplicatePhone:
			writeError(w, http.StatusConflict, ErrDuplicatePhone)
		default:
			writeError(w, http.StatusInternalServerError, ErrInternalServer)
		}
		return
	}

	writeJSON(w, http.StatusCreated, courier)
}

func (h *Handler) updateCourier(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateCourierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	courier, err := h.courierService.UpdateCourier(r.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrInvalidInput, service.ErrInvalidStatus:
			writeError(w, http.StatusBadRequest, err.Error())
		case service.ErrCourierNotFound:
			writeError(w, http.StatusNotFound, ErrCourierNotFound)
		case service.ErrDuplicatePhone:
			writeError(w, http.StatusConflict, ErrDuplicatePhone)
		default:
			writeError(w, http.StatusInternalServerError, ErrInternalServer)
		}
		return
	}

	writeJSON(w, http.StatusOK, courier)
}

// Вспомогательные функции для ответов
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
