.PHONY: build clean install test lint fmt help

# Binary name
BINARY_NAME=git-guardian-mcp
VERSION=1.0.0

# Build directories
BUILD_DIR=./build
DIST_DIR=./dist

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v
	@echo "✓ Build complete: ./$(BINARY_NAME)"

build-all: ## Build for multiple platforms
	@echo "Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@echo "✓ Multi-platform builds complete in $(DIST_DIR)/"

clean: ## Remove build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@echo "✓ Clean complete"

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "✓ Tests complete"

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Running basic checks..."; \
		$(GOVET) ./...; \
	fi
	@echo "✓ Lint complete"

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✓ Format complete"

install: build ## Build and install hooks in current directory
	@echo "Installing hooks..."
	@chmod +x scripts/install-hooks.sh
	@./scripts/install-hooks.sh
	@echo "✓ Installation complete"

uninstall: ## Uninstall hooks from current directory
	@echo "Uninstalling hooks..."
	@chmod +x scripts/uninstall-hooks.sh
	@./scripts/uninstall-hooks.sh
	@echo "✓ Uninstall complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOGET) -v ./...
	@echo "✓ Dependencies downloaded"

run: build ## Build and run the MCP server
	@echo "Starting MCP server..."
	./$(BINARY_NAME)

version: ## Show version
	@echo "$(BINARY_NAME) version $(VERSION)"

# Development helpers
dev-setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "✓ Development environment ready"

check: fmt lint test ## Run all checks (format, lint, test)
	@echo "✓ All checks passed"

