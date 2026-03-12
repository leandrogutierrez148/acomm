package products_specifications

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ProductsSpecificationsMCPHandler struct {
	repo interfaces.IProductsSpecificationsRepository
}

func NewProductsSpecificationsMCPHandler(repo interfaces.IProductsSpecificationsRepository) *ProductsSpecificationsMCPHandler {
	return &ProductsSpecificationsMCPHandler{repo: repo}
}

func (h *ProductsSpecificationsMCPHandler) RegisterTools(srv *server.MCPServer) {
	// get_product_specifications
	getSpecsTool := mcp.NewTool("get_product_specifications",
		mcp.WithDescription("Returns all product specifications or specifications filtered by product ID"),
		mcp.WithNumber("product_id", mcp.Description("Filter specifications by product ID (optional)")),
	)
	srv.AddTool(getSpecsTool, h.handleGetProductSpecifications)

	// create_product_specification
	createSpecTool := mcp.NewTool("create_product_specification",
		mcp.WithDescription("Creates a new specification for a product"),
		mcp.WithNumber("product_id", mcp.Required(), mcp.Description("Product ID")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Specification key")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Specification value")),
	)
	srv.AddTool(createSpecTool, h.handleCreateProductSpecification)
}

func (h *ProductsSpecificationsMCPHandler) handleGetProductSpecifications(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	productID := request.GetInt("product_id", 0)

	var specs []models.ProductSpecification
	var err error

	if productID > 0 {
		specs, err = h.repo.FindByProductID(uint(productID))
	} else {
		specs, err = h.repo.FindAll()
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get specifications: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(specs, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *ProductsSpecificationsMCPHandler) handleCreateProductSpecification(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	productID := request.GetInt("product_id", 0)
	key := request.GetString("key", "")
	value := request.GetString("value", "")

	if productID == 0 || key == "" || value == "" {
		return mcp.NewToolResultError("product_id, key, and value are required"), nil
	}

	spec := &models.ProductSpecification{
		ProductID: uint(productID),
		Key:       key,
		Value:     value,
	}

	if err := h.repo.Create(spec); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create specification: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}
