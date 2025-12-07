package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	deploymentfrequency "github.com/ClappFormOrg/AI-CO/monitoring/dora-metrics/deployment-frequency"
	leadtime "github.com/ClappFormOrg/AI-CO/monitoring/dora-metrics/lead-time"
)

var (
	// Deployment frequency metrics
	deploymentFrequencyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "platform_deployments_total",
			Help: "Total number of deployments",
		},
		[]string{"namespace", "status"},
	)

	// Lead time metrics
	leadTimeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "platform_lead_time_seconds",
			Help:    "Lead time from commit to deployment in seconds",
			Buckets: prometheus.ExponentialBuckets(60, 2, 10), // 1 min to ~17 hours
		},
		[]string{"namespace", "deployment"},
	)

	// Change failure rate metrics
	deploymentFailuresCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "platform_deployment_failures_total",
			Help: "Total number of failed deployments",
		},
		[]string{"namespace"},
	)

	// MTTR metrics
	recoveryTimeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "platform_recovery_time_seconds",
			Help:    "Time to recover from failures in seconds",
			Buckets: prometheus.ExponentialBuckets(60, 2, 8), // 1 min to ~2 hours
		},
		[]string{"namespace", "incident"},
	)
)

func init() {
	prometheus.MustRegister(deploymentFrequencyGauge)
	prometheus.MustRegister(leadTimeHistogram)
	prometheus.MustRegister(deploymentFailuresCounter)
	prometheus.MustRegister(recoveryTimeHistogram)
}

func main() {
	var (
		kubeconfig = flag.String("kubeconfig", "", "Path to kubeconfig file")
		namespace  = flag.String("namespace", "nl-appstore-registry", "Kubernetes namespace")
		port       = flag.Int("port", 9090, "Port to serve metrics on")
	)
	flag.Parse()

	// Create Kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to build kubeconfig: %v", err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create clientset: %v", err))
	}

	// Create collectors
	deploymentFreqCollector := deploymentfrequency.NewCollector(clientset, *namespace)
	leadTimeCollector := leadtime.NewCollector("", "")

	// Start metrics collection in background
	go collectMetrics(context.Background(), deploymentFreqCollector, leadTimeCollector, *namespace)

	// Serve metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	fmt.Printf("Starting Prometheus metrics exporter on port %d\n", *port)
	fmt.Printf("Metrics available at http://localhost:%d/metrics\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func collectMetrics(ctx context.Context, deploymentFreqCollector *deploymentfrequency.DeploymentFrequencyCollector, leadTimeCollector *leadtime.LeadTimeCollector, namespace string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Collect deployment frequency (daily)
			if count, err := deploymentFreqCollector.CollectDaily(ctx); err == nil {
				deploymentFrequencyGauge.WithLabelValues(namespace, "success").Set(float64(count.Count))
			}

			// Note: Lead time collection requires Git integration
			// This is a placeholder - actual implementation would need Git API access
			// leadTimeHistogram.WithLabelValues(namespace, "deployment-name").Observe(leadTimeSeconds)
		}
	}
}

