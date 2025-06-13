# GitHub Repository Settings for Declarate

This document describes recommended GitHub repository settings for optimal CI/CD workflow.

## Branch Protection Rules

### Main/Master Branch
Set up the following branch protection rules for the main/master branch:

1. **Require pull request reviews before merging**
   - Required number of reviews: 1
   - Dismiss stale reviews when new commits are pushed: ✓
   - Require review from code owners: ✓

2. **Require status checks to pass before merging**
   - Require branches to be up to date before merging: ✓
   - Required status checks:
     - `test` (from CI workflow)
     - `lint` (from CI workflow)
     - `quality` (from Code Quality workflow)

3. **Require conversation resolution before merging**: ✓

4. **Restrict pushes that create files larger than 100 MB**: ✓

5. **Require signed commits**: ✓ (recommended)

## Repository Settings

### General
- Default branch: `main` or `master`
- Allow merge commits: ✓
- Allow squash merging: ✓
- Allow rebase merging: ✓
- Automatically delete head branches: ✓

### Security
- Enable vulnerability alerts: ✓
- Enable security updates: ✓
- Enable secret scanning: ✓
- Enable push protection for secrets: ✓

### Code security and analysis
- Enable Dependabot alerts: ✓
- Enable Dependabot security updates: ✓
- Enable CodeQL analysis: ✓

## Secrets Configuration

Add the following secrets in repository settings (if needed):

1. `CODECOV_TOKEN` - for code coverage reporting
2. `DOCKER_USERNAME` and `DOCKER_PASSWORD` - if using Docker registry
3. Any database connection strings for integration tests

## Labels

Create the following labels for better issue and PR management:

- `bug` (red) - Something isn't working
- `enhancement` (blue) - New feature or request
- `documentation` (yellow) - Improvements or additions to documentation
- `good first issue` (green) - Good for newcomers
- `help wanted` (purple) - Extra attention is needed
- `question` (pink) - Further information is requested
- `dependencies` (orange) - Pull requests that update a dependency file
- `github-actions` (gray) - Pull requests that update GitHub Actions
- `go` (blue) - Go related issues
- `tests` (green) - Test related changes
- `performance` (orange) - Performance improvements
- `security` (red) - Security related issues

## Notifications

Configure notifications to ensure maintainers are alerted about:
- Failed CI/CD workflows
- Security alerts
- New issues and PRs
- Review requests
