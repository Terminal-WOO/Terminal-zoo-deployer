# HTTPS Complete Setup

Complete setup voor HTTPS op Scaleway Kubernetes.

## Vereisten Check

```bash
# Check of je een domain hebt
# Check of ingress controller bestaat
kubectl get ingressclass
```

---

## Stap 1: Installeer Ingress Controller

Scaleway Kubernetes gebruikt meestal **Traefik** of **NGINX Ingress Controller**.

### Optie A: NGINX Ingress Controller (Aanbevolen)

```bash
# Installeer NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Wacht tot ready
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=300s

# Check status
kubectl get pods -n ingress-nginx
```

### Optie B: Traefik (Als Scaleway dit gebruikt)

```bash
# Check of Traefik al bestaat
kubectl get pods --all-namespaces | grep traefik

# Als niet bestaat, installeer via Helm of Scaleway addon
```

---

## Stap 2: Kies HTTPS Methode

### Methode 1: Met Domain + Let's Encrypt (Gratis) ⭐

**Vereisten:**
- Domain naam (bijv. `app.nl-appstore-registry.nl`)
- DNS A record naar je LoadBalancer/node IP

**Setup:**

```bash
# 1. Installeer cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Wacht tot ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=cert-manager -n cert-manager --timeout=5m

# 2. Maak ClusterIssuer
cat > k8s/cluster-issuer.yaml <<'EOF'
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@example.com  # VERVANG MET JE EMAIL
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF

kubectl apply -f k8s/cluster-issuer.yaml

# 3. Update ingress.yaml met je domain (vervang 'your-domain.com')
# 4. Configureer DNS A record
# 5. Apply ingress
kubectl apply -f k8s/ingress.yaml

# 6. Check certificaat (kan 2-5 minuten duren)
kubectl get certificate -n nl-appstore-registry -w
```

### Methode 2: LoadBalancer met TLS Termination

**Vereisten:**
- LoadBalancer service
- SSL certificaat (zelf-signed of gekocht)

**Setup:**

```bash
# 1. Switch naar LoadBalancer
make k8s-apply-loadbalancer

# 2. Get LoadBalancer IP
kubectl get svc -n nl-appstore-registry

# 3. Configureer TLS in Scaleway Console:
#    - Ga naar Load Balancers
#    - Selecteer je Load Balancer
#    - SSL/TLS tab
#    - Upload certificaat of gebruik Scaleway managed SSL
```

### Methode 3: Zelf-Signed Certificaat (Alleen Testing)

**⚠️ Waarschuwing**: Browser geeft security warning!

```bash
# 1. Get node IP
NODE_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1)

# 2. Genereer certificaat
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /tmp/tls.key \
  -out /tmp/tls.crt \
  -subj "/CN=$NODE_IP" \
  -addext "subjectAltName=IP:$NODE_IP"

# 3. Maak secret
kubectl create secret tls app-tls-secret \
  --cert=/tmp/tls.crt \
  --key=/tmp/tls.key \
  -n nl-appstore-registry

# 4. Apply ingress met TLS
kubectl apply -f k8s/ingress-https-ip.yaml

# 5. Get ingress IP
kubectl get ingress -n nl-appstore-registry
```

---

## Stap 3: Verifieer HTTPS

```bash
# Test HTTPS (zelf-signed geeft warning)
curl -k https://<your-domain-or-ip>

# Check certificaat details
kubectl get certificate -n nl-appstore-registry
kubectl describe certificate -n nl-appstore-registry
```

---

## Welke Methode?

**Vraag**: Heb je een domain naam beschikbaar?

- ✅ **Ja** → Methode 1 (Let's Encrypt, gratis)
- ❌ **Nee, maar wel budget** → Methode 2 (LoadBalancer + SSL, ~€10/maand)
- ❌ **Nee, alleen testing** → Methode 3 (zelf-signed, browser warning)

---

**Laat me weten welke methode je wilt, dan configureer ik het direct!** 🔒

