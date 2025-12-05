# Deployment Documentatie

Deze documentatie beschrijft hoe je de App Store applicatie naar Scaleway productie deployt.

## Overzicht

De applicatie bestaat uit twee hoofdcomponenten:
- **Frontend**: Nuxt.js applicatie (poort 3000) - neutrale branding, geen specifieke organisatie logo's
- **Backend**: Go Kubernetes API server (poort 8080)

Beide componenten worden gedeployed naar een Kubernetes cluster op Scaleway.

**Belangrijke configuratie**:
- Alle Docker builds zijn standaard geconfigureerd voor `linux/amd64` platform
- Kubernetes namespace: `nl-appstore-registry`
- Frontend image: `rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest`
- Backend image: `rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest`

## Vereisten

### Scaleway Account
- Scaleway Kubernetes Kapsule cluster (of K8S)
- Scaleway Container Registry (SCR)
- Domain name voor de applicatie
- Scaleway account met voldoende quota

### Lokale Tools
- `kubectl` geconfigureerd voor je Scaleway cluster
- `docker` voor lokale builds
- `make` voor build commando's

> **💡 Kubeconfig Setup**: Zorg dat je kubeconfig geconfigureerd hebt voordat je kubectl gebruikt! Zie **[KUBECONFIG_SETUP.md](KUBECONFIG_SETUP.md)**

## Voorbereiding

### 1. Scaleway Container Registry Setup

1. Maak een Container Registry aan in Scaleway Console
2. Noteer de registry URL (bijv. `rg.nl-ams.scw.cloud`)
3. Genereer een Secret Key voor de registry

### 2. Kubernetes Cluster Setup

1. Maak een Kubernetes Kapsule cluster aan in Scaleway
2. Download en configureer `kubectl`:
   ```bash
   # Download kubeconfig van Scaleway Console
   # Of gebruik Scaleway CLI:
   scw k8s kubeconfig install <cluster-id>
   ```

3. Verifieer connectie:
   ```bash
   kubectl get nodes
   ```

### 3. Environment Variabelen Configureren

Kopieer `.env.example` naar `.env` en vul de waarden in:

```bash
cp .env.example .env
# Bewerk .env met je productie waarden
```

Belangrijke variabelen:
- `NUXT_EXTERNAL_API_BASE`: Backend API URL
- `NUXT_EXTERNAL_API_AUTH`: API authenticatie token
- `SCR_REGISTRY`: Scaleway Container Registry URL
- `SCR_NAMESPACE`: Je registry namespace
- `SCR_SECRET_KEY`: Je Scaleway secret key

### 4. Kubernetes Secrets Aanmaken

#### Container Registry Secret
```bash
make k8s-create-scr-secret
# Of handmatig:
kubectl create secret docker-registry scr-secret \
  --docker-server=rg.nl-ams.scw.cloud \
  --docker-username=nologin \
  --docker-password=YOUR_SECRET_KEY \
  -n nl-appstore-registry
```

#### Application Secrets
```bash
kubectl create secret generic app-secrets \
  --from-literal=NUXT_EXTERNAL_API_AUTH="Bearer YOUR_TOKEN" \
  -n nl-appstore-registry
```

## Deployment Stappen

### Stap 1: Build en Push Docker Images

```bash
# Build en push beide images naar Scaleway Container Registry
make push-all

# Of individueel:
make push-frontend
make push-backend
```

### Stap 2: Update Kubernetes Manifests

Bewerk de volgende bestanden met je specifieke waarden:

1. **k8s/configmap.yaml**: Update `NUXT_EXTERNAL_API_BASE` indien nodig
2. **k8s/ingress.yaml**: Update domain names
3. **k8s/frontend-deployment.yaml**: Update image registry URL
4. **k8s/backend-deployment.yaml**: Update image registry URL

### Stap 3: Deploy naar Kubernetes

```bash
# Apply alle Kubernetes manifests
make k8s-apply

# Check status
make k8s-status
```

### Stap 4: Verifieer Deployment

```bash
# Check pod status
kubectl get pods -n nl-appstore-registry

# Check logs
make k8s-logs-frontend
make k8s-logs-backend

# Test endpoints
kubectl port-forward -n nl-appstore-registry service/nuxt-frontend 3000:3000
# Open http://localhost:3000 in browser
```

## SSL/TLS Certificaten

### Optie 1: Scaleway Managed SSL
1. Upload certificaat in Scaleway Console
2. Update `k8s/ingress.yaml` met certificaat ID:
   ```yaml
   annotations:
     ingress.scaleway.com/ssl-certificate-id: "your-cert-id"
   ```

### Optie 2: Cert-Manager met Let's Encrypt
1. Installeer cert-manager:
   ```bash
   kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
   ```

2. Update `k8s/ingress.yaml` met cert-manager annotations (al geconfigureerd)

## CI/CD Pipeline

De applicatie heeft een GitHub Actions workflow (`.github/workflows/deploy.yml`) die automatisch:
1. Docker images bouwt bij push naar `main` branch
2. Images pusht naar Scaleway Container Registry
3. Deployt naar Kubernetes met Scaleway CLI

### Quick Start

Voor snelle setup, volg: **[QUICK_START.md](QUICK_START.md)**

### GitHub Secrets Configureren

Configureer de volgende secrets in GitHub (Settings → Secrets → Actions):

**Verplichte Secrets:**
- `SCALEWAY_ACCESS_KEY`: Je Scaleway API access key
- `SCALEWAY_SECRET_KEY`: Je Scaleway API secret key
- `SCALEWAY_K8S_CLUSTER_ID`: Je Kubernetes cluster ID (UUID)
- `SCR_NAMESPACE`: Je Container Registry namespace naam
- `SCR_USERNAME`: Meestal `nologin`
- `SCR_SECRET_KEY`: Je Container Registry secret key

Zie **[AUTOMATIC_DEPLOYMENT_SETUP.md](AUTOMATIC_DEPLOYMENT_SETUP.md)** voor gedetailleerde instructies.

## Monitoring en Health Checks

### Health Check Endpoints

- Frontend: `http://your-domain.com/api/health`
- Backend: `http://api.your-domain.com/health`
- Readiness: `http://api.your-domain.com/ready`

### Logs Bekijken

```bash
# Frontend logs
make k8s-logs-frontend

# Backend logs
make k8s-logs-backend

# Alle pods
kubectl logs -f -l app=nuxt-frontend -n nl-appstore-registry
kubectl logs -f -l app=go-backend -n nl-appstore-registry
```

## Troubleshooting

### Pods starten niet
```bash
# Check pod events
kubectl describe pod <pod-name> -n nl-appstore-registry

# Check logs
kubectl logs <pod-name> -n nl-appstore-registry
```

### Image pull errors
```bash
# Verifieer image pull secret
kubectl get secret scr-secret -n nl-appstore-registry

# Test image pull
kubectl run test-pod --image=<your-image> --restart=Never -n nl-appstore-registry
```

### Ingress werkt niet
```bash
# Check ingress status
kubectl describe ingress app-ingress -n nl-appstore-registry

# Check ingress controller
kubectl get pods -n kube-system | grep ingress
```

## Rollback

```bash
# Rollback naar vorige versie
kubectl rollout undo deployment/nuxt-frontend -n nl-appstore-registry
kubectl rollout undo deployment/go-backend -n nl-appstore-registry

# Check rollout history
kubectl rollout history deployment/nuxt-frontend -n nl-appstore-registry
```

## Scaling

```bash
# Scale frontend
kubectl scale deployment/nuxt-frontend --replicas=3 -n nl-appstore-registry

# Scale backend
kubectl scale deployment/go-backend --replicas=3 -n nl-appstore-registry
```

## Onderhoud

### Updates Deployen

1. Push nieuwe code naar `main` branch
2. CI/CD pipeline bouwt automatisch nieuwe images
3. Update image tags in deployment files:
   ```bash
   kubectl set image deployment/nuxt-frontend nuxt-frontend=<new-image> -n nl-appstore-registry
   kubectl set image deployment/go-backend go-backend=<new-image> -n nl-appstore-registry
   ```

### Configuratie Updates

```bash
# Update configmap
kubectl apply -f k8s/configmap.yaml

# Restart pods om nieuwe config te laden
make k8s-rollout-restart
```

## Veiligheid

- Containers draaien als non-root users
- Read-only root filesystem waar mogelijk
- Resource limits zijn geconfigureerd
- Network policies kunnen worden toegevoegd voor extra isolatie

## Kosten Optimalisatie

De configuratie is geoptimaliseerd voor minimale kosten:

### Huidige Configuratie
- **1 replica** per service (minimaal voor beschikbaarheid)
- **Lage resource requests**: Frontend 100m CPU/128Mi RAM, Backend 50m CPU/64Mi RAM
- **Beperkte resource limits**: Voorkomt over-allocation
- **ImagePullPolicy: IfNotPresent**: Bespaart op registry pull kosten
- **HPA geconfigureerd**: Scale alleen op bij hoge load (min 1, max 2-3)

### Scaleway Cluster Setup Tips

1. **Kies de kleinste node pool**:
   - **DEV1-S** (2 vCPU, 4GB RAM) - ~€15/maand
   - Of **DEV1-M** (4 vCPU, 8GB RAM) - ~€30/maand voor meer headroom
   - Start met 1 node, voeg nodes toe indien nodig

2. **Gebruik spot instances** (indien beschikbaar):
   ```bash
   # Spot instances zijn tot 70% goedkoper
   # Configureer in node pool settings
   ```

3. **Cluster autoscaling**:
   - Configureer cluster autoscaler met minimum 1 node
   - Maximum 2-3 nodes voor kostenbeheersing

4. **Monitoring en optimalisatie**:
   ```bash
   # Check resource usage
   kubectl top pods -n nl-appstore-registry
   kubectl top nodes
   
   # Pas resources aan indien nodig
   kubectl edit deployment nuxt-frontend -n nl-appstore-registry
   kubectl edit deployment go-backend -n nl-appstore-registry
   ```

5. **HPA activeren** (optioneel, voor automatische scaling):
   ```bash
   kubectl apply -f k8s/hpa.yaml
   ```

### Geschatte Maandelijkse Kosten

Met minimale configuratie:
- **Kubernetes cluster** (1x DEV1-S node): ~€15/maand
- **Container Registry**: ~€1-2/maand (afhankelijk van storage)
- **Load Balancer**: ~€10/maand (indien gebruikt)
- **Totaal**: ~€26-27/maand voor basis setup

### Verdere Kostenbesparing

- **Gebruik Ingress zonder Load Balancer**: Gebruik NodePort of Ingress met gratis SSL
- **Schakel monitoring uit** indien niet nodig
- **Gebruik lokale storage** in plaats van managed volumes waar mogelijk
- **Cleanup oude images** regelmatig uit registry
- **Overweeg serverless** voor minder kritieke componenten

## Support

Voor vragen of problemen:
1. Check logs: `make k8s-logs-frontend` en `make k8s-logs-backend`
2. Check pod status: `kubectl get pods -n nl-appstore-registry`
3. Check events: `kubectl get events -n nl-appstore-registry --sort-by='.lastTimestamp'`

