package item_images

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lgutierrez148/acomm/internal/http/api"
	"github.com/lgutierrez148/acomm/internal/inbound"
	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/lgutierrez148/acomm/internal/outbound"
)

type ItemImagesHandler struct {
	repo interfaces.IItemImagesRepository
}

func NewItemImagesHandler(repo interfaces.IItemImagesRepository) *ItemImagesHandler {
	return &ItemImagesHandler{repo: repo}
}

func (h *ItemImagesHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	images, err := h.repo.FindAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetItemImagesResponse{Images: mapToItemImagesResponse(images)})
}

func (h *ItemImagesHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	image, err := h.repo.FindByID(uint(id))
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "image not found")
		return
	}
	api.OKResponse(w, outbound.GetItemImageResponse{Image: mapToItemImageResponse(image)})
}

func (h *ItemImagesHandler) HandleGetByItemID(w http.ResponseWriter, r *http.Request) {
	itemIDStr := r.PathValue("item_id")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid item_id")
		return
	}

	images, err := h.repo.FindByItemID(uint(itemID))
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetItemImagesResponse{Images: mapToItemImagesResponse(images)})
}

func (h *ItemImagesHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	itemIDStr := r.PathValue("item_id")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid item_id in path")
		return
	}

	var req inbound.CreateItemImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	image := req.ToDomain()
	image.ItemID = uint(itemID)

	if err := h.repo.Create(image); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.CreateItemImageResponse{Image: mapToItemImageResponse(image)})
}

func (h *ItemImagesHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	// Fetch existing to retain ItemID
	existing, err := h.repo.FindByID(uint(id))
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "image not found")
		return
	}

	var req inbound.UpdateItemImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	image := req.ToDomain()
	image.ID = uint(id)
	image.ItemID = existing.ItemID

	if err := h.repo.Update(image); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.UpdateItemImageResponse{Image: mapToItemImageResponse(image)})
}

func (h *ItemImagesHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, nil)
}

func mapToItemImageResponse(i *models.ItemImage) outbound.ItemImage {
	return outbound.ItemImage{
		ID:        i.ID,
		ItemID:    i.ItemID,
		Url:       i.Url,
		Label:     i.Label,
		Text:      i.Text,
		IsMain:    i.IsMain,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

func mapToItemImagesResponse(images []models.ItemImage) []outbound.ItemImage {
	resp := make([]outbound.ItemImage, len(images))
	for i := range images {
		resp[i] = mapToItemImageResponse(&images[i])
	}
	return resp
}
