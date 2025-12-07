package events

import (
	"time"
)

// EventType represents the type of event
type EventType string

const (
	// Deployment Events
	EventTypeDeploymentRequested EventType = "DeploymentRequested"
	EventTypeDeploymentStarted   EventType = "DeploymentStarted"
	EventTypeDeploymentCompleted EventType = "DeploymentCompleted"
	EventTypeDeploymentFailed    EventType = "DeploymentFailed"
	EventTypeDeploymentRolledBack EventType = "DeploymentRolledBack"

	// Infrastructure Events
	EventTypeClusterCreated      EventType = "ClusterCreated"
	EventTypeClusterUpdated      EventType = "ClusterUpdated"
	EventTypeNamespaceCreated    EventType = "NamespaceCreated"
	EventTypeResourceQuotaExceeded EventType = "ResourceQuotaExceeded"

	// Platform Events
	EventTypePlatformHealthCheck EventType = "PlatformHealthCheck"
	EventTypeSLOViolation       EventType = "SLOViolation"
	EventTypeErrorBudgetDepleted EventType = "ErrorBudgetDepleted"
)

// Event represents a platform event
type Event struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	Timestamp   time.Time              `json:"timestamp"`
	Source      string                 `json:"source"`
	Data        map[string]interface{} `json:"data"`
	Metadata    map[string]string     `json:"metadata"`
}

// DeploymentEventData contains data for deployment events
type DeploymentEventData struct {
	Namespace     string `json:"namespace"`
	Deployment    string `json:"deployment"`
	Image         string `json:"image"`
	Replicas      int32  `json:"replicas"`
	Status        string `json:"status"`
	Error         string `json:"error,omitempty"`
	Duration      int64  `json:"duration_seconds,omitempty"`
}

// ClusterEventData contains data for cluster events
type ClusterEventData struct {
	ClusterName   string `json:"cluster_name"`
	Region        string `json:"region"`
	Action        string `json:"action"`
	Status        string `json:"status"`
}

// PlatformEventData contains data for platform events
type PlatformEventData struct {
	Service       string  `json:"service"`
	SLO           string  `json:"slo,omitempty"`
	CurrentValue  float64 `json:"current_value,omitempty"`
	TargetValue   float64 `json:"target_value,omitempty"`
	ErrorBudget   float64 `json:"error_budget,omitempty"`
}

// NewEvent creates a new event
func NewEvent(eventType EventType, source string, data map[string]interface{}) *Event {
	return &Event{
		ID:        generateEventID(),
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    source,
		Data:      data,
		Metadata:  make(map[string]string),
	}
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string (simplified implementation)
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

