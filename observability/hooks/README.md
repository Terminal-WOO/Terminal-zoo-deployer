# Observability Hooks

Dit directory bevat observability hooks voor platform services.

## Overzicht

Observability hooks worden geïntegreerd in platform services om automatisch metrics, logs en traces te verzamelen.

## Hook Types

### Deployment Hooks

**Purpose**: Track deployment events en metrics

**File**: `deployment-hook.go`

**Metrics**:
- `platform_deployments_total` - Total deployments
- `platform_deployment_duration_seconds` - Deployment duration
- `platform_deployment_failures_total` - Failed deployments

**Usage**:
```go
hook := hooks.NewDeploymentHook(namespace, deployment)
hook.OnDeploymentStart(ctx)
// ... deployment logic ...
hook.OnDeploymentSuccess(ctx)
```

---

### API Hooks

**Purpose**: Track API requests en performance

**File**: `api-hook.go`

**Metrics**:
- `platform_api_requests_total` - Total API requests
- `platform_api_request_duration_seconds` - Request duration
- `platform_api_errors_total` - API errors

**Usage**:
```go
hook := hooks.NewAPIHook("go-backend", "GET", "/deployments")
hook.OnRequestStart()
// ... handle request ...
hook.OnRequestSuccess(200)
```

---

## Integration

### Backend Integration

**Location**: `go/internal/server/server.go`

**Example**:
```go
func (h *Handler) handleDeploymentCreation() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        hook := hooks.NewDeploymentHook(namespace, deployment)
        hook.OnDeploymentStart(r.Context())
        
        // ... deployment logic ...
        
        if err != nil {
            hook.OnDeploymentFailure(r.Context(), err.Error())
            return
        }
        
        hook.OnDeploymentSuccess(r.Context())
    }
}
```

### Frontend Integration

**Location**: `server/api/` (Nuxt server API)

**Example**:
```typescript
export default defineEventHandler(async (event) => {
  const hook = new APIHook('nuxt-frontend', event.method, event.path)
  hook.onRequestStart()
  
  try {
    // ... handle request ...
    hook.onRequestSuccess(200)
  } catch (error) {
    hook.onRequestError(500, error.type)
  }
})
```

---

## Metrics Export

Alle hooks exporteren metrics in Prometheus format via `/metrics` endpoint.

**Metrics Endpoint**: `http://localhost:8080/metrics` (backend)  
**Metrics Endpoint**: `http://localhost:3000/metrics` (frontend, planned)

---

## Referenties

- [Prometheus Client Go](https://github.com/prometheus/client_golang)
- [Effective Platform Engineering - Chapter 5: Evolutionary Observability]

---

**Status**: In ontwikkeling  
**Eigenaar**: Platform Engineering Team

