# Scaleway Kosten Optimalisatie Guide

Deze guide helpt je om een minimalistisch en kostenefficiënt Kubernetes cluster op te zetten op Scaleway.

## Minimale Cluster Setup

### 1. Cluster Aanmaken

1. Ga naar Scaleway Console → Kubernetes → Create Cluster
2. Kies **Kapsule** (managed Kubernetes)
3. Selecteer **Amsterdam** region (meestal goedkoper)
4. Kies **Kubernetes 1.31+** (laatste stabiele versie)

### 2. Node Pool Configuratie

**Minimale Setup (Beschikbaar):**
- **Pool Name**: `default-pool`
- **Node Type**: `DEV1-M` (4 vCPU, 8GB RAM)
  - ⚠️ Meestal het minimum beschikbare type
- **Node Count**: `1` (start klein)
- **Auto-scaling**: 
  - Min: `1`
  - Max: `2` (voor kostenbeheersing)
- **Spot Instances**: `Disabled` (voor stabiliteit, kan later aan)

**Geschatte Kosten**: ~€30/maand voor 1 node

**Voordeel**: Meer headroom voor groei zonder direct op te hoeven schalen!

### 3. Container Registry Setup

1. Maak een **Container Registry** aan in Scaleway Console
2. Kies **Amsterdam** region
3. Kies **Small** storage (10GB) - genoeg voor start
4. **Geschatte Kosten**: ~€1-2/maand

### 4. Load Balancer (Optioneel)

**Optie A: Zonder Load Balancer (Goedkoopst)**
- Gebruik **NodePort** service type
- Toegang via node IP + poort
- **Kosten**: €0/maand

**Optie B: Met Load Balancer**
- Gebruik **Ingress** met Scaleway Load Balancer
- **Kosten**: ~€10/maand
- Alleen nodig voor productie met custom domain

**Aanbeveling**: Start zonder Load Balancer, voeg later toe indien nodig.

## Resource Optimalisatie

### Huidige Configuratie

De applicatie is geconfigureerd met minimale resources:

**Frontend (Nuxt):**
- Replicas: 1
- CPU Request: 100m (0.1 core)
- Memory Request: 128Mi
- CPU Limit: 500m (0.5 core)
- Memory Limit: 512Mi

**Backend (Go):**
- Replicas: 1
- CPU Request: 50m (0.05 core)
- Memory Request: 64Mi
- CPU Limit: 200m (0.2 core)
- Memory Limit: 256Mi

**Totaal Resource Gebruik:**
- CPU: ~150m requests, 700m limits
- Memory: ~192Mi requests, 768Mi limits

Dit past ruimschoots op een **DEV1-S** node (2 vCPU, 4GB RAM).

## Kosten Breakdown

### Maandelijkse Kosten (Met DEV1-M)

| Component | Kosten |
|-----------|--------|
| Kubernetes Cluster (1x DEV1-M) | ~€30 |
| Container Registry (10GB) | ~€1-2 |
| Load Balancer (optioneel) | ~€10 |
| **Totaal (zonder LB)** | **~€31-32** |
| **Totaal (met LB)** | **~€41-42** |

**Opmerking**: DEV1-M is meestal het minimum beschikbare type. Dit geeft je meer headroom voor groei!

### Jaarlijkse Kosten

- **Zonder Load Balancer**: ~€192-204/jaar
- **Met Load Balancer**: ~€312-324/jaar

## Kostenbesparende Tips

### 1. Gebruik Spot Instances (Geavanceerd)

Spot instances zijn tot 70% goedkoper maar kunnen worden gestopt:
```bash
# Configureer spot instances in node pool
# Alleen voor development/staging, niet voor productie
```

### 2. Schaal Alleen Op Bij Nodig

De HPA is geconfigureerd maar niet actief standaard:
```bash
# Activeer alleen als je verwacht variabele load te hebben
make k8s-apply-hpa
```

### 3. Monitor Resource Gebruik

```bash
# Check actueel gebruik
kubectl top pods -n terminal-zoo
kubectl top nodes

# Als gebruik laag is, verlaag resources verder
kubectl edit deployment nuxt-frontend -n terminal-zoo
kubectl edit deployment go-backend -n terminal-zoo
```

### 4. Cleanup Oude Images

```bash
# Verwijder oude images uit registry regelmatig
# Via Scaleway Console of CLI
scw registry image list
scw registry image delete <image-id>
```

### 5. Gebruik Node Affinity (Geavanceerd)

Forceer pods op dezelfde node om kosten te besparen:
```yaml
# Voeg toe aan deployment spec
affinity:
  podAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - nuxt-frontend
            - go-backend
        topologyKey: kubernetes.io/hostname
```

### 6. Schakel Ongebruikte Services Uit

- Monitoring tools (indien niet nodig)
- Log aggregation (gebruik kubectl logs)
- Extra ingress controllers

## Scaling Strategie

### Wanneer Schalen?

**Schaal op wanneer:**
- CPU gebruik consistent >70%
- Memory gebruik consistent >80%
- Response times stijgen
- Pods worden gekilled door OOM

**Schaal af wanneer:**
- CPU gebruik consistent <30%
- Memory gebruik consistent <50%
- Je kosten wilt besparen

### Handmatig Schalen

```bash
# Scale frontend
kubectl scale deployment/nuxt-frontend --replicas=2 -n terminal-zoo

# Scale backend
kubectl scale deployment/go-backend --replicas=2 -n terminal-zoo
```

### Automatisch Schalen (HPA)

```bash
# Activeer HPA
make k8s-apply-hpa

# Check HPA status
kubectl get hpa -n terminal-zoo

# Deactiveer HPA om kosten te besparen
make k8s-delete-hpa
```

## Cluster Autoscaling

Configureer cluster autoscaling voor automatische node management:

```bash
# Via Scaleway Console:
# Kubernetes → Your Cluster → Node Pools → Edit
# Enable Auto-scaling:
# - Min nodes: 1
# - Max nodes: 2 (voor kostenbeheersing)
```

## Monitoring Kosten

### Check Resource Usage

```bash
# Pod resource usage
kubectl top pods -n terminal-zoo

# Node resource usage
kubectl top nodes

# Detailed resource usage
kubectl describe node <node-name>
```

### Kosten Tracking

1. **Scaleway Console**: Billing → Usage
2. **Set Budget Alerts**: Configureer alerts bij bepaalde kosten
3. **Monthly Review**: Review maandelijks en pas aan indien nodig

## Troubleshooting Kosten

### Pods Starten Niet (Resource Issues)

```bash
# Check beschikbare resources
kubectl describe node

# Als niet genoeg resources, verlaag requests/limits
kubectl edit deployment nuxt-frontend -n terminal-zoo
kubectl edit deployment go-backend -n terminal-zoo
```

### Hoge Kosten

1. Check aantal nodes: `kubectl get nodes`
2. Check aantal pods: `kubectl get pods -A`
3. Check resource usage: `kubectl top nodes`
4. Scale down indien mogelijk
5. Overweeg kleinere node types

## Best Practices

1. **Start Klein**: Begin met 1 node, schaal op indien nodig
2. **Monitor Regelmatig**: Check resource usage wekelijks
3. **Set Limits**: Zorg dat resource limits zijn gezet
4. **Cleanup Regelmatig**: Verwijder oude images en resources
5. **Use Spot Instances**: Voor development/staging omgevingen
6. **Review Maandelijks**: Pas configuratie aan op basis van gebruik

## Support

Voor vragen over Scaleway kosten:
- Scaleway Billing Dashboard
- Scaleway Support (via console)
- Check [Scaleway Pricing](https://www.scaleway.com/en/pricing/)

