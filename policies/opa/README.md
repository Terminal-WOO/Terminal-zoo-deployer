# Open Policy Agent (OPA) Policies

Dit directory bevat OPA (Open Policy Agent) policies voor het App Store platform.

## Overzicht

OPA policies worden gebruikt voor policy-as-code implementatie. Policies zijn geschreven in Rego, de policy language van OPA.

## Policy Bestanden

### deployment.rego
Validatie policies voor Kubernetes deployments:
- Namespace validatie
- Resource limits (CPU, memory)
- Security constraints (non-root, read-only filesystem)
- Required labels

**Usage**:
```bash
opa eval --data policies/opa/deployment.rego --input deployment.yaml 'data.platform.deployment.allow'
```

### security.rego
Security policies voor deployments:
- Image security (registry, tags)
- Secrets management
- Network policy compliance

**Usage**:
```bash
opa eval --data policies/opa/security.rego --input deployment.yaml 'data.platform.security.allow'
```

## Policy Testing

### Test Policies
```bash
# Test deployment policy
opa test policies/opa/

# Test with specific input
opa eval --data policies/opa/deployment.rego --input test-deployment.yaml 'data.platform.deployment.allow'
```

## Policy Integration

### CI/CD Integration
Policies worden gecontroleerd in CI/CD pipeline:
- Pre-deployment validation
- Policy compliance check
- Error reporting

### Kubernetes Integration
Policies kunnen geïntegreerd worden met:
- OPA Gatekeeper (Kubernetes admission controller)
- OPA sidecar voor runtime policy enforcement

## Policy Development

### Best Practices
1. **Test First**: Write tests before policies
2. **Document**: Document policy rationale
3. **Version**: Version policies with code
4. **Review**: Review policies like code

### Policy Structure
```rego
package platform.<category>

# Default deny
default allow = false

# Allow if checks pass
allow {
    validate_check1
    validate_check2
}

# Validation functions
validate_check1 { ... }
validate_check2 { ... }

# Error messages
errors[msg] {
    not validate_check1
    msg := "Error message"
}
```

## Referenties

- [OPA Documentation](https://www.openpolicyagent.org/docs/)
- [Rego Language](https://www.openpolicyagent.org/docs/latest/policy-language/)
- [OPA Gatekeeper](https://open-policy-agent.github.io/gatekeeper/)

---

**Status**: In ontwikkeling  
**Eigenaar**: Platform Engineering Team, Security Team

