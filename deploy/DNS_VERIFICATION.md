# DNS Verificatie voor nlappstore.nl

## LoadBalancer IP

Je ingress controller heeft het volgende externe IP adres:

**LoadBalancer IP: `51.158.210.162`**

## DNS Records die Geconfigureerd Moeten Zijn

Bij je DNS provider (waar je nlappstore.nl hebt geregistreerd):

| Type | Name | Value | TTL |
|------|------|-------|-----|
| A | @ (of leeg) | `51.158.210.162` | 300 |
| A | www | `51.158.210.162` | 300 |
| A | api | `51.158.210.162` | 300 |

## Verificatie Commando's

### Test DNS Lokaal

```bash
# Test met dig
dig nlappstore.nl
dig www.nlappstore.nl
dig api.nlappstore.nl

# Test met nslookup
nslookup nlappstore.nl
nslookup www.nlappstore.nl
nslookup api.nlappstore.nl

# Test met Google DNS
nslookup nlappstore.nl 8.8.8.8
```

**Verwacht resultaat**: Alle moeten naar `51.158.210.162` wijzen

### Test via Online Tools

- **DNS Checker**: https://dnschecker.org/#A/nlappstore.nl
- **What's My DNS**: https://www.whatsmydns.net/#A/nlappstore.nl

Voer `nlappstore.nl` in en check of het naar `51.158.210.162` wijst.

## Certificaat Status Checken

```bash
# Check certificaat
kubectl get certificate -n nl-appstore-registry

# Check details
kubectl describe certificate app-tls-secret -n nl-appstore-registry

# Check challenges (Let's Encrypt verificatie)
kubectl get challenges -n nl-appstore-registry
kubectl describe challenge -n nl-appstore-registry
```

## Troubleshooting

### DNS wijst niet naar juiste IP

**Probleem**: DNS records zijn niet correct geconfigureerd of nog niet gepropageerd.

**Oplossing**:
1. Check DNS records bij je provider
2. Wacht 5-60 minuten voor propagation
3. Test met verschillende DNS servers (8.8.8.8, 1.1.1.1)

### Challenge blijft "pending"

**Probleem**: Let's Encrypt kan de domain niet bereiken.

**Oplossing**:
1. Verifieer DNS wijst naar `51.158.210.162`
2. Check firewall rules (poorten 80 en 443 moeten open zijn)
3. Check ingress controller draait: `kubectl get pods -n ingress-nginx`

### Certificaat wordt niet aangemaakt

**Probleem**: Cert-manager kan certificaat niet aanvragen.

**Oplossing**:
```bash
# Check cert-manager logs
kubectl logs -n cert-manager -l app.kubernetes.io/name=cert-manager

# Check ClusterIssuer
kubectl get clusterissuer letsencrypt-prod
kubectl describe clusterissuer letsencrypt-prod

# Check email in ClusterIssuer (moet geldig zijn)
```

---

## Na Succesvolle DNS Configuratie

1. **DNS propagation**: 5-60 minuten
2. **Certificaat aanvraag**: 2-5 minuten na DNS propagation
3. **HTTPS actief**: Automatisch na certificaat

**Test HTTPS**:
```bash
curl https://nlappstore.nl/api/health
```

---

**Check je DNS configuratie en laat weten als je vragen hebt!** 🌐

