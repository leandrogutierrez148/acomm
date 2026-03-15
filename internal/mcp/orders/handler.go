package orders

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lgutierrez148/acomm/internal/inbound"
	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type OrdersMCPHandler struct {
	repo interfaces.IOrdersRepository
}

func NewOrdersMCPHandler(repo interfaces.IOrdersRepository) *OrdersMCPHandler {
	return &OrdersMCPHandler{repo: repo}
}

func (h *OrdersMCPHandler) RegisterTools(srv *server.MCPServer) {
	// get_orders
	getOrdersTool := mcp.NewTool("get_orders",
		mcp.WithDescription("Returns all orders"),
	)
	srv.AddTool(getOrdersTool, h.handleGetOrders)

	// get_order_by_id
	getOrderByIDTool := mcp.NewTool("get_order_by_id",
		mcp.WithDescription("Returns a specific order by its ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Order ID")),
	)
	srv.AddTool(getOrderByIDTool, h.handleGetOrderByID)

	// create_order
	createOrderTool := mcp.NewTool("create_order",
		mcp.WithDescription("Creates a new order."),
		mcp.WithString("customer_email", mcp.Required(), mcp.Description("Customer's email address")),
		mcp.WithString("customer_name", mcp.Required(), mcp.Description("Customer's full name")),
		mcp.WithString("customer_address", mcp.Required(), mcp.Description("Customer's shipping address")),
		mcp.WithString("customer_phone", mcp.Required(), mcp.Description("Customer's phone number")),
		mcp.WithArray("items", mcp.Required(), mcp.Description("List of items in the order"), mcp.Items(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"item_id":  map[string]any{"type": "number"},
				"quantity": map[string]any{"type": "number"},
				"price":    map[string]any{"type": "number"},
			},
			"required": []string{"item_id", "quantity", "price"},
		})),
	)
	srv.AddTool(createOrderTool, h.handleCreateOrder)
}

func (h *OrdersMCPHandler) handleGetOrders(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orders, err := h.repo.FindAll()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get orders: %v", err)), nil
	}
	responseJSON, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}
	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *OrdersMCPHandler) handleGetOrderByID(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetInt("id", 0)
	if id == 0 {
		return mcp.NewToolResultError("invalid or missing order id"), nil
	}

	order, err := h.repo.FindByID(id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get order: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *OrdersMCPHandler) handleCreateOrder(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	email, err := request.RequireString("customer_email")
	if err != nil {
		return mcp.NewToolResultError("customer_email is required"), nil
	}
	name, err := request.RequireString("customer_name")
	if err != nil {
		return mcp.NewToolResultError("customer_name is required"), nil
	}
	address, err := request.RequireString("customer_address")
	if err != nil {
		return mcp.NewToolResultError("customer_address is required"), nil
	}
	phone, err := request.RequireString("customer_phone")
	if err != nil {
		return mcp.NewToolResultError("customer_phone is required"), nil
	}

	args := request.GetArguments()
	rawItems, ok := args["items"].([]any)
	if !ok {
		return mcp.NewToolResultError("items array is required"), nil
	}

	var parsedItems []inbound.ItemOrderRequest
	for i, rawItem := range rawItems {
		itemMap, ok := rawItem.(map[string]any)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("item at index %d is not an object", i)), nil
		}

		var reqItem inbound.ItemOrderRequest
		if itemID, ok := itemMap["item_id"].(float64); ok {
			reqItem.ItemID = int(itemID)
		} else if itemIDStr, ok := itemMap["item_id"].(string); ok {
			fmt.Sscanf(itemIDStr, "%d", &reqItem.ItemID)
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("item_id is required for item at index %d", i)), nil
		}

		if quantity, ok := itemMap["quantity"].(float64); ok {
			reqItem.Quantity = int(quantity)
		} else if quantityStr, ok := itemMap["quantity"].(string); ok {
			fmt.Sscanf(quantityStr, "%d", &reqItem.Quantity)
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("quantity is required for item at index %d", i)), nil
		}

		if price, ok := itemMap["price"].(float64); ok {
			reqItem.Price = price
		} else if priceStr, ok := itemMap["price"].(string); ok {
			fmt.Sscanf(priceStr, "%f", &reqItem.Price)
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("price is required for item at index %d", i)), nil
		}

		parsedItems = append(parsedItems, reqItem)
	}

	req := inbound.CreateOrderRequest{
		CustomerEmail:   email,
		CustomerName:    name,
		CustomerAddress: address,
		CustomerPhone:   phone,
		Items:           parsedItems,
	}

	order := req.ToDomain()

	if err := h.repo.Create(order); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create order: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}
