package specifications

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lgutierrez148/acomm/app/http/api"
	inbound "github.com/lgutierrez148/acomm/app/inbound/specifications"
	outbound "github.com/lgutierrez148/acomm/app/outbound/specifications"
	"github.com/lgutierrez148/acomm/interfaces"
	"github.com/lgutierrez148/acomm/models"
)

type SpecificationsHandler struct {
	repo interfaces.ISpecificationsRepository
}

func NewSpecificationsHandler(repo interfaces.ISpecificationsRepository) *SpecificationsHandler {
	return &SpecificationsHandler{repo: repo}
}

func (h *SpecificationsHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	specs, err := h.repo.FindAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetSpecificationsResponse{Specifications: mapToSpecificationsResponse(specs)})
}

func (h *SpecificationsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
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
	api.OKResponse(w, outbound.GetSpecificationResponse{Specification: mapToSpecificationResponse(spec)})
}

func (h *SpecificationsHandler) HandleGetByProductID(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid product_id")
		return
	}

	specs, err := h.repo.FindByProductID(uint(productID))
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetSpecificationsResponse{Specifications: mapToSpecificationsResponse(specs)})
}

func (h *SpecificationsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req inbound.CreateSpecificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	spec := req.ToDomain()
	if err := h.repo.Create(spec); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.CreateSpecificationResponse{Specification: mapToSpecificationResponse(spec)})
}

func (h *SpecificationsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req inbound.UpdateSpecificationRequest
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
	api.OKResponse(w, outbound.UpdateSpecificationResponse{Specification: mapToSpecificationResponse(spec)})
}

func (h *SpecificationsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
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

func mapToSpecificationResponse(s *models.Specification) outbound.Specification {
	return outbound.Specification{
		ID:        s.ID,
		ProductID: s.ProductID,
		Key:       s.Key,
		Value:     s.Value,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func mapToSpecificationsResponse(specs []models.Specification) []outbound.Specification {
	resp := make([]outbound.Specification, len(specs))
	for i := range specs {
		resp[i] = mapToSpecificationResponse(&specs[i])
	}
	return resp
}
