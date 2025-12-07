# ADR 0003: Scaleway als Cloud Provider

## Status
Geaccepteerd

## Context
Het platform heeft een cloud provider nodig voor:
- Kubernetes cluster hosting
- Container registry voor Docker images
- Load balancing en networking
- DNS management
- TLS certificaat management

Mogelijke providers: AWS, Azure, GCP, Scaleway, DigitalOcean, etc.

## Beslissing
We kiezen voor **Scaleway** als cloud provider met:
- **Kubernetes**: Scaleway Kubernetes Kapsule (managed Kubernetes)
- **Container Registry**: Scaleway Container Registry (SCR)
- **Load Balancer**: Scaleway Load Balancer via Ingress
- **DNS**: Externe DNS provider (voor nlappstore.nl)
- **TLS**: Cert-manager met Let's Encrypt

## Gevolgen

### Positief
- ✅ **EU Data Residency**: Scaleway heeft datacenters in EU (Nederland, Frankrijk)
- ✅ **Cost-Effective**: Lagere kosten dan AWS/Azure/GCP
- ✅ **Managed Services**: Kubernetes Kapsule en Container Registry zijn managed
- ✅ **Simplicity**: Minder complex dan multi-cloud setup
- ✅ **Compliance**: EU datacenters helpen met GDPR compliance
- ✅ **Performance**: Lage latency voor Nederlandse gebruikers

### Negatief
- ⚠️ **Vendor Lock-in**: Platform is afhankelijk van Scaleway
- ⚠️ **Ecosystem**: Kleinere ecosystem dan AWS/Azure/GCP
- ⚠️ **Documentation**: Minder uitgebreide documentatie dan grote providers
- ⚠️ **Multi-Region**: Beperkte multi-region opties

## Implementatie Details

### Kubernetes Cluster
- **Service**: Scaleway Kubernetes Kapsule
- **Region**: Amsterdam (nl-ams)
- **Cluster Name**: `k8s-ams-nl-appstore`
- **Namespace**: `nl-appstore-registry`
- **Kubeconfig**: `kubeconfig-k8s-ams-nl-appstore.yaml`

### Container Registry
- **Service**: Scaleway Container Registry
- **Registry URL**: `rg.nl-ams.scw.cloud`
- **Namespace**: `nl-appstore-registry`
- **Authentication**: Secret key via `SCW_SECRET_KEY`

### Networking
- **Ingress**: NGINX Ingress Controller
- **Load Balancer**: Scaleway Load Balancer (automatisch via Ingress)
- **Domain**: `nlappstore.nl`, `www.nlappstore.nl`, `api.nlappstore.nl`
- **TLS**: Cert-manager met Let's Encrypt ClusterIssuer

### Deployment
- Deployment script: `deploy/setup-from-env.sh`
- Docker builds: `--platform linux/amd64` voor Scaleway compatibiliteit
- Image tags: `rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest`

## Alternatieven Overwogen

### AWS (EKS)
- **Waarom niet**: Hogere kosten, US datacenters (GDPR concerns), meer complexiteit

### Azure (AKS)
- **Waarom niet**: Hogere kosten, Microsoft ecosystem lock-in

### GCP (GKE)
- **Waarom niet**: Hogere kosten, US datacenters (GDPR concerns)

### DigitalOcean
- **Waarom niet**: Minder managed services, kleinere Kubernetes ecosystem

### Multi-Cloud
- **Waarom niet**: Te complex voor huidige scope, operational overhead

## Migration Strategy

### Als we moeten migreren:
1. Kubernetes manifests zijn provider-agnostic (standaard Kubernetes)
2. Container images kunnen naar andere registry gepusht worden
3. Ingress configuratie kan aangepast worden voor andere providers
4. Application code is provider-agnostic

### Abstraction Layers:
- Kubernetes API abstraheert cloud provider details
- Container registry is standaard Docker registry
- Ingress is standaard Kubernetes Ingress

## Referenties
- [Scaleway Documentation](https://www.scaleway.com/en/docs/)
- [Scaleway Kubernetes Kapsule](https://www.scaleway.com/en/docs/containers/kubernetes/)
- [Effective Platform Engineering - Chapter 7: Platform Control Plane Foundations]

---

**Datum**: 2025-01-XX  
**Auteur**: Platform Engineering Team  
**Reviewers**: DevOps Team, Security Team

