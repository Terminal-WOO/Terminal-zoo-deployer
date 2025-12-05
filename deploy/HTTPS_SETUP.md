# HTTPS Setup Guide

Complete guide voor het configureren van HTTPS voor je applicatie.

## Vereisten

Voor HTTPS heb je **een van de volgende** nodig:
1. ✅ **Domain naam** (aanbevolen) - voor Let's Encrypt of Scaleway SSL
2. ⚠️ **Zelf-signed certificaat** (alleen voor testing, niet voor productie)

---

## Optie 1: HTTPS met Domain + Let's Encrypt (Gratis, Aanbevolen)

### Vereisten
- Domain naam (bijv. `app.nl-appstore-registry.nl`)
- Domain DNS moet naar je LoadBalancer IP wijzen
- Cert-manager geïnstalleerd in cluster

### Stap 1: Installeer Cert-Manager

```bash
# Installeer cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Wacht tot cert-manager klaar is
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=cert-manager -n cert-manager --timeout=5m
```

### Stap 2: Maak ClusterIssuer voor Let's Encrypt

```bash
# Maak ClusterIssuer bestand
cat > k8s/cluster-issuer.yaml <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com  # Vervang met je email
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF

# Apply ClusterIssuer
kubectl apply -f k8s/cluster-issuer.yaml
```

### Stap 3: Update Ingress met Domain

```bash
# Update ingress.yaml met je domain
# Vervang 'your-domain.com' met je echte domain
```

### Stap 4: Configureer DNS

```bash
# Get LoadBalancer IP
LB_IP=$(kubectl get svc -n nl-appstore-registry -o jsonpath='{.items[?(@.spec.type=="LoadBalancer")].status.loadBalancer.ingress[0].ip}')

# Of als je NodePort gebruikt, gebruik node IP:
NODE_IP=$(kubectl get nodes -o wide | grep -v NAME | awk '{print $7}' | head -1)

# Configureer DNS A record:
# Type: A
# Name: @ (of je subdomain)
# Value: <LB_IP of NODE_IP>
# TTL: 300
```

### Stap 5: Apply Ingress

```bash
kubectl apply -f k8s/ingress.yaml

# Check certificaat status
kubectl get certificate -n nl-appstore-registry
kubectl describe certificate -n nl-appstore-registry
```

---

## Optie 2: HTTPS met Domain + Scaleway Managed SSL

### Vereisten
- Domain naam
- Scaleway SSL certificaat

### Stap 1: Upload SSL Certificaat in Scaleway

1. Ga naar **Scaleway Console** → **Load Balancers** → **SSL Certificates**
2. Upload je certificaat (of koop via Scaleway)
3. Noteer het **Certificate ID**

### Stap 2: Update Ingress

```yaml
# k8s/ingress.yaml
annotations:
  ingress.scaleway.com/ssl-certificate-id: "your-cert-id-here"
  ingress.scaleway.com/ssl-redirect: "true"
```

### Stap 3: Apply Ingress

```bash
kubectl apply -f k8s/ingress.yaml
```

---

## Optie 3: HTTPS zonder Domain (Zelf-Signed, Alleen Testing)

⚠️ **Waarschuwing**: Zelf-signed certificaten geven browser warnings. Niet geschikt voor productie!

### Stap 1: Genereer Zelf-Signed Certificaat

```bash
# Genereer certificaat
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout tls.key \
  -out tls.crt \
  -subj "/CN=51.15.49.35" \
  -addext "subjectAltName=IP:51.15.49.35"

# Maak Kubernetes secret
kubectl create secret tls app-tls-secret \
  --cert=tls.crt \
  --key=tls.key \
  -n nl-appstore-registry
```

### Stap 2: Update Ingress voor IP

```yaml
# Maak k8s/ingress-ip.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress-ip
  namespace: nl-appstore-registry
  annotations:
    ingress.scaleway.com/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
    - secretName: app-tls-secret
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: nuxt-frontend-np
                port:
                  number: 3000
```

### Stap 3: Apply Ingress

```bash
kubectl apply -f k8s/ingress-ip.yaml
```

---

## Optie 4: HTTPS via LoadBalancer met TLS Termination

### Stap 1: Switch naar LoadBalancer Services

```bash
# Apply LoadBalancer services
make k8s-apply-loadbalancer

# Wacht op externe IPs
kubectl get svc -n nl-appstore-registry -w
```

### Stap 2: Configureer TLS in LoadBalancer

Via Scaleway Console:
1. Ga naar **Load Balancers**
2. Selecteer je Load Balancer
3. Configureer **SSL/TLS termination**
4. Upload certificaat

---

## Quick Setup: Met Domain (Aanbevolen)

Als je een domain hebt:

```bash
# 1. Installeer cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# 2. Maak ClusterIssuer (pas email aan)
kubectl apply -f k8s/cluster-issuer.yaml

# 3. Update ingress.yaml met je domain
# Vervang 'your-domain.com' met je echte domain

# 4. Configureer DNS A record naar je LoadBalancer/node IP

# 5. Apply ingress
kubectl apply -f k8s/ingress.yaml

# 6. Check certificaat
kubectl get certificate -n nl-appstore-registry
```

---

## Troubleshooting

### Certificaat wordt niet aangemaakt

```bash
# Check cert-manager logs
kubectl logs -n cert-manager -l app.kubernetes.io/name=cert-manager

# Check certificate status
kubectl describe certificate -n nl-appstore-registry

# Check challenge
kubectl get challenges -n nl-appstore-registry
```

### DNS niet correct geconfigureerd

```bash
# Test DNS
dig your-domain.com
nslookup your-domain.com

# Moet naar je LoadBalancer/node IP wijzen
```

### Ingress werkt niet

```bash
# Check ingress status
kubectl get ingress -n nl-appstore-registry
kubectl describe ingress -n nl-appstore-registry

# Check ingress controller
kubectl get pods -n ingress-nginx
```

---

## Verificatie

Na setup:

```bash
# Test HTTPS
curl -k https://your-domain.com/api/health

# Of in browser (accepteer certificaat warning als zelf-signed)
https://your-domain.com
```

---

**Heb je een domain naam beschikbaar?** Laat me weten welke, dan kan ik de ingress direct configureren! 🌐

