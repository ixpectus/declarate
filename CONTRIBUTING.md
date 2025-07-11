# Contributing to Declarate

Thank you for your interest in contributing to Declarate! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive feedback
- Respect differing viewpoints and experiences

## Getting Started

### Prerequisites

- Go 1.20 or later
- Docker (for running PostgreSQL tests)
- Git

### Setting up the development environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/declarate.git
   cd declarate
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/ixpectus/declarate.git
   ```

4. Install development dependencies:
   ```bash
   make deps-dev
   ```

5. Start PostgreSQL for tests:
   ```bash
   make run-postgres
   ```

6. Run tests to verify everything works:
   ```bash
   make test
   ```

## Development Workflow

### 1. Create a feature branch

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-new-command` for new features
- `fix/command-parsing-bug` for bug fixes
- `docs/update-readme` for documentation updates

### 2. Make your changes

- Write clean, readable code
- Follow the existing code style
- Add tests for new functionality
- Update documentation as needed

### 3. Test your changes

```bash
# Run unit tests
make test-unit

# Run all tests
make test

# Run tests with coverage
make test-cover

# Run linter
make lint

# Format code
make fmt
```

### 4. Commit your changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add new echo command for debugging

- Implement echo command with variable substitution
- Add comprehensive tests for echo functionality
- Update documentation with echo command examples"
```

Use conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for test additions/changes
- `refactor:` for code refactoring
- `perf:` for performance improvements

## Code Style

### Go Code Style

1. **Follow standard Go conventions**:
   - Use `gofmt` for formatting
   - Use `goimports` for import organization
   - Follow effective Go guidelines

2. **Naming conventions**:
   - Use descriptive names for functions, variables, and types
   - Use camelCase for unexported functions and variables
   - Use PascalCase for exported functions and types

3. **Comments**:
   - Add comments for exported functions and types
   - Use complete sentences in comments
   - Explain complex logic with inline comments

4. **Error handling**:
   - Always handle errors explicitly
   - Provide meaningful error messages
   - Use error wrapping when appropriate

### Example:

```go
// ProcessRequest handles HTTP request processing and validation.
// It returns the processed response body or an error if processing fails.
func ProcessRequest(req *http.Request) (*string, error) {
    if req == nil {
        return nil, fmt.Errorf("request cannot be nil")
    }
    
    // Process the request
    result, err := doProcessing(req)
    if err != nil {
        return nil, fmt.Errorf("failed to process request: %w", err)
    }
    
    return result, nil
}
```

### YAML Test Files

1. Use consistent indentation (2 spaces)
2. Use descriptive test names
3. Include comments for complex test scenarios

```yaml
- name: test user creation with validation
  method: POST
  path: /users
  request: |
    {
      "name": "{{$username}}",
      "email": "{{$email}}"
    }
  response: |
    {
      "id": "$matchRegexp(^[0-9]+$)",
      "name": "{{$username}}",
      "status": "active"
    }
  variables:
    user_id: "id"  # Extract user ID for subsequent tests
```

## Testing

### Unit Tests

- Write unit tests for all new functionality
- Aim for good test coverage (>80%)
- Use table-driven tests for multiple test cases
- Mock external dependencies

### Integration Tests

- Add integration tests for new commands
- Test with real services when possible
- Use Docker containers for dependencies

### Test Structure

```go
func TestCommandName_Method(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:    "invalid input",
            input:   invalidInput,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := MethodUnderTest(tt.input)
            
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Submitting Changes

### 1. Push your branch

```bash
git push origin feature/your-feature-name
```

### 2. Create a Pull Request

1. Go to your fork on GitHub
2. Click "New Pull Request"
3. Select your feature branch
4. Fill out the PR template with:
   - Clear description of changes
   - Link to related issues
   - Testing notes

### 3. Code Review Process

1. Automated checks must pass (CI, linting, tests)
2. At least one maintainer review required
3. Address review feedback promptly
4. Keep the PR updated with upstream changes

### 4. Merging

- Maintainers will merge approved PRs
- Use squash merge for feature branches
- Maintain clean commit history

## Reporting Issues

### Bug Reports

Use the bug report template and include:

1. **Clear description** of the issue
2. **Steps to reproduce** the problem
3. **Expected vs actual behavior**
4. **Environment details** (OS, Go version, etc.)
5. **Sample YAML** that reproduces the issue
6. **Error output** or logs

### Feature Requests

Use the feature request template and include:

1. **Problem description** - what problem does this solve?
2. **Proposed solution** - how should it work?
3. **Example usage** - show YAML configuration examples
4. **Alternatives considered** - what other approaches did you consider?

### Questions

Use the question template for:

- Usage questions
- Best practices
- Configuration help
- Integration guidance

## Development Tips

### Useful Commands

```bash
# Run specific tests
go test -v ./commands/echo -run TestEcho

# Run tests with verbose output
make test-unit -v

# Check code coverage for specific package
go test -cover ./commands/echo

# Run linter on specific files
golangci-lint run ./commands/echo/...

# Format specific files
gofmt -s -w ./commands/echo/
```

### Debugging

1. Use `fmt.Printf` for quick debugging
2. Use delve debugger for complex issues:
   ```bash
   dlv test ./commands/echo -- -test.run TestEcho
   ```

3. Use test files in `tests/` directory for integration testing

### Adding New Commands

1. Create new package in `commands/` directory
2. Implement `contract.Doer` interface
3. Implement `contract.CommandBuilder` interface
4. Add comprehensive tests
5. Update documentation
6. Add to default builders in `defaults/suite.go`

## Release Process

Releases are handled by maintainers:

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create git tag
4. GitHub Actions will build and create release

## Getting Help

- Check existing issues and PRs
- Ask questions in new issues using question template
- Join discussions in existing issues
- Read the documentation in `docs/` directory

Thank you for contributing to Declarate! 🎉
