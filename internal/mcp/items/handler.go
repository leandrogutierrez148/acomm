package items

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ItemsMCPHandler struct {
	repo interfaces.IItemsRepository
}

func NewItemsMCPHandler(repo interfaces.IItemsRepository) *ItemsMCPHandler {
	return &ItemsMCPHandler{repo: repo}
}

func (h *ItemsMCPHandler) RegisterTools(srv *server.MCPServer) {
	// get_items
	getItemsTool := mcp.NewTool("get_items",
		mcp.WithDescription("Returns all items or items filtered by product ID"),
		mcp.WithNumber("product_id", mcp.Description("Filter items by product ID (optional)")),
	)
	srv.AddTool(getItemsTool, h.handleGetItems)

	// get_item_by_id
	getItemByIDTool := mcp.NewTool("get_item_by_id",
		mcp.WithDescription("Returns a specific item by its ID"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Item ID")),
	)
	srv.AddTool(getItemByIDTool, h.handleGetItemByID)

	// create_item
	createItemTool := mcp.NewTool("create_item",
		mcp.WithDescription("Creates a new item"),
		mcp.WithNumber("product_id", mcp.Required(), mcp.Description("Product ID")),
		mcp.WithString("sku", mcp.Required(), mcp.Description("SKU")),
		mcp.WithNumber("price", mcp.Required(), mcp.Description("Price")),
		mcp.WithNumber("stock", mcp.Description("Stock quantity (default 0)")),
	)
	srv.AddTool(createItemTool, h.handleCreateItem)
}

func (h *ItemsMCPHandler) handleGetItems(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	productID := request.GetInt("product_id", 0)

	var items []models.Item
	var err error

	if productID > 0 {
		items, err = h.repo.FindByProductID(uint(productID))
	} else {
		items, err = h.repo.FindAll()
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get items: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *ItemsMCPHandler) handleGetItemByID(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetInt("id", 0)
	if id == 0 {
		return mcp.NewToolResultError("item id is required"), nil
	}

	item, err := h.repo.FindByID(uint(id))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get item: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *ItemsMCPHandler) handleCreateItem(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	productID := request.GetInt("product_id", 0)
	sku := request.GetString("sku", "")
	price := request.GetFloat("price", 0)
	stock := request.GetInt("stock", 0)

	if productID == 0 || sku == "" || price <= 0 {
		return mcp.NewToolResultError("product_id, sku, and price are required"), nil
	}

	item := &models.Item{
		ProductID: uint(productID),
		SKU:       sku,
		Price:     price,
		Stock:     stock,
	}

	if err := h.repo.Create(item); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create item: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}
