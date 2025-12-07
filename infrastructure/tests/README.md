# Infrastructure Tests

Dit directory bevat tests voor infrastructure code.

## Test Types

### 1. Manifest Validation Tests

**Purpose**: Validate Kubernetes manifests

**File**: `validate-manifests.sh`

**Usage**:
```bash
./infrastructure/tests/validate-manifests.sh
```

**What it tests**:
- kubectl dry-run validation
- kubeval validation (if available)
- Network policies validation

---

### 2. Terraform Tests (Planned)

**Purpose**: Test Terraform modules

**Status**: 🔄 Gepland

**Tools**:
- `terraform validate`
- `terraform plan`
- `terratest` (Go testing framework)

---

### 3. Integration Tests (Planned)

**Purpose**: Test infrastructure components together

**Status**: 🔄 Gepland

**Tools**:
- Kubernetes test framework
- End-to-end testing tools

---

## Running Tests

### Local Testing

```bash
# Validate manifests
./infrastructure/tests/validate-manifests.sh

# Run all tests (when implemented)
make infra-test
```

### CI/CD Testing

Tests run automatically in:
- `.github/workflows/infrastructure-ci.yml` - On PR and push
- `.github/workflows/platform-ci.yml` - As part of platform CI

---

## Test Best Practices

1. **Test Before Deploy**: Always run tests before deploying
2. **Test in CI/CD**: Automate tests in CI/CD pipeline
3. **Test Locally**: Run tests locally before committing
4. **Test Incrementally**: Test changes incrementally

---

## Referenties

- [Kubernetes Testing](https://kubernetes.io/docs/concepts/cluster-administration/testing/)
- [Terraform Testing](https://www.terraform.io/docs/cloud/guides/testing/)

---

**Status**: In ontwikkeling  
**Eigenaar**: Platform Engineering Team

