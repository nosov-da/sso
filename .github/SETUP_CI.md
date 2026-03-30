# CI/CD Setup Guide

This guide helps you set up the complete CI/CD pipeline for the SSO project on GitHub.

## 📋 Prerequisites

- GitHub repository access (you own or have admin rights)
- Optional: Docker Hub account for image hosting
- Optional: Codecov account for coverage reports

---

## 🚀 Quick Start (5 minutes)

### 1. Push Code to GitHub
```bash
cd /path/to/sso
git remote add origin https://github.com/YOUR_USERNAME/sso.git
git branch -M main
git push -u origin main
```

### 2. GitHub Actions Auto-Discovery
Once code is pushed, GitHub automatically:
- ✅ Detects `.github/workflows/` directory
- ✅ Enables GitHub Actions (usually auto-enabled)
- ✅ Runs workflows on push/PR events

### 3. Verify Workflows
1. Go to your GitHub repository
2. Click **Actions** tab
3. See workflows running automatically

**That's it!** Your CI is now active.

---

## 🔧 Optional Configurations

### Enable Codecov (Coverage Reports)

No setup required! Codecov automatically:
1. Detects your GitHub Actions
2. Processes coverage data
3. Posts results to PRs
4. Generates badges

**To view coverage:**
- Visit: https://codecov.io/gh/YOUR_USERNAME/sso

### Enable Docker Hub Integration

1. **Get Docker Hub Access Token:**
   ```
   https://hub.docker.com/settings/security
   ```

2. **Add GitHub Secrets:**
   ```
   Settings > Secrets and variables > Actions > New repository secret
   
   DOCKER_USERNAME = your_docker_username
   DOCKER_PASSWORD = your_access_token
   ```

3. **Verify Docker Build:**
   - Push to `main` branch
   - Go to Actions tab
   - Docker Build workflow should run

### Enable Dependabot / Renovate

**Option 1: GitHub Dependabot** (Built-in)
1. Go to **Settings** > **Code security and analysis**
2. Enable **Dependabot alerts**
3. Enable **Dependabot security updates**
4. Enable **Dependabot version updates** (set schedule)

**Option 2: Renovate** (More features)
1. Commit `renovate.json` to your repo (already provided)
2. Install Renovate app: https://github.com/apps/renovate
3. Select your repository
4. Renovate will auto-create update PRs

### Setup Branch Protection

Follow [.github/BRANCH_PROTECTION.md](./.github/BRANCH_PROTECTION.md) guide

---

## 📊 Workflow Overview

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| **test.yml** | Push, PR | Run tests, lint, coverage |
| **docker.yml** | Push main, Tags, PR | Build & push Docker images |
| **security.yml** | Push, PR, Daily | Security scanning |
| **release.yml** | Tag push | Auto-build multiplatform binaries |
| **review.yml** | PR creation | PR checklist & validation |

---

## 🔍 Monitoring CI/CD

### View Workflow Runs
```
GitHub > Actions > Select Workflow > View Runs
```

### View Status Badges
Add to `README.md`:
```markdown
[![Tests & Quality](https://github.com/YOUR_USERNAME/sso/actions/workflows/test.yml/badge.svg)](https://github.com/YOUR_USERNAME/sso/actions/workflows/test.yml)
[![Docker](https://github.com/YOUR_USERNAME/sso/actions/workflows/docker.yml/badge.svg)](https://github.com/YOUR_USERNAME/sso/actions/workflows/docker.yml)
[![Security](https://github.com/YOUR_USERNAME/sso/actions/workflows/security.yml/badge.svg)](https://github.com/YOUR_USERNAME/sso/actions/workflows/security.yml)
```

### Debug Failed Workflows
1. Click on failed workflow run
2. Click on failed job
3. Check "Run" section for logs
4. Fix locally and push again

---

## 🎯 First Deployment Steps

### 1. Create First Release
```bash
# Tag the release
git tag -a v1.0.0 -m "Initial release"
git push origin v1.0.0
```

Result: GitHub automatically:
- ✅ Builds binaries for 5 platforms
- ✅ Generates checksums
- ✅ Creates Release page with downloads

### 2. Deploy With Docker
```bash
# Image automatically built and available at:
# Docker Hub: docker pull YOUR_USERNAME/sso:v1.0.0
# GitHub Registry: docker pull ghcr.io/YOUR_USERNAME/sso:v1.0.0

# Run the image
docker run -e SSO_JWT_SECRET="your-secret" ghcr.io/YOUR_USERNAME/sso:v1.0.0
```

---

## 📚 File Structure

```
.github/
├── workflows/
│   ├── test.yml          # Tests, lint, coverage
│   ├── docker.yml        # Docker builds
│   ├── security.yml      # Security scanning
│   ├── release.yml       # Release binaries
│   └── review.yml        # PR automation
├── CI-CD.md              # This guide (detailed)
└── BRANCH_PROTECTION.md  # Branch protection rules

renovate.json            # Dependency updates
.gitignore               # Git ignore patterns
Dockerfile              # Docker configuration
Makefile                # Local build commands
```

---

## 🆘 Common Issues & Fixes

### Issue: "workflow syntax is invalid"
**Solution:**
```bash
# Validate locally
npm install -g ajv-cli
ajv validate -s node_modules/ajv/package.json -d .github/workflows/test.yml
```

### Issue: Tests fail in CI but pass locally
**Solution:**
```bash
# Test with Go 1.24 (same as CI)
go version  # Check local version
# If different, use Docker:
docker run -v $PWD:/app -w /app golang:1.24 go test ./...
```

### Issue: Docker build fails
**Solution:**
```bash
# Test Docker build locally
docker build -t sso:test .
# Debug with
DOCKER_BUILDKIT=1 docker build --progress=plain .
```

### Issue: "Permission denied" on git push
**Solution:**
```bash
# Use SSH instead of HTTPS
git remote set-url origin git@github.com:YOUR_USERNAME/sso.git
# Or use GitHub CLI
gh auth login
```

---

## ✅ Verification Checklist

- [ ] All workflows are visible in Actions tab
- [ ] Test workflow passes on first push
- [ ] PR creates automatic comments
- [ ] Status checks show up on PRs
- [ ] Coverage is tracked (if Codecov enabled)
- [ ] Docker images build (if Docker Hub secrets added)
- [ ] Release workflow works (try creating tag)
- [ ] Branch protection rules are active

---

## 🔐 Security Considerations

1. **Secrets Management**
   ```
   ❌ Never commit secrets to repo
   ✅ Use GitHub Secrets for:
      - Docker Hub credentials
      - Database passwords
      - API keys
   ```

2. **Access Control**
   ```
   ✅ Use branch protection rules
   ✅ Require PR reviews before merge
   ✅ Enable 2FA on your account
   ✅ Regularly audit collaborators
   ```

3. **Dependency Security**
   ```
   ✅ Enable Dependabot alerts
   ✅ Review dependency updates
   ✅ Run security scans (included)
   ✅ Keep Go updated
   ```

---

## 📖 Next Steps

1. **Read detailed docs:**
   - [CI-CD.md](./.github/CI-CD.md) - Detailed workflow info
   - [BRANCH_PROTECTION.md](./.github/BRANCH_PROTECTION.md) - protection setup

2. **Customize workflows:**
   - Edit `.github/workflows/*.yml` for your needs
   - Add custom job steps
   - Adjust triggers/schedules

3. **Monitor and improve:**
   - Watch Actions tab
   - Check coverage reports
   - Update workflows based on feedback

---

## 📞 Support

If workflows fail:
1. Check GitHub Actions logs
2. Review workflow file syntax
3. Test commands locally
4. Check secrets are configured
5. Verify branch/tag patterns match

---

## 🔗 Links

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Workflow Syntax Reference](https://docs.github.com/en/actions/reference)
- [Using Secrets in GitHub Actions](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions)
- [Codecov Documentation](https://docs.codecov.io/)
- [Docker Documentation](https://docs.docker.com/)
