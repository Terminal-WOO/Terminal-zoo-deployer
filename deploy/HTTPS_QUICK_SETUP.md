# HTTPS Quick Setup

Snelle setup voor HTTPS op je applicatie.

## Optie 1: Met Domain (Aanbevolen - Gratis)

### Stap 1: Installeer Cert-Manager

```bash
# Installeer cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Wacht tot ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=cert-manager -n cert-manager --timeout=5m
```

### Stap 2: Maak ClusterIssuer

```bash
# Maak bestand
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

# Apply
kubectl apply -f k8s/cluster-issuer.yaml
```

### Stap 3: Update Ingress met Je Domain

```bash
# Bewerk k8s/ingress.yaml
# Vervang 'your-domain.com' met je echte domain
```

### Stap 4: Configureer DNS

```bash
# Get LoadBalancer IP (als je LoadBalancer gebruikt)
kubectl get svc -n nl-appstore-registry -o jsonpath='{.items[?(@.spec.type=="LoadBalancer")].status.loadBalancer.ingress[0].ip}'

# Of node IP (als je NodePort gebruikt)
kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1

# Configureer DNS A record naar dit IP
```

### Stap 5: Apply Ingress

```bash
kubectl apply -f k8s/ingress.yaml

# Check certificaat (kan 2-5 minuten duren)
kubectl get certificate -n nl-appstore-registry -w
```

---

## Optie 2: Zonder Domain (Zelf-Signed Certificaat)

⚠️ **Waarschuwing**: Browser geeft security warning. Alleen voor testing!

### Stap 1: Genereer Certificaat

```bash
# Get node IP
NODE_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1)

# Genereer zelf-signed certificaat
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /tmp/tls.key \
  -out /tmp/tls.crt \
  -subj "/CN=$NODE_IP" \
  -addext "subjectAltName=IP:$NODE_IP,DNS:$NODE_IP"

# Maak Kubernetes secret
kubectl create secret tls app-tls-secret \
  --cert=/tmp/tls.crt \
  --key=/tmp/tls.key \
  -n nl-appstore-registry \
  --dry-run=client -o yaml | kubectl apply -f -
```

### Stap 2: Check Ingress Controller

```bash
# Check of ingress controller bestaat
kubectl get ingressclass

# Als niet bestaat, installeer nginx ingress controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml
```

### Stap 3: Apply HTTPS Ingress

```bash
kubectl apply -f k8s/ingress-https-ip.yaml

# Check ingress
kubectl get ingress -n nl-appstore-registry
```

### Stap 4: Get HTTPS URL

```bash
# Get ingress IP
INGRESS_IP=$(kubectl get ingress app-ingress-https -n nl-appstore-registry -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

# Of als leeg, gebruik node IP
if [ -z "$INGRESS_IP" ]; then
  INGRESS_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1)
fi

echo "HTTPS URL: https://$INGRESS_IP"
echo "⚠️ Browser zal security warning geven - accepteer certificaat"
```

---

## Optie 3: LoadBalancer met TLS Termination

### Stap 1: Switch naar LoadBalancer

```bash
# Apply LoadBalancer services
make k8s-apply-loadbalancer

# Wacht op externe IPs
kubectl get svc -n nl-appstore-registry -w
```

### Stap 2: Configureer TLS in Scaleway Console

1. Ga naar **Load Balancers** in Scaleway Console
2. Selecteer je Load Balancer
3. Ga naar **SSL/TLS** tab
4. Upload certificaat of gebruik Scaleway managed SSL

---

## Welke Optie Kies Je?

- **Heb je een domain?** → Optie 1 (gratis, Let's Encrypt)
- **Geen domain, maar wel budget?** → Optie 3 (LoadBalancer + SSL)
- **Alleen testing?** → Optie 2 (zelf-signed, browser warning)

---

**Laat me weten welke optie je wilt, dan configureer ik het direct!** 🔒

