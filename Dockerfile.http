FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the HTTP server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/http-server ./cmd/http

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/http-server .
COPY .env .

EXPOSE 8484

CMD ["/app/http-server"]
