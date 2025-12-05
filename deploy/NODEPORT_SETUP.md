# NodePort Setup - Stap-voor-Stap

Complete guide voor het deployen met NodePort (zonder domain, gratis).

## Stap 1: Namespace Aanmaken (Al Gedaan!)

```bash
# Verifieer namespace bestaat
kubectl get namespace nl-appstore-registry

# Je zou moeten zien:
# NAME                    STATUS   AGE
# nl-appstore-registry    Active   Xm
```

---

## Stap 2: ConfigMap en Secrets Aanmaken

### ConfigMap

```bash
# Apply configmap
kubectl apply -f k8s/configmap.yaml

# Verifieer
kubectl get configmap app-config -n nl-appstore-registry
```

### Secrets (Als Je Die Nog Niet Hebt)

```bash
# Maak app-secrets aan (voor NUXT_EXTERNAL_API_AUTH)
kubectl create secret generic app-secrets \
  --from-literal=NUXT_EXTERNAL_API_AUTH="Bearer YOUR_TOKEN" \
  -n nl-appstore-registry

# Maak registry secret aan (voor image pull)
kubectl create secret docker-registry scr-secret \
  --docker-server=rg.nl-ams.scw.cloud \
  --docker-username=nologin \
  --docker-password=<jouw-api-secret-key> \
  -n nl-appstore-registry
```

---

## Stap 3: Deployments Aanmaken

```bash
# Apply deployments
kubectl apply -f k8s/frontend-deployment.yaml
kubectl apply -f k8s/backend-deployment.yaml

# Check status
kubectl get pods -n nl-appstore-registry

# Wacht tot pods Running zijn
kubectl wait --for=condition=ready pod -l app=nuxt-frontend -n nl-appstore-registry --timeout=5m
kubectl wait --for=condition=ready pod -l app=go-backend -n nl-appstore-registry --timeout=5m
```

---

## Stap 4: NodePort Services Deployen

```bash
# Apply NodePort services
make k8s-apply-nodeport

# Of handmatig:
kubectl apply -f k8s/frontend-service-nodeport.yaml
kubectl apply -f k8s/backend-service-nodeport.yaml

# Check services
kubectl get svc -n nl-appstore-registry
```

Je zou moeten zien:
```
NAME                 TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
nuxt-frontend-np     NodePort   10.x.x.x       <none>        3000:30080/TCP   Xs
go-backend-np        NodePort   10.x.x.x       <none>        8080:30088/TCP   Xs
```

---

## Stap 5: Node IP Adres Vinden

```bash
# Via Makefile
make k8s-get-node-ip

# Of handmatig:
kubectl get nodes -o wide

# Noteer de EXTERNAL-IP (of INTERNAL-IP als EXTERNAL-IP leeg is)
```

**Belangrijk**: Noteer het IP adres - je hebt dit nodig voor toegang!

---

## Stap 6: Firewall Rules Configureren

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

## Stap 7: Test Toegang

### Test Frontend

```bash
# Vervang <node-ip> met je node IP
curl http://<node-ip>:30080/api/health

# Verwacht: {"status":"healthy"}
```

### Test Backend

```bash
curl http://<node-ip>:30088/health

# Verwacht: {"status":"healthy"}
```

### Test in Browser

- **Frontend**: `http://<node-ip>:30080`
- **Backend**: `http://<node-ip>:30088`

---

## Stap 8: ConfigMap Updaten (Als Nodig)

Als je frontend extern moet communiceren met backend (via IP):

```bash
# Update configmap
kubectl edit configmap app-config -n nl-appstore-registry

# Verander:
# NUXT_EXTERNAL_API_BASE: "http://go-backend:8080"
# Naar:
# NUXT_EXTERNAL_API_BASE: "http://<node-ip>:30088"

# Of gebruik service naam (werkt binnen cluster):
# NUXT_EXTERNAL_API_BASE: "http://go-backend-np:8080"
```

**Aanbeveling**: Gebruik service naam (`go-backend:8080`) als beide pods in hetzelfde cluster draaien - dit werkt beter!

---

## Complete Workflow

```bash
# 1. Namespace (al gedaan)
kubectl get namespace nl-appstore-registry

# 2. ConfigMap
kubectl apply -f k8s/configmap.yaml

# 3. Secrets (als nodig)
kubectl create secret generic app-secrets \
  --from-literal=NUXT_EXTERNAL_API_AUTH="Bearer YOUR_TOKEN" \
  -n nl-appstore-registry

kubectl create secret docker-registry scr-secret \
  --docker-server=rg.nl-ams.scw.cloud \
  --docker-username=nologin \
  --docker-password=<api-secret-key> \
  -n nl-appstore-registry

# 4. Deployments
kubectl apply -f k8s/frontend-deployment.yaml
kubectl apply -f k8s/backend-deployment.yaml

# 5. Wacht op pods
kubectl wait --for=condition=ready pod -l app=nuxt-frontend -n nl-appstore-registry --timeout=5m
kubectl wait --for=condition=ready pod -l app=go-backend -n nl-appstore-registry --timeout=5m

# 6. NodePort services
make k8s-apply-nodeport

# 7. Get node IP
NODE_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1)
echo "Node IP: $NODE_IP"
echo "Frontend: http://$NODE_IP:30080"
echo "Backend: http://$NODE_IP:30088"

# 8. Test
curl http://$NODE_IP:30080/api/health
curl http://$NODE_IP:30088/health
```

---

## Troubleshooting

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

### ImagePullBackOff Error

```bash
# Check image pull secret
kubectl get secret scr-secret -n nl-appstore-registry

# Als niet bestaat, maak aan:
kubectl create secret docker-registry scr-secret \
  --docker-server=rg.nl-ams.scw.cloud \
  --docker-username=nologin \
  --docker-password=<api-secret-key> \
  -n nl-appstore-registry

# Restart pods
kubectl rollout restart deployment/nuxt-frontend -n nl-appstore-registry
kubectl rollout restart deployment/go-backend -n nl-appstore-registry
```

### Kan Niet Verbinden Met NodePort

**Probleem**: Connection refused of timeout

**Oplossing**:
1. Check firewall rules (poorten moeten open zijn)
2. Check node IP is correct: `kubectl get nodes -o wide`
3. Check service draait: `kubectl get svc -n nl-appstore-registry`
4. Check pods draaien: `kubectl get pods -n nl-appstore-registry`

### Service Heeft Geen NodePort

```bash
# Check service details
kubectl describe svc nuxt-frontend-np -n nl-appstore-registry

# Als NodePort niet zichtbaar is, check of service correct is aangemaakt
kubectl get svc nuxt-frontend-np -n nl-appstore-registry -o yaml
```

---

## Verificatie Checklist

- [ ] Namespace bestaat: `kubectl get namespace nl-appstore-registry`
- [ ] ConfigMap bestaat: `kubectl get configmap -n nl-appstore-registry`
- [ ] Secrets bestaan: `kubectl get secrets -n nl-appstore-registry`
- [ ] Pods draaien: `kubectl get pods -n nl-appstore-registry` (status: Running)
- [ ] Services bestaan: `kubectl get svc -n nl-appstore-registry` (type: NodePort)
- [ ] Node IP bekend: `kubectl get nodes -o wide`
- [ ] Firewall rules geconfigureerd (poorten 30080, 30088 open)
- [ ] Health checks werken: `curl http://<node-ip>:30080/api/health`

---

## Volgende Stappen

Na succesvolle NodePort setup:

1. ✅ **Test applicatie**: Open `http://<node-ip>:30080` in browser
2. ✅ **Monitor logs**: `make k8s-logs-frontend` en `make k8s-logs-backend`
3. ✅ **Check status**: `make k8s-status`
4. ✅ **Later domain toevoegen**: Zie `DEPLOYMENT_WITHOUT_DOMAIN.md`

---

**Tip**: Gebruik `make k8s-status` om snel alles te checken! 🚀

