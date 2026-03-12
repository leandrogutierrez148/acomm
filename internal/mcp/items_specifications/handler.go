package items_specifications

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ItemsSpecificationsMCPHandler struct {
	repo interfaces.IItemsSpecificationsRepository
}

func NewItemsSpecificationsMCPHandler(repo interfaces.IItemsSpecificationsRepository) *ItemsSpecificationsMCPHandler {
	return &ItemsSpecificationsMCPHandler{repo: repo}
}

func (h *ItemsSpecificationsMCPHandler) RegisterTools(srv *server.MCPServer) {
	// get_item_specifications
	getSpecsTool := mcp.NewTool("get_item_specifications",
		mcp.WithDescription("Returns all item specifications or specifications filtered by item ID"),
		mcp.WithNumber("item_id", mcp.Description("Filter specifications by item ID (optional)")),
	)
	srv.AddTool(getSpecsTool, h.handleGetItemSpecifications)

	// create_item_specification
	createSpecTool := mcp.NewTool("create_item_specification",
		mcp.WithDescription("Creates a new specification for an item"),
		mcp.WithNumber("item_id", mcp.Required(), mcp.Description("Item ID")),
		mcp.WithString("key", mcp.Required(), mcp.Description("Specification key")),
		mcp.WithString("value", mcp.Required(), mcp.Description("Specification value")),
	)
	srv.AddTool(createSpecTool, h.handleCreateItemSpecification)
}

func (h *ItemsSpecificationsMCPHandler) handleGetItemSpecifications(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	itemID := request.GetInt("item_id", 0)

	var specs []models.ItemSpecification
	var err error

	if itemID > 0 {
		specs, err = h.repo.FindByItemID(uint(itemID))
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

func (h *ItemsSpecificationsMCPHandler) handleCreateItemSpecification(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	itemID := request.GetInt("item_id", 0)
	key := request.GetString("key", "")
	value := request.GetString("value", "")

	if itemID == 0 || key == "" || value == "" {
		return mcp.NewToolResultError("item_id, key, and value are required"), nil
	}

	spec := &models.ItemSpecification{
		ItemID: uint(itemID),
		Key:    key,
		Value:  value,
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
