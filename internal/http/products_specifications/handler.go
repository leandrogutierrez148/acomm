package products_specifications

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

type ProductsSpecificationsHandler struct {
	repo interfaces.IProductsSpecificationsRepository
}

func NewProductsSpecificationsHandler(repo interfaces.IProductsSpecificationsRepository) *ProductsSpecificationsHandler {
	return &ProductsSpecificationsHandler{repo: repo}
}

func (h *ProductsSpecificationsHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	specs, err := h.repo.FindAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetProductSpecificationsResponse{Specifications: mapToProductSpecificationsResponse(specs)})
}

func (h *ProductsSpecificationsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
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
	api.OKResponse(w, outbound.GetProductSpecificationResponse{Specification: mapToProductSpecificationResponse(spec)})
}

func (h *ProductsSpecificationsHandler) HandleGetByProductID(w http.ResponseWriter, r *http.Request) {
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
	api.OKResponse(w, outbound.GetProductSpecificationsResponse{Specifications: mapToProductSpecificationsResponse(specs)})
}

func (h *ProductsSpecificationsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req inbound.CreateProductSpecificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	spec := req.ToDomain()
	if err := h.repo.Create(spec); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.CreateProductSpecificationResponse{Specification: mapToProductSpecificationResponse(spec)})
}

func (h *ProductsSpecificationsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req inbound.UpdateProductSpecificationRequest
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
	api.OKResponse(w, outbound.UpdateProductSpecificationResponse{Specification: mapToProductSpecificationResponse(spec)})
}

func (h *ProductsSpecificationsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
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

func mapToProductSpecificationResponse(s *models.ProductSpecification) outbound.ProductSpecification {
	return outbound.ProductSpecification{
		ID:        s.ID,
		ProductID: s.ProductID,
		Key:       s.Key,
		Value:     s.Value,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func mapToProductSpecificationsResponse(specs []models.ProductSpecification) []outbound.ProductSpecification {
	resp := make([]outbound.ProductSpecification, len(specs))
	for i := range specs {
		resp[i] = mapToProductSpecificationResponse(&specs[i])
	}
	return resp
}
