package items

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

func TestNewItemsHandler(t *testing.T) {
	repo := mocks.NewMockIItemsRepository(t)
	handler := NewItemsHandler(repo)
	assert.NotNil(t, handler)
	assert.Equal(t, repo, handler.repo)
}

func TestHandleGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		items := []models.Item{
			{ID: 1, SKU: "SKU1", Price: 100},
			{ID: 2, SKU: "SKU2", Price: 200},
		}

		repo.EXPECT().FindAll().Return(items, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetItemsResponse{Items: mapToItemsResponse(items)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		repo.EXPECT().FindAll().Return(nil, errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		item := &models.Item{ID: 1, SKU: "SKU1", Price: 100}

		repo.EXPECT().FindByID(uint(1)).Return(item, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetItemResponse{Item: mapToItemResponse(item)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/invalid", nil)
		req.SetPathValue("id", "invalid")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		repo.EXPECT().FindByID(uint(1)).Return(nil, errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleGetByProductID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		items := []models.Item{
			{ID: 1, ProductID: 5, SKU: "SKU1", Price: 100},
		}
		repo.EXPECT().FindByProductID(uint(5)).Return(items, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/product/5", nil)
		req.SetPathValue("product_id", "5")

		handler.HandleGetByProductID(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetItemsResponse{Items: mapToItemsResponse(items)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("invalid product_id", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/product/invalid", nil)
		req.SetPathValue("product_id", "invalid")

		handler.HandleGetByProductID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		repo.EXPECT().FindByProductID(uint(5)).Return(nil, errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/items/product/5", nil)
		req.SetPathValue("product_id", "5")

		handler.HandleGetByProductID(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		reqBody := inbound.CreateItemRequest{
			ProductID: 1,
			SKU:       "SKU1",
			Price:     100,
		}

		bodyBytes, _ := json.Marshal(reqBody)

		itemToCreate := reqBody.ToDomain()

		repo.EXPECT().Create(itemToCreate).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/items", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("bad json", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/items", bytes.NewReader([]byte("invalid json")))

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		reqBody := inbound.CreateItemRequest{SKU: "SKU1", Price: 100}
		bodyBytes, _ := json.Marshal(reqBody)

		repo.EXPECT().Create(reqBody.ToDomain()).Return(errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/items", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		reqBody := inbound.UpdateItemRequest{SKU: "SKU_UPDATED", Price: 150}
		bodyBytes, _ := json.Marshal(reqBody)

		item := reqBody.ToDomain()
		item.ID = 1
		repo.EXPECT().Update(item).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/items/1", bytes.NewReader(bodyBytes))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/items/invalid", bytes.NewReader([]byte("{}")))
		req.SetPathValue("id", "invalid")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("bad json", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/items/1", bytes.NewReader([]byte("invalid json")))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		reqBody := inbound.UpdateItemRequest{SKU: "SKU_UPDATED", Price: 150}
		bodyBytes, _ := json.Marshal(reqBody)

		item := reqBody.ToDomain()
		item.ID = 1
		repo.EXPECT().Update(item).Return(errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/items/1", bytes.NewReader(bodyBytes))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		repo.EXPECT().Delete(uint(1)).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/items/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/items/invalid", nil)
		req.SetPathValue("id", "invalid")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		repo.EXPECT().Delete(uint(1)).Return(errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/items/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}
