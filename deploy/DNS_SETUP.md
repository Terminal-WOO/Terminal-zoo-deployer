# DNS Setup voor nlappstore.nl

## LoadBalancer IP

Je ingress controller heeft het volgende externe IP adres:

**LoadBalancer IP: `51.158.210.162`**

## DNS Configuratie

Configureer de volgende DNS records bij je domain provider:

### A Records

| Type | Name | Value | TTL |
|------|------|-------|-----|
| A | @ | `51.158.210.162` | 300 |
| A | www | `51.158.210.162` | 300 |
| A | api | `51.158.210.162` | 300 |

### Via DNS Provider

1. **Log in op je DNS provider** (waar je nlappstore.nl hebt geregistreerd)
2. **Ga naar DNS Management**
3. **Voeg A records toe**:
   - **@** (of leeg) → `51.158.210.162`
   - **www** → `51.158.210.162`
   - **api** → `51.158.210.162`

### Verificatie

Na DNS configuratie (kan 5-60 minuten duren):

```bash
# Test DNS
dig nlappstore.nl
dig www.nlappstore.nl
dig api.nlappstore.nl

# Moeten allemaal naar 51.158.210.162 wijzen
```

### Check Certificaat Status

```bash
# Check certificaat wordt aangemaakt
kubectl get certificate -n nl-appstore-registry

# Check details
kubectl describe certificate app-tls-secret -n nl-appstore-registry

# Check challenges (Let's Encrypt verificatie)
kubectl get challenges -n nl-appstore-registry
```

### Na DNS Configuratie

1. **Wacht 5-60 minuten** voor DNS propagation
2. **Certificaat wordt automatisch aangemaakt** (kan 2-5 minuten duren)
3. **Test HTTPS**:
   ```bash
   curl https://nlappstore.nl/api/health
   ```

---

## Toegang na Setup

- **Frontend**: `https://nlappstore.nl`
- **Backend API**: `https://api.nlappstore.nl`
- **WWW**: `https://www.nlappstore.nl` (redirect naar hoofddomein)

---

**Configureer DNS nu, dan is HTTPS binnen 10-15 minuten actief!** 🌐🔒

