package products

import (
	"net/http"
	"strconv"

	"github.com/lgutierrez148/acomm/internal/http/api"
	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/lgutierrez148/acomm/internal/outbound"
	"github.com/shopspring/decimal"
)

type ProductsHandler struct {
	repo interfaces.IProductsRepository
}

func NewProductsHandler(r interfaces.IProductsRepository) *ProductsHandler {
	return &ProductsHandler{
		repo: r,
	}
}

// HandleGet retrieves all products from the repository and returns them as a JSON response.
func (h *ProductsHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	prods, err := h.repo.GetAllProducts()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve products: "+err.Error())
		return
	}

	products := mapToProductsResponse(prods)

	api.OKResponse(w, outbound.GetProductsResponse{
		Products: products,
	})
}

// HandleGetPaginated retrieves products with pagination from the repository and returns them as a JSON response.
func (h *ProductsHandler) HandleGetPaginated(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	OffsetStr := query.Get("Offset")
	sizeStr := query.Get("limit")

	// Defaults
	offset := 0
	limit := 10

	if p, err := strconv.Atoi(OffsetStr); err == nil && p > 0 {
		offset = p
	}
	if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
		limit = s
	}

	prods, count, err := h.repo.GetProductsPaginated(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := mapToProductsResponse(prods)

	api.OKResponse(w, outbound.GetProductsPagedResponse{
		Products: products,
		Pagination: outbound.Pagination{
			Offset:     offset,
			Limit:      limit,
			TotalCount: count,
		},
	})
}

// mapToProductsResponse maps a slice of models.Product to a slice of Product response objects.
func mapToProductsResponse(products []models.Product) []outbound.Product {
	resp := make([]outbound.Product, len(products))
	for i, p := range products {
		price := 0.0
		if len(p.Items) > 0 {
			price = p.Items[0].Price
		}

		resp[i] = outbound.Product{
			Code:     strconv.FormatUint(uint64(p.ID), 10),
			Price:    price,
			Category: p.Category.Name,
		}
	}
	return resp
}

// mapToProductResponse maps a single models.Product to a Product response object.
func mapToProductResponse(prod *models.Product) outbound.Product {
	vars := make([]outbound.Variant, len(prod.Items))
	for i, v := range prod.Items {
		vars[i] = outbound.Variant{
			Name:  v.SKU,
			SKU:   v.SKU,
			Price: v.Price,
		}
	}

	price := 0.0
	if len(prod.Items) > 0 {
		price = prod.Items[0].Price
	}

	return outbound.Product{
		Code:     strconv.FormatUint(uint64(prod.ID), 10),
		Price:    price,
		Category: prod.Category.Name,
		Variants: vars,
	}
}

// HandleSearchPaginated retrieves products based on search criteria with pagination and returns them as a JSON response.
func (h *ProductsHandler) HandleSearchPaginated(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	categoryStr := query.Get("category")
	maxPriceStr := query.Get("maxPrice")

	maxPrice := decimal.Zero

	if mp, err := decimal.NewFromString(maxPriceStr); err == nil && mp.Compare(decimal.Zero) > 0 {
		maxPrice = mp
	}

	offsetStr := query.Get("offset")
	limitStr := query.Get("limit")

	// Defaults
	offset := 0
	limit := 10

	if o, err := strconv.Atoi(offsetStr); err == nil && o > 0 {
		offset = o
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	prods, count, err := h.repo.SearchProductsPaginated(offset, limit, categoryStr, maxPrice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := mapToProductsResponse(prods)

	api.OKResponse(w, outbound.GetProductsPagedResponse{
		Products: products,
		Pagination: outbound.Pagination{
			Offset:     offset,
			Limit:      limit,
			TotalCount: count,
		},
	})
}

// HandleGetByCode retrieves a product by its code and returns it as a JSON response.
func (h *ProductsHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	prod, err := h.repo.GetProductByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if prod == nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	api.OKResponse(w, outbound.GetProductResponse{Product: mapToProductResponse(prod)})
}
