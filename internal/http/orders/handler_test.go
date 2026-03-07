package orders

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lgutierrez148/acomm/internal/inbound"
	"github.com/lgutierrez148/acomm/internal/mocks"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/lgutierrez148/acomm/internal/outbound"
	"github.com/stretchr/testify/assert"
)

func TestNewOrdersHandler(t *testing.T) {
	repo := mocks.NewMockIOrdersRepository(t)
	handler := NewOrdersHandler(repo)
	assert.NotNil(t, handler)
	assert.Equal(t, repo, handler.repo)
}

func TestHandleGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		ords := []models.Order{
			{ID: 1, CustomerEmail: "test@test.com"},
			{ID: 2, CustomerEmail: "test2@test.com"},
		}

		repo.EXPECT().FindAll().Return(ords, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetOrdersResponse{Orders: mapToOrdersResponse(ords)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		repo.EXPECT().FindAll().Return(nil, errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		ord := &models.Order{ID: 1, CustomerEmail: "test@test.com"}

		repo.EXPECT().FindByID(1).Return(ord, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetOrderResponse{Order: mapToOrderResponse(ord)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders/invalid", nil)
		req.SetPathValue("id", "invalid")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

func TestHandleGetByCustomerEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		ords := []models.Order{
			{ID: 1, CustomerEmail: "test@test.com"},
		}

		repo.EXPECT().FindByCustomerEmail("test@test.com").Return(ords, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders/email/test@test.com", nil)
		req.SetPathValue("email", "test@test.com")

		handler.HandleGetByCustomerEmail(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		expectedResp := outbound.GetOrdersResponse{Orders: mapToOrdersResponse(ords)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})
}

func TestHandleCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		reqBody := inbound.CreateOrderRequest{
			CustomerEmail: "new@test.com",
			CustomerName:  "Test",
		}

		bodyBytes, _ := json.Marshal(reqBody)
		ordToCreate := reqBody.ToDomain()

		repo.EXPECT().Create(ordToCreate).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
