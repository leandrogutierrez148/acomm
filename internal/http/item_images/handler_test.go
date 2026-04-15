package item_images

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

func TestNewItemImagesHandler(t *testing.T) {
	repo := mocks.NewMockIItemImagesRepository(t)
	handler := NewItemImagesHandler(repo)
	assert.NotNil(t, handler)
	assert.Equal(t, repo, handler.repo)
}

func TestHandleGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		lbl := "label"
		images := []models.ItemImage{
			{ID: 1, ItemID: 1, Url: "http://example.com/1.jpg", Label: &lbl, IsMain: true},
		}

		repo.EXPECT().FindAll().Return(images, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/item-images", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetItemImagesResponse{Images: mapToItemImagesResponse(images)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("error", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		repo.EXPECT().FindAll().Return(nil, errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/item-images", nil)

		handler.HandleGetAll(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		lbl := "label"
		image := &models.ItemImage{ID: 1, ItemID: 1, Url: "http://example.com/1.jpg", Label: &lbl, IsMain: true}

		repo.EXPECT().FindByID(uint(1)).Return(image, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/item-images/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetItemImageResponse{Image: mapToItemImageResponse(image)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/item-images/invalid", nil)
		req.SetPathValue("id", "invalid")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		repo.EXPECT().FindByID(uint(1)).Return(nil, errors.New("not found"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/item-images/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleGetByID(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
}

func TestHandleGetByItemID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		lbl := "label"
		images := []models.ItemImage{
			{ID: 1, ItemID: 2, Url: "http://example.com/1.jpg", Label: &lbl},
		}

		repo.EXPECT().FindByItemID(uint(2)).Return(images, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/item-images/item/2", nil)
		req.SetPathValue("item_id", "2")

		handler.HandleGetByItemID(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedResp := outbound.GetItemImagesResponse{Images: mapToItemImagesResponse(images)}
		expectedJSON, _ := json.Marshal(expectedResp)
		assert.JSONEq(t, string(expectedJSON), recorder.Body.String())
	})

	t.Run("invalid item_id", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/item-images/item/invalid", nil)
		req.SetPathValue("item_id", "invalid")

		handler.HandleGetByItemID(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		repo.EXPECT().FindByItemID(uint(2)).Return(nil, errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/item-images/item/2", nil)
		req.SetPathValue("item_id", "2")

		handler.HandleGetByItemID(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		lbl := "label"
		reqBody := inbound.CreateItemImageRequest{
			Url:    "http://example.com/img.jpg",
			Label:  &lbl,
			IsMain: true,
		}

		bodyBytes, _ := json.Marshal(reqBody)

		imgToCreate := reqBody.ToDomain()
		imgToCreate.ItemID = 1 // Path variable

		repo.EXPECT().Create(imgToCreate).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/item-images/item/1", bytes.NewReader(bodyBytes))
		req.SetPathValue("item_id", "1")

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("invalid item_id", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/item-images/item/invalid", bytes.NewReader([]byte("{}")))
		req.SetPathValue("item_id", "invalid")

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("bad json", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/item-images/item/1", bytes.NewReader([]byte("invalid json")))
		req.SetPathValue("item_id", "1")

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		reqBody := inbound.CreateItemImageRequest{Url: "http://example.com/img.jpg"}
		bodyBytes, _ := json.Marshal(reqBody)

		imgToCreate := reqBody.ToDomain()
		imgToCreate.ItemID = 1
		repo.EXPECT().Create(imgToCreate).Return(errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/item-images/item/1", bytes.NewReader(bodyBytes))
		req.SetPathValue("item_id", "1")

		handler.HandleCreate(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		existing := &models.ItemImage{ID: 1, ItemID: 2, Url: "http://old.com/img.jpg"}
		repo.EXPECT().FindByID(uint(1)).Return(existing, nil)

		reqBody := inbound.UpdateItemImageRequest{Url: "http://new.com/img.jpg", IsMain: false}
		bodyBytes, _ := json.Marshal(reqBody)

		updated := reqBody.ToDomain()
		updated.ID = 1
		updated.ItemID = existing.ItemID
		repo.EXPECT().Update(updated).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/item-images/1", bytes.NewReader(bodyBytes))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/item-images/invalid", bytes.NewReader([]byte("{}")))
		req.SetPathValue("id", "invalid")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("find error", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		repo.EXPECT().FindByID(uint(1)).Return(nil, errors.New("not found"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/item-images/1", bytes.NewReader([]byte("{}")))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("bad json", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		existing := &models.ItemImage{ID: 1, ItemID: 2}
		repo.EXPECT().FindByID(uint(1)).Return(existing, nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/item-images/1", bytes.NewReader([]byte("invalid json")))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error on update", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		existing := &models.ItemImage{ID: 1, ItemID: 2}
		repo.EXPECT().FindByID(uint(1)).Return(existing, nil)

		reqBody := inbound.UpdateItemImageRequest{Url: "http://new.com/img.jpg"}
		bodyBytes, _ := json.Marshal(reqBody)

		updated := reqBody.ToDomain()
		updated.ID = 1
		updated.ItemID = existing.ItemID
		repo.EXPECT().Update(updated).Return(errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/item-images/1", bytes.NewReader(bodyBytes))
		req.SetPathValue("id", "1")

		handler.HandleUpdate(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandleDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		repo.EXPECT().Delete(uint(1)).Return(nil)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/item-images/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/item-images/invalid", nil)
		req.SetPathValue("id", "invalid")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("db error", func(t *testing.T) {
		repo := mocks.NewMockIItemImagesRepository(t)
		handler := NewItemImagesHandler(repo)

		repo.EXPECT().Delete(uint(1)).Return(errors.New("db error"))

		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/item-images/1", nil)
		req.SetPathValue("id", "1")

		handler.HandleDelete(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}
