# Automatische Deployment naar Scaleway Setup

Deze guide legt uit hoe je automatische deployment naar Scaleway configureert via GitHub Actions.

## Overzicht

De CI/CD pipeline (`/.github/workflows/deploy.yml`) doet automatisch:
1. ✅ Build Docker images bij push naar `main` branch
2. ✅ Push images naar Scaleway Container Registry
3. ✅ Deploy naar Kubernetes cluster
4. ✅ Verifieer deployment status

## Vereisten

### 1. Scaleway Account Setup

Je hebt nodig:
- **Scaleway Access Key**: API access key
- **Scaleway Secret Key**: API secret key  
- **Kubernetes Cluster ID**: ID van je Kapsule cluster
- **Container Registry**: Registry namespace naam

### 2. GitHub Secrets Configureren

Ga naar je GitHub repository → **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

Voeg de volgende secrets toe:

#### Verplichte Secrets

| Secret Name | Beschrijving | Voorbeeld |
|------------|--------------|-----------|
| `SCALEWAY_ACCESS_KEY` | Je Scaleway API access key | `SCWXXXXXXXXXXXXXXXXX` |
| `SCALEWAY_SECRET_KEY` | Je Scaleway API secret key | `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |
| `SCALEWAY_K8S_CLUSTER_ID` | Kubernetes cluster ID | `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |
| `SCR_NAMESPACE` | Container Registry namespace | `my-registry` |
| `SCR_USERNAME` | Registry username (meestal `nologin`) | `nologin` |
| `SCR_SECRET_KEY` | Registry secret key | `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |

#### Optionele Secrets

| Secret Name | Beschrijving | Wanneer nodig |
|------------|--------------|---------------|
| `KUBECONFIG` | Base64 encoded kubeconfig | Alleen als je Scaleway CLI niet gebruikt |

## Stap-voor-Stap Setup

### Stap 1: Scaleway API Keys Aanmaken

1. Log in op [Scaleway Console](https://console.scaleway.com)
2. Ga naar **IAM** → **API Keys**
3. Klik op **Generate API Key**
4. Geef een naam (bijv. "GitHub Actions")
5. Kopieer de **Access Key** en **Secret Key**
6. **BELANGRIJK**: Bewaar de Secret Key veilig, deze wordt maar 1x getoond!

### Stap 2: Kubernetes Cluster ID Vinden

1. Ga naar **Kubernetes** → **Kapsule**
2. Klik op je cluster
3. Kopieer de **Cluster ID** (UUID formaat)

Of via CLI:
```bash
scw k8s cluster list
```

### Stap 3: Container Registry Info

1. Ga naar **Container Registry**
2. Noteer je **Registry URL** (bijv. `rg.nl-ams.scw.cloud`)
3. Noteer je **Namespace** naam
4. Ga naar **Secrets** → **Generate new secret key**
5. Kopieer de secret key

### Stap 4: GitHub Secrets Toevoegen

1. Ga naar je GitHub repository
2. **Settings** → **Secrets and variables** → **Actions**
3. Klik **New repository secret** voor elk secret:

```bash
# Voeg deze toe:
SCALEWAY_ACCESS_KEY=<jouw-access-key>
SCALEWAY_SECRET_KEY=<jouw-secret-key>
SCALEWAY_K8S_CLUSTER_ID=<cluster-id>
SCR_NAMESPACE=<registry-namespace>
SCR_USERNAME=nologin
SCR_SECRET_KEY=<registry-secret-key>
```

### Stap 5: Kubernetes Manifests Aanpassen

Update de volgende bestanden met je specifieke waarden:

#### `k8s/configmap.yaml`
```yaml
data:
  NUXT_EXTERNAL_API_BASE: "http://go-backend:8080"  # Of externe URL
```

#### `k8s/frontend-deployment.yaml` en `k8s/backend-deployment.yaml`
```yaml
# Frontend
image: rg.nl-ams.scw.cloud/<jouw-registry>/ai-co:latest

# Backend (gebruikt dezelfde image naam)
image: rg.nl-ams.scw.cloud/<jouw-registry>/ai-co:latest
```

#### `k8s/ingress.yaml` (optioneel)
```yaml
spec:
  rules:
    - host: jouw-domein.com  # Update met je eigen domain
```

### Stap 6: Eerste Deployment

1. **Push naar main branch**:
   ```bash
   git add .
   git commit -m "Configure automatic deployment"
   git push origin main
   ```

2. **Check GitHub Actions**:
   - Ga naar **Actions** tab in GitHub
   - Je ziet de workflow draaien
   - Klik op de run om logs te zien

3. **Verifieer deployment**:
   ```bash
   kubectl get pods -n terminal-zoo
   kubectl get services -n terminal-zoo
   ```

## Workflow Details

### Triggers

De workflow draait automatisch wanneer:
- Code wordt gepusht naar `main` branch
- Handmatig via **Actions** → **Run workflow**

### Build Proces

1. **Frontend Build**:
   - Builds Nuxt.js applicatie
   - Tagt image met commit SHA
   - Pusht naar Container Registry

2. **Backend Build**:
   - Builds Go binary
   - Tagt image met commit SHA
   - Pusht naar Container Registry

3. **Deployment**:
   - Configureert kubectl met Scaleway CLI
   - Update image tags in deployments
   - Apply alle Kubernetes manifests
   - Wacht op rollout completion
   - Verifieert deployment status

### Image Tagging

Images worden getagged als:
- `main-<commit-sha>` (bijv. `main-a1b2c3d`)
- `latest` (alleen voor main branch)

## Troubleshooting

### Workflow Fails: Authentication Error

**Probleem**: `scaleway: authentication failed`

**Oplossing**:
1. Verifieer `SCALEWAY_ACCESS_KEY` en `SCALEWAY_SECRET_KEY` zijn correct
2. Check of de keys niet zijn verlopen
3. Verifieer de keys hebben de juiste permissions

### Workflow Fails: Cluster Not Found

**Probleem**: `cluster not found`

**Oplossing**:
1. Verifieer `SCALEWAY_K8S_CLUSTER_ID` is correct
2. Check of het cluster bestaat en actief is
3. Verifieer de API keys hebben toegang tot het cluster

### Workflow Fails: Image Pull Error

**Probleem**: `ImagePullBackOff` in Kubernetes

**Oplossing**:
1. Verifieer `SCR_SECRET_KEY` is correct
2. Check of image pull secret bestaat:
   ```bash
   kubectl get secret scr-secret -n terminal-zoo
   ```
3. Maak secret aan indien nodig:
   ```bash
   make k8s-create-scr-secret
   ```

### Workflow Fails: Rollout Timeout

**Probleem**: `rollout status timeout`

**Oplossing**:
1. Check pod logs:
   ```bash
   kubectl logs -l app=nuxt-frontend -n terminal-zoo
   kubectl logs -l app=go-backend -n terminal-zoo
   ```
2. Check pod status:
   ```bash
   kubectl describe pod -l app=nuxt-frontend -n terminal-zoo
   ```
3. Verifieer resource limits zijn niet te laag
4. Check of health checks werken

### Deployment Succeeds But App Doesn't Work

1. **Check ingress**:
   ```bash
   kubectl get ingress -n terminal-zoo
   kubectl describe ingress app-ingress -n terminal-zoo
   ```

2. **Check services**:
   ```bash
   kubectl get svc -n terminal-zoo
   kubectl describe svc nuxt-frontend -n terminal-zoo
   ```

3. **Test endpoints**:
   ```bash
   kubectl port-forward -n terminal-zoo svc/nuxt-frontend 3000:3000
   curl http://localhost:3000/api/health
   ```

## Security Best Practices

1. **Gebruik GitHub Secrets**: Nooit credentials in code
2. **Rotate Keys Regelmatig**: Update API keys elke 3-6 maanden
3. **Minimale Permissions**: Geef alleen benodigde permissions aan API keys
4. **Review Workflow Logs**: Check regelmatig voor security issues
5. **Use Branch Protection**: Bescherm main branch met required reviews

## Advanced: Custom Deployment

### Deploy Alleen Frontend of Backend

Pas de workflow aan om alleen specifieke componenten te deployen:

```yaml
jobs:
  deploy-frontend:
    if: contains(github.event.head_commit.message, '[deploy-frontend]')
    # ... deployment steps
```

### Deploy naar Staging

Voeg een staging environment toe:

```yaml
on:
  push:
    branches:
      - main      # Production
      - staging   # Staging
```

### Manual Deployment

Deploy handmatig via GitHub Actions UI:
1. Ga naar **Actions** → **Build and Deploy**
2. Klik **Run workflow**
3. Selecteer branch en klik **Run workflow**

## Monitoring

### GitHub Actions Status

- **Groen**: Deployment succesvol
- **Geel**: Deployment bezig
- **Rood**: Deployment gefaald (check logs)

### Kubernetes Status

```bash
# Check deployment status
kubectl get deployments -n terminal-zoo

# Check pod status
kubectl get pods -n terminal-zoo

# Check recent events
kubectl get events -n terminal-zoo --sort-by='.lastTimestamp'
```

## Rollback

Als deployment faalt, rollback naar vorige versie:

```bash
# Via kubectl
kubectl rollout undo deployment/nuxt-frontend -n terminal-zoo
kubectl rollout undo deployment/go-backend -n terminal-zoo

# Of via GitHub Actions (re-run vorige succesvolle workflow)
```

## Support

Voor problemen:
1. Check GitHub Actions logs
2. Check Kubernetes events: `kubectl get events -n terminal-zoo`
3. Check pod logs: `kubectl logs -l app=nuxt-frontend -n terminal-zoo`
4. Scaleway Support: [Scaleway Documentation](https://www.scaleway.com/en/docs/)


