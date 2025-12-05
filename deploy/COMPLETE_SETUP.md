# Complete Setup - Vanaf Images Pushen Tot Deployment

Complete stap-voor-stap guide om alles te deployen, inclusief images pushen.

## Prerequisites Check

Voordat je begint, check:

- [ ] Docker draait: `docker ps` (moet werken zonder error)
- [ ] Scaleway CLI geconfigureerd: `scw config get`
- [ ] kubectl werkt: `kubectl get nodes`
- [ ] Je bent ingelogd op Scaleway Container Registry

---

## Stap 1: Docker Starten (Als Nodig)

### Mac

```bash
# Start Docker Desktop applicatie
# Of via terminal:
open -a Docker

# Wacht tot Docker gestart is
docker ps
```

### Linux

```bash
# Start Docker service
sudo systemctl start docker

# Check status
sudo systemctl status docker
```

---

## Stap 2: Login op Scaleway Container Registry

```bash
# Login met je API Secret Key
docker login rg.nl-ams.scw.cloud \
  -u nologin \
  -p <jouw-api-secret-key>

# Verifieer
docker login rg.nl-ams.scw.cloud
```

---

## Stap 3: Images Bouwen en Pushen

### Frontend Image

```bash
# Build frontend image
make build-frontend

# Push naar registry
make push-frontend

# Of beide in één keer:
make push-all
```

### Backend Image

```bash
# Build backend image (gebruikt go/Makefile intern)
make build-backend

# Push naar registry
make push-backend
```

### Beide Images

```bash
# Build en push beide images
make push-all

# Dit doet:
# 1. Build frontend
# 2. Push frontend
# 3. Build backend (inclusief Go binary)
# 4. Push backend
```

**Verwachte output**: Images worden gepusht naar `rg.nl-ams.scw.cloud/nl-appstore-registry/`

---

## Stap 4: ConfigMap Aanmaken

```bash
# Apply configmap
kubectl apply -f k8s/configmap.yaml

# Verifieer
kubectl get configmap app-config -n nl-appstore-registry
```

---

## Stap 5: Secrets Aanmaken

### App Secrets (Voor Frontend API Auth)

```bash
# Maak app-secrets aan
kubectl create secret generic app-secrets \
  --from-literal=NUXT_EXTERNAL_API_AUTH="Bearer YOUR_TOKEN_HERE" \
  -n nl-appstore-registry

# Verifieer
kubectl get secret app-secrets -n nl-appstore-registry
```

**⚠️ Belangrijk**: Vervang `YOUR_TOKEN_HERE` met je echte API token!

### Registry Secret (Voor Image Pull)

```bash
# Maak registry secret aan voor image pull
kubectl create secret docker-registry scr-secret \
  --docker-server=rg.nl-ams.scw.cloud \
  --docker-username=nologin \
  --docker-password=<jouw-api-secret-key> \
  -n nl-appstore-registry

# Verifieer
kubectl get secret scr-secret -n nl-appstore-registry
```

**⚠️ Belangrijk**: Vervang `<jouw-api-secret-key>` met je Scaleway API Secret Key!

---

## Stap 6: Deployments Deployen

```bash
# Apply frontend deployment
kubectl apply -f k8s/frontend-deployment.yaml

# Apply backend deployment
kubectl apply -f k8s/backend-deployment.yaml

# Check status
kubectl get pods -n nl-appstore-registry

# Wacht tot pods Running zijn (kan 2-5 minuten duren)
kubectl wait --for=condition=ready pod -l app=nuxt-frontend -n nl-appstore-registry --timeout=5m
kubectl wait --for=condition=ready pod -l app=go-backend -n nl-appstore-registry --timeout=5m
```

---

## Stap 7: NodePort Services Deployen

```bash
# Apply NodePort services
make k8s-apply-nodeport

# Check services
kubectl get svc -n nl-appstore-registry

# Je zou moeten zien:
# nuxt-frontend-np     NodePort   10.x.x.x   <none>   3000:30080/TCP
# go-backend-np        NodePort   10.x.x.x   <none>   8080:30088/TCP
```

---

## Stap 8: Node IP Vinden

```bash
# Via Makefile
make k8s-get-node-ip

# Of handmatig:
kubectl get nodes -o wide

# Noteer het IP adres (EXTERNAL-IP of INTERNAL-IP)
```

---

## Stap 9: Firewall Rules Configureren

### Via Scaleway Console

1. **Ga naar Instances → Security Groups**
2. **Zoek de security group** van je Kubernetes nodes
3. **Voeg inbound rules toe**:
   - **TCP poort 30080** (voor frontend)
   - **TCP poort 30088** (voor backend)
   - **Source**: `0.0.0.0/0` (of specifiek IP voor security)

### Via Scaleway CLI

```bash
# List security groups
scw instance security-group list

# Voeg rules toe (vervang <sg-id> met je security group ID)
scw instance security-group rule add-inbound \
  security-group-id=<sg-id> \
  action=accept \
  protocol=TCP \
  port=30080 \
  ip-range=0.0.0.0/0

scw instance security-group rule add-inbound \
  security-group-id=<sg-id> \
  action=accept \
  protocol=TCP \
  port=30088 \
  ip-range=0.0.0.0/0
```

---

## Stap 10: Test Toegang

```bash
# Vervang <node-ip> met je node IP
NODE_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1)

# Test frontend health
curl http://$NODE_IP:30080/api/health
# Verwacht: {"status":"healthy"}

# Test backend health
curl http://$NODE_IP:30088/health
# Verwacht: {"status":"healthy"}
```

---

## Complete Script (Alles in Één Keer)

Kopieer en plak dit script (pas de waarden aan):

```bash
#!/bin/bash

# Configuratie
SCR_SECRET_KEY="<jouw-api-secret-key>"
API_TOKEN="Bearer <jouw-api-token>"

# 1. Docker login
echo "=== Docker Login ==="
docker login rg.nl-ams.scw.cloud -u nologin -p "$SCR_SECRET_KEY"

# 2. Build en push images
echo "=== Building and Pushing Images ==="
make push-all

# 3. ConfigMap
echo "=== Creating ConfigMap ==="
kubectl apply -f k8s/configmap.yaml

# 4. Secrets
echo "=== Creating Secrets ==="
kubectl create secret generic app-secrets \
  --from-literal=NUXT_EXTERNAL_API_AUTH="$API_TOKEN" \
  -n nl-appstore-registry || echo "Secret already exists"

kubectl create secret docker-registry scr-secret \
  --docker-server=rg.nl-ams.scw.cloud \
  --docker-username=nologin \
  --docker-password="$SCR_SECRET_KEY" \
  -n nl-appstore-registry || echo "Registry secret already exists"

# 5. Deployments
echo "=== Deploying Applications ==="
kubectl apply -f k8s/frontend-deployment.yaml
kubectl apply -f k8s/backend-deployment.yaml

# 6. Wacht op pods
echo "=== Waiting for Pods ==="
kubectl wait --for=condition=ready pod -l app=nuxt-frontend -n nl-appstore-registry --timeout=5m
kubectl wait --for=condition=ready pod -l app=go-backend -n nl-appstore-registry --timeout=5m

# 7. NodePort services
echo "=== Deploying NodePort Services ==="
make k8s-apply-nodeport

# 8. Get node IP
echo "=== Node IP Address ==="
NODE_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1)
echo "Frontend: http://$NODE_IP:30080"
echo "Backend: http://$NODE_IP:30088"

# 9. Test
echo "=== Testing ==="
curl http://$NODE_IP:30080/api/health
curl http://$NODE_IP:30088/health

echo ""
echo "✅ Setup compleet!"
echo "Frontend: http://$NODE_IP:30080"
echo "Backend: http://$NODE_IP:30088"
```

---

## Troubleshooting

### Docker Build Fails

```bash
# Check Docker draait
docker ps

# Check Docker build logs
docker build -t test-image . 2>&1 | tail -20
```

### Image Push Fails

```bash
# Verifieer login
docker login rg.nl-ams.scw.cloud

# Check registry bestaat
scw registry namespace list | grep nl-appstore-registry

# Als namespace niet bestaat, maak aan:
scw registry namespace create name=nl-appstore-registry region=nl-ams
```

### Pods Starten Niet

```bash
# Check pod status
kubectl get pods -n nl-appstore-registry

# Check logs
kubectl logs -l app=nuxt-frontend -n nl-appstore-registry
kubectl logs -l app=go-backend -n nl-appstore-registry

# Check events
kubectl get events -n nl-appstore-registry --sort-by='.lastTimestamp'
```

### ImagePullBackOff

```bash
# Check image pull secret
kubectl get secret scr-secret -n nl-appstore-registry

# Check image bestaat in registry
scw registry image list namespace-id=<namespace-id>

# Of test image pull handmatig
docker pull rg.nl-ams.scw.cloud/nl-appstore-registry/ai-co:latest
```

---

## Verificatie Checklist

Na setup:

- [ ] Images gepusht: `scw registry image list namespace-id=<id>`
- [ ] ConfigMap bestaat: `kubectl get configmap -n nl-appstore-registry`
- [ ] Secrets bestaan: `kubectl get secrets -n nl-appstore-registry`
- [ ] Pods draaien: `kubectl get pods -n nl-appstore-registry` (status: Running)
- [ ] Services bestaan: `kubectl get svc -n nl-appstore-registry`
- [ ] Node IP bekend
- [ ] Firewall rules geconfigureerd
- [ ] Health checks werken

---

## Quick Commands

```bash
# Alles checken
make k8s-status

# Logs bekijken
make k8s-logs-frontend
make k8s-logs-backend

# Pods herstarten
make k8s-rollout-restart
```

---

**Tip**: Start Docker eerst, dan kun je alles in één keer uitvoeren met het script hierboven! 🚀

