package brands

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

type BrandsHandler struct {
	repo interfaces.IBrandsRepository
}

func NewBrandsHandler(r interfaces.IBrandsRepository) *BrandsHandler {
	return &BrandsHandler{repo: r}
}

func (h *BrandsHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	brands, err := h.repo.FindAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetBrandsResponse{Brands: mapToBrandsResponse(brands)})
}

func (h *BrandsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	brand, err := h.repo.FindByID(uint(id))
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.GetBrandResponse{Brand: mapToBrandResponse(brand)})
}

func (h *BrandsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req inbound.CreateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	brand := req.ToDomain()
	if err := h.repo.Create(brand); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.CreateBrandResponse{Brand: mapToBrandResponse(brand)})
}

func (h *BrandsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req inbound.UpdateBrandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	brand := req.ToDomain()
	brand.ID = uint(id)
	if err := h.repo.Update(brand); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, outbound.UpdateBrandResponse{Brand: mapToBrandResponse(brand)})
}

func (h *BrandsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
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

func mapToBrandResponse(b *models.Brand) outbound.Brand {
	return outbound.Brand{
		ID:                 b.ID,
		Name:               b.Name,
		Title:              b.Title,
		ImageURL:           b.ImageURL,
		MetaTagDescription: b.MetaTagDescription,
		IsActive:           b.IsActive,
	}
}

func mapToBrandsResponse(brands []models.Brand) []outbound.Brand {
	resp := make([]outbound.Brand, len(brands))
	for i := range brands {
		resp[i] = mapToBrandResponse(&brands[i])
	}
	return resp
}
