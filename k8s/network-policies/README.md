# Zero-Trust Network Policies

Dit directory bevat Kubernetes Network Policies voor zero-trust networking.

## Overzicht

Zero-trust networking betekent: "Never trust, always verify". Alle network traffic wordt gecontroleerd op basis van policies, niet op basis van network location.

## Network Policies

### default-deny-all.yaml
Default deny-all policy voor alle pods in de namespace. Specifieke policies zullen toegestane traffic definiëren.

**Principle**: Deny by default, allow explicitly.

### frontend-policy.yaml
Network policy voor de frontend (Nuxt.js) pods:
- **Ingress**: Toegang van ingress controller en backend
- **Egress**: DNS, HTTPS naar externe APIs, HTTP naar backend

### backend-policy.yaml
Network policy voor de backend (Go) pods:
- **Ingress**: Toegang van frontend en ingress controller
- **Egress**: DNS, HTTPS naar externe services, Kubernetes API

## Zero-Trust Principles

### 1. Never Trust, Always Verify
- Alle network traffic wordt gecontroleerd
- Geen impliciete trust op basis van network location
- Verificatie op basis van pod labels en selectors

### 2. Least Privilege
- Pods krijgen alleen toegang tot wat ze nodig hebben
- Geen brede network access
- Specifieke port en protocol definities

### 3. Micro-Segmentation
- Services zijn gesegmenteerd per applicatie
- Frontend en backend hebben verschillende policies
- Isolatie tussen services

### 4. Continuous Monitoring
- Network policies worden gemonitord
- Violations worden gedetecteerd
- Policies worden geëvalueerd en verbeterd

## Deployment

### Apply Network Policies
```bash
# Apply default deny-all first
kubectl apply -f k8s/network-policies/default-deny-all.yaml

# Apply specific policies
kubectl apply -f k8s/network-policies/frontend-policy.yaml
kubectl apply -f k8s/network-policies/backend-policy.yaml
```

### Verify Policies
```bash
# List network policies
kubectl get networkpolicies -n nl-appstore-registry

# Describe specific policy
kubectl describe networkpolicy frontend-network-policy -n nl-appstore-registry
```

## Policy Testing

### Test Network Connectivity
```bash
# Test frontend to backend connectivity
kubectl exec -it <frontend-pod> -n nl-appstore-registry -- curl http://go-backend:8080/health

# Test backend to Kubernetes API
kubectl exec -it <backend-pod> -n nl-appstore-registry -- curl https://kubernetes.default.svc
```

## Best Practices

1. **Start with Deny-All**: Begin met default-deny-all policy
2. **Add Policies Incrementally**: Voeg policies toe per service
3. **Test Thoroughly**: Test network connectivity na policy changes
4. **Document Policies**: Document waarom policies nodig zijn
5. **Review Regularly**: Review policies regelmatig voor optimalisatie

## Troubleshooting

### Common Issues

**Issue**: Pods kunnen niet communiceren
- **Solution**: Check network policies, verify pod labels match selectors

**Issue**: DNS resolution fails
- **Solution**: Ensure DNS egress policy allows UDP port 53 to kube-system

**Issue**: External API calls fail
- **Solution**: Ensure HTTPS egress policy allows TCP port 443

## Referenties

- [Kubernetes Network Policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [Zero Trust Architecture](https://www.nist.gov/publications/zero-trust-architecture)
- [Effective Platform Engineering - Chapter 4: Governance, Compliance, and Trust]

---

**Status**: In ontwikkeling  
**Eigenaar**: Platform Engineering Team, Security Team

