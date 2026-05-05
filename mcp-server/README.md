# Acomm MCP Server

This directory contains the Go application responsible for managing the e-commerce catalog and ordering system. It implements a Model Context Protocol (MCP) server to securely expose system capabilities to AI agents.

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

1. **`cmd/`**: Application entry points.
   - `http/main.go`: The main application entry point for the REST API (if applicable).
   - `mcp/main.go`: The entry point for the Model Context Protocol server.
2. **`internal/`**: Contains the core application logic.
   - `mcp/`: MCP tool implementations (e.g., `search_products`) exposing system capabilities to AI agents.
   - `http/`: REST API handlers and routing details.
   - `inbound/` & `outbound/`: DTOs and mappers for domain models.
   - `models/`: Core domain models and structs.
   - `interfaces/`: Defines interfaces for repositories and decoupled services.
   - `repositories/`: Concrete repository implementations for database operations.
   - `database/`: Database connection and configuration utilities.
3. **`sql/`**: Database migration and setup scripts.

## Development Setup

Ensure you have Go installed before proceeding.

### Dependencies
- [`mockery`](https://vektra.github.io/mockery/latest/) – used to generate interface mocks for testing.

```bash
# Generate mocks for interfaces
mockery
```

### Useful Commands

You can run standard Go commands locally to develop:

- `go mod tidy`: Install all dependencies.
- `go test ./...`: Run the unit and integration test suites.
- `go run cmd/mcp/main.go`: Start the MCP server.
