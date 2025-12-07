package deploymentfrequency

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentFrequencyCollector collects deployment frequency metrics
type DeploymentFrequencyCollector struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewCollector creates a new deployment frequency collector
func NewCollector(clientset *kubernetes.Clientset, namespace string) *DeploymentFrequencyCollector {
	return &DeploymentFrequencyCollector{
		clientset: clientset,
		namespace: namespace,
	}
}

// DeploymentCount represents deployment count for a time period
type DeploymentCount struct {
	Period    string    `json:"period"`
	Count     int       `json:"count"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// Collect collects deployment frequency metrics for the specified time period
func (c *DeploymentFrequencyCollector) Collect(ctx context.Context, startTime, endTime time.Time) (*DeploymentCount, error) {
	// Get deployments in the specified time range
	deployments, err := c.clientset.AppsV1().Deployments(c.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Count deployments that were created or updated in the time range
	count := 0
	for _, deployment := range deployments.Items {
		// Check if deployment was created or updated in the time range
		creationTime := deployment.CreationTimestamp.Time
		if creationTime.After(startTime) && creationTime.Before(endTime) {
			count++
			continue
		}

		// Check if deployment was updated in the time range
		for _, condition := range deployment.Status.Conditions {
			if condition.LastTransitionTime.Time.After(startTime) && condition.LastTransitionTime.Time.Before(endTime) {
				count++
				break
			}
		}
	}

	return &DeploymentCount{
		Period:    endTime.Sub(startTime).String(),
		Count:     count,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// CollectDaily collects daily deployment frequency
func (c *DeploymentFrequencyCollector) CollectDaily(ctx context.Context) (*DeploymentCount, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return c.Collect(ctx, startOfDay, now)
}

// CollectWeekly collects weekly deployment frequency
func (c *DeploymentFrequencyCollector) CollectWeekly(ctx context.Context) (*DeploymentCount, error) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	return c.Collect(ctx, startOfWeek, now)
}

// CollectMonthly collects monthly deployment frequency
func (c *DeploymentFrequencyCollector) CollectMonthly(ctx context.Context) (*DeploymentCount, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return c.Collect(ctx, startOfMonth, now)
}

