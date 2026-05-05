FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the MCP Server
RUN go build -o mcp-server ./cmd/mcp

FROM alpine:latest
WORKDIR /app

# Install CA certificates for external requests if needed
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/mcp-server .

# Run the SSE server
EXPOSE 8080
CMD ["./mcp-server"]
