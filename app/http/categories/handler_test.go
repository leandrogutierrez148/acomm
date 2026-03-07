package categories

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	outbound "github.com/lgutierrez148/acomm/app/outbound/categories"
	"github.com/lgutierrez148/acomm/mocks"
	"github.com/lgutierrez148/acomm/models"
	"github.com/stretchr/testify/assert"
)

func TestNewCategoriesHandler(t *testing.T) {
	t.Run("succesful creation of categories handler", func(t *testing.T) {
		repo := mocks.NewMockICategoriesRepository(t)
		handler := NewCategoriesHandler(repo)

		if handler == nil {
			t.Fatal("Expected CategoriesHandler to be created, got nil")
		}

		assert.Equal(t, repo, handler.repo, "Expected repository to be %v, got %v", repo, handler.repo)
	})
}

func TestGetEndpoint(t *testing.T) {
	t.Run("succesful retrieving categories", func(t *testing.T) {
		repo := mocks.NewMockICategoriesRepository(t)
		handler := NewCategoriesHandler(repo)

		prods := []models.Category{
			{Code: "C001", Name: "Category1"},
			{Code: "C002", Name: "Category1"},
		}

		// Mock the GetAllCategories method to return a sample category list
		repo.EXPECT().GetAllCategories().Return(prods, nil)

		// Create a response recorder to capture the response
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/categories", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		handler.HandleGet(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code 200 OK, got %d", recorder.Code)
		}

		categories := mapToCategoryResponse(prods)

		resp := outbound.GetCategoriesResponse{
			Categories: categories,
		}

		expected, _ := json.Marshal(resp)

		assert.JSONEq(t, string(expected), recorder.Body.String(), "Response body does not match expected")
	})

	t.Run("error retrieving categories", func(t *testing.T) {
		repo := mocks.NewMockICategoriesRepository(t)
		handler := NewCategoriesHandler(repo)

		// Mock the GetAllCategories method to return an error
		repo.EXPECT().GetAllCategories().Return(nil, errors.New("database error"))

		// Create a response recorder to capture the response
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/categories", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		handler.HandleGet(recorder, req)

		assert.Equal(t, recorder.Code, http.StatusInternalServerError, "Expected status code 500 Bad Request, got %d", recorder.Code)
	})
}

func TestCreateEndpoint(t *testing.T) {
	t.Run("succesful creating category", func(t *testing.T) {
		repo := mocks.NewMockICategoriesRepository(t)
		handler := NewCategoriesHandler(repo)

		category := models.Category{Code: "C001", Name: "Category1"}

		// Mock the CreateCategory method to return nil error
		repo.EXPECT().CreateCategory(&category).Return(nil)

		// Create a request with the category data
		body, _ := json.Marshal(category)
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		handler.HandleCreate(recorder, req)

		assert.Equal(t, recorder.Code, http.StatusCreated, "Expected status code 201 Bad Request, got %d", recorder.Code)
	})

	t.Run("error creating category", func(t *testing.T) {
		repo := mocks.NewMockICategoriesRepository(t)
		handler := NewCategoriesHandler(repo)

		category := models.Category{Code: "C001", Name: "Category1"}

		// Mock the CreateCategory method to return an error
		repo.EXPECT().CreateCategory(&category).Return(errors.New("database error"))

		// Create a request with the category data
		body, _ := json.Marshal(category)
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		handler.HandleCreate(recorder, req)

		assert.Equal(t, recorder.Code, http.StatusInternalServerError, "Expected status code 500 Bad Request, got %d", recorder.Code)
	})

	t.Run("bad request creating category", func(t *testing.T) {
		repo := mocks.NewMockICategoriesRepository(t)
		handler := NewCategoriesHandler(repo)

		// Create a request with invalid JSON data
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer([]byte("invalid json")))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		handler.HandleCreate(recorder, req)

		assert.Equal(t, recorder.Code, http.StatusBadRequest, "Expected status code 400 Bad Request, got %d", recorder.Code)
	})

	t.Run("bad request creating category with empty body", func(t *testing.T) {
		repo := mocks.NewMockICategoriesRepository(t)
		handler := NewCategoriesHandler(repo)

		// Create a request with an empty body
		recorder := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/categories", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		handler.HandleCreate(recorder, req)

		assert.Equal(t, recorder.Code, http.StatusBadRequest, "Expected status code 400 Bad Request, got %d", recorder.Code)
	})
}
