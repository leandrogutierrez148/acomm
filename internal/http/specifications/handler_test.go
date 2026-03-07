package specifications

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

func TestNewSpecificationsHandler(t *testing.T) {
	repo := mocks.NewMockISpecificationsRepository(t)
	handler := NewSpecificationsHandler(repo)
	assert.NotNil(t, handler)
	assert.Equal(t, repo, handler.repo)
}

func TestHandleGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockISpecificationsRepository(t)
		handler := NewSpecificationsHandler(repo)

		specs := []models.Specification{
			{ID: 1, Key: "Color", Value: "Red"},
			{ID: 2, Key: "Size", Value: "Large"},
		}

		repo.EXPECT().FindAll().Return(specs, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/specifications", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetSpecificationsResponse{Specifications: mapToSpecificationsResponse(specs)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewMockISpecificationsRepository(t)
		handler := NewSpecificationsHandler(repo)

		repo.EXPECT().FindAll().Return(nil, errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/specifications", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockISpecificationsRepository(t)
		handler := NewSpecificationsHandler(repo)

		spec := &models.Specification{ID: 1, Key: "Color", Value: "Red"}

		repo.EXPECT().FindByID(uint(1)).Return(spec, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/specifications/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetSpecificationResponse{Specification: mapToSpecificationResponse(spec)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockISpecificationsRepository(t)
		handler := NewSpecificationsHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/specifications/invalid", nil)
		req.SetPathValue("id", "invalid")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

func TestHandleGetByProductID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockISpecificationsRepository(t)
		handler := NewSpecificationsHandler(repo)

		specs := []models.Specification{
			{ID: 1, ProductID: 1, Key: "Color", Value: "Red"},
		}

		repo.EXPECT().FindByProductID(uint(1)).Return(specs, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/specifications/product/1", nil)
		req.SetPathValue("product_id", "1")

		handler.HandleGetByProductID(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetSpecificationsResponse{Specifications: mapToSpecificationsResponse(specs)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})
}

func TestHandleCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockISpecificationsRepository(t)
		handler := NewSpecificationsHandler(repo)

		reqBody := inbound.CreateSpecificationRequest{
			ProductID: 1,
			Key:       "Color",
			Value:     "Red",
		}

		bodyBytes, _ := json.Marshal(reqBody)

		specToCreate := reqBody.ToDomain()

		repo.EXPECT().Create(specToCreate).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/specifications", bytes.NewReader(bodyBytes))

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
