package categories

import (
	"encoding/json"
	"net/http"

	"github.com/lgutierrez148/acomm/app/http/api"
	inbound "github.com/lgutierrez148/acomm/app/inbound/categories"
	outbound "github.com/lgutierrez148/acomm/app/outbound/categories"
	"github.com/lgutierrez148/acomm/interfaces"
	"github.com/lgutierrez148/acomm/models"
)

type CategoriesHandler struct {
	repo interfaces.ICategoriesRepository
}

func NewCategoriesHandler(r interfaces.ICategoriesRepository) *CategoriesHandler {
	return &CategoriesHandler{
		repo: r,
	}
}

// HandleGet retrieves all categories from the repository and returns them as a JSON response.
func (h *CategoriesHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	prods, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	categories := mapToCategoryResponse(prods)

	api.OKResponse(w, outbound.GetCategoriesResponse{
		Categories: categories,
	})
}

// mapToCategoryResponse maps the models.Category slice to a slice of Category response objects.
func mapToCategoryResponse(categories []models.Category) []outbound.Category {
	resp := make([]outbound.Category, len(categories))
	for i, c := range categories {
		resp[i] = outbound.Category{
			Name: c.Name,
			Code: c.Code,
		}
	}
	return resp
}

// HandleCreate creates a new category based on the request body and returns a success response.
func (h *CategoriesHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "request body is required", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req inbound.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	category := req.ToDomain()
	if err := h.repo.CreateCategory(category); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKCreatedResponse(w, outbound.CreateCategoryResponse{Message: "Category created successfully"})
}
