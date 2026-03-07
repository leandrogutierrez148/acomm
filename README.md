# Acomm (Agent Commerce) — shopping infrastructure for AI agents.

This repository contains a Go application for managing an e-commerce catalog and ordering system. It includes functionalities for CRUD operations, an HTTP REST API, and an MCP (Model Context Protocol) server tailored for AI agent integrations.

## Managed Entities

The application manages the following core domain entities:
- **Products**: Main product catalogue entries.
- **Categories**: Taxonomies for organizing products.
- **Brands**: Brand information for items.
- **Specifications**: Details and characteristics linked to items or products.
- **Items**: Individual variants or stock-keeping units (formerly mapped as variations).
- **Orders**: Customer orders tracking purchased items.

## Project Structure

The project is structured following clean architecture and bounded contexts:

1. **cmd/**: Application entry points.
   - `http/main.go`: The main application entry point, serves the REST API.
   - `mcp/main.go`: The entry point for the Model Context Protocol server.

2. **app/**: Contains the application/delivery logic.
   - `http/`: REST API handlers and routing details.
   - `mcp/`: MCP tool implementations exposing system capabilities to AI agents.
   - `inbound/`: Request DTOs and mappers to domain models.
   - `outbound/`: Response DTOs and serializations to the external world.

3. **models/**: Contains the core domain models and structs.
4. **interfaces/**: Defines interfaces for repositories and decoupled services.
5. **repositories/**: Concrete repository implementations for database operations.
6. **database/**: Database connection and configuration utilities.
7. **sql/**: Database migration and setup scripts.
8. `.env`: Environment variables file for configuration.

## Setup & Tools

### Development Dependencies
Ensure you have Go and Docker installed before proceeding.
- [`mockery`](https://vektra.github.io/mockery/latest/) – used to generate interface mocks for testing.
- [`husky`](https://github.com/automation-co/husky) – used to manage Git hooks.

```bash
# Generate mocks for interfaces
mockery

# Install git hooks
husky install
```

### Useful Commands

You can use the provided Makefile to manage the environment:

- `make tidy`: Will install all dependencies.
- `make docker-up`: Will start the required database and infrastructure services via Docker Compose.
- `make test`: Will run the unit and integration test suites.
- `make run`: Will start the HTTP application.
- `make docker-down`: Will stop the Docker containers.
