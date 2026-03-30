#!/bin/bash

# ci-setup.sh - Quick CI/CD Setup Script for SSO Project
# This script provides instructions for setting up CI/CD

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          SSO Project - GitHub Actions CI/CD Setup              ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 Workflow Files Created:${NC}"
echo "   ✓ .github/workflows/test.yml         - Tests & Quality Checks"
echo "   ✓ .github/workflows/docker.yml       - Docker Build & Push"
echo "   ✓ .github/workflows/security.yml     - Security Scanning"
echo "   ✓ .github/workflows/release.yml      - Automated Releases"
echo "   ✓ .github/workflows/review.yml       - PR Automation"
echo ""

echo -e "${BLUE}📚 Documentation Files Created:${NC}"
echo "   ✓ .github/SETUP_CI.md                - Quick Start Guide"
echo "   ✓ .github/CI-CD.md                   - Detailed Workflow Info"
echo "   ✓ .github/BRANCH_PROTECTION.md       - Branch Protection Setup"
echo "   ✓ renovate.json                      - Dependency Updates Config"
echo ""

echo -e "${YELLOW}🚀 Next Steps:${NC}"
echo ""

echo -e "${GREEN}1. Push to GitHub${NC}"
echo "   git add .github renovate.json"
echo "   git commit -m 'chore: add GitHub Actions CI/CD workflows'"
echo "   git push origin main"
echo ""

echo -e "${GREEN}2. Verify Workflows${NC}"
echo "   → Go to: https://github.com/YOUR_USERNAME/sso/actions"
echo "   → You should see workflows running automatically"
echo ""

echo -e "${GREEN}3. (Optional) Setup GitHub Secrets${NC}"
echo "   For Docker Hub integration:"
echo "   → Settings → Secrets and variables → Actions → New repository secret"
echo "   → Name: DOCKER_USERNAME"
echo "   → Name: DOCKER_PASSWORD"
echo ""

echo -e "${GREEN}4. (Optional) Enable Dependabot${NC}"
echo "   → Settings → Code security and analysis"
echo "   → Enable: Dependabot alerts"
echo "   → Enable: Dependabot security updates"
echo ""

echo -e "${GREEN}5. Setup Branch Protection (Recommended)${NC}"
echo "   → See .github/BRANCH_PROTECTION.md for detailed guide"
echo ""

echo -e "${YELLOW}📊 Workflow Details:${NC}"
echo ""
echo "   Test:       Runs on: push to main/develop, PRs"
echo "              Does:    tests, lint, coverage, build"
echo ""
echo "   Docker:     Runs on: push to main, tags, PRs"
echo "              Does:    build Docker image, push to registries"
echo ""
echo "   Security:   Runs on: push, PRs, daily schedule"
echo "              Does:    gosec, dependency check, trivy scan"
echo ""
echo "   Release:    Runs on: git tags (v*)"
echo "              Does:    build multi-platform binaries, create release"
echo ""
echo "   Review:     Runs on: PRs"
echo "              Does:    checklist, title validation"
echo ""

echo -e "${YELLOW}📖 Read the Guides:${NC}"
echo "   → .github/SETUP_CI.md         - How to set up CI"
echo "   → .github/CI-CD.md            - Detailed workflow info"
echo "   → .github/BRANCH_PROTECTION.md - Branch protection rules"
echo ""

echo -e "${GREEN}✅ All CI/CD files are ready!${NC}"
echo ""
echo "For detailed information, see .github/SETUP_CI.md"
