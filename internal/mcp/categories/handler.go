package categories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lgutierrez148/acomm/internal/interfaces"
	"github.com/lgutierrez148/acomm/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// CategoriesMCPHandler handles MCP tools related to categories.
type CategoriesMCPHandler struct {
	catRepo interfaces.ICategoriesRepository
}

func NewCategoriesMCPHandler(repo interfaces.ICategoriesRepository) *CategoriesMCPHandler {
	return &CategoriesMCPHandler{
		catRepo: repo,
	}
}

// RegisterTools registers the categories tools onto the provided MCP server.
func (h *CategoriesMCPHandler) RegisterTools(srv *server.MCPServer) {
	// 3. get_categories
	getCategoriesTool := mcp.NewTool("get_categories",
		mcp.WithDescription("Returns a list of all available categories."),
	)
	srv.AddTool(getCategoriesTool, h.handleGetCategories)

	// 4. create_category
	createCategoryTool := mcp.NewTool("create_category",
		mcp.WithDescription("Creates a new category in the catalog."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the category")),
		mcp.WithString("code", mcp.Required(), mcp.Description("Unique code for the category")),
	)
	srv.AddTool(createCategoryTool, h.handleCreateCategory)
}

func (h *CategoriesMCPHandler) handleGetCategories(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cats, err := h.catRepo.GetAllCategories()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get categories: %v", err)), nil
	}

	type CategoryResponse struct {
		Categories []models.Category `json:"categories"`
	}

	resp := CategoryResponse{Categories: cats}
	responseJSON, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *CategoriesMCPHandler) handleCreateCategory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err1 := request.RequireString("name")
	code, err2 := request.RequireString("code")

	if err1 != nil || err2 != nil || name == "" || code == "" {
		return mcp.NewToolResultError("Category name and code are required"), nil
	}

	category := &models.Category{
		Name: name,
		Code: code,
	}

	if err := h.catRepo.CreateCategory(category); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create category: %v", err)), nil
	}

	responseMsg := fmt.Sprintf("Category created successfully: %s (%s)", name, code)
	return mcp.NewToolResultText(responseMsg), nil
}
