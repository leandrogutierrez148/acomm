FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the MCP server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/mcp-server ./cmd/mcp

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/mcp-server .
COPY .env .

CMD ["/app/mcp-server"]
