# Platform Product Domains

Dit document definieert de platform product domains zoals beschreven in "Effective Platform Engineering". Deze domains helpen bij het organiseren van platform capabilities en het identificeren van ownership en verantwoordelijkheden.

## Overzicht

Het App Store platform is georganiseerd in drie hoofddomains:

1. **Developer Tools** - Tools en interfaces voor ontwikkelaars
2. **Infrastructure & Operations** - Platform infrastructuur en operaties
3. **Governance** - Policies, compliance en security

## Domain 1: Developer Tools

### Doel
Biedt ontwikkelaars de tools en interfaces die ze nodig hebben om AI-toepassingen te ontdekken, evalueren, deployen en beheren.

### Capabilities

#### Application Catalog
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: Catalogus van AI-toepassingen met metadata, criteria en reviews
- **Componenten**:
  - Frontend: `app/pages/apps/`, `app/stores/appCatalog.ts`
  - Backend: `server/api/apps/`
  - Data: `app/assets/userApps.json`, `server/assets/userApps.json`

#### Deployment Management
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: Self-service deployment van applicaties naar Kubernetes
- **Componenten**:
  - Frontend: `app/pages/deploy/new.vue`, `app/stores/deploy.ts`
  - Backend: `go/internal/server/server.go` (deployment handlers)
  - API: `server/api/deploy/`

#### Developer Dashboard
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: Dashboard voor ontwikkelaars om deployments te monitoren
- **Componenten**:
  - Frontend: `app/pages/dashboard/index.vue`
  - Stores: `app/stores/deployments.ts`

#### Community Features
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: Community forums en threads voor kennis delen
- **Componenten**:
  - Frontend: `app/pages/community/`, `app/stores/community.ts`

### Toekomstige Capabilities
- Developer Portal integratie
- SDK voor programmatische toegang
- CLI tools voor deployment
- IDE plugins

---

## Domain 2: Infrastructure & Operations

### Doel
Biedt de onderliggende infrastructuur en operationele capabilities die het platform ondersteunen.

### Capabilities

#### Kubernetes Control Plane
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: Kubernetes-native deployment en management
- **Componenten**:
  - Backend: `go/pkg/kube/client/` (Kubernetes client)
  - Kubernetes manifests: `k8s/`
  - Cluster management: `server/api/clusters/`

#### Container Registry Integration
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: Integratie met Scaleway Container Registry
- **Componenten**:
  - Docker builds: `Dockerfile`, `go/Dockerfile.multi-stage`
  - Deployment script: `deploy/setup-from-env.sh`
  - Registry secrets: `k8s/secrets.yaml.template`

#### Deployment Automation
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: Geautomatiseerde deployment workflows
- **Componenten**:
  - Setup script: `deploy/setup-from-env.sh`
  - Kubernetes deployments: `k8s/frontend-deployment.yaml`, `k8s/backend-deployment.yaml`
  - Services: `k8s/frontend-service-*.yaml`, `k8s/backend-service-*.yaml`

#### Ingress & TLS Management
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: Ingress routing en TLS certificaat management
- **Componenten**:
  - Ingress: `k8s/ingress.yaml`
  - Cert-manager: `k8s/cluster-issuer.yaml`
  - Domain: `nlappstore.nl`, `api.nlappstore.nl`

#### Observability (Basis)
- **Status**: 🔄 In ontwikkeling
- **Beschrijving**: Monitoring en logging capabilities
- **Componenten**:
  - Health checks: `server/api/health.get.ts`, `server/api/ready.get.ts`
  - Backend logging: `go/pkg/log/`

### Toekomstige Capabilities
- Geavanceerde observability platform (Grafana/LGTM stack)
- SLOs als code
- Distributed tracing
- Cost monitoring en optimization
- Multi-cluster federation
- Service mesh (Istio/Linkerd)

---

## Domain 3: Governance

### Doel
Zorgt voor compliance, security en governance zonder developer friction te veroorzaken.

### Capabilities

#### Application Criteria & Review
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: Transparante criteria voor AI-toepassingen
- **Componenten**:
  - Frontend: `app/pages/criteria/index.vue`, `app/components/ToelatingsChecklist.vue`
  - Stores: `app/stores/criteria.ts`, `app/stores/submissions.ts`
  - Review queue: `app/pages/review/queue.vue`

#### Authentication & Authorization
- **Status**: ✅ Geïmplementeerd
- **Beschrijving**: User authenticatie en autorisatie
- **Componenten**:
  - Frontend: `app/pages/auth/`, `app/middleware/auth.ts`
  - Backend: `server/api/auth/`
  - Stores: `app/stores/user.ts`

#### Security Policies
- **Status**: 🔄 Basis geïmplementeerd
- **Beschrijving**: Security headers en policies
- **Componenten**:
  - CSP headers: `nuxt.config.ts` (Content-Security-Policy)
  - CORS: `go/cmd/server/cors.go`
  - Security store: `app/stores/security.ts`

#### Audit & Compliance
- **Status**: ✅ Basis geïmplementeerd
- **Beschrijving**: Audit logging en compliance tracking
- **Componenten**:
  - Audit store: `app/stores/audits.ts`
  - Backend logging: `go/pkg/log/`

### Toekomstige Capabilities
- Policy-as-code met Open Policy Agent (OPA)
- Software supply chain security (SBOM, Cosign)
- Zero-trust networking policies
- Automated compliance checks
- Risk assessment automation

---

## Domain Boundaries

### Developer Tools ↔ Infrastructure & Operations
- **Interface**: REST API tussen frontend en backend
- **Contract**: API endpoints gedefinieerd in `server/api/`
- **Ownership**: Developer Tools team gebruikt Infrastructure & Operations capabilities via API

### Developer Tools ↔ Governance
- **Interface**: Criteria en review workflows
- **Contract**: Criteria definitie in `app/stores/criteria.ts`
- **Ownership**: Governance team definieert criteria, Developer Tools team implementeert UI

### Infrastructure & Operations ↔ Governance
- **Interface**: Security policies en compliance checks
- **Contract**: Security policies in Kubernetes manifests en code
- **Ownership**: Governance team definieert policies, Infrastructure team implementeert

## Domain Ownership

### Developer Tools Domain
- **Owner**: Frontend Team
- **Stakeholders**: Developers, Product Managers
- **Success Metrics**: Developer satisfaction, Self-service adoption

### Infrastructure & Operations Domain
- **Owner**: Platform Engineering Team
- **Stakeholders**: DevOps Engineers, SREs
- **Success Metrics**: Platform availability, Deployment success rate

### Governance Domain
- **Owner**: Security & Compliance Team
- **Stakeholders**: Security Engineers, Compliance Officers
- **Success Metrics**: Policy compliance rate, Security incident reduction

## Evolution Strategy

### Incremental Design
- Domains evolueren op basis van user feedback
- Backlog management per domain
- Cross-domain coördinatie via platform roadmap

### Architectural Fitness Functions
- Per domain worden fitness functions gedefinieerd
- Automatische validatie van domain boundaries
- Continuous monitoring van domain health

---

**Laatste update**: 2025-01-XX  
**Status**: Actief  
**Eigenaar**: Platform Engineering Team

