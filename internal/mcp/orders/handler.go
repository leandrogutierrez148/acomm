package orders

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lgutierrez148/acomm/internal/inbound"
	"github.com/lgutierrez148/acomm/internal/interfaces"
	internalmcp "github.com/lgutierrez148/acomm/internal/mcp"
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
	createOrderJSONSchema := internalmcp.MustGetToolSchema(internalmcp.CreateOrderSchemaName)

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
		mcp.WithDescription("Creates a new order. Input payload for order_json must follow this JSON Schema:\n"+createOrderJSONSchema),
		mcp.WithString("order_json", mcp.Required(), mcp.Description("JSON string payload for the order. Schema:\n"+createOrderJSONSchema)),
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
	orderJSON := request.GetString("order_json", "")
	if orderJSON == "" {
		return mcp.NewToolResultError("order_json is required"), nil
	}

	var req inbound.CreateOrderRequest

	if err := json.Unmarshal([]byte(orderJSON), &req); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to parse order_json: %v", err)), nil
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
