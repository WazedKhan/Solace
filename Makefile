BINARY_NAME := solace
BUILD_DIR := bin
MAIN_PACKAGE_PATH := ./cmd/server

DB_URL=postgres://solace:strong-password@localhost:5432/solace?sslmode=disable

.PHONY: all build run test clean tidy help migrate-up migrate-down migrate-create

all: tidy test build

build:
	@echo "📦 Building the binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "🧹 Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@go clean

tidy:
	@go mod tidy

test:
	@go test ./...

migrate-up:
	@echo "⬆️ Running migrations up..."
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	@echo "⬇️ Running migrations down..."
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
