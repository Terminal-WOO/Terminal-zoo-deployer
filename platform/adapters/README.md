# Platform Adapters

Dit directory bevat adapter implementations voor integraties met externe systemen.

## Overzicht

Adapters gebruiken het adapter pattern om platform events te integreren met externe systemen zoals CI/CD, observability, en issue tracking.

## Adapter Pattern

### Benefits

- **Loose Coupling**: Platform is niet direct gekoppeld aan externe systemen
- **Easy to Swap**: Implementaties kunnen eenvoudig gewisseld worden
- **Testability**: Adapters kunnen gemakkelijk getest worden
- **Extensibility**: Nieuwe adapters kunnen eenvoudig toegevoegd worden

---

## Adapter Types

### CI Hooks Adapter

**Purpose**: Integrate met CI/CD systems

**Adapters**:
- `ci/github-actions-adapter.go` - GitHub Actions integration
- `ci/gitlab-ci-adapter.go` - GitLab CI integration (planned)
- `ci/jenkins-adapter.go` - Jenkins integration (planned)

**Events**:
- `DeploymentRequested` → Trigger CI workflow
- `DeploymentCompleted` → Update CI status
- `DeploymentFailed` → Update CI status

**Location**: `platform/adapters/ci/`

---

### Observability Hooks Adapter

**Purpose**: Integrate met observability platforms

**Adapters**:
- `observability/prometheus-adapter.go` - Prometheus metrics
- `observability/grafana-adapter.go` - Grafana integration (planned)
- `observability/datadog-adapter.go` - Datadog integration (planned)

**Events**:
- `DeploymentCompleted` → Update metrics
- `DeploymentFailed` → Update metrics
- `SLOViolation` → Create alerts

**Location**: `platform/adapters/observability/`

---

### Issue Tracking Adapter

**Purpose**: Integrate met issue tracking systems

**Adapters**:
- `issues/github-issues-adapter.go` - GitHub Issues integration
- `issues/jira-adapter.go` - Jira integration (planned)
- `issues/linear-adapter.go` - Linear integration (planned)

**Events**:
- `DeploymentFailed` → Create issue
- `SLOViolation` → Create issue
- `ErrorBudgetDepleted` → Create issue

**Location**: `platform/adapters/issues/`

---

### Configuration Management Adapter

**Purpose**: Integrate met configuration management systems

**Adapters**:
- `config/vault-adapter.go` - HashiCorp Vault (planned)
- `config/aws-secrets-adapter.go` - AWS Secrets Manager (planned)
- `config/k8s-secrets-adapter.go` - Kubernetes Secrets (planned)

**Location**: `platform/adapters/config/`

---

## Adapter Interface

### Common Interface

```go
type Adapter interface {
    // Register registers the adapter
    Register(ctx context.Context) error
    
    // OnEvent handles platform events
    OnEvent(ctx context.Context, event *Event) error
}
```

---

## Usage

### Registering Adapters

```go
// Create adapters
githubAdapter := ci.NewGitHubActionsAdapter(webhookURL)
prometheusAdapter := observability.NewPrometheusAdapter()
issuesAdapter := issues.NewGitHubIssuesAdapter(owner, repo, token)

// Register with event bus
eventBus.Subscribe(ctx, EventTypeDeploymentRequested, githubAdapter.OnDeploymentRequested)
eventBus.Subscribe(ctx, EventTypeDeploymentCompleted, prometheusAdapter.OnDeploymentCompleted)
eventBus.Subscribe(ctx, EventTypeDeploymentFailed, issuesAdapter.OnDeploymentFailed)
```

---

## Referenties

- [Adapter Pattern](https://refactoring.guru/design-patterns/adapter)
- [Effective Platform Engineering - Chapter 9: Architecture Changes to Support Scale]

---

**Status**: In ontwikkeling  
**Eigenaar**: Platform Engineering Team

