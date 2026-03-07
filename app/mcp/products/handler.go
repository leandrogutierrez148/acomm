package products

import (
	"context"
	"encoding/json"
	"fmt"

	outbound "github.com/lgutierrez148/acomm/app/outbound/products"
	"github.com/lgutierrez148/acomm/interfaces"
	"github.com/lgutierrez148/acomm/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/shopspring/decimal"
)

// ProductsMCPHandler handles MCP tools related to the product catalog.
type ProductsMCPHandler struct {
	prodRepo interfaces.IProductsRepository
}

func NewProductsMCPHandler(repo interfaces.IProductsRepository) *ProductsMCPHandler {
	return &ProductsMCPHandler{
		prodRepo: repo,
	}
}

// RegisterTools registers the catalog tools onto the provided MCP server.
func (h *ProductsMCPHandler) RegisterTools(srv *server.MCPServer) {
	// 1. search_products
	searchProductsTool := mcp.NewTool("search_products",
		mcp.WithDescription("Returns a list of products filtered by category and maximum price, with pagination."),
		mcp.WithString("category", mcp.Description("Product category to filter by (optional)")),
		mcp.WithString("maxPrice", mcp.Description("Maximum price of the product (optional)")),
		mcp.WithNumber("offset", mcp.Description("Starting index of the products (default 0)")),
		mcp.WithNumber("limit", mcp.Description("Number of products to return (default 10)")),
	)
	srv.AddTool(searchProductsTool, h.handleSearchProducts)

	// 2. get_product_by_code
	getProductByCodeTool := mcp.NewTool("get_product_by_code",
		mcp.WithDescription("Returns a specific product based on its code."),
		mcp.WithString("code", mcp.Required(), mcp.Description("Unique code of the product")),
	)
	srv.AddTool(getProductByCodeTool, h.handleGetProductByCode)
}

func (h *ProductsMCPHandler) handleSearchProducts(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	category := request.GetString("category", "")
	maxPriceStr := request.GetString("maxPrice", "")

	maxPrice := decimal.Zero
	if maxPriceStr != "" {
		if mp, err := decimal.NewFromString(maxPriceStr); err == nil && mp.Compare(decimal.Zero) > 0 {
			maxPrice = mp
		}
	}

	offset := request.GetInt("offset", 0)
	limit := request.GetInt("limit", 10)

	prods, count, err := h.prodRepo.SearchProductsPaginated(offset, limit, category, maxPrice)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search products: %v", err)), nil
	}

	response := outbound.GetProductsPagedResponse{
		Products: h.mapToProductsResponse(prods),
		Pagination: outbound.Pagination{
			Offset:     offset,
			Limit:      limit,
			TotalCount: count,
		},
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *ProductsMCPHandler) handleGetProductByCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	code, err := request.RequireString("code")
	if err != nil || code == "" {
		return mcp.NewToolResultError("product code is required"), nil
	}

	prod, err := h.prodRepo.GetProductByCode(code)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get product: %v", err)), nil
	}

	if prod == nil {
		return mcp.NewToolResultError("Product not found"), nil
	}

	resp := h.mapToProductResponse(prod)
	responseJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

// Helpers isolated to the handler

func (h *ProductsMCPHandler) mapToProductsResponse(products []models.Product) []outbound.Product {
	resp := make([]outbound.Product, len(products))
	for i, p := range products {
		price := 0.0
		if len(p.Items) > 0 {
			price = p.Items[0].Price
		}

		resp[i] = outbound.Product{
			Code:     fmt.Sprintf("%d", p.ID),
			Price:    price,
			Category: p.Category.Name,
		}
	}
	return resp
}

func (h *ProductsMCPHandler) mapToProductResponse(prod *models.Product) outbound.Product {
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

	resp := outbound.Product{
		Code:     fmt.Sprintf("%d", prod.ID),
		Price:    price,
		Category: prod.Category.Name,
		Variants: vars,
	}
	return resp
}
