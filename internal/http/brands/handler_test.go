package brands

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

func TestNewBrandsHandler(t *testing.T) {
	repo := mocks.NewMockIBrandsRepository(t)
	handler := NewBrandsHandler(repo)
	assert.NotNil(t, handler)
	assert.Equal(t, repo, handler.repo)
}

func TestHandleGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIBrandsRepository(t)
		handler := NewBrandsHandler(repo)

		brands := []models.Brand{
			{ID: 1, Name: "Brand1"},
			{ID: 2, Name: "Brand2"},
		}

		repo.EXPECT().FindAll().Return(brands, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/brands", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetBrandsResponse{Brands: mapToBrandsResponse(brands)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewMockIBrandsRepository(t)
		handler := NewBrandsHandler(repo)

		repo.EXPECT().FindAll().Return(nil, errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/brands", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIBrandsRepository(t)
		handler := NewBrandsHandler(repo)

		brand := &models.Brand{ID: 1, Name: "Brand1"}

		repo.EXPECT().FindByID(uint(1)).Return(brand, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/brands/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetBrandResponse{Brand: mapToBrandResponse(brand)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIBrandsRepository(t)
		handler := NewBrandsHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/brands/invalid", nil)
		req.SetPathValue("id", "invalid")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

func TestHandleCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIBrandsRepository(t)
		handler := NewBrandsHandler(repo)

		reqBody := inbound.CreateBrandRequest{
			Name: "Brand1",
		}

		bodyBytes, _ := json.Marshal(reqBody)

		brandToCreate := reqBody.ToDomain()

		repo.EXPECT().Create(brandToCreate).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/brands", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
