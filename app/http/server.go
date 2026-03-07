package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lgutierrez148/acomm/app/http/brands"
	"github.com/lgutierrez148/acomm/app/http/categories"
	"github.com/lgutierrez148/acomm/app/http/items"
	"github.com/lgutierrez148/acomm/app/http/orders"
	"github.com/lgutierrez148/acomm/app/http/products"
	"github.com/lgutierrez148/acomm/app/http/specifications"
)

// HTTPServer represents the HTTP server wrapper.
type HTTPServer struct {
	srv      *http.Server
	products *products.ProductsHandler
	cats     *categories.CategoriesHandler
	brands   *brands.BrandsHandler
	items    *items.ItemsHandler
	orders   *orders.OrdersHandler
	specs    *specifications.SpecificationsHandler
}

// NewHTTPServer initializes a new HTTP server with the required handlers.
func NewHTTPServer(
	port string,
	productsHandler *products.ProductsHandler,
	categoriesHandler *categories.CategoriesHandler,
	brandsHandler *brands.BrandsHandler,
	itemsHandler *items.ItemsHandler,
	ordersHandler *orders.OrdersHandler,
	specsHandler *specifications.SpecificationsHandler,
) *HTTPServer {
	mux := http.NewServeMux()

	// Catalog endpoints
	mux.HandleFunc("GET /products", productsHandler.HandleSearchPaginated)
	mux.HandleFunc("GET /products/{code}", productsHandler.HandleGetByCode)

	// Categories endpoints
	mux.HandleFunc("GET /categories", categoriesHandler.HandleGet)
	mux.HandleFunc("POST /categories", categoriesHandler.HandleCreate)

	// Brands endpoints
	mux.HandleFunc("GET /brands", brandsHandler.HandleGetAll)
	mux.HandleFunc("GET /brands/{id}", brandsHandler.HandleGetByID)
	mux.HandleFunc("POST /brands", brandsHandler.HandleCreate)
	mux.HandleFunc("PUT /brands/{id}", brandsHandler.HandleUpdate)
	mux.HandleFunc("DELETE /brands/{id}", brandsHandler.HandleDelete)

	// Items endpoints
	mux.HandleFunc("GET /items", itemsHandler.HandleGetAll)
	mux.HandleFunc("GET /items/{id}", itemsHandler.HandleGetByID)
	mux.HandleFunc("GET /items/product/{product_id}", itemsHandler.HandleGetByProductID)
	mux.HandleFunc("POST /items", itemsHandler.HandleCreate)
	mux.HandleFunc("PUT /items/{id}", itemsHandler.HandleUpdate)
	mux.HandleFunc("DELETE /items/{id}", itemsHandler.HandleDelete)

	// Orders endpoints
	mux.HandleFunc("GET /orders", ordersHandler.HandleGetAll)
	mux.HandleFunc("GET /orders/{id}", ordersHandler.HandleGetByID)
	mux.HandleFunc("GET /orders/customer/{email}", ordersHandler.HandleGetByCustomerEmail)
	mux.HandleFunc("POST /orders", ordersHandler.HandleCreate)
	mux.HandleFunc("PUT /orders/{id}", ordersHandler.HandleUpdate)
	mux.HandleFunc("DELETE /orders/{id}", ordersHandler.HandleDelete)

	// Specifications endpoints
	mux.HandleFunc("GET /specifications", specsHandler.HandleGetAll)
	mux.HandleFunc("GET /specifications/{id}", specsHandler.HandleGetByID)
	mux.HandleFunc("GET /specifications/product/{product_id}", specsHandler.HandleGetByProductID)
	mux.HandleFunc("POST /specifications", specsHandler.HandleCreate)
	mux.HandleFunc("PUT /specifications/{id}", specsHandler.HandleUpdate)
	mux.HandleFunc("DELETE /specifications/{id}", specsHandler.HandleDelete)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		Handler: mux,
	}

	return &HTTPServer{
		srv:      srv,
		products: productsHandler,
		cats:     categoriesHandler,
		brands:   brandsHandler,
		items:    itemsHandler,
		orders:   ordersHandler,
		specs:    specsHandler,
	}
}

// Start runs the HTTP server in a blocking manner and handles graceful shutdown when the context is canceled.
func (h *HTTPServer) Start(ctx context.Context) error {
	serverErrCh := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("Starting server on http://%s", h.srv.Addr)
		if err := h.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- fmt.Errorf("server failed: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Context canceled, shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := h.srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}
		log.Println("Server stopped gracefully")
		return nil
	case err := <-serverErrCh:
		return err
	}
}
