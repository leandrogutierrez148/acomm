FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy mod files
COPY go.mod ./
RUN go mod download

# Copy source
COPY . .

# Build HTTP server
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-http ./cmd/http/main.go

FROM alpine:3.19

WORKDIR /

COPY --from=builder /auth-http /auth-http

EXPOSE 8081

CMD ["/auth-http"]
