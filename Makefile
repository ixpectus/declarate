# Declarate Makefile

.PHONY: help build test test-cover lint fmt clean run-postgres start-postgres stop-postgres deps deps-dev install

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the main binary
	@echo "Building declarate..."
	@mkdir -p build
	@go build -o build/main cmd/example/main.go

build-all: ## Build all binaries
	@echo "Building all binaries..."
	@mkdir -p build/bin
	@go build -o build/bin/declarate cmd/example/main.go
	@go build -o build/bin/converter cmd/converter/main.go
	@go build -o build/bin/formatter cmd/formatter/main.go
	@go build -o build/bin/server cmd/server/main.go

# Test targets
test: allure-results build ## Run tests
	@echo "Running tests..."
	@go test -v ./tests/... -count=1 -run TestSuite

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v ./... -short

test-cover: allure-results build ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -cover -coverpkg github.com/ixpectus/declarate/... -coverprofile cover.out -v ./tests/... -count=1 -run TestExample
	@go tool cover -html=cover.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	@go test -race -short ./...

# Code quality targets
lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w .
	@goimports -w .

fmt-check: ## Check if code is formatted
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -s -l .)" || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Security targets
security: ## Run security checks
	@echo "Running security checks..."
	@gosec ./...
	@govulncheck ./...

# Dependency management
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

deps-dev: ## Install development dependencies
	@echo "Installing development dependencies..."
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/kisielk/errcheck@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/gordonklaus/ineffassign@latest
	@go install github.com/client9/misspell/cmd/misspell@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest

tidy: ## Clean up dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

# Docker targets for PostgreSQL
run-postgres: ## Start PostgreSQL in Docker
	@echo "Starting PostgreSQL..."
	@docker run --name pg -d -e POSTGRES_HOST_AUTH_METHOD=trust -p 5440:5432 postgres:10.21

start-postgres: ## Start existing PostgreSQL container
	@echo "Starting PostgreSQL container..."
	@docker start pg

stop-postgres: ## Stop PostgreSQL container
	@echo "Stopping PostgreSQL container..."
	@docker stop pg

remove-postgres: ## Remove PostgreSQL container
	@echo "Removing PostgreSQL container..."
	@docker rm -f pg

# Utility targets
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf build/
	@rm -rf allure-results/
	@rm -f cover.out coverage.html
	@rm -f persistent persistent.idx
	@rm -rf tests/persistent tests/persistent.idx

allure-results:
	@mkdir -p allure-results

install: build ## Install binary to $GOPATH/bin
	@echo "Installing declarate..."
	@go install ./cmd/example

# Git hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	@cp scripts/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed successfully!"

uninstall-hooks: ## Uninstall git hooks
	@echo "Uninstalling git hooks..."
	@rm -f .git/hooks/pre-commit
	@echo "Git hooks uninstalled!"

# CI targets
ci-check: ## Run all CI checks locally (same as GitHub Actions)
	@echo "Running CI checks locally..."
	@./scripts/ci-check.sh

ci-test: deps lint test-race test-cover ## Run all CI tests
	@echo "All CI checks passed!"

ci-build: ## Build for multiple platforms (used in CI)
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -o dist/declarate-linux-amd64 ./cmd/example/
	@GOOS=linux GOARCH=arm64 go build -o dist/declarate-linux-arm64 ./cmd/example/
	@GOOS=darwin GOARCH=amd64 go build -o dist/declarate-darwin-amd64 ./cmd/example/
	@GOOS=darwin GOARCH=arm64 go build -o dist/declarate-darwin-arm64 ./cmd/example/
	go tool cover -html=cover.out -o cover.html


run-polling: 
	go run cmd/example/main.go -dir ./tests/yaml_poll/poll_long.yaml -progress_bar

