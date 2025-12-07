package mttr

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
)

// MTTRCollector collects Mean Time to Recovery metrics
type MTTRCollector struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewCollector creates a new MTTR collector
func NewCollector(clientset *kubernetes.Clientset, namespace string) *MTTRCollector {
	return &MTTRCollector{
		clientset: clientset,
		namespace: namespace,
	}
}

// RecoveryTime represents recovery time for a deployment failure
type RecoveryTime struct {
	DeploymentName string        `json:"deployment_name"`
	FailureTime    time.Time    `json:"failure_time"`
	RecoveryTime   time.Time    `json:"recovery_time"`
	Duration       time.Duration `json:"duration"`
	DurationSeconds int64       `json:"duration_seconds"`
	DurationMinutes float64     `json:"duration_minutes"`
}

// MTTR represents Mean Time to Recovery for a time period
type MTTR struct {
	Period          string        `json:"period"`
	RecoveryTimes   []RecoveryTime `json:"recovery_times"`
	MeanTimeSeconds int64         `json:"mean_time_seconds"`
	MeanTimeMinutes float64       `json:"mean_time_minutes"`
	MeanTimeHours   float64       `json:"mean_time_hours"`
	StartTime       time.Time     `json:"start_time"`
	EndTime         time.Time     `json:"end_time"`
}

// Collect collects MTTR metrics for the specified time period
func (c *MTTRCollector) Collect(ctx context.Context, startTime, endTime time.Time) (*MTTR, error) {
	// Get all deployments in the namespace
	deployments, err := c.clientset.AppsV1().Deployments(c.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var recoveryTimes []RecoveryTime

	for _, deployment := range deployments.Items {
		var failureTime *time.Time
		var recoveryTime *time.Time

		// Check deployment conditions for failures and recoveries
		for _, condition := range deployment.Status.Conditions {
			if condition.Type == appsv1.DeploymentReplicaFailure {
				if condition.Status == "True" {
					// Failure detected
					if failureTime == nil || condition.LastTransitionTime.Time.Before(*failureTime) {
						t := condition.LastTransitionTime.Time
						failureTime = &t
					}
				} else if condition.Status == "False" && failureTime != nil {
					// Recovery detected
					if recoveryTime == nil || condition.LastTransitionTime.Time.After(*recoveryTime) {
						t := condition.LastTransitionTime.Time
						recoveryTime = &t
					}
				}
			}

			// Check if deployment is available (recovered)
			if condition.Type == appsv1.DeploymentAvailable && condition.Status == "True" {
				if failureTime != nil && recoveryTime == nil {
					t := condition.LastTransitionTime.Time
					recoveryTime = &t
				}
			}
		}

		// If we have both failure and recovery times in the range, calculate recovery time
		if failureTime != nil && recoveryTime != nil {
			if failureTime.After(startTime) && failureTime.Before(endTime) {
				duration := recoveryTime.Sub(*failureTime)
				recoveryTimes = append(recoveryTimes, RecoveryTime{
					DeploymentName:  deployment.Name,
					FailureTime:     *failureTime,
					RecoveryTime:    *recoveryTime,
					Duration:        duration,
					DurationSeconds: int64(duration.Seconds()),
					DurationMinutes: duration.Minutes(),
				})
			}
		}
	}

	// Calculate mean time to recovery
	var totalSeconds int64
	for _, rt := range recoveryTimes {
		totalSeconds += rt.DurationSeconds
	}

	meanTimeSeconds := int64(0)
	meanTimeMinutes := 0.0
	meanTimeHours := 0.0

	if len(recoveryTimes) > 0 {
		meanTimeSeconds = totalSeconds / int64(len(recoveryTimes))
		meanDuration := time.Duration(meanTimeSeconds) * time.Second
		meanTimeMinutes = meanDuration.Minutes()
		meanTimeHours = meanDuration.Hours()
	}

	return &MTTR{
		Period:          endTime.Sub(startTime).String(),
		RecoveryTimes:   recoveryTimes,
		MeanTimeSeconds: meanTimeSeconds,
		MeanTimeMinutes: meanTimeMinutes,
		MeanTimeHours:   meanTimeHours,
		StartTime:       startTime,
		EndTime:         endTime,
	}, nil
}

// CollectWeekly collects weekly MTTR
func (c *MTTRCollector) CollectWeekly(ctx context.Context) (*MTTR, error) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	return c.Collect(ctx, startOfWeek, now)
}

// CollectMonthly collects monthly MTTR
func (c *MTTRCollector) CollectMonthly(ctx context.Context) (*MTTR, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return c.Collect(ctx, startOfMonth, now)
}

