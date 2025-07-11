# GitHub Repository Settings

This document outlines the recommended GitHub repository settings for optimal CI/CD operation.

## Branch Protection Rules

### For `main` (or `master`) branch:

1. **General Settings**:
   - ✅ Require a pull request before merging
   - ✅ Require approvals: 1
   - ✅ Dismiss stale PR approvals when new commits are pushed
   - ✅ Require review from code owners (if CODEOWNERS file is present)

2. **Status Checks**:
   - ✅ Require status checks to pass before merging
   - ✅ Require branches to be up to date before merging
   - **Required status checks**:
     - `test` (from ci.yml)
     - `lint` (from ci.yml)
     - `quality` (from code-quality.yml)
     - `test-matrix (ubuntu-latest, 1.20)` (at minimum)

3. **Additional Restrictions**:
   - ✅ Require conversation resolution before merging
   - ✅ Require signed commits (optional, recommended for security)
   - ✅ Include administrators (applies rules to admins too)
   - ✅ Allow force pushes: ❌ (disabled)
   - ✅ Allow deletions: ❌ (disabled)

## Repository Settings

### General

- **Features**:
  - ✅ Wikis: Enable if you want community documentation
  - ✅ Issues: Enable
  - ✅ Sponsorships: Enable if accepting donations
  - ✅ Preserve this repository: Enable for important projects
  - ✅ Discussions: Optional, good for community engagement

- **Pull Requests**:
  - ✅ Allow merge commits
  - ✅ Allow squash merging (recommended for clean history)
  - ✅ Allow rebase merging
  - ✅ Always suggest updating PR branches
  - ✅ Allow auto-merge
  - ✅ Automatically delete head branches

### Security & Analysis

1. **Security Features**:
   - ✅ Dependency graph: Enable
   - ✅ Dependabot alerts: Enable
   - ✅ Dependabot security updates: Enable
   - ✅ Code scanning: Enable (will be populated by workflows)
   - ✅ Secret scanning: Enable

2. **Code Scanning**:
   - Will be automatically configured by the security workflow
   - Results will appear in Security tab

### Actions

1. **General**:
   - ✅ Allow all actions and reusable workflows
   - ✅ Allow actions created by GitHub
   - ✅ Allow actions by Marketplace verified creators

2. **Artifact and Log Retention**:
   - Set to 90 days (default) or organization policy

3. **Fork Pull Request Workflows**:
   - ✅ Run workflows from fork pull requests
   - ✅ Send write tokens to workflows from fork pull requests: ❌
   - ✅ Send secrets to workflows from fork pull requests: ❌

### Secrets and Variables

#### Repository Secrets (if needed):
- `CODECOV_TOKEN`: For coverage reporting (optional, works without for public repos)

#### Repository Variables:
- None required for basic setup

### Collaborators & Teams

- Configure appropriate access levels:
  - **Admin**: Project maintainers
  - **Write**: Core contributors
  - **Read**: General contributors

### Notifications

#### Email Notifications:
- ✅ Actions workflow runs: For maintainers
- ✅ Dependabot alerts: For maintainers
- ✅ Vulnerability alerts: For maintainers

#### GitHub App Notifications:
- Configure Slack/Discord integrations if desired

## Repository Topics

Add relevant topics for discoverability:
- `go`
- `golang`
- `testing`
- `api-testing`
- `declarative`
- `yaml`
- `automation`
- `ci-cd`
- `quality-assurance`

## About Section

### Description:
"Declarative testing library for APIs, CLIs and anything else using simple YAML syntax"

### Website:
Link to documentation site if available

### Topics:
Add the topics listed above

## Social Preview

Upload a social preview image (1280×640px) that represents the project visually.

## License

Ensure LICENSE file is present and properly configured in repository settings.

## Security Policy

Consider adding a SECURITY.md file with vulnerability reporting instructions.

## Code of Conduct

Consider adding a CODE_OF_CONDUCT.md file for community guidelines.

## Issue Templates Configuration

The issue templates are already configured in `.github/ISSUE_TEMPLATE/`. 

To disable blank issues (force users to use templates):
1. Go to Settings → Features → Issues
2. Uncheck "Issues" temporarily
3. Re-check "Issues"
4. The templates will now be enforced

## Quick Setup Checklist

After pushing the `feature/github-actions` branch and merging:

- [ ] Configure branch protection rules for main branch
- [ ] Enable required status checks
- [ ] Configure security features (Dependabot, code scanning)
- [ ] Set up notifications for maintainers
- [ ] Add repository topics
- [ ] Review and adjust Actions permissions
- [ ] Add CODEOWNERS file if needed
- [ ] Configure integrations (Slack, Discord, etc.)
- [ ] Set up project boards if using GitHub Projects
- [ ] Review collaborator access levels

## Monitoring

Regularly review:
- Security alerts and Dependabot PRs
- Action workflow success rates
- Code coverage trends
- Performance of status checks

The CI/CD pipeline will provide ongoing monitoring of code quality, security, and functionality.
