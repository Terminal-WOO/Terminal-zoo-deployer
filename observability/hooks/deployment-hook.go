package hooks

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Deployment metrics
	deploymentCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "platform_deployments_total",
			Help: "Total number of deployments",
		},
		[]string{"namespace", "deployment", "status"},
	)

	deploymentDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "platform_deployment_duration_seconds",
			Help:    "Deployment duration in seconds",
			Buckets: prometheus.ExponentialBuckets(1, 2, 10), // 1s to ~17 minutes
		},
		[]string{"namespace", "deployment"},
	)

	deploymentFailures = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "platform_deployment_failures_total",
			Help: "Total number of failed deployments",
		},
		[]string{"namespace", "deployment", "reason"},
	)
)

// DeploymentHook provides observability hooks for deployments
type DeploymentHook struct {
	startTime time.Time
	namespace string
	deployment string
}

// NewDeploymentHook creates a new deployment hook
func NewDeploymentHook(namespace, deployment string) *DeploymentHook {
	return &DeploymentHook{
		startTime: time.Now(),
		namespace: namespace,
		deployment: deployment,
	}
}

// OnDeploymentStart is called when a deployment starts
func (h *DeploymentHook) OnDeploymentStart(ctx context.Context) {
	h.startTime = time.Now()
	deploymentCounter.WithLabelValues(h.namespace, h.deployment, "started").Inc()
}

// OnDeploymentSuccess is called when a deployment succeeds
func (h *DeploymentHook) OnDeploymentSuccess(ctx context.Context) {
	duration := time.Since(h.startTime).Seconds()
	deploymentDuration.WithLabelValues(h.namespace, h.deployment).Observe(duration)
	deploymentCounter.WithLabelValues(h.namespace, h.deployment, "success").Inc()
}

// OnDeploymentFailure is called when a deployment fails
func (h *DeploymentHook) OnDeploymentFailure(ctx context.Context, reason string) {
	duration := time.Since(h.startTime).Seconds()
	deploymentDuration.WithLabelValues(h.namespace, h.deployment).Observe(duration)
	deploymentCounter.WithLabelValues(h.namespace, h.deployment, "failed").Inc()
	deploymentFailures.WithLabelValues(h.namespace, h.deployment, reason).Inc()
}

// OnDeploymentRollback is called when a deployment is rolled back
func (h *DeploymentHook) OnDeploymentRollback(ctx context.Context) {
	deploymentCounter.WithLabelValues(h.namespace, h.deployment, "rolled_back").Inc()
}

