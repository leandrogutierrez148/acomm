package items

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	inbound "github.com/lgutierrez148/acomm/app/inbound/items"
	outbound "github.com/lgutierrez148/acomm/app/outbound/items"
	"github.com/lgutierrez148/acomm/mocks"
	"github.com/lgutierrez148/acomm/models"
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
}

func TestHandleCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsRepository(t)
		handler := NewItemsHandler(repo)

		reqBody := inbound.CreateItemRequest{
			ProductID: 1,
			SKU:       "SKU1",
			Price:     100,
			Stock:     10,
		}

		bodyBytes, _ := json.Marshal(reqBody)

		itemToCreate := reqBody.ToDomain()

		repo.EXPECT().Create(itemToCreate).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/items", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
