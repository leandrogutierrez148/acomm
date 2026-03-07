package products

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lgutierrez148/acomm/internal/mocks"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/lgutierrez148/acomm/internal/outbound"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewProductsHandler(t *testing.T) {
	t.Run("succesful creation of catalog handler", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		if handler == nil {
			t.Fatal("Expected ProductsHandler to be created, got nil")
		}

		if handler.repo != repo {
			t.Errorf("Expected repository to be %v, got %v", repo, handler.repo)
		}
	})
}

func TestGetEndpoint(t *testing.T) {
	t.Run("succesful retrieving catalog", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		prods := []models.Product{
			{ID: 1, Items: []models.Item{{Price: 100}}},
			{ID: 2, Items: []models.Item{{Price: 20}}},
		}

		// Mock the GetAllProducts method to return a sample product list
		repo.EXPECT().GetAllProducts().Return(prods, nil)

		// Create a response recorder to capture the response
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/products", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		handler.HandleGet(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code 200 OK, got %d", recorder.Code)
		}

		products := mapToProductsResponse(prods)

		resp := outbound.GetProductsResponse{
			Products: products,
		}

		expected, _ := json.Marshal(resp)

		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})

	t.Run("error retrieving catalog", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		// Mock the GetAllProducts method to return an error
		repo.EXPECT().GetAllProducts().Return(nil, errors.New("database error"))

		// Create a response recorder to capture the response
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/products", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		handler.HandleGet(recorder, req)
		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 500 Internal Server Error, got %d", recorder.Code)
		}
	})
}

func TestGetPaginatedEndpoint(t *testing.T) {
	t.Run("succesful retrieving paginated catalog", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		prods := []models.Product{
			{ID: 1, Items: []models.Item{{Price: 100}}},
			{ID: 2, Items: []models.Item{{Price: 20}}},
		}

		// Mock the GetProductsPaginated method to return a sample product list
		repo.EXPECT().GetProductsPaginated(0, 10).Return(prods, int64(len(prods)), nil)

		// Create a response recorder to capture the response
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/products?offset=0&limit=10", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		handler.HandleGetPaginated(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code 200 OK, got %d", recorder.Code)
		}

		products := mapToProductsResponse(prods)

		resp := outbound.GetProductsPagedResponse{
			Products: products,
			Pagination: outbound.Pagination{
				Offset:     0,
				Limit:      10,
				TotalCount: int64(len(prods)),
			},
		}

		expected, _ := json.Marshal(resp)

		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})

	t.Run("error retrieving paginated catalog", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		// Mock the GetPaginatedProducts method to return an error
		repo.EXPECT().
			GetProductsPaginated(0, 10).
			Return(nil, 0, errors.New("database error"))

		// Create a response recorder to capture the response
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/products?offset=0&limit=10", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		handler.HandleGetPaginated(recorder, req)
		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 500 Internal Server Error, got %d", recorder.Code)
		}
	})
}

func TestSearchPaginatedEndpoint(t *testing.T) {
	t.Run("succesful searching paginated catalog", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		prods := []models.Product{
			{ID: 1, Items: []models.Item{{Price: 100}}},
			{ID: 2, Items: []models.Item{{Price: 20}}},
		}

		// Mock the SearchProductsPaginated method to return a sample product list
		repo.EXPECT().
			SearchProductsPaginated(0, 10, "Category1", decimal.NewFromInt(50)).
			Return(prods, int64(len(prods)), nil)

		// Create a response recorder to capture the response
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/products?offset=0&limit=10&category=Category1&maxPrice=50", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		handler.HandleSearchPaginated(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code 200 OK, got %d", recorder.Code)
		}

		products := mapToProductsResponse(prods)

		resp := outbound.GetProductsPagedResponse{
			Products: products,
			Pagination: outbound.Pagination{
				Offset:     0,
				Limit:      10,
				TotalCount: int64(len(prods)),
			},
		}

		expected, _ := json.Marshal(resp)

		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})

	t.Run("error searching paginated catalog", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		// Mock the SearchPaginatedProducts method to return an error
		repo.EXPECT().
			SearchProductsPaginated(0, 10, "Category1", decimal.NewFromInt(50)).
			Return(nil, 0, errors.New("database error"))

		// Create a response recorder
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/products?offset=0&limit=10&category=Category1&maxPrice=50", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		handler.HandleSearchPaginated(recorder, req)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Internal Server Error, got %d", recorder.Code)
	})
}

func TestGetByCodeEndpoint(t *testing.T) {
	t.Run("succesful retrieving product by code", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		prod := &models.Product{
			ID:       1,
			Category: models.Category{Name: "Category1"},
			Items: []models.Item{
				{SKU: "SKU001", Price: 90},
				{SKU: "SKU002", Price: 80},
			},
		}

		// Mock the GetProductByCode method to return a sample product
		repo.EXPECT().GetProductByCode("1").Return(prod, nil)

		// Create a response recorder to capture the response
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/products/", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.SetPathValue("code", "1")

		handler.HandleGetByCode(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code 200 OK, got %d", recorder.Code)
		}

		expected := outbound.GetProductResponse{Product: mapToProductResponse(prod)}
		expectedJSON, _ := json.Marshal(expected)

		assert.JSONEq(t, string(expectedJSON), recorder.Body.String(), "Response body does not match expected")
	})

	t.Run("error retrieving product by code", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		// Mock the GetProductByCode method to return an error
		repo.EXPECT().GetProductByCode("1").Return(nil, errors.New("product not found"))

		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/products/", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.SetPathValue("code", "1")

		handler.HandleGetByCode(recorder, req)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Not Found, got %d", recorder.Code)
	})

	t.Run("product not found", func(t *testing.T) {
		repo := mocks.NewMockIProductsRepository(t)
		handler := NewProductsHandler(repo)

		// Mock the GetProductByCode method to return nil
		repo.EXPECT().GetProductByCode("1").Return(nil, nil)

		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/products/", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.SetPathValue("code", "1")

		handler.HandleGetByCode(recorder, req)
		assert.Equal(t, http.StatusNotFound, recorder.Code, "Expected status code 404 Not Found, got %d", recorder.Code)
	})
}
