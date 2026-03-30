# GitHub Actions CI/CD Configuration

This project includes comprehensive GitHub Actions workflows for continuous integration and continuous deployment.

## 📋 Workflows

### 1. **Tests & Quality** (`.github/workflows/test.yml`)
**Triggers:** Push to `main`/`develop`, Pull Requests

Performs:
- ✅ Run tests on Go 1.23 and 1.24
- ✅ Code format checking
- ✅ Test coverage analysis
- ✅ Upload coverage to Codecov
- ✅ Linting with golangci-lint
- ✅ Go vet analysis
- ✅ Build binary artifact

**Required Secrets:** None (unless using Codecov)

---

### 2. **Docker Build & Push** (`.github/workflows/docker.yml`)
**Triggers:** Push to `main`, Tags (v*), Pull Requests

Builds and pushes Docker images to:
- Docker Hub (optional, requires secrets)
- GitHub Container Registry (uses GITHUB_TOKEN)

**Features:**
- Multi-platform builds (Linux, ARM64)
- Semantic versioning support
- Build caching
- Automatic metadata extraction

**Required Secrets:**
```
DOCKER_USERNAME    # Your Docker Hub username (optional)
DOCKER_PASSWORD    # Your Docker Hub token (optional)
```

---

### 3. **Security Scanning** (`.github/workflows/security.yml`)
**Triggers:** Push to `main`/`develop`, Pull Requests, Daily schedule (3 AM UTC)

Performs:
- 🔒 Gosec security scanning
- 🔒 Dependency vulnerability checks (govulncheck)
- 🔒 Trivy filesystem scanning

Results uploaded to GitHub Security tab.

**Required Secrets:** None

---

### 4. **Release** (`.github/workflows/release.yml`)
**Triggers:** Tag push (v*)

Automatically:
- 📦 Builds binaries for multiple platforms
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
- 📦 Generate checksums
- 📦 Create GitHub Release with files

**Required Secrets:** None (uses GITHUB_TOKEN)

---

### 5. **Code Review** (`.github/workflows/review.yml`)
**Triggers:** Pull Requests

Performs:
- 📝 PR checklist generation
- 📝 Title format validation (feat:, fix:, etc.)
- 📝 Description length check

**Required Secrets:** None

---

## 🚀 Setup Instructions

### 1. Codecov Integration (Optional)

To enable coverage reports:

```bash
# No additional setup needed - uses default Codecov service
# Results are auto-uploaded from GitHub Actions
```

### 2. Docker Hub Integration (Optional)

To push to Docker Hub:

1. Create a Docker Hub account at https://hub.docker.com
2. Create an access token: https://hub.docker.com/settings/security
3. Add secrets to GitHub:

```bash
# Go to: Settings > Secrets and variables > Actions
# Add:
DOCKER_USERNAME = your_docker_username
DOCKER_PASSWORD = your_docker_access_token
```

### 3. GitHub Container Registry (Auto)

Uses `GITHUB_TOKEN` automatically.

Images available at: `ghcr.io/YOUR_USERNAME/sso:tag`

---

## 📊 Status Badges

Add to your `README.md`:

```markdown
[![Tests & Quality](https://github.com/YOUR_USERNAME/sso/actions/workflows/test.yml/badge.svg)](https://github.com/YOUR_USERNAME/sso/actions/workflows/test.yml)
[![Docker Build & Push](https://github.com/YOUR_USERNAME/sso/actions/workflows/docker.yml/badge.svg)](https://github.com/YOUR_USERNAME/sso/actions/workflows/docker.yml)
[![Security Scan](https://github.com/YOUR_USERNAME/sso/actions/workflows/security.yml/badge.svg)](https://github.com/YOUR_USERNAME/sso/actions/workflows/security.yml)
```

---

## 🔄 PR Workflow

1. **Create Feature Branch**
   ```bash
   git checkout -b feat/my-feature
   ```

2. **Make Changes & Commit**
   ```bash
   # Code...
   git add .
   git commit -m "feat: add awesome feature"
   ```

3. **Push & Create PR**
   ```bash
   git push origin feat/my-feature
   ```

4. **GitHub Actions Automatically:**
   - ✅ Runs tests
   - ✅ Checks code quality
   - ✅ Scans security
   - ✅ Posts checklist comment

5. **Review & Merge**
   - Address review comments
   - All checks must pass
   - Merge to main/develop

---

## 🏷️ Release Process

1. **Create Release Tag**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions Automatically:**
   - ✅ Builds binaries for all platforms
   - ✅ Generates checksums
   - ✅ Creates GitHub Release
   - ✅ Uploads artifacts

3. **Release Published**
   - Available at: https://github.com/YOUR_USERNAME/sso/releases/tag/v1.0.0

---

## 📌 Branch Protection Rules

Recommended GitHub Settings > Branches > Branch Protection:

```
✅ Require a pull request before merging
✅ Require status checks to pass before merging
   - Tests & Quality (test.yml)
   - Lint
   - Build
✅ Require branches to be up to date before merging
✅ Dismiss stale pull request approvals when new commits are pushed
✅ Require code review from at least 1 approving review
```

---

## 🎯 Best Practices

1. **Use descriptive commit messages**
   ```
   feat: add user registration endpoint
   fix: resolve JWT token expiration bug
   docs: update API documentation
   ```

2. **Keep PRs focused**
   - One feature per PR
   - Well-tested changes
   - Clear description

3. **Monitor CI/CD**
   - Check workflow runs: Actions tab
   - Address failures immediately
   - Keep tests passing

4. **Security**
   - Review dependency updates
   - Address security warnings
   - Keep Go updated

---

## 🆘 Troubleshooting

### Tests Failing
```bash
# Run locally before pushing
go test -v ./...
```

### Format Issues
```bash
# Auto-format code
go fmt ./...
```

### Linting Issues
```bash
# Run linter locally
golangci-lint run ./...
```

### Docker Build Fails
```bash
# Test Docker build locally
docker build -t sso:latest .
```

---

## 📖 Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Testing Best Practices](https://golang.org/doc/effective_go#testing)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices)
- [Semantic Versioning](https://semver.org/)
