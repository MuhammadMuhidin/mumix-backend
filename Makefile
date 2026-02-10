APP_NAME := app

.PHONY: dev build run clean

# Jalankan mode development
dev:
	@echo "Running in dev mode..."
	GIN_MODE=debug go run .

# Build binary untuk production
build:
	@echo "Building production binary..."
	GIN_MODE=release go build -ldflags="-s -w" -o $(APP_NAME)

# Jalankan binary hasil build
run:
	@echo "Running production binary..."
	GIN_MODE=release ./$(APP_NAME)

# Hapus binary
clean:
	rm -f $(APP_NAME)
	@echo "Cleaned up."
	