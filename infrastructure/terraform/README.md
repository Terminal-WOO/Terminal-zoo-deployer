# Terraform Infrastructure

Dit directory bevat Terraform modules voor infrastructure as code.

## Overzicht

Terraform wordt gebruikt voor het beheren van cloud infrastructure (Scaleway resources).

## Structure (Planned)

```
infrastructure/terraform/
├── account/          # Account-level resources
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── cluster/          # Cluster-level resources
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── namespace/        # Namespace-level resources
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── modules/          # Reusable modules
│   ├── kubernetes-cluster/
│   ├── container-registry/
│   └── namespace/
└── environments/     # Environment-specific configs
    ├── production/
    └── staging/
```

## Current Status

**Status**: 🔄 Gepland

**Current Implementation**: 
- ✅ Kubernetes manifests (k8s/)
- ✅ Manual Scaleway setup
- 🔄 Terraform modules (to be implemented)

## Future Implementation

### Account Module
- Scaleway account configuration
- IAM roles and policies
- Billing setup

### Cluster Module
- Kubernetes cluster (Scaleway Kapsule)
- Cluster add-ons
- Cluster monitoring

### Namespace Module
- Namespace creation
- Resource quotas
- Service accounts

## Terraform Best Practices

1. **Modularity**: Use reusable modules
2. **Versioning**: Pin provider versions
3. **State Management**: Use remote state
4. **Testing**: Test modules before use
5. **Documentation**: Document all modules

## Referenties

- [Terraform Documentation](https://www.terraform.io/docs/)
- [Scaleway Provider](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs)

---

**Status**: Gepland  
**Eigenaar**: Platform Engineering Team

