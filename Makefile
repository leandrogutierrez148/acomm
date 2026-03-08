tidy ::
	@go mod tidy && go mod vendor

run-http::
	@go run cmd/http/main.go

run-mcp::
	@go run cmd/mcp/main.go

test::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

up::
	docker compose up -d

down::
	docker compose down

clear::
	docker compose down -v --remove-orphans
