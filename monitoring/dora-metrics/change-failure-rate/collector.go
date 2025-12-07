package changefailurerate

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
)

// ChangeFailureRateCollector collects change failure rate metrics
type ChangeFailureRateCollector struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewCollector creates a new change failure rate collector
func NewCollector(clientset *kubernetes.Clientset, namespace string) *ChangeFailureRateCollector {
	return &ChangeFailureRateCollector{
		clientset: clientset,
		namespace: namespace,
	}
}

// FailureRate represents the change failure rate for a time period
type FailureRate struct {
	Period           string    `json:"period"`
	TotalDeployments int       `json:"total_deployments"`
	FailedDeployments int      `json:"failed_deployments"`
	FailureRate      float64   `json:"failure_rate"` // Percentage
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
}

// Collect collects change failure rate metrics for the specified time period
func (c *ChangeFailureRateCollector) Collect(ctx context.Context, startTime, endTime time.Time) (*FailureRate, error) {
	// Get all deployments in the namespace
	deployments, err := c.clientset.AppsV1().Deployments(c.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	totalDeployments := 0
	failedDeployments := 0

	for _, deployment := range deployments.Items {
		// Check if deployment was created or updated in the time range
		creationTime := deployment.CreationTimestamp.Time
		wasInRange := (creationTime.After(startTime) && creationTime.Before(endTime))

		// Check deployment conditions for failures
		for _, condition := range deployment.Status.Conditions {
			if condition.Type == appsv1.DeploymentReplicaFailure && condition.Status == "True" {
				if wasInRange || (condition.LastTransitionTime.Time.After(startTime) && condition.LastTransitionTime.Time.Before(endTime)) {
					failedDeployments++
					break
				}
			}
		}

		// Count deployments in the time range
		if wasInRange {
			totalDeployments++
		}
	}

	// Calculate failure rate
	failureRate := 0.0
	if totalDeployments > 0 {
		failureRate = (float64(failedDeployments) / float64(totalDeployments)) * 100
	}

	return &FailureRate{
		Period:           endTime.Sub(startTime).String(),
		TotalDeployments: totalDeployments,
		FailedDeployments: failedDeployments,
		FailureRate:      failureRate,
		StartTime:        startTime,
		EndTime:          endTime,
	}, nil
}

// CollectWeekly collects weekly change failure rate
func (c *ChangeFailureRateCollector) CollectWeekly(ctx context.Context) (*FailureRate, error) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	return c.Collect(ctx, startOfWeek, now)
}

// CollectMonthly collects monthly change failure rate
func (c *ChangeFailureRateCollector) CollectMonthly(ctx context.Context) (*FailureRate, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return c.Collect(ctx, startOfMonth, now)
}

