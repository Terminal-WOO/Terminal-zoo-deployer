# ADR 0005: Self-Service Deployment Model

## Status
Geaccepteerd

## Context
Het platform moet ontwikkelaars in staat stellen om AI-toepassingen te deployen. Er zijn verschillende deployment modellen mogelijk:
- Centralized deployment (alleen platform team deployt)
- Self-service deployment (ontwikkelaars deployen zelf)
- Hybrid model (mix van beide)

## Beslissing
We kiezen voor een **self-service deployment model** waarbij:
- Ontwikkelaars kunnen zelf applicaties deployen via de platform UI
- Platform biedt self-service capabilities zonder handmatige interventie
- Governance en compliance zijn ingebouwd (niet als gatekeeper)
- Platform abstraheert Kubernetes complexiteit

## Gevolgen

### Positief
- ✅ **Developer Autonomy**: Ontwikkelaars zijn niet afhankelijk van platform team
- ✅ **Faster Delivery**: Snellere time-to-market voor applicaties
- ✅ **Scalability**: Platform team kan schalen zonder per-deployment werk
- ✅ **Reduced Cognitive Load**: Platform abstraheert infrastructuur complexiteit
- ✅ **Compliance Built-In**: Policies zijn geautomatiseerd, niet handmatig

### Negatief
- ⚠️ **Platform Complexity**: Platform moet robuust genoeg zijn voor self-service
- ⚠️ **Error Handling**: Betere error messages en recovery nodig
- ⚠️ **Security Risks**: Meer attack surface door self-service
- ⚠️ **Resource Management**: Mogelijk resource waste zonder governance

## Implementatie Details

### Self-Service Workflow
1. **Application Selection**: Developer selecteert applicatie uit catalog
2. **Deployment Configuration**: Developer configureert deployment (namespace, resources, etc.)
3. **Validation**: Platform valideert configuratie tegen policies
4. **Deployment**: Platform deployt naar Kubernetes automatisch
5. **Monitoring**: Developer kan deployment status monitoren
6. **Management**: Developer kan deployment updaten, restarten, of verwijderen

### Platform Components

#### Frontend
- **Deployment UI**: `app/pages/deploy/new.vue`
- **Deployment Store**: `app/stores/deploy.ts`
- **Dashboard**: `app/pages/dashboard/index.vue` voor monitoring
- **Deployment Management**: `app/stores/deployments.ts`

#### Backend
- **Deployment API**: `go/internal/server/server.go` (deployment handlers)
- **Kubernetes Client**: `go/pkg/kube/client/` voor cluster interactie
- **Validation**: Configuratie validatie voor deployment specs
- **Error Handling**: Duidelijke error messages voor developers

#### API Endpoints
- `POST /deployments` - Create deployment
- `GET /deployments/{namespace}` - List deployments
- `GET /deployments/{namespace}/{name}` - Get deployment details
- `PUT /deployments/{namespace}/{name}` - Update deployment
- `DELETE /deployments/{namespace}/{name}` - Delete deployment
- `POST /deployments/{namespace}/{name}/restart` - Restart deployment

### Governance & Compliance

#### Built-In Policies
- Namespace validatie (bestaat namespace?)
- Resource limits (CPU, memory)
- Security policies (non-root, read-only filesystem)
- Network policies (ingress/egress rules)

#### Future Enhancements
- Policy-as-code met Open Policy Agent (OPA)
- Automated compliance checks
- Resource quota management
- Cost tracking per deployment

### Developer Experience

#### Abstraction Layers
- **Kubernetes Abstraction**: Developers hoeven geen Kubernetes kennis te hebben
- **Simple Configuration**: UI-based configuratie, geen YAML nodig
- **Clear Feedback**: Duidelijke status en error messages
- **Self-Service Monitoring**: Developers kunnen zelf deployment status checken

#### Documentation
- Deployment guides in platform
- Error troubleshooting guides
- Best practices documentatie

## Alternatieven Overwogen

### Centralized Deployment
- **Waarom niet**: Bottleneck voor platform team, langzame delivery, niet schaalbaar

### Hybrid Model
- **Waarom niet**: Inconsistentie, developers weten niet wanneer ze zelf kunnen deployen

### Fully Automated (No Developer Input)
- **Waarom niet**: Te weinig flexibiliteit, developers willen controle over configuratie

## Success Metrics

### Developer Experience
- **Time-to-first-deployment**: < 15 minuten
- **Self-service adoption rate**: > 80% van deployments
- **Developer satisfaction**: > 4.0/5.0

### Platform Performance
- **Deployment success rate**: > 95%
- **Average deployment time**: < 5 minuten
- **Error recovery time**: < 10 minuten

## Future Enhancements

1. **GitOps Integration**: Deployments via Git commits
2. **CI/CD Integration**: Automatische deployments vanuit CI/CD
3. **Advanced Policies**: Meer geavanceerde policy-as-code
4. **Cost Optimization**: Automatische resource optimization
5. **Multi-Cluster**: Deployments naar meerdere clusters

## Referenties
- [Effective Platform Engineering - Chapter 1: What is Platform Engineering]
- [Effective Platform Engineering - Chapter 4: Governance, Compliance, and Trust]
- [Team Topologies - Platform Team Model](https://teamtopologies.com/)

---

**Datum**: 2025-01-XX  
**Auteur**: Platform Engineering Team  
**Reviewers**: Developer Experience Team, Security Team

