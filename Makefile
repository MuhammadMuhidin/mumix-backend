APP_NAME := app

.PHONY: dev build run clean

# Development (lokal / Codespaces)
dev:
	@echo "Running in dev mode..."
	GIN_MODE=debug go run ./cmd/api

# Build production binary
build:
	@echo "Building production binary..."
	GIN_MODE=release go build -ldflags="-s -w" -o $(APP_NAME) ./cmd/api

# Run binary locally (prod-like)
run:
	@echo "Running production binary..."
	GIN_MODE=release ./$(APP_NAME)

# Cleanup
clean:
	rm -f $(APP_NAME)
	@echo "Cleaned up."