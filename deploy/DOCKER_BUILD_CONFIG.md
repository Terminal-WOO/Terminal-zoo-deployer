# Docker Build Configuratie

## Platform Configuratie

Alle Docker builds zijn standaard geconfigureerd voor `linux/amd64` platform om compatibiliteit met Scaleway Kubernetes clusters te garanderen.

## Frontend Build

### Dockerfile Configuratie

De frontend `Dockerfile` gebruikt multi-stage builds met expliciete platform specificatie:

```dockerfile
# Build stage
FROM --platform=linux/amd64 node:24.6.0-alpine3.22 AS builder
# ... build steps ...

# Production stage
FROM --platform=linux/amd64 node:24.6.0-alpine3.22 AS production
# ... production steps ...
```

### Build Commando

```bash
# Standaard build (gebruikt linux/amd64 automatisch)
make build-frontend

# Of handmatig met expliciete platform:
docker build --platform linux/amd64 -t <image-name> .
```

## Backend Build

### Dockerfile Configuratie

De backend gebruikt `go/Dockerfile.multi-stage` met platform specificatie:

```dockerfile
# Stage 1: Build
FROM --platform=linux/amd64 golang:1.25-alpine AS builder
# ... build steps ...

# Stage 2: Runtime
FROM --platform=linux/amd64 alpine:3.22.1
# ... runtime steps ...
```

### Build Commando

```bash
# Standaard build (gebruikt linux/amd64 automatisch)
make build-backend

# Of handmatig:
docker build --platform linux/amd64 -t <image-name> -f go/Dockerfile.multi-stage ./go
```

## Makefile Configuratie

Alle `docker build` commando's in de `Makefile` bevatten expliciet `--platform linux/amd64`:

```makefile
build-frontend:
	docker build --platform linux/amd64 -t $(FRONTEND_IMAGE):$(IMAGE_TAG) .

build-backend:
	docker build --platform linux/amd64 -t $(BACKEND_IMAGE):$(IMAGE_TAG) -f go/Dockerfile.multi-stage ./go
```

## Waarom linux/amd64?

1. **Scaleway Compatibiliteit**: Scaleway Kubernetes clusters gebruiken standaard AMD64 architecture
2. **Consistentie**: Voorkomt platform-specifieke build problemen
3. **Performance**: AMD64 images zijn geoptimaliseerd voor productie omgevingen

## Troubleshooting

### Platform Mismatch Errors

Als je een error krijgt zoals:
```
ERROR: failed to solve: failed to compute cache key: no match for platform in manifest
```

**Oplossing**: Zorg dat je `--platform linux/amd64` gebruikt in alle build commando's.

### Build op Apple Silicon (M1/M2)

Als je op Apple Silicon bouwt maar voor linux/amd64 target:

```bash
# Docker Desktop doet automatisch cross-platform builds
# Zorg dat BuildKit enabled is:
export DOCKER_BUILDKIT=1

# Build commando's werken automatisch met --platform flag
make build-frontend
make build-backend
```

### Verificatie

Controleer dat je image het juiste platform heeft:

```bash
# Inspect image platform
docker inspect <image-name> | grep Architecture

# Of gebruik docker buildx
docker buildx imagetools inspect <image-name>
```

## Best Practices

1. **Altijd platform specificeren**: Gebruik altijd `--platform linux/amd64` voor productie builds
2. **Consistentie**: Gebruik dezelfde platform configuratie in alle build scripts
3. **CI/CD**: Zorg dat CI/CD pipelines ook `--platform linux/amd64` gebruiken
4. **Testing**: Test builds lokaal voordat je pusht naar registry

