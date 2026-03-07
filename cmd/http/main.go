package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	http_server "github.com/lgutierrez148/acomm/app/http"
	"github.com/lgutierrez148/acomm/app/http/brands"
	"github.com/lgutierrez148/acomm/app/http/categories"
	"github.com/lgutierrez148/acomm/app/http/items"
	"github.com/lgutierrez148/acomm/app/http/orders"
	"github.com/lgutierrez148/acomm/app/http/products"
	"github.com/lgutierrez148/acomm/app/http/specifications"
	"github.com/lgutierrez148/acomm/database"
	"github.com/lgutierrez148/acomm/repositories"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		host,
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize databases and repos
	prodRepo := repositories.NewProductsRepository(db)
	catRepo := repositories.NewCategoriesRepository(db)
	brandsRepo := repositories.NewBrandsRepository(db)
	itemsRepo := repositories.NewItemsRepository(db)
	ordersRepo := repositories.NewOrdersRepository(db)
	specsRepo := repositories.NewSpecificationsRepository(db)

	// Initialize HTTP handlers
	prodsHandler := products.NewProductsHandler(prodRepo)
	catsHandler := categories.NewCategoriesHandler(catRepo)
	brandsHandler := brands.NewBrandsHandler(brandsRepo)
	itemsHandler := items.NewItemsHandler(itemsRepo)
	ordersHandler := orders.NewOrdersHandler(ordersRepo)
	specsHandler := specifications.NewSpecificationsHandler(specsRepo)

	// Initialize HTTP Server
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8484" // Default port if not in .env
	}

	httpSrv := http_server.NewHTTPServer(port, prodsHandler, catsHandler, brandsHandler, itemsHandler, ordersHandler, specsHandler)

	// Start the server blocking until context is done
	if err := httpSrv.Start(ctx); err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
