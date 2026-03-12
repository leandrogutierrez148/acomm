package items_specifications

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

type ItemsSpecificationsHandler struct {
	repo interfaces.IItemsSpecificationsRepository
}

func NewItemsSpecificationsHandler(repo interfaces.IItemsSpecificationsRepository) *ItemsSpecificationsHandler {
	return &ItemsSpecificationsHandler{repo: repo}
}

func (h *ItemsSpecificationsHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	specs, err := h.repo.FindAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetItemSpecificationsResponse{Specifications: mapToItemSpecificationsResponse(specs)})
}

func (h *ItemsSpecificationsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	spec, err := h.repo.FindByID(uint(id))
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "specification not found")
		return
	}
	api.OKResponse(w, outbound.GetItemSpecificationResponse{Specification: mapToItemSpecificationResponse(spec)})
}

func (h *ItemsSpecificationsHandler) HandleGetByItemID(w http.ResponseWriter, r *http.Request) {
	itemIDStr := r.PathValue("item_id")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid item_id")
		return
	}

	specs, err := h.repo.FindByItemID(uint(itemID))
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetItemSpecificationsResponse{Specifications: mapToItemSpecificationsResponse(specs)})
}

func (h *ItemsSpecificationsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req inbound.CreateItemSpecificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	spec := req.ToDomain()
	if err := h.repo.Create(spec); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.CreateItemSpecificationResponse{Specification: mapToItemSpecificationResponse(spec)})
}

func (h *ItemsSpecificationsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req inbound.UpdateItemSpecificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	spec := req.ToDomain()
	spec.ID = uint(id)
	if err := h.repo.Update(spec); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.UpdateItemSpecificationResponse{Specification: mapToItemSpecificationResponse(spec)})
}

func (h *ItemsSpecificationsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
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

func mapToItemSpecificationResponse(s *models.ItemSpecification) outbound.ItemSpecification {
	return outbound.ItemSpecification{
		ID:        s.ID,
		ItemID:    s.ItemID,
		Key:       s.Key,
		Value:     s.Value,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func mapToItemSpecificationsResponse(specs []models.ItemSpecification) []outbound.ItemSpecification {
	resp := make([]outbound.ItemSpecification, len(specs))
	for i := range specs {
		resp[i] = mapToItemSpecificationResponse(&specs[i])
	}
	return resp
}
