package products_specifications

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

func TestNewProductsSpecificationsHandler(t *testing.T) {
	repo := mocks.NewMockIProductsSpecificationsRepository(t)
	handler := NewProductsSpecificationsHandler(repo)
	assert.NotNil(t, handler)
}

func TestHandleGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIProductsSpecificationsRepository(t)
		handler := NewProductsSpecificationsHandler(repo)

		specs := []models.ProductSpecification{
			{ID: 1, ProductID: 1, Key: "Color", Value: "Red", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, ProductID: 2, Key: "Size", Value: "M", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		repo.EXPECT().FindAll().Return(specs, nil)

		req := httptest.NewRequest(http.MethodGet, "/products-specifications", nil)
		rr := httptest.NewRecorder()

		handler.HandleGetAll(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp outbound.GetProductSpecificationsResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Specifications, 2)
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewMockIProductsSpecificationsRepository(t)
		handler := NewProductsSpecificationsHandler(repo)

		repo.EXPECT().FindAll().Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/products-specifications", nil)
		rr := httptest.NewRecorder()

		handler.HandleGetAll(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestHandleGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIProductsSpecificationsRepository(t)
		handler := NewProductsSpecificationsHandler(repo)

		spec := &models.ProductSpecification{ID: 1, ProductID: 1, Key: "Color", Value: "Red"}
		repo.EXPECT().FindByID(uint(1)).Return(spec, nil)

		req := httptest.NewRequest(http.MethodGet, "/products-specifications/1", nil)
		req.SetPathValue("id", "1")
		rr := httptest.NewRecorder()

		handler.HandleGetByID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp outbound.GetProductSpecificationResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, spec.Key, resp.Specification.Key)
	})

	t.Run("invalid_id", func(t *testing.T) {
		repo := mocks.NewMockIProductsSpecificationsRepository(t)
		handler := NewProductsSpecificationsHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/products-specifications/abc", nil)
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()

		handler.HandleGetByID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestHandleGetByProductID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIProductsSpecificationsRepository(t)
		handler := NewProductsSpecificationsHandler(repo)

		specs := []models.ProductSpecification{
			{ID: 1, ProductID: 1, Key: "Color", Value: "Red"},
		}
		repo.EXPECT().FindByProductID(uint(1)).Return(specs, nil)

		req := httptest.NewRequest(http.MethodGet, "/products-specifications/product/1", nil)
		req.SetPathValue("product_id", "1")
		rr := httptest.NewRecorder()

		handler.HandleGetByProductID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp outbound.GetProductSpecificationsResponse
		err := json.NewDecoder(rr.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Specifications, 1)
	})
}

func TestHandleCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIProductsSpecificationsRepository(t)
		handler := NewProductsSpecificationsHandler(repo)

		reqBody := inbound.CreateProductSpecificationRequest{
			ProductID: 1,
			Key:       "Color",
			Value:     "Red",
		}
		body, _ := json.Marshal(reqBody)

		repo.EXPECT().Create(mock.AnythingOfType("*models.ProductSpecification")).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/products-specifications", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.HandleCreate(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
