# Deployment Zonder Domeinnaam - Alleen IP Adres

Deze guide legt uit hoe je de applicatie deployt zonder custom domain, alleen met IP adres.

## Opties Zonder Domain

Je hebt 3 opties:

1. **NodePort** (Goedkoopst - €0/maand) ✅ Aanbevolen voor start
2. **LoadBalancer** (Duurder - ~€10/maand, maar makkelijker)
3. **Port Forward** (Alleen voor testing - gratis)

---

## Optie 1: NodePort (Aanbevolen - Goedkoopst)

### Voordelen
- ✅ **Gratis** (geen extra kosten)
- ✅ Eenvoudig te gebruiken
- ✅ Direct toegankelijk via node IP

### Nadelen
- ⚠️ Toegang via `<node-ip>:<poort>` (minder gebruiksvriendelijk)
- ⚠️ Poort moet open zijn in firewall

### Setup

1. **Vind je Node IP Adres**
   ```bash
   # Via kubectl
   kubectl get nodes -o wide
   # Noteer de EXTERNAL-IP kolom
   
   # Of via Scaleway Console
   # Kubernetes → Je cluster → Nodes → Noteer IP adres
   ```

2. **Apply NodePort Services**
   ```bash
   kubectl apply -f k8s/frontend-service-nodeport.yaml
   kubectl apply -f k8s/backend-service-nodeport.yaml
   ```

3. **Check Services**
   ```bash
   kubectl get svc -n terminal-zoo
   # Je ziet:
   # nuxt-frontend-np   NodePort   10.x.x.x   3000:30080/TCP
   # go-backend-np      NodePort   10.x.x.x   8080:30088/TCP
   ```

4. **Toegang tot Applicatie**
   - **Frontend**: `http://<node-ip>:30080`
   - **Backend**: `http://<node-ip>:30088`

### Firewall Rules

Zorg dat poorten open zijn:

```bash
# Via Scaleway Console:
# Instances → Security Groups → Je security group
# Voeg rules toe:
# - Inbound TCP 30080 (voor frontend)
# - Inbound TCP 30088 (voor backend)
```

---

## Optie 2: LoadBalancer (Makkelijker, Maar Duurder)

### Voordelen
- ✅ **Makkelijker** - Krijgt automatisch extern IP
- ✅ Geen poort nummers nodig
- ✅ Professionele setup

### Nadelen
- ⚠️ **Kosten**: ~€10/maand per Load Balancer
- ⚠️ Je hebt 2 Load Balancers nodig (frontend + backend) = ~€20/maand

### Setup

1. **Apply LoadBalancer Services**
   ```bash
   kubectl apply -f k8s/frontend-service-loadbalancer.yaml
   kubectl apply -f k8s/backend-service-loadbalancer.yaml
   ```

2. **Wacht op Extern IP**
   ```bash
   kubectl get svc -n terminal-zoo -w
   # Wacht tot EXTERNAL-IP kolom een IP adres krijgt
   ```

3. **Noteer IP Adressen**
   ```bash
   kubectl get svc -n terminal-zoo
   # Je ziet:
   # nuxt-frontend-lb   LoadBalancer   10.x.x.x   <extern-ip>   80/TCP
   # go-backend-lb      LoadBalancer   10.x.x.x   <extern-ip>   8080/TCP
   ```

4. **Toegang tot Applicatie**
   - **Frontend**: `http://<extern-ip-frontend>`
   - **Backend**: `http://<extern-ip-backend>`

### Kosten

- Frontend Load Balancer: ~€10/maand
- Backend Load Balancer: ~€10/maand
- **Totaal**: ~€20/maand extra

---

## Optie 3: Port Forward (Alleen Testing)

### Voordelen
- ✅ **Gratis**
- ✅ Geen configuratie nodig
- ✅ Perfect voor lokale testing

### Nadelen
- ⚠️ Alleen lokaal toegankelijk
- ⚠️ Niet voor productie

### Setup

```bash
# Frontend port forward
kubectl port-forward -n terminal-zoo svc/nuxt-frontend 3000:3000

# In andere terminal, backend port forward
kubectl port-forward -n terminal-zoo svc/go-backend 8080:8080

# Toegang:
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
```

---

## Aanbevolen Setup: NodePort

Voor kostenbesparing, gebruik **NodePort**:

### Stap 1: Update Kubernetes Manifests

Gebruik NodePort services in plaats van ingress:

```bash
# Apply NodePort services
kubectl apply -f k8s/frontend-service-nodeport.yaml
kubectl apply -f k8s/backend-service-nodeport.yaml

# Verwijder ingress (niet nodig zonder domain)
# kubectl delete -f k8s/ingress.yaml
```

### Stap 2: Vind Node IP

```bash
# Via kubectl
kubectl get nodes -o wide

# Of via Scaleway Console
# Kubernetes → Je cluster → Nodes → Noteer IP
```

### Stap 3: Configureer Firewall

```bash
# Via Scaleway Console:
# Instances → Security Groups → Je security group
# Voeg inbound rules toe:
# - TCP 30080 (frontend)
# - TCP 30088 (backend)
```

### Stap 4: Test Toegang

```bash
# Test frontend
curl http://<node-ip>:30080/api/health

# Test backend
curl http://<node-ip>:30088/health
```

---

## Update Frontend Configuratie

Als je backend via IP adres bereikbaar is, update `k8s/configmap.yaml`:

```yaml
data:
  # Frontend environment variables
  NUXT_EXTERNAL_API_BASE: "http://<node-ip>:30088"  # Of LoadBalancer IP
  # Of gebruik service naam als binnen cluster:
  # NUXT_EXTERNAL_API_BASE: "http://go-backend:8080"
```

---

## Makefile Updates

Voeg commando's toe aan Makefile:

```makefile
# Get node IP
k8s-get-node-ip:
	@kubectl get nodes -o wide | grep -v NAME | awk '{print $$7}' | head -1

# Get LoadBalancer IPs
k8s-get-lb-ips:
	@echo "Frontend:"
	@kubectl get svc nuxt-frontend-lb -n terminal-zoo -o jsonpath='{.status.loadBalancer.ingress[0].ip}' || echo "N/A"
	@echo "Backend:"
	@kubectl get svc go-backend-lb -n terminal-zoo -o jsonpath='{.status.loadBalancer.ingress[0].ip}' || echo "N/A"
```

---

## Kosten Vergelijking

| Optie | Maandelijkse Kosten | Moeilijkheid |
|-------|-------------------|--------------|
| **NodePort** | €0 | Eenvoudig |
| **LoadBalancer** | ~€20 (2x €10) | Zeer eenvoudig |
| **Port Forward** | €0 | Alleen testing |

**Aanbeveling**: Start met **NodePort** (gratis), upgrade naar LoadBalancer later als nodig.

---

## Later Domain Toevoegen

Als je later een domain krijgt:

1. **Koop domain** (bijv. via Namecheap, Cloudflare, etc.)
2. **Configureer DNS**:
   ```
   A record: @ → <node-ip> of <loadbalancer-ip>
   ```
3. **Update ingress.yaml** met je domain
4. **Apply ingress**:
   ```bash
   kubectl apply -f k8s/ingress.yaml
   ```

---

## Troubleshooting

### Kan Niet Verbinden Met NodePort

**Probleem**: `Connection refused` of timeout

**Oplossing**:
1. Check firewall rules (poorten moeten open zijn)
2. Check node IP is correct
3. Check service draait: `kubectl get svc -n terminal-zoo`
4. Check pods draaien: `kubectl get pods -n terminal-zoo`

### LoadBalancer Krijgt Geen IP

**Probleem**: EXTERNAL-IP blijft `<pending>`

**Oplossing**:
1. Wacht 2-5 minuten (kan even duren)
2. Check Scaleway Console → Load Balancers
3. Check of je quota hebt voor Load Balancers
4. Check service configuratie: `kubectl describe svc -n terminal-zoo`

---

## Quick Start: NodePort Setup

```bash
# 1. Apply NodePort services
kubectl apply -f k8s/frontend-service-nodeport.yaml
kubectl apply -f k8s/backend-service-nodeport.yaml

# 2. Get node IP
NODE_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1)
echo "Node IP: $NODE_IP"

# 3. Test toegang
curl http://$NODE_IP:30080/api/health
curl http://$NODE_IP:30088/health

# 4. Open in browser
echo "Frontend: http://$NODE_IP:30080"
echo "Backend: http://$NODE_IP:30088"
```

---

**Tip**: Start met NodePort (gratis), upgrade naar LoadBalancer of domain later als je budget hebt! 💰


