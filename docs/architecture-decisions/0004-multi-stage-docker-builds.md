# ADR 0004: Multi-Stage Docker Builds

## Status
Geaccepteerd

## Context
Docker images moeten gebouwd worden voor zowel frontend als backend. Er zijn verschillende benaderingen:
- Single-stage builds (eenvoudig, maar grote images)
- Multi-stage builds (kleinere images, maar complexer)
- Buildpacks (automatische builds, maar minder controle)
- External build tools (CI/CD builds, maar meer dependencies)

## Beslissing
We kiezen voor **multi-stage Docker builds** waarbij:
- **Frontend**: Multi-stage build met builder en production stage
- **Backend**: Multi-stage build met Go builder en Alpine runtime
- **Platform**: Standaard `linux/amd64` voor Scaleway compatibiliteit
- **Optimization**: Minimale production images met alleen runtime dependencies

## Gevolgen

### Positief
- ✅ **Smaller Images**: Production images bevatten alleen runtime dependencies
- ✅ **Security**: Minder attack surface door kleinere images
- ✅ **Performance**: Snellere image pulls en container starts
- ✅ **Cost**: Lagere storage en transfer costs
- ✅ **Build Optimization**: Build dependencies worden niet in production image opgenomen

### Negatief
- ⚠️ **Complexity**: Meer complexe Dockerfiles
- ⚠️ **Build Time**: Mogelijk langere build tijd door meerdere stages
- ⚠️ **Debugging**: Moeilijker om build issues te debuggen

## Implementatie Details

### Frontend Multi-Stage Build
```dockerfile
# Builder stage
FROM --platform=linux/amd64 node:24.6.0-alpine3.22 AS builder
# ... build steps ...

# Production stage
FROM --platform=linux/amd64 node:24.6.0-alpine3.22 AS production
# ... copy built artifacts, install only production dependencies ...
```

**Location**: `Dockerfile`  
**Size Reduction**: ~70% kleiner dan single-stage build

### Backend Multi-Stage Build
```dockerfile
# Builder stage
FROM --platform=linux/amd64 golang:1.25-alpine AS builder
# ... build Go binary ...

# Runtime stage
FROM --platform=linux/amd64 alpine:3.22.1
# ... copy binary, install only runtime dependencies ...
```

**Location**: `go/Dockerfile.multi-stage`  
**Size Reduction**: ~90% kleiner dan single-stage build (geen Go toolchain in production)

### Platform Specification
- Alle builds gebruiken expliciet `--platform linux/amd64`
- Zorgt voor consistentie tussen lokale builds (mogelijk Apple Silicon) en productie (Scaleway AMD64)
- Voorkomt platform mismatch errors in Kubernetes

### Build Commands
```bash
# Frontend
docker build --platform linux/amd64 -t <image> .

# Backend
docker build --platform linux/amd64 -t <image> -f go/Dockerfile.multi-stage ./go
```

## Alternatieven Overwogen

### Single-Stage Builds
- **Waarom niet**: Te grote images, security risico's, hogere costs

### Buildpacks
- **Waarom niet**: Minder controle over build proces, mogelijk niet optimaal voor onze use case

### External Build Tools (Bazel, Buildah)
- **Waarom niet**: Meer complexiteit, extra tooling vereist

### Distroless Images
- **Overweging**: Mogelijk toekomstige optimalisatie voor backend
- **Huidige status**: Alpine is voldoende voor nu

## Image Size Comparison

### Frontend
- Single-stage: ~500MB
- Multi-stage: ~150MB
- **Reduction**: ~70%

### Backend
- Single-stage (met Go toolchain): ~800MB
- Multi-stage (Alpine + binary): ~50MB
- **Reduction**: ~94%

## Best Practices Gevolgd

1. **Minimal Base Images**: Alpine Linux voor minimale footprint
2. **Layer Caching**: Build dependencies eerst, code later voor betere caching
3. **Non-Root User**: Images draaien als non-root user (security)
4. **Platform Specification**: Expliciete platform voor consistentie
5. **.dockerignore**: Exclude onnodige files van build context

## Referenties
- [Docker Multi-Stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Effective Platform Engineering - Chapter 6: Building Software-Defined Platforms]
- [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)

---

**Datum**: 2025-01-XX  
**Auteur**: Platform Engineering Team  
**Reviewers**: DevOps Team

