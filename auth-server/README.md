# Acomm Auth Server

This is the identity and security microservice for the Acomm platform, built with Golang. It manages user registration, authentication, and distributes JWT tokens that are used by clients to securely access the commercial agent backend.

## Architecture

The Auth Server connects to its own dedicated PostgreSQL database (`acomm_auth`) to store user credentials securely. It is designed to be fully decoupled from the core agentic reasoning engine.

### Endpoints

- `POST /api/v1/auth/register`: Register a new user account.
- `POST /api/v1/auth/login`: Authenticate and receive a JWT token.
- `GET /health`: Service health check.

## Development

Make sure you have Go installed.

### Environment Variables

The auth server requires the following environment variables:

- `POSTGRES_HOST`: Database host (default: `localhost`)
- `POSTGRES_PORT`: Database port (default: `5434`)
- `POSTGRES_USER`: Database user (default: `postgres`)
- `POSTGRES_PASSWORD`: Database password
- `POSTGRES_DB`: Database name (default: `acomm_auth`)
- `HTTP_PORT`: Port to run the server on (default: `8081`)

### Running Locally

To run the server locally (outside of Docker):

```bash
# Export necessary environment variables or rely on the root .env file
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5434
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=password
export POSTGRES_DB=acomm_auth
export HTTP_PORT=8081

# Run the server
go run cmd/http/main.go
```
