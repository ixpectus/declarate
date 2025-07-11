# GitHub Actions CI/CD Setup

This document describes the GitHub Actions workflows set up for the Declarate project.

## Workflows Overview

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Triggered on**: Push and PR to main/master/develop branches

**Jobs**:
- **Test**: Runs full test suite with PostgreSQL service
  - Sets up Go 1.20
  - Downloads and verifies dependencies
  - Builds the project
  - Runs tests with race detector and coverage
  - Uploads coverage to Codecov
  
- **Lint**: Code quality checks
  - Runs golangci-lint with project configuration
  
- **Security**: Security vulnerability scanning
  - Runs Gosec security scanner
  - Uploads results to GitHub Security tab

### 2. Compatibility Workflow (`.github/workflows/compatibility.yml`)

**Triggered on**: Push and PR to main/master branches

**Purpose**: Tests compatibility across multiple Go versions and operating systems

**Matrix**:
- **OS**: Ubuntu, Windows, macOS
- **Go versions**: 1.20, 1.21, 1.22

### 3. Code Quality Workflow (`.github/workflows/code-quality.yml`)

**Triggered on**: 
- Push/PR to main/master/develop branches
- Daily schedule (00:00 UTC)

**Checks**:
- Go formatting (`gofmt`)
- Import organization (`goimports`)
- Go vet analysis
- Static analysis (`staticcheck`)
- Error checking (`errcheck`)
- Ineffective assignments (`ineffassign`)
- Spelling (`misspell`)
- Module tidiness (`go mod tidy`)
- Vulnerability scanning (`govulncheck`)
- Dependency vulnerabilities (Nancy)

### 4. Release Workflow (`.github/workflows/release.yml`)

**Triggered on**: Git tags matching `v*` pattern

**Actions**:
- Runs full test suite
- Builds binaries for multiple platforms:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64) 
  - Windows (amd64)
- Generates changelog from git commits
- Creates GitHub release with binaries

### 5. Dependabot (`.github/dependabot.yml`)

**Automated dependency management**:
- **Go modules**: Daily updates
- **GitHub Actions**: Weekly updates
- Auto-assigns PRs to maintainers
- Uses conventional commit format

## Repository Settings

### Required Status Checks

The following checks must pass before merging PRs:
- `test` (CI workflow)
- `lint` (CI workflow) 
- `quality` (Code Quality workflow)

### Branch Protection

Recommended settings for main/master branch:
- Require PR reviews (minimum 1)
- Require status checks to pass
- Require branches to be up to date
- Dismiss stale reviews on new commits
- Require conversation resolution

## Local Development

### Prerequisites

Install development tools:
```bash
make deps-dev
```

### Available Make Targets

```bash
make help                 # Show all available targets
make build               # Build main binary
make test                # Run tests with PostgreSQL
make test-unit          # Run unit tests only
make test-cover         # Run tests with coverage
make lint               # Run linter
make fmt                # Format code
make clean              # Clean build artifacts
make ci-test            # Run all CI checks locally
```

### Running Tests Locally

1. Start PostgreSQL:
   ```bash
   make run-postgres
   ```

2. Run tests:
   ```bash
   make test
   ```

3. Run specific test:
   ```bash
   go test -v ./commands/echo -run TestEcho
   ```

### Code Quality Checks

Before submitting PRs, run:
```bash
make fmt              # Format code
make lint             # Check code quality
make test-race        # Test with race detector
make test-cover       # Check test coverage
```

## IDE Setup

### VS Code

The repository includes VS Code settings (`.vscode/`):
- Go extension configuration
- Recommended extensions
- Code formatting on save
- Lint on save

### Recommended Extensions

- Go (golang.go)
- YAML (redhat.vscode-yaml)
- Markdown lint (davidanson.vscode-markdownlint)
- Code spell checker (streetsidesoftware.code-spell-checker)
- GitHub integration (github.vscode-pull-request-github)

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed contribution guidelines.

## Security

- Secrets are scanned automatically
- Vulnerabilities are reported daily
- Security updates are automated via Dependabot
- SARIF results are uploaded to GitHub Security tab

## Maintenance

### Adding New Workflows

1. Create workflow file in `.github/workflows/`
2. Follow existing patterns for consistency
3. Test locally if possible
4. Document any new required secrets or settings

### Updating Dependencies

- Go modules: Automated by Dependabot
- GitHub Actions: Automated by Dependabot
- Development tools: Update in `make deps-dev` target

### Release Process

1. Ensure all tests pass on main branch
2. Create and push git tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. GitHub Actions will automatically create release

## Troubleshooting

### Common Issues

1. **Tests failing in CI but passing locally**
   - Check PostgreSQL service status
   - Verify environment variables
   - Check for race conditions

2. **Linter errors**
   - Run `make fmt` to fix formatting
   - Run `make lint` locally to see issues
   - Check `.golangci.yml` configuration

3. **Build failures on specific platforms**
   - Check platform-specific code
   - Verify build constraints
   - Test with GOOS/GOARCH environment variables

### Getting Help

- Check workflow logs in GitHub Actions tab
- Review existing issues and PRs
- Create new issue with detailed information
