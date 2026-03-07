package brands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type BrandsMCPHandler struct {
	repo interfaces.IBrandsRepository
}

func NewBrandsMCPHandler(repo interfaces.IBrandsRepository) *BrandsMCPHandler {
	return &BrandsMCPHandler{repo: repo}
}

func (h *BrandsMCPHandler) RegisterTools(srv *server.MCPServer) {
	// get_brands
	getBrandsTool := mcp.NewTool("get_brands",
		mcp.WithDescription("Returns all brands"),
	)
	srv.AddTool(getBrandsTool, h.handleGetBrands)

	// get_brand_by_id
	getBrandByIDTool := mcp.NewTool("get_brand_by_id",
		mcp.WithDescription("Returns a specific brand by its ID"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Brand ID")),
	)
	srv.AddTool(getBrandByIDTool, h.handleGetBrandByID)

	// create_brand
	createBrandTool := mcp.NewTool("create_brand",
		mcp.WithDescription("Creates a new brand"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Brand name")),
		mcp.WithString("description", mcp.Description("Brand description (optional)")),
	)
	srv.AddTool(createBrandTool, h.handleCreateBrand)
}

func (h *BrandsMCPHandler) handleGetBrands(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	brands, err := h.repo.FindAll()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get brands: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(brands, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *BrandsMCPHandler) handleGetBrandByID(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetInt("id", 0)
	if id == 0 {
		return mcp.NewToolResultError("brand id is required"), nil
	}

	brand, err := h.repo.FindByID(uint(id))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get brand: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(brand, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *BrandsMCPHandler) handleCreateBrand(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")
	description := request.GetString("description", "")

	if name == "" {
		return mcp.NewToolResultError("name is required"), nil
	}

	brand := &models.Brand{
		Name:               name,
		MetaTagDescription: description,
	}

	if err := h.repo.Create(brand); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create brand: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(brand, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}
