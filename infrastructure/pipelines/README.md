# Infrastructure Pipelines

Dit directory bevat reusable pipeline templates en workflows voor infrastructure deployment.

## Pipeline Types

### Namespace-Level Pipeline

**Purpose**: Deploy applications to Kubernetes namespaces

**File**: `namespace-deployment.yml`

**Usage**:
```yaml
jobs:
  deploy:
    uses: ./.github/workflows/infrastructure/pipelines/namespace-deployment.yml
    with:
      namespace: nl-appstore-registry
      environment: production
      image_tag: main-abc1234
```

**Steps**:
1. Validate Kubernetes manifests
2. Configure kubectl
3. Apply namespace
4. Apply ConfigMaps
5. Update image tags
6. Apply deployments
7. Wait for rollout
8. Verify deployment

---

### Cluster-Level Pipeline (Planned)

**Purpose**: Manage Kubernetes cluster configuration

**Status**: 🔄 Gepland

**Components**:
- Cluster updates
- Add-on management
- Cluster monitoring

---

### Account-Level Pipeline (Planned)

**Purpose**: Manage cloud account baseline

**Status**: 🔄 Gepland

**Components**:
- Account security scanning
- Cost optimization
- Resource cleanup

---

## Pipeline Best Practices

### 1. Validation First
Always validate infrastructure code before deploying:
- Syntax validation
- Policy validation
- Dry-run validation

### 2. Idempotency
Pipelines should be idempotent:
- Can be run multiple times safely
- No side effects on re-run
- State convergence

### 3. Rollback Support
Pipelines should support rollback:
- Version tracking
- Rollback procedures
- Health checks

### 4. Observability
Pipelines should be observable:
- Logging
- Metrics
- Alerts

---

## Pipeline Templates

### Reusable Workflows

**Location**: `.github/workflows/reusable/`

**Templates**:
- `deploy-k8s.yml` - Kubernetes deployment
- `validate-manifests.yml` - Manifest validation
- `terraform-plan.yml` - Terraform planning
- `terraform-apply.yml` - Terraform apply

---

## Testing Pipelines

### Test Infrastructure

**Location**: `infrastructure/tests/`

**Test Types**:
- Unit tests voor modules
- Integration tests voor pipelines
- End-to-end tests voor deployments

---

## Referenties

- [GitHub Actions Reusable Workflows](https://docs.github.com/en/actions/using-workflows/reusing-workflows)
- [Effective Platform Engineering - Chapter 6: Building Software-Defined Platforms]

---

**Status**: In ontwikkeling  
**Eigenaar**: Platform Engineering Team

