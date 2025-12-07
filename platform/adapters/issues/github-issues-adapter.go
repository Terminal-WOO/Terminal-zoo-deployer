package issues

import (
	"context"
	"fmt"
)

// GitHubIssuesAdapter adapts platform events to GitHub Issues
type GitHubIssuesAdapter struct {
	owner      string
	repo       string
	githubToken string
}

// NewGitHubIssuesAdapter creates a new GitHub Issues adapter
func NewGitHubIssuesAdapter(owner, repo, githubToken string) *GitHubIssuesAdapter {
	return &GitHubIssuesAdapter{
		owner:      owner,
		repo:       repo,
		githubToken: githubToken,
	}
}

// OnDeploymentFailed creates a GitHub issue when deployment fails
func (a *GitHubIssuesAdapter) OnDeploymentFailed(ctx context.Context, data map[string]interface{}) error {
	// Create GitHub issue for failed deployment
	// In production, this would use GitHub API
	
	title := fmt.Sprintf("Deployment failed: %s/%s", data["namespace"], data["deployment"])
	body := fmt.Sprintf(`
## Deployment Failure

**Namespace**: %s
**Deployment**: %s
**Error**: %s
**Timestamp**: %s

Please investigate and fix the deployment issue.
`, data["namespace"], data["deployment"], data["error"], data["timestamp"])
	
	fmt.Printf("GitHub Issues: Creating issue - %s\n%s\n", title, body)
	
	// In production: POST to GitHub API to create issue
	_ = title
	_ = body
	
	return nil
}

// OnSLOViolation creates a GitHub issue when SLO is violated
func (a *GitHubIssuesAdapter) OnSLOViolation(ctx context.Context, data map[string]interface{}) error {
	title := fmt.Sprintf("SLO Violation: %s", data["slo"])
	body := fmt.Sprintf(`
## SLO Violation

**Service**: %s
**SLO**: %s
**Current Value**: %.2f
**Target Value**: %.2f
**Error Budget**: %.2f

Please investigate and improve service performance.
`, data["service"], data["slo"], data["current_value"], data["target_value"], data["error_budget"])
	
	fmt.Printf("GitHub Issues: Creating issue - %s\n%s\n", title, body)
	
	return nil
}

