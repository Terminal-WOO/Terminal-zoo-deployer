package leadtime

import (
	"context"
	"time"
)

// LeadTimeCollector collects lead time metrics (commit to deployment)
type LeadTimeCollector struct {
	// Git repository information
	gitRepoURL string
	gitBranch  string
	
	// Deployment tracking
	deploymentStartTime time.Time
	deploymentEndTime   time.Time
}

// NewCollector creates a new lead time collector
func NewCollector(gitRepoURL, gitBranch string) *LeadTimeCollector {
	return &LeadTimeCollector{
		gitRepoURL: gitRepoURL,
		gitBranch:  gitBranch,
	}
}

// LeadTime represents the time from commit to deployment
type LeadTime struct {
	CommitHash      string        `json:"commit_hash"`
	CommitTime      time.Time     `json:"commit_time"`
	DeploymentTime  time.Time     `json:"deployment_time"`
	LeadTimeSeconds int64         `json:"lead_time_seconds"`
	LeadTimeMinutes float64       `json:"lead_time_minutes"`
	LeadTimeHours   float64       `json:"lead_time_hours"`
}

// CalculateLeadTime calculates lead time for a deployment
func (c *LeadTimeCollector) CalculateLeadTime(ctx context.Context, commitHash string, commitTime, deploymentTime time.Time) (*LeadTime, error) {
	leadTimeDuration := deploymentTime.Sub(commitTime)
	
	return &LeadTime{
		CommitHash:      commitHash,
		CommitTime:      commitTime,
		DeploymentTime:  deploymentTime,
		LeadTimeSeconds: int64(leadTimeDuration.Seconds()),
		LeadTimeMinutes: leadTimeDuration.Minutes(),
		LeadTimeHours:   leadTimeDuration.Hours(),
	}, nil
}

// AverageLeadTime calculates average lead time for multiple deployments
func AverageLeadTime(leadTimes []*LeadTime) *LeadTime {
	if len(leadTimes) == 0 {
		return nil
	}

	var totalSeconds int64
	for _, lt := range leadTimes {
		totalSeconds += lt.LeadTimeSeconds
	}

	avgSeconds := totalSeconds / int64(len(leadTimes))
	avgDuration := time.Duration(avgSeconds) * time.Second

	return &LeadTime{
		LeadTimeSeconds: avgSeconds,
		LeadTimeMinutes: avgDuration.Minutes(),
		LeadTimeHours:   avgDuration.Hours(),
	}
}

