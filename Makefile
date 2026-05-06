.DEFAULT_GOAL := help
.PHONY: help up down clear infra-up infra-down tidy-mcp test-mcp run-mcp tidy-auth test-auth run-auth install-backend run-backend format-backend install-frontend run-frontend format-frontend

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- Docker Compose ---
up: ## Start all infrastructure and applications with Docker Compose
	docker compose up --build -d

down: ## Stop all Docker Compose containers
	docker compose down

clear: ## Stop Docker Compose containers and remove orphaned volumes
	docker compose down -v --remove-orphans
