package observability

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusAdapter adapts platform events to Prometheus metrics
type PrometheusAdapter struct {
	deploymentCounter *prometheus.CounterVec
	deploymentDuration *prometheus.HistogramVec
}

// NewPrometheusAdapter creates a new Prometheus adapter
func NewPrometheusAdapter() *PrometheusAdapter {
	return &PrometheusAdapter{
		deploymentCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "platform_deployments_total",
				Help: "Total number of deployments",
			},
			[]string{"namespace", "deployment", "status"},
		),
		deploymentDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "platform_deployment_duration_seconds",
				Help:    "Deployment duration in seconds",
				Buckets: prometheus.ExponentialBuckets(1, 2, 10),
			},
			[]string{"namespace", "deployment"},
		),
	}
}

// Register registers Prometheus metrics
func (a *PrometheusAdapter) Register() error {
	if err := prometheus.Register(a.deploymentCounter); err != nil {
		return fmt.Errorf("failed to register deployment counter: %w", err)
	}
	if err := prometheus.Register(a.deploymentDuration); err != nil {
		return fmt.Errorf("failed to register deployment duration: %w", err)
	}
	return nil
}

// OnDeploymentCompleted is called when a deployment completes
func (a *PrometheusAdapter) OnDeploymentCompleted(ctx context.Context, data map[string]interface{}) error {
	namespace := fmt.Sprintf("%v", data["namespace"])
	deployment := fmt.Sprintf("%v", data["deployment"])
	status := fmt.Sprintf("%v", data["status"])
	
	a.deploymentCounter.WithLabelValues(namespace, deployment, status).Inc()
	
	if duration, ok := data["duration_seconds"].(int64); ok {
		a.deploymentDuration.WithLabelValues(namespace, deployment).Observe(float64(duration))
	}
	
	return nil
}

// OnDeploymentFailed is called when a deployment fails
func (a *PrometheusAdapter) OnDeploymentFailed(ctx context.Context, data map[string]interface{}) error {
	namespace := fmt.Sprintf("%v", data["namespace"])
	deployment := fmt.Sprintf("%v", data["deployment"])
	
	a.deploymentCounter.WithLabelValues(namespace, deployment, "failed").Inc()
	
	return nil
}

