package items

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lgutierrez148/acomm/app/http/api"
	inbound "github.com/lgutierrez148/acomm/app/inbound/items"
	outbound "github.com/lgutierrez148/acomm/app/outbound/items"
	"github.com/lgutierrez148/acomm/interfaces"
	"github.com/lgutierrez148/acomm/models"
)

type ItemsHandler struct {
	repo interfaces.IItemsRepository
}

func NewItemsHandler(r interfaces.IItemsRepository) *ItemsHandler {
	return &ItemsHandler{repo: r}
}

func (h *ItemsHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.FindAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetItemsResponse{Items: mapToItemsResponse(items)})
}

func (h *ItemsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetItemResponse{Item: mapToItemResponse(item)})
}

func (h *ItemsHandler) HandleGetByProductID(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid product_id")
		return
	}

	items, err := h.repo.FindByProductID(uint(productID))
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetItemsResponse{Items: mapToItemsResponse(items)})
}

func (h *ItemsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req inbound.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	item := req.ToDomain()
	if err := h.repo.Create(item); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.CreateItemResponse{Item: mapToItemResponse(item)})
}

func (h *ItemsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req inbound.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	item := req.ToDomain()
	item.ID = uint(id)
	if err := h.repo.Update(item); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.UpdateItemResponse{Item: mapToItemResponse(item)})
}

func (h *ItemsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
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

// ─── mappers ────────────────────────────────────────────────────────────────

func mapToItemResponse(it *models.Item) outbound.Item {
	return outbound.Item{
		ID:        it.ID,
		ProductID: it.ProductID,
		SKU:       it.SKU,
		Price:     it.Price,
		Stock:     it.Stock,
		CreatedAt: it.CreatedAt,
		UpdatedAt: it.UpdatedAt,
	}
}

func mapToItemsResponse(items []models.Item) []outbound.Item {
	resp := make([]outbound.Item, len(items))
	for i := range items {
		resp[i] = mapToItemResponse(&items[i])
	}
	return resp
}
