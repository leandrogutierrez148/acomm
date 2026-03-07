package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lgutierrez148/acomm/internal/database"
	"github.com/lgutierrez148/acomm/internal/mcp"
	"github.com/lgutierrez148/acomm/internal/mcp/brands"
	"github.com/lgutierrez148/acomm/internal/mcp/categories"
	"github.com/lgutierrez148/acomm/internal/mcp/items"
	"github.com/lgutierrez148/acomm/internal/mcp/orders"
	"github.com/lgutierrez148/acomm/internal/mcp/products"
	"github.com/lgutierrez148/acomm/internal/mcp/specifications"
	"github.com/lgutierrez148/acomm/internal/repositories"
)

func main() {
	// Try to load .env, but don't fail if it's not found as MCP servers might be run by hosts missing the file
	// or we might pass environment variables via the client configuration side.
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}

	// Initialize database connection
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	// Initialize database connection
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		host,
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize repositories
	prodRepo := repositories.NewProductsRepository(db)
	catRepo := repositories.NewCategoriesRepository(db)
	brandsRepo := repositories.NewBrandsRepository(db)
	itemsRepo := repositories.NewItemsRepository(db)
	ordersRepo := repositories.NewOrdersRepository(db)
	specsRepo := repositories.NewSpecificationsRepository(db)

	// Initialize handlers
	productsHandler := products.NewProductsMCPHandler(prodRepo)
	categoriesHandler := categories.NewCategoriesMCPHandler(catRepo)
	brandsHandler := brands.NewBrandsMCPHandler(brandsRepo)
	itemsHandler := items.NewItemsMCPHandler(itemsRepo)
	ordersHandler := orders.NewOrdersMCPHandler(ordersRepo)
	specsHandler := specifications.NewSpecificationsMCPHandler(specsRepo)

	mcpSrv := mcp.NewMCPServer(
		productsHandler,
		categoriesHandler,
		brandsHandler,
		itemsHandler,
		ordersHandler,
		specsHandler,
	)

	// Since MCP relies on stdio for its protocol, all logging must go to standard error to avoid
	// corrupting the standard output which is strictly for MCP JSON-RPC messages.
	log.SetOutput(os.Stderr)

	mcpPort := os.Getenv("MCP_PORT")
	if mcpPort != "" {
		if err := mcpSrv.ServeSSE(mcpPort); err != nil {
			log.Fatalf("SSE Server error: %v", err)
		}
	} else {
		if err := mcpSrv.ServeStdio(); err != nil {
			log.Fatalf("Stdio Server error: %v", err)
		}
	}
}
