# ADR 0002: Frontend-Backend Split Architecture

## Status
Geaccepteerd

## Context
Het platform heeft zowel een user interface (voor ontwikkelaars) als backend APIs (voor deployment management) nodig. Er zijn verschillende architecturen mogelijk:
- Monolithic applicatie (alles in één codebase)
- Frontend-backend split (gescheiden services)
- Microservices architecture (veel kleine services)
- Server-side rendering (SSR) vs Client-side rendering (CSR)

## Beslissing
We kiezen voor een **frontend-backend split** waarbij:
- **Frontend**: Nuxt.js applicatie (Vue 3, SSR-capable)
- **Backend**: Go HTTP server met Kubernetes client
- Communicatie via REST API
- Frontend en backend kunnen onafhankelijk deployed worden

## Gevolgen

### Positief
- ✅ **Separation of Concerns**: Frontend en backend teams kunnen onafhankelijk werken
- ✅ **Technology Choice**: Beste tool voor elke taak (Vue voor UI, Go voor backend)
- ✅ **Scalability**: Frontend en backend kunnen onafhankelijk geschaald worden
- ✅ **Deployment Flexibility**: Frontend en backend kunnen apart geüpdatet worden
- ✅ **Developer Experience**: Frontend developers werken met Vue/Nuxt, backend met Go
- ✅ **SSR Support**: Nuxt.js biedt SSR voor betere SEO en performance

### Negatief
- ⚠️ **Complexity**: Twee codebases om te beheren
- ⚠️ **API Contract**: API contract tussen frontend en backend moet beheerd worden
- ⚠️ **Deployment Coordination**: Frontend en backend moeten compatibel blijven

## Implementatie Details

### Frontend Stack
- **Framework**: Nuxt.js 3 (Vue 3)
- **UI Library**: PrimeVue
- **State Management**: Pinia
- **Styling**: Tailwind CSS, PrimeFlex
- **Location**: `app/` directory
- **Port**: 3000

### Backend Stack
- **Language**: Go 1.25+
- **HTTP Server**: Standard library `net/http`
- **Kubernetes Client**: `k8s.io/client-go`
- **Location**: `go/` directory
- **Port**: 8080

### API Communication
- **Protocol**: REST over HTTP/HTTPS
- **Authentication**: Bearer token via `Authorization` header
- **API Base**: Configurable via `NUXT_EXTERNAL_API_BASE`
- **API Auth**: Configurable via `NUXT_EXTERNAL_API_AUTH`

### Deployment
- Frontend en backend zijn gescheiden Docker images
- Frontend image: `rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest`
- Backend image: `rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest` (same name, different Dockerfile)
- Ingress routes frontend naar `/` en backend naar `/api`

## Alternatieven Overwogen

### Monolithic Application
- **Waarom niet**: Minder flexibiliteit, moeilijker te schalen, technology lock-in

### Microservices
- **Waarom niet**: Te complex voor huidige scope, operational overhead

### Backend-Only (Server-Side Rendering)
- **Waarom niet**: Minder goede developer experience voor frontend developers, minder interactieve UI mogelijkheden

### Frontend-Only (SPA met Backend-as-a-Service)
- **Waarom niet**: Minder controle over backend, vendor lock-in risico

## Referenties
- [Nuxt.js Documentation](https://nuxt.com/docs)
- [Go Documentation](https://go.dev/doc/)
- [Effective Platform Engineering - Chapter 2: Software-Defined Products]

---

**Datum**: 2025-01-XX  
**Auteur**: Platform Engineering Team  
**Reviewers**: Frontend Team, Backend Team

