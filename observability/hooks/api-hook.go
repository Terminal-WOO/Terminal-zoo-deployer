package hooks

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// API request metrics
	apiRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "platform_api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"service", "method", "endpoint", "status"},
	)

	apiRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "platform_api_request_duration_seconds",
			Help:    "API request duration in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 12), // 1ms to ~4s
		},
		[]string{"service", "method", "endpoint"},
	)

	apiErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "platform_api_errors_total",
			Help: "Total number of API errors",
		},
		[]string{"service", "method", "endpoint", "error_type"},
	)
)

// APIHook provides observability hooks for API requests
type APIHook struct {
	service string
	method  string
	endpoint string
	startTime time.Time
}

// NewAPIHook creates a new API hook
func NewAPIHook(service, method, endpoint string) *APIHook {
	return &APIHook{
		service: service,
		method: method,
		endpoint: endpoint,
		startTime: time.Now(),
	}
}

// OnRequestStart is called when an API request starts
func (h *APIHook) OnRequestStart() {
	h.startTime = time.Now()
}

// OnRequestSuccess is called when an API request succeeds
func (h *APIHook) OnRequestSuccess(statusCode int) {
	duration := time.Since(h.startTime).Seconds()
	apiRequestDuration.WithLabelValues(h.service, h.method, h.endpoint).Observe(duration)
	apiRequestsTotal.WithLabelValues(h.service, h.method, h.endpoint, "success").Inc()
}

// OnRequestError is called when an API request fails
func (h *APIHook) OnRequestError(statusCode int, errorType string) {
	duration := time.Since(h.startTime).Seconds()
	apiRequestDuration.WithLabelValues(h.service, h.method, h.endpoint).Observe(duration)
	apiRequestsTotal.WithLabelValues(h.service, h.method, h.endpoint, "error").Inc()
	apiErrorsTotal.WithLabelValues(h.service, h.method, h.endpoint, errorType).Inc()
}

