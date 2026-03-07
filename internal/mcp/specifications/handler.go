package specifications

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type SpecificationsMCPHandler struct {
	repo interfaces.ISpecificationsRepository
}

func NewSpecificationsMCPHandler(repo interfaces.ISpecificationsRepository) *SpecificationsMCPHandler {
	return &SpecificationsMCPHandler{repo: repo}
}

func (h *SpecificationsMCPHandler) RegisterTools(srv *server.MCPServer) {
	// get_specifications
	getSpecificationsTool := mcp.NewTool("get_specifications",
		mcp.WithDescription("Returns all specifications or specifications filtered by product ID"),
		mcp.WithNumber("product_id", mcp.Description("Filter specifications by product ID (optional)")),
	)
	srv.AddTool(getSpecificationsTool, h.handleGetSpecifications)

	// create_specification
	createSpecificationTool := mcp.NewTool("create_specification",
		mcp.WithDescription("Creates a new specification for a product"),
		mcp.WithNumber("product_id", mcp.Required(), mcp.Description("Product ID")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Specification key")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Specification value")),
	)
	srv.AddTool(createSpecificationTool, h.handleCreateSpecification)
}

func (h *SpecificationsMCPHandler) handleGetSpecifications(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	productID := request.GetInt("product_id", 0)

	var specs []models.Specification
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

func (h *SpecificationsMCPHandler) handleCreateSpecification(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	productID := request.GetInt("product_id", 0)
	key := request.GetString("key", "")
	value := request.GetString("value", "")

	if productID == 0 || key == "" || value == "" {
		return mcp.NewToolResultError("product_id, key, and value are required"), nil
	}

	spec := &models.Specification{
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
