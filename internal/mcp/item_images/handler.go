package item_images

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/leandrogutierrez148/acomm/internal/interfaces"
	"github.com/leandrogutierrez148/acomm/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ItemImagesMCPHandler struct {
	repo interfaces.IItemImagesRepository
}

func NewItemImagesMCPHandler(repo interfaces.IItemImagesRepository) *ItemImagesMCPHandler {
	return &ItemImagesMCPHandler{repo: repo}
}

func (h *ItemImagesMCPHandler) RegisterTools(srv *server.MCPServer) {
	// get_item_images
	getImagesTool := mcp.NewTool("get_item_images",
		mcp.WithDescription("Returns all item images or images filtered by item ID."),
		mcp.WithNumber("item_id", mcp.Description("Filter images by item ID (optional)")),
	)
	srv.AddTool(getImagesTool, h.handleGetItemImages)

	// create_item_image
	createImageTool := mcp.NewTool("create_item_image",
		mcp.WithDescription("Creates a new image for an item."),
		mcp.WithNumber("item_id", mcp.Required(), mcp.Description("Item ID")),
		mcp.WithString("url", mcp.Required(), mcp.Description("Image URL")),
		mcp.WithString("label", mcp.Description("Image Label (optional)")),
		mcp.WithString("text", mcp.Description("Image Text (optional)")),
		mcp.WithBoolean("is_main", mcp.Description("Is Main Image (optional)")),
	)
	srv.AddTool(createImageTool, h.handleCreateItemImage)
}

func (h *ItemImagesMCPHandler) handleGetItemImages(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	itemID := request.GetInt("item_id", 0)

	var images []models.ItemImage
	var err error

	if itemID > 0 {
		images, err = h.repo.FindByItemID(uint(itemID))
	} else {
		images, err = h.repo.FindAll()
	}

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get images: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}

func (h *ItemImagesMCPHandler) handleCreateItemImage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	itemID := request.GetInt("item_id", 0)
	url := request.GetString("url", "")

	args := request.GetArguments()
	valLabel, okLabel := args["label"]
	valText, okText := args["text"]
	isMain := request.GetBool("is_main", false)

	if itemID == 0 || url == "" {
		return mcp.NewToolResultError("item_id and url are required"), nil
	}

	image := &models.ItemImage{
		ItemID: uint(itemID),
		Url:    url,
		IsMain: isMain,
	}

	if okLabel {
		lbl := valLabel.(string)
		image.Label = &lbl
	}
	if okText {
		txt := valText.(string)
		image.Text = &txt
	}

	if err := h.repo.Create(image); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create image: %v", err)), nil
	}

	responseJSON, err := json.MarshalIndent(image, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to encode response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(responseJSON)), nil
}
