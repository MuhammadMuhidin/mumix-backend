.PHONY: build run test clean

APP_NAME=mumix-backend
GO=go

build:
	$(GO) build -o bin/$(APP_NAME) ./cmd/server

run:
	$(GO) run ./cmd/server

test:
	$(GO) test ./...

clean:
	rm -rf bin/

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy
