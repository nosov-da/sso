# GitHub Branch Protection Configuration

This guide helps setup recommended branch protection rules for the SSO project.

## 🔐 Setup Branch Protection Rules

### Step 1: Navigate to Settings
1. Go to your GitHub repository
2. Click **Settings** tab
3. Go to **Branches** (left sidebar)
4. Click **Add rule** under "Branch protection rules"

### Step 2: Configure Main Branch

**Branch name pattern:** `main`

#### Require pull request reviews before merging
- ✅ Enable "Require a pull request before merging"
- ✅ Set "Require approvals": 1
- ✅ Enable "Dismiss stale pull request approvals when new commits are pushed"
- ✅ Enable "Require review from Code Owners" (optional)

#### Require status checks to pass before merging
- ✅ Enable "Require status checks to pass before merging"
- ✅ Enable "Require branches to be up to date before merging"

**Required status checks:**
```
test
lint
build
```

#### Require conversation resolution
- ✅ Enable "Require conversations to be resolved before merging"

#### Restrict who can push to matching branches
- ⚠️ Optional: Limit who can push (usually admins only)

#### Do not allow bypassing the above settings
- ✅ Enable (recommended)

---

### Step 3: Configure Develop Branch

**Branch name pattern:** `develop`

**Same settings as main, but:**
- ✅ Require approvals: 1 (or 0 for faster development)
- ⚠️ Enable "Allow specified actors to bypass required pull requests"
  - Add team leads who can fast-track hotfixes

---

## 📝 Code Owners Configuration

Create `.github/CODEOWNERS` file:

```
# File/directory owners

# Root
* @YOUR_USERNAME

# Domain layer
/domain/ @YOUR_USERNAME
/domain/user/ @YOUR_USERNAME

# Application layer
/application/ @YOUR_USERNAME

# Infrastructure
/infrastructure/ @YOUR_USERNAME

# Internal APIs
/internal/api/ @YOUR_USERNAME

# Configuration and DevOps
/.github/ @YOUR_USERNAME
Dockerfile @YOUR_USERNAME
docker-compose.yml @YOUR_USERNAME

# Documentation
*.md @YOUR_USERNAME
docs/ @YOUR_USERNAME
```

---

## ✅ Default Branch Settings

1. Go to **Settings** > **General**
2. Set **Default branch** to: `main`
3. Enable **Always suggest to update branch protection rules**

---

## 🔔 Status Check Requirements

The workflows define these status checks:

| Workflow | Status Check | Purpose |
|----------|-------------|---------|
| test.yml | test | Run unit tests |
| test.yml | lint | Check code quality |
| test.yml | build | Build binary |
| docker.yml | docker | Build Docker image |
| security.yml | security | Security scanning |

---

## 📋 Pull Request Template

Create `.github/pull_request_template.md`:

```markdown
## Description
<!-- Describe your changes in detail -->

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to change)
- [ ] Documentation update

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests passed locally with my changes

## Testing
<!-- Describe how you tested your changes -->

## Related Issues
<!-- Link to related issues: Closes #123 -->

## Screenshots (if applicable)
<!-- Add screenshots for UI changes -->
```

---

## 🚫 Branch Protection Rules - Do Not Allow

The following should be prevented:

- ❌ Direct commits to `main` (only PRs)
- ❌ Force pushes to `main` or `develop`
- ❌ Deleting branches
- ❌ Merging without passing status checks
- ❌ Merging without reviews

---

## 🔄 Recommended Workflow

```
1. Create feature branch from develop
   git checkout -b feat/my-feature develop

2. Make changes and push
   git push origin feat/my-feature

3. Create Pull Request
   - GitHub runs automated CI
   - Request reviewers

4. Address feedback
   - Make requested changes
   - Push again (CI runs automatically)

5. Get approval
   - Minimum 1 approval required

6. Merge to develop
   - Use "Squash and merge" or "Create a merge commit"
   - Delete branch after merge

7. Create release when ready
   - Create PR from develop to main
   - Same review process
   - Merge to main
   - Tag release (v1.0.0)
```

---

## 📊 Monitoring

### View Status Checks
1. Go to **Actions** tab
2. Filter by branch or workflow
3. Click on individual runs to see details

### View Branch Protections
1. Go to **Settings** > **Branches**
2. Click **Edit** on rule to modify
3. Click **Delete** to remove

### View PR Status
1. On any Pull Request
2. Check status checks at bottom
3. Click details to see logs

---

## 🆘 Troubleshooting

### "Required status checks are pending"
- Wait for CI to complete
- Check Actions tab for results

### "This branch has conflicts that must be resolved"
```bash
git fetch origin develop
git rebase origin/develop
# Fix conflicts
git push -f origin feat/my-feature
```

### "At least 1 approving review is required"
- Request review from team member
- They'll review and approve

### Cannot force push
- Feature branch protection enabled
- Use normal push instead

---

## 🔐 Security Best Practices

1. **Never push secrets to GitHub**
   - Use GitHub Secrets for sensitive data
   - Scan commits: `git log -p | grep -i password`

2. **Enable 2FA**
   - Account Settings > Security
   - Require 2FA for all organization members

3. **Review Dependencies**
   - Check Dependabot alerts
   - Keep Go and dependencies updated

4. **Monitor Access**
   - Review collaborators
   - Remove inactive members
   - Use teams for permissions

---

## 📚 Additional Resources

- [GitHub Branch Protection Rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/managing-a-branch-protection-rule)
- [CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
- [Required Status Checks](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)
