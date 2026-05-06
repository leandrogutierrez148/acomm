# Acomm (Agent Commerce)

Acomm is a modern shopping infrastructure designed specifically for AI agents. This repository contains the complete stack for the Acomm platform, separated into bounded contexts and decoupled services.

## Architecture

![Acomm Architecture](./public/architecture.png)

The project is structured as a monorepo with four main components:

1. **[`auth-server`](./auth-server/)**: A Go microservice managing user identity, security, and JWT generation, connected to its own dedicated PostgreSQL database.
2. **[`mcp-server`](./mcp-server/)**: A Model Context Protocol (MCP) server written in Go. It manages the core e-commerce catalog (products, categories, brands, orders) and exposes these capabilities as tools that an AI agent can natively consume.
3. **[`backend`](./backend/)**: The agentic "brain" built with Python, FastAPI, and LangGraph. It acts as the routing and reasoning layer, connecting to the MCP server to fulfill user requests and maintaining conversation state in PostgreSQL.
4. **[`frontend`](./frontend/)**: A Streamlit-based web interface for users to interact with the Acomm purchasing agent.

## Quickstart

The easiest way to run the entire stack locally is using Docker Compose. Make sure you have Docker installed and your `.env` file configured. We provide a `Makefile` to simplify common operations.

```bash
# Start all services (frontend, backend, mcp-server, and databases)
make up

# To stop all services gracefully
make down

# To stop all services and wipe the databases (remove volumes)
make clear
```

Once the containers are running, you can access:
- **Frontend UI**: http://localhost:8501
- **Backend API**: http://localhost:8000
- **Auth Server API**: http://localhost:8081
- **MCP Server**: http://localhost:8080

![Acomm Workflow](./public/acomm-agent.gif)

## Development

For detailed development instructions, please refer to the README files in each respective project directory:
- [Auth Server README](./auth-server/README.md)
- [MCP Server README](./mcp-server/README.md)
- [Backend README](./backend/README.md)
- [Frontend README](./frontend/README.md)
