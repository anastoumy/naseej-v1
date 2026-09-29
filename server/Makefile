.PHONY: run generate seed test lint build clean

# Start the server
run:
	go run cmd/server/main.go

# Generate GraphQL code
generate:
	go run github.com/99designs/gqlgen generate

# Seed development data
seed:
	go run cmd/seed/main.go

# Run all tests
test:
	go test -race ./...

# Run linter
lint:
	golangci-lint run

# Build production binary
build:
	CGO_ENABLED=0 go build -o bin/server cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Start PostgreSQL via Docker
db-up:
	docker-compose up -d

# Stop PostgreSQL
db-down:
	docker-compose down

# Tidy dependencies
tidy:
	go mod tidy
