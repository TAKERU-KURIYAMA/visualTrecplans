.PHONY: help build run test clean migrate-up migrate-down seed lint fmt

# Default target
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo "  make seed           - Seed database with sample data"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/server ./cmd/main.go

# Run the application
run:
	@echo "Running application..."
	@go run ./cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf ./tmp
	@rm -rf ./logs

# Run database migrations
migrate-up:
	@echo "Running migrations..."
	@migrate -path ./migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

# Rollback database migrations
migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path ./migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down 1

# Seed database
seed:
	@echo "Seeding database..."
	@go run ./cmd/seeder

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Development helpers
dev:
	@echo "Starting development server with hot reload..."
	@air

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t visualtrecplans-backend .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env visualtrecplans-backend

# Database helpers
db-create:
	@echo "Creating database..."
	@createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

db-drop:
	@echo "Dropping database..."
	@dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

db-reset: db-drop db-create migrate-up seed
	@echo "Database reset complete!"