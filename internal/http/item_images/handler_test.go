package item_images

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
}
