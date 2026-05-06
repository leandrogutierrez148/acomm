package orders

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leandrogutierrez148/acomm/internal/inbound"
	"github.com/leandrogutierrez148/acomm/internal/mocks"
	"github.com/leandrogutierrez148/acomm/internal/models"
	"github.com/leandrogutierrez148/acomm/internal/outbound"
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

	t.Run("empty id", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders/", nil)
		req.SetPathValue("id", "")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
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

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		repo.EXPECT().FindByID(1).Return(nil, errors.New("not found"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
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

	t.Run("empty email", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders/email/", nil)
		req.SetPathValue("email", "")

		handler.HandleGetByCustomerEmail(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		repo.EXPECT().FindByCustomerEmail("test@test.com").Return(nil, errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders/email/test@test.com", nil)
		req.SetPathValue("email", "test@test.com")

		handler.HandleGetByCustomerEmail(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
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

	t.Run("bad json", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/orders", bytes.NewReader([]byte("invalid json")))

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		reqBody := inbound.CreateOrderRequest{CustomerEmail: "new@test.com"}
		bodyBytes, _ := json.Marshal(reqBody)

		repo.EXPECT().Create(reqBody.ToDomain()).Return(errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		reqBody := inbound.UpdateOrderRequest{
			CustomerEmail: "updated@test.com",
			Status:        "shipped",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		ord := reqBody.ToDomain()
		ord.ID = 1
		repo.EXPECT().Update(ord).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/orders/1", bytes.NewReader(bodyBytes))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("empty id", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/orders/", bytes.NewReader([]byte("{}")))
		req.SetPathValue("id", "")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/orders/invalid", bytes.NewReader([]byte("{}")))
		req.SetPathValue("id", "invalid")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("bad json", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/orders/1", bytes.NewReader([]byte("invalid json")))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		reqBody := inbound.UpdateOrderRequest{CustomerEmail: "updated@test.com", Status: "shipped"}
		bodyBytes, _ := json.Marshal(reqBody)

		ord := reqBody.ToDomain()
		ord.ID = 1
		repo.EXPECT().Update(ord).Return(errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/orders/1", bytes.NewReader(bodyBytes))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		repo.EXPECT().Delete(1).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/orders/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("empty id", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/orders/", nil)
		req.SetPathValue("id", "")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/orders/invalid", nil)
		req.SetPathValue("id", "invalid")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIOrdersRepository(t)
		handler := NewOrdersHandler(repo)

		repo.EXPECT().Delete(1).Return(errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/orders/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestMapToOrderResponseWithItems(t *testing.T) {
	ord := &models.Order{
		ID:            1,
		CustomerEmail: "test@test.com",
		Items: []models.ItemOrder{
			{ItemID: 10, OrderID: 1},
		},
	}

	repo := mocks.NewMockIOrdersRepository(t)
	handler := NewOrdersHandler(repo)
	repo.EXPECT().FindByID(1).Return(ord, nil)

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orders/1", nil)
	req.SetPathValue("id", "1")

	handler.HandleGetByID(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	expectedResp := outbound.GetOrderResponse{Order: mapToOrderResponse(ord)}
	expectedJSON, _ := json.Marshal(expectedResp)
	assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
}
