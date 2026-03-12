package items_specifications

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lgutierrez148/acomm/internal/inbound"
	"github.com/lgutierrez148/acomm/internal/mocks"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/lgutierrez148/acomm/internal/outbound"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewItemsSpecificationsHandler(t *testing.T) {
	repo := mocks.NewMockIItemsSpecificationsRepository(t)
	handler := NewItemsSpecificationsHandler(repo)
	assert.NotNil(t, handler)
}

func TestHandleGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsSpecificationsRepository(t)
		handler := NewItemsSpecificationsHandler(repo)

		specs := []models.ItemSpecification{
			{ID: 1, ItemID: 1, Key: "Color", Value: "Red", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, ItemID: 2, Key: "Size", Value: "M", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		repo.EXPECT().FindAll().Return(specs, nil)

		req := httptest.NewRequest(http.MethodGet, "/items-specifications", nil)
		rr := httptest.NewRecorder()

		handler.HandleGetAll(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp outbound.GetItemSpecificationsResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Specifications, 2)
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewMockIItemsSpecificationsRepository(t)
		handler := NewItemsSpecificationsHandler(repo)

		repo.EXPECT().FindAll().Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/items-specifications", nil)
		rr := httptest.NewRecorder()

		handler.HandleGetAll(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestHandleGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsSpecificationsRepository(t)
		handler := NewItemsSpecificationsHandler(repo)

		spec := &models.ItemSpecification{ID: 1, ItemID: 1, Key: "Color", Value: "Red"}
		repo.EXPECT().FindByID(uint(1)).Return(spec, nil)

		req := httptest.NewRequest(http.MethodGet, "/items-specifications/1", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		handler.HandleGetByID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp outbound.GetItemSpecificationResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, spec.Key, resp.Specification.Key)
	})

	t.Run("invalid_id", func(t *testing.T) {
		repo := mocks.NewMockIItemsSpecificationsRepository(t)
		handler := NewItemsSpecificationsHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/items-specifications/abc", nil)
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()

		handler.HandleGetByID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestHandleGetByItemID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsSpecificationsRepository(t)
		handler := NewItemsSpecificationsHandler(repo)

		specs := []models.ItemSpecification{
			{ID: 1, ItemID: 1, Key: "Color", Value: "Red"},
		}
		repo.EXPECT().FindByItemID(uint(1)).Return(specs, nil)

		req := httptest.NewRequest(http.MethodGet, "/items-specifications/item/1", nil)
		req.SetPathValue("item_id", "1")
		rr := httptest.NewRecorder()

		handler.HandleGetByItemID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp outbound.GetItemSpecificationsResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Specifications, 1)
	})
}

func TestHandleCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemsSpecificationsRepository(t)
		handler := NewItemsSpecificationsHandler(repo)

		reqBody := inbound.CreateItemSpecificationRequest{
			ItemID: 1,
			Key:    "Color",
			Value:  "Red",
		}
		body, _ := json.Marshal(reqBody)

		repo.EXPECT().Create(mock.AnythingOfType("*models.ItemSpecification")).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/items-specifications", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.HandleCreate(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
