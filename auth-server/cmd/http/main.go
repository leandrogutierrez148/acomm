package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/leandrogutierrez148/acomm/auth-server/internal/database"
	"github.com/leandrogutierrez148/acomm/auth-server/internal/http/auth"
	"github.com/leandrogutierrez148/acomm/auth-server/internal/repositories"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Error loading .env file: %s", err)
	}

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		host,
		port,
	)
	defer close()

	usersRepo := repositories.NewUsersRepository(db)
	authHandler := auth.NewAuthHandler(usersRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Simple CORS Middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	serverPort := os.Getenv("HTTP_PORT")
	if serverPort == "" {
		serverPort = "8081"
	}

	log.Printf("Starting auth-server on port %s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, corsMiddleware(mux)); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
