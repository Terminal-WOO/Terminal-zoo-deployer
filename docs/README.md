# Platform Engineering Documentatie

Dit directory bevat alle documentatie voor het App Store platform, georganiseerd volgens de principes uit "Effective Platform Engineering".

## 📋 Navigatie

### 🗺️ Start Hier

- **[Platform Engineering Overview](PLATFORM_ENGINEERING_OVERVIEW.md)** ⭐ - Hoofdnavigatie met alle links
- **[Executive Summary](PLATFORM_ENGINEERING_SUMMARY.md)** - Snelle samenvatting voor managers
- **[Platform Engineering Modules](../deploy/PLATFORM_ENGINEERING_MODULES.md)** - Volledige module beschrijvingen

---

## 📚 Documentatie per Categorie

### Foundation (Fase 1)

**Platform Fundamentals**:
- [Platform Visie](platform-vision.md) - Visie, missie, kernwaarden en strategische doelen
- [Platform Product Domains](platform-domains.md) - Developer Tools, Infrastructure & Operations, Governance
- [Platform Roadmap](platform-roadmap.md) - Strategische roadmap Q1-Q4 2025

**Architecture**:
- [Platform Lifecycle](platform-lifecycle.md) - SDLC (Design → Code → Build → Release → Operate)
- [Architectural Fitness Functions](fitness-functions.md) - 9 fitness functions voor architectuur validatie
- [Architecture Decisions](architecture-decisions/README.md) - ADR proces en 5 beslissingen

**Metrics**:
- [Metrics Framework](metrics-framework.md) - DORA metrics, platform value metrics, cognitive load

---

### Building (Fase 2)

**Governance**:
- [Governance Framework](governance-framework.md) - 4 governance lagen en workflows
- [Developer Autonomy](developer-autonomy.md) - Autonomy framework met metrics

**Observability**:
- [Observability Guide](observability-guide.md) - Observability principes, SLOs, LGTM stack

**Infrastructure**:
- [Infrastructure Guide](infrastructure-guide.md) - Infrastructure as code, pipelines, GitOps

---

### Scaling (Fase 3)

**Scaling**:
- [Scaling Architecture](scaling-architecture.md) - Event-driven automation, federated control planes, adapters

**Product Evolution**:
- [Product Evolution](product-evolution.md) - Platform als product strategie
- [Developer Experience](developer-experience.md) - DevEx principes en improvements
- [Culture Principles](culture-principles.md) - DevOps en Team Topologies

---

## 🏗️ Architecture Decisions

Alle belangrijke architectuur beslissingen zijn gedocumenteerd in Architecture Decision Records (ADRs):

- [ADR 0001](architecture-decisions/0001-kubernetes-native-architecture.md) - Kubernetes-Native Architecture
- [ADR 0002](architecture-decisions/0002-frontend-backend-split.md) - Frontend-Backend Split
- [ADR 0003](architecture-decisions/0003-scaleway-cloud-provider.md) - Scaleway als Cloud Provider
- [ADR 0004](architecture-decisions/0004-multi-stage-docker-builds.md) - Multi-Stage Docker Builds
- [ADR 0005](architecture-decisions/0005-self-service-deployment-model.md) - Self-Service Deployment Model

Zie [ADR Index](architecture-decisions/README.md) voor volledige lijst en proces.

---

## 📊 Code Directories

### Monitoring & Observability
- `monitoring/dora-metrics/` - DORA metrics collection
- `monitoring/platform-value/` - Platform value metrics
- `observability/slos/` - SLOs als code
- `observability/dashboards/` - Grafana dashboards
- `observability/hooks/` - Observability hooks

### Governance & Policies
- `policies/opa/` - Open Policy Agent policies
- `policies/supply-chain/` - Supply chain security
- `k8s/network-policies/` - Zero-trust network policies

### Infrastructure
- `infrastructure/pipelines/` - Infrastructure pipelines
- `infrastructure/tests/` - Infrastructure tests
- `infrastructure/terraform/` - Terraform modules

### Platform Services
- `platform/events/` - Event-driven automation
- `platform/federation/` - Federated control planes
- `platform/adapters/` - Adapter implementations
- `platform/intelligent-tools/` - AI/automation tooling
- `platform/developer-portal/` - Developer portal

---

## 🎯 Quick Links

### Voor Developers
- [Quick Start Guide](../deploy/QUICK_START.md) - Snel aan de slag
- [Platform Visie](platform-vision.md) - Begrijp het platform
- [Developer Experience](developer-experience.md) - DevEx improvements

### Voor Platform Engineers
- [Platform Engineering Modules](../deploy/PLATFORM_ENGINEERING_MODULES.md) - Volledige modules
- [Architecture Decisions](architecture-decisions/README.md) - ADRs
- [Platform Roadmap](platform-roadmap.md) - Wat bouwen we?

### Voor Managers
- [Executive Summary](PLATFORM_ENGINEERING_SUMMARY.md) - Snelle samenvatting
- [Platform Roadmap](platform-roadmap.md) - Strategische planning
- [Success Metrics](PLATFORM_ENGINEERING_OVERVIEW.md#success-metrics-overzicht) - Metrics overzicht

---

## 📈 Module Status

**✅ Voltooid**: 8 van 10 modules (80%)
- Fase 1: Foundation (3/3) ✅
- Fase 2: Building (3/5) ✅
- Fase 3: Scaling (2/2) ✅

**🔄 In Ontwikkeling**: 2 modules
- Module 2.4: Platform Control Plane Foundations
- Module 2.5: Control Plane Services & Extensions

Zie [Platform Engineering Overview](PLATFORM_ENGINEERING_OVERVIEW.md) voor volledige status.

---

## 🔗 Externe Referenties

- **Boek**: [Effective Platform Engineering](https://www.manning.com/books/effective-platform-engineering) (Manning Publications)
- **GitHub Companion**: [github.com/effective-platform-engineering/companion-code](https://github.com/effective-platform-engineering/companion-code)
- **Manning Book Page**: [manning.com/books/effective-platform-engineering](https://manning.com/books/effective-platform-engineering)

---

**Laatste update**: 2025-01-XX  
**Status**: Actief  
**Eigenaar**: Platform Engineering Team

