# Event-Driven Automation

Dit directory bevat de event-driven automation implementatie voor het platform.

## Overzicht

Event-driven automation maakt het platform schaalbaar door asynchrone, event-driven operaties te gebruiken in plaats van synchrone API calls.

## Components

### Event Types

**File**: `types.go`

**Event Types**:
- Deployment events (Requested, Started, Completed, Failed, RolledBack)
- Infrastructure events (ClusterCreated, NamespaceCreated, etc.)
- Platform events (HealthCheck, SLOViolation, ErrorBudgetDepleted)

---

### Event Bus

**File**: `bus.go`

**Implementation**: In-memory event bus (can be extended to message queue)

**Features**:
- Event publishing
- Event subscription
- Event storage (with limit)
- Async event handling

**Future**: Message queue (RabbitMQ, NATS, Kafka) voor production

---

### Event Handlers

**Location**: `handlers/`

**Handlers**:
- `deployment-handler.go` - Handles deployment events
- `infrastructure-handler.go` - Handles infrastructure events (planned)
- `platform-handler.go` - Handles platform events (planned)

---

## Usage

### Publishing Events

```go
eventBus := events.NewInMemoryEventBus(1000)

event := events.NewEvent(
    events.EventTypeDeploymentRequested,
    "deployment-service",
    map[string]interface{}{
        "namespace": "nl-appstore-registry",
        "deployment": "my-app",
        "image": "my-image:latest",
    },
)

eventBus.Publish(ctx, event)
```

### Subscribing to Events

```go
handler := func(ctx context.Context, event *events.Event) error {
    // Handle event
    return nil
}

eventBus.Subscribe(ctx, events.EventTypeDeploymentCompleted, handler)
```

---

## Benefits

### Scalability
- Asynchronous operations
- Non-blocking
- Better throughput

### Resilience
- Decoupled components
- Fault isolation
- Retry logic

### Extensibility
- Easy to add new handlers
- Plugin architecture
- Adapter pattern support

---

## Referenties

- [Event-Driven Architecture](https://martinfowler.com/articles/201701-event-driven.html)
- [Effective Platform Engineering - Chapter 9: Architecture Changes to Support Scale]

---

**Status**: In ontwikkeling  
**Eigenaar**: Platform Engineering Team

