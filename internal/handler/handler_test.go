package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"service-courier/internal/mocks"
	"service-courier/internal/model"
	"service-courier/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateCourier_Success(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{
		"name": "Test Courier",
		"phone": "1234567890",
		"status": "available",
		"transport_type": "car"
	}`

	expectedCourier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}

	mockCourierService.On("CreateCourier", mock.Anything, mock.AnythingOfType("*model.CreateCourierRequest")).
		Return(expectedCourier, nil)

	req := httptest.NewRequest("POST", "/courier", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response model.Courier
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCourier.ID, response.ID)
	assert.Equal(t, expectedCourier.Name, response.Name)
	assert.Equal(t, expectedCourier.Phone, response.Phone)
	assert.Equal(t, expectedCourier.Status, response.Status)
	assert.Equal(t, expectedCourier.TransportType, response.TransportType)

	mockCourierService.AssertExpectations(t)
}

func TestHandler_CreateCourier_InvalidJSON(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{invalid json`
	req := httptest.NewRequest("POST", "/courier", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockCourierService.AssertNotCalled(t, "CreateCourier")
}

func TestHandler_GetCourier_Success(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	expectedCourier := &model.Courier{
		ID:            1,
		Name:          "Test Courier",
		Phone:         "1234567890",
		Status:        "available",
		TransportType: "car",
	}

	mockCourierService.On("GetCourier", mock.Anything, int64(1)).
		Return(expectedCourier, nil)

	req := httptest.NewRequest("GET", "/courier/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.Courier
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCourier.ID, response.ID)
	mockCourierService.AssertExpectations(t)
}

func TestHandler_GetAllCouriers_Success(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	expectedCouriers := []*model.Courier{
		{ID: 1, Name: "Courier 1", Phone: "111", Status: "available", TransportType: "car"},
		{ID: 2, Name: "Courier 2", Phone: "222", Status: "busy", TransportType: "on_foot"},
	}

	mockCourierService.On("GetAllCouriers", mock.Anything).
		Return(expectedCouriers, nil)

	req := httptest.NewRequest("GET", "/couriers", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []model.Courier
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, expectedCouriers[0].Name, response[0].Name)

	mockCourierService.AssertExpectations(t)
}

func TestHandler_AssignDelivery_Success(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{"order_id": "test-order-123"}`

	expectedResponse := &model.AssignResponse{
		CourierID:        1,
		OrderID:          "test-order-123",
		TransportType:    "car",
		DeliveryDeadline: time.Now().Add(30 * time.Minute),
	}

	mockDeliveryService.On("AssignCourier", mock.Anything, mock.AnythingOfType("*model.AssignRequest")).
		Return(expectedResponse, nil)

	req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.AssignResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.CourierID, response.CourierID)
	assert.Equal(t, expectedResponse.OrderID, response.OrderID)

	mockDeliveryService.AssertExpectations(t)
}

func TestHandler_UpdateCourier_Success(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{
		"id": 1,
		"name": "Updated Name",
		"phone": "9999999999",
		"status": "busy",
		"transport_type": "car"
	}`

	expectedCourier := &model.Courier{
		ID:            1,
		Name:          "Updated Name",
		Phone:         "9999999999",
		Status:        "busy",
		TransportType: "car",
	}

	mockCourierService.On("UpdateCourier", mock.Anything, mock.AnythingOfType("*model.UpdateCourierRequest")).
		Return(expectedCourier, nil)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.Courier
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCourier.ID, response.ID)
	assert.Equal(t, expectedCourier.Name, response.Name)
	assert.Equal(t, expectedCourier.Phone, response.Phone)
	assert.Equal(t, expectedCourier.Status, response.Status)
	assert.Equal(t, expectedCourier.TransportType, response.TransportType)

	mockCourierService.AssertExpectations(t)
}

func TestHandler_UpdateCourier_NotFound(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{
		"id": 999,
		"name": "Non-existent Courier"
	}`

	mockCourierService.On("UpdateCourier", mock.Anything, mock.AnythingOfType("*model.UpdateCourierRequest")).
		Return(nil, service.ErrCourierNotFound)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Courier not found")

	mockCourierService.AssertExpectations(t)
}

func TestHandler_UpdateCourier_InvalidStatus(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{
		"id": 1,
		"status": "invalid_status"
	}`

	mockCourierService.On("UpdateCourier", mock.Anything, mock.AnythingOfType("*model.UpdateCourierRequest")).
		Return(nil, service.ErrInvalidStatus)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid status")

	mockCourierService.AssertExpectations(t)
}

func TestHandler_UpdateCourier_DuplicatePhone(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{
		"id": 1,
		"phone": "1234567890"
	}`

	mockCourierService.On("UpdateCourier", mock.Anything, mock.AnythingOfType("*model.UpdateCourierRequest")).
		Return(nil, service.ErrDuplicatePhone)

	req := httptest.NewRequest("PUT", "/courier", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "already exists")

	mockCourierService.AssertExpectations(t)
}

func TestHandler_UpdateCourier_InvalidJSON(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{invalid json`
	req := httptest.NewRequest("PUT", "/courier", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockCourierService.AssertNotCalled(t, "UpdateCourier")
}

func TestHandler_UnassignCourier_Success(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{"order_id": "test-order-123"}`

	expectedResponse := &model.UnassignResponse{
		OrderID:   "test-order-123",
		Status:    "unassigned",
		CourierID: 1,
	}

	mockDeliveryService.On("UnassignCourier", mock.Anything, mock.AnythingOfType("*model.UnassignRequest")).
		Return(expectedResponse, nil)

	req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.UnassignResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.OrderID, response.OrderID)
	assert.Equal(t, expectedResponse.Status, response.Status)
	assert.Equal(t, expectedResponse.CourierID, response.CourierID)

	mockDeliveryService.AssertExpectations(t)
}

func TestHandler_UnassignCourier_NotFound(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{"order_id": "non-existent-order"}`

	mockDeliveryService.On("UnassignCourier", mock.Anything, mock.AnythingOfType("*model.UnassignRequest")).
		Return(nil, service.ErrDeliveryNotFound)

	req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Delivery not found")

	mockDeliveryService.AssertExpectations(t)
}

func TestHandler_UnassignCourier_InvalidInput(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{"order_id": ""}`

	mockDeliveryService.On("UnassignCourier", mock.Anything, mock.AnythingOfType("*model.UnassignRequest")).
		Return(nil, service.ErrInvalidInput)

	req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid input")

	mockDeliveryService.AssertExpectations(t)
}

func TestHandler_UnassignCourier_InvalidJSON(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{invalid json`
	req := httptest.NewRequest("POST", "/delivery/unassign", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockDeliveryService.AssertNotCalled(t, "UnassignCourier")
}

func TestHandler_GetCourier_InvalidID(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	req := httptest.NewRequest("GET", "/courier/not-a-number", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid courier ID")

	mockCourierService.AssertNotCalled(t, "GetCourier")
}

func TestHandler_GetCourier_NotFound(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	mockCourierService.On("GetCourier", mock.Anything, int64(999)).
		Return(nil, service.ErrCourierNotFound)

	req := httptest.NewRequest("GET", "/courier/999", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Courier not found")

	mockCourierService.AssertExpectations(t)
}

func TestHandler_CreateCourier_InvalidStatus(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{
		"name": "Test Courier",
		"phone": "1234567890",
		"status": "invalid_status",
		"transport_type": "car"
	}`

	mockCourierService.On("CreateCourier", mock.Anything, mock.AnythingOfType("*model.CreateCourierRequest")).
		Return(nil, service.ErrInvalidStatus)

	req := httptest.NewRequest("POST", "/courier", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid status")

	mockCourierService.AssertExpectations(t)
}

func TestHandler_CreateCourier_DuplicatePhone(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{
		"name": "Test Courier",
		"phone": "1234567890",
		"status": "available",
		"transport_type": "car"
	}`

	mockCourierService.On("CreateCourier", mock.Anything, mock.AnythingOfType("*model.CreateCourierRequest")).
		Return(nil, service.ErrDuplicatePhone)

	req := httptest.NewRequest("POST", "/courier", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "already exists")

	mockCourierService.AssertExpectations(t)
}

func TestHandler_AssignDelivery_NoAvailableCouriers(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	reqBody := `{"order_id": "test-order-123"}`

	mockDeliveryService.On("AssignCourier", mock.Anything, mock.AnythingOfType("*model.AssignRequest")).
		Return(nil, service.ErrNoAvailableCouriers)

	req := httptest.NewRequest("POST", "/delivery/assign", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "No available couriers")

	mockDeliveryService.AssertExpectations(t)
}

func TestHandler_HandlePing(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	req := httptest.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "pong", response["message"])
}

func TestHandler_HandleHealthCheck(t *testing.T) {
	t.Parallel()
	mockCourierService := &mocks.CourierService{}
	mockDeliveryService := &mocks.DeliveryService{}

	router := NewHandler(mockCourierService, mockDeliveryService)

	req := httptest.NewRequest("HEAD", "/healthcheck", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}
