# GitHub Actions Status Badges

Add these badges to your main README.md to show the status of your CI/CD pipelines:

## Basic CI Badge
```markdown
[![CI](https://github.com/ixpectus/declarate/actions/workflows/ci.yml/badge.svg)](https://github.com/ixpectus/declarate/actions/workflows/ci.yml)
```

## Code Quality Badge
```markdown
[![Code Quality](https://github.com/ixpectus/declarate/actions/workflows/code-quality.yml/badge.svg)](https://github.com/ixpectus/declarate/actions/workflows/code-quality.yml)
```

## Go Compatibility Badge
```markdown
[![Go Compatibility](https://github.com/ixpectus/declarate/actions/workflows/compatibility.yml/badge.svg)](https://github.com/ixpectus/declarate/actions/workflows/compatibility.yml)
```

## Security Badge
```markdown
[![Security](https://github.com/ixpectus/declarate/actions/workflows/ci.yml/badge.svg?event=schedule)](https://github.com/ixpectus/declarate/actions/workflows/ci.yml)
```

## All Badges Combined
```markdown
[![CI](https://github.com/ixpectus/declarate/actions/workflows/ci.yml/badge.svg)](https://github.com/ixpectus/declarate/actions/workflows/ci.yml)
[![Code Quality](https://github.com/ixpectus/declarate/actions/workflows/code-quality.yml/badge.svg)](https://github.com/ixpectus/declarate/actions/workflows/code-quality.yml)
[![Go Compatibility](https://github.com/ixpectus/declarate/actions/workflows/compatibility.yml/badge.svg)](https://github.com/ixpectus/declarate/actions/workflows/compatibility.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ixpectus/declarate)](https://goreportcard.com/report/github.com/ixpectus/declarate)
[![codecov](https://codecov.io/gh/ixpectus/declarate/branch/main/graph/badge.svg)](https://codecov.io/gh/ixpectus/declarate)
```

## Workflow Status Descriptions

| Workflow | Description | Trigger |
|----------|-------------|---------|
| **CI** | Main continuous integration pipeline with tests, linting, and security | Push/PR to main branches |
| **Code Quality** | Comprehensive code quality checks and static analysis | Push/PR + daily schedule |
| **Go Compatibility** | Cross-platform and cross-version compatibility testing | Push/PR to main branches |
| **Release** | Automated release creation with multi-platform binaries | Git tags (v*) |

## Setting Up External Services

### Codecov
1. Sign up at [codecov.io](https://codecov.io)
2. Add your repository
3. No additional configuration needed - coverage is uploaded automatically

### Go Report Card
1. Visit [goreportcard.com](https://goreportcard.com)
2. Enter your repository URL
3. The badge will be automatically available

## Branch Protection Settings

Recommended settings for main/master branch:

```yaml
protection_rules:
  required_status_checks:
    strict: true
    contexts:
      - "test"
      - "lint" 
      - "quality"
  enforce_admins: false
  required_pull_request_reviews:
    dismiss_stale_reviews: true
    require_code_owner_reviews: true
    required_approving_review_count: 1
  restrictions: null
```

## Notifications

Set up GitHub notifications for:
- Failed CI runs
- Security alerts
- Dependency updates

Configure in Repository Settings > Notifications.
