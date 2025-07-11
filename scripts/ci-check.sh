#!/bin/bash

# Pre-push CI/CD validation script
# This script runs the same checks that will run in GitHub Actions

set -e

echo "🚀 Running pre-push CI/CD validation..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2 passed${NC}"
    else
        echo -e "${RED}❌ $2 failed${NC}"
        exit 1
    fi
}

# Check if required tools are installed
echo "🔧 Checking required tools..."

# Check Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi

# Check golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${YELLOW}⚠️  golangci-lint not found. Installing...${NC}"
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

echo -e "${GREEN}✅ All required tools are available${NC}"

# 1. Format check
echo ""
echo "📝 Checking code formatting..."
if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    echo -e "${RED}❌ Code is not formatted correctly:${NC}"
    gofmt -s -l .
    echo -e "${YELLOW}Run 'make fmt' to fix formatting${NC}"
    exit 1
fi
print_status 0 "Code formatting"

# 2. Go mod tidy check
echo ""
echo "📦 Checking go mod tidy..."
go mod tidy
if [ -n "$(git status --porcelain go.mod go.sum)" ]; then
    echo -e "${RED}❌ go.mod/go.sum not tidy${NC}"
    echo -e "${YELLOW}Run 'go mod tidy' and commit changes${NC}"
    exit 1
fi
print_status 0 "Go modules"

# 3. Go vet
echo ""
echo "🔍 Running go vet..."
go vet ./...
print_status $? "Go vet"

# 4. Linting
echo ""
echo "🔍 Running golangci-lint..."
golangci-lint run --timeout=5m
print_status $? "Linting"

# 5. Build
echo ""
echo "🔨 Building project..."
go build -v ./...
print_status $? "Build"

# 6. Unit tests
echo ""
echo "🧪 Running unit tests..."
go test -v -short ./...
print_status $? "Unit tests"

# 7. Race detector tests
echo ""
echo "🏃 Running race detector tests..."
go test -race -short ./...
print_status $? "Race detector tests"

# 8. Security scan (if gosec is available)
if command -v gosec &> /dev/null; then
    echo ""
    echo "🔒 Running security scan..."
    gosec -quiet ./...
    print_status $? "Security scan"
else
    echo -e "${YELLOW}⚠️  gosec not found, skipping security scan${NC}"
fi

# 9. Check for vulnerabilities (if govulncheck is available)
if command -v govulncheck &> /dev/null; then
    echo ""
    echo "🛡️  Checking for vulnerabilities..."
    govulncheck ./...
    print_status $? "Vulnerability check"
else
    echo -e "${YELLOW}⚠️  govulncheck not found, skipping vulnerability check${NC}"
fi

echo ""
echo -e "${GREEN}🎉 All checks passed! Ready to push to GitHub.${NC}"
echo ""
echo "📋 Summary of checks performed:"
echo "  ✅ Code formatting (gofmt)"
echo "  ✅ Module tidiness (go mod tidy)"
echo "  ✅ Static analysis (go vet)"
echo "  ✅ Linting (golangci-lint)"
echo "  ✅ Build verification"
echo "  ✅ Unit tests"
echo "  ✅ Race condition detection"
if command -v gosec &> /dev/null; then
    echo "  ✅ Security scan (gosec)"
fi
if command -v govulncheck &> /dev/null; then
    echo "  ✅ Vulnerability check (govulncheck)"
fi
echo ""
echo "🚀 Push away! GitHub Actions will run the same checks."
