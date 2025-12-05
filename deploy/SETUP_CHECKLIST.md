# Automatische Deployment Setup Checklist

Gebruik deze checklist om automatische deployment naar Scaleway te configureren.

## Pre-Flight Checks

- [ ] Scaleway account aangemaakt
- [ ] Kubernetes Kapsule cluster aangemaakt en actief
- [ ] Container Registry aangemaakt
- [ ] GitHub repository heeft Actions enabled

## Scaleway Setup

### API Keys
- [ ] Scaleway Access Key aangemaakt
- [ ] Scaleway Secret Key gekopieerd en veilig opgeslagen
- [ ] API keys hebben juiste permissions (Kubernetes + Container Registry)

### Kubernetes Cluster
- [ ] Cluster ID genoteerd
- [ ] Cluster is actief en bereikbaar
- [ ] kubectl werkt lokaal met cluster
- [ ] Test: `kubectl get nodes` werkt

### Container Registry
- [ ] Registry namespace naam genoteerd
- [ ] Registry secret key aangemaakt
- [ ] Registry URL genoteerd (bijv. `rg.nl-ams.scw.cloud`)

## GitHub Secrets Configuratie

Ga naar: **Repository Settings** → **Secrets and variables** → **Actions**

- [ ] `SCALEWAY_ACCESS_KEY` toegevoegd
- [ ] `SCALEWAY_SECRET_KEY` toegevoegd
- [ ] `SCALEWAY_K8S_CLUSTER_ID` toegevoegd
- [ ] `SCR_NAMESPACE` toegevoegd
- [ ] `SCR_USERNAME` toegevoegd (meestal `nologin`)
- [ ] `SCR_SECRET_KEY` toegevoegd

## Kubernetes Manifests Configuratie

- [ ] `k8s/configmap.yaml` - Registry URL en namespace geüpdatet
- [ ] `k8s/frontend-deployment.yaml` - Image URL geüpdatet
- [ ] `k8s/backend-deployment.yaml` - Image URL geüpdatet
- [ ] `k8s/ingress.yaml` - Domain names geüpdatet (optioneel)

## Eerste Deployment

### Lokale Test (Optioneel)
- [ ] Docker images gebouwd: `make build-frontend build-backend`
- [ ] Images gepusht naar registry: `make push-all`
- [ ] Kubernetes manifests getest: `make k8s-apply`
- [ ] Pods draaien: `make k8s-status`

### GitHub Actions Test
- [ ] Code gepusht naar `main` branch
- [ ] GitHub Actions workflow gestart
- [ ] Build jobs succesvol
- [ ] Deploy job succesvol
- [ ] Pods draaien in cluster: `kubectl get pods -n terminal-zoo`

## Post-Deployment Verificatie

- [ ] Frontend pod draait: `kubectl get pods -l app=nuxt-frontend -n terminal-zoo`
- [ ] Backend pod draait: `kubectl get pods -l app=go-backend -n terminal-zoo`
- [ ] Health checks werken: `kubectl logs -l app=nuxt-frontend -n terminal-zoo | grep health`
- [ ] Services zijn beschikbaar: `kubectl get svc -n terminal-zoo`
- [ ] Ingress werkt (indien geconfigureerd): `kubectl get ingress -n terminal-zoo`

## Monitoring Setup

- [ ] GitHub Actions notifications ingeschakeld
- [ ] Kubernetes events monitoring: `kubectl get events -n terminal-zoo`
- [ ] Logs toegankelijk: `make k8s-logs-frontend` en `make k8s-logs-backend`

## Security Checklist

- [ ] API keys zijn veilig opgeslagen (GitHub Secrets)
- [ ] Geen credentials in code
- [ ] Kubernetes secrets aangemaakt: `make k8s-create-scr-secret`
- [ ] Application secrets aangemaakt: `kubectl create secret generic app-secrets ...`
- [ ] Branch protection rules ingesteld (optioneel maar aanbevolen)

## Troubleshooting Resources

- [ ] `deploy/AUTOMATIC_DEPLOYMENT_SETUP.md` gelezen
- [ ] `deploy/README.md` deployment sectie gelezen
- [ ] Scaleway documentatie geraadpleegd
- [ ] GitHub Actions logs gecontroleerd bij problemen

## Quick Test Commands

```bash
# Test kubectl connectie
kubectl get nodes

# Test registry access
docker login rg.nl-ams.scw.cloud -u nologin -p <secret-key>

# Test deployment
make k8s-status

# Check logs
make k8s-logs-frontend
make k8s-logs-backend
```

## Volgende Stappen

Na succesvolle setup:
- [ ] Monitor eerste paar deployments
- [ ] Setup alerts voor failed deployments
- [ ] Documenteer custom configuraties
- [ ] Overweeg staging environment
- [ ] Review kosten maandelijks

---

**Status**: ⬜ Niet gestart | 🟡 In progress | ✅ Voltooid

**Laatste update**: _________________

**Notities**:
_________________________________________________
_________________________________________________


