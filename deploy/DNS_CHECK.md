# DNS Check voor nlappstore.nl

## Huidige Status

**LoadBalancer IP**: `51.158.210.162`

**Probleem**: DNS records zijn nog niet actief of niet correct geconfigureerd.

Cert-manager kan de domain niet bereiken:
```
lookup nlappstore.nl: no such host
```

## DNS Records die Geconfigureerd Moeten Zijn

Bij je DNS provider (waar je nlappstore.nl hebt geregistreerd):

### A Records

| Type | Name | Value | TTL |
|------|------|-------|-----|
| A | @ (of leeg/root) | `51.158.210.162` | 300 |
| A | www | `51.158.210.162` | 300 |
| A | api | `51.158.210.162` | 300 |

### Hoe te Configureren

**Voorbeelden per DNS provider:**

#### TransIP / Openprovider / SIDN
1. Log in op je DNS beheer panel
2. Ga naar "DNS Records" of "DNS Beheer"
3. Voeg A record toe:
   - **Hostname**: `@` of leeg (voor hoofddomein)
   - **Type**: A
   - **Value**: `51.158.210.162`
   - **TTL**: 300
4. Herhaal voor `www` en `api`

#### Cloudflare
1. Log in op Cloudflare
2. Selecteer je domain `nlappstore.nl`
3. Ga naar **DNS** → **Records**
4. Voeg A records toe:
   - **Type**: A
   - **Name**: `@` (voor nlappstore.nl)
   - **IPv4 address**: `51.158.210.162`
   - **Proxy**: Uit (grijs wolkje)
   - **TTL**: Auto
5. Herhaal voor `www` en `api`

## Verificatie

### Test DNS Online

Gebruik deze tools om te checken of DNS correct is geconfigureerd:

1. **DNS Checker**: https://dnschecker.org/#A/nlappstore.nl
   - Voer `nlappstore.nl` in
   - Check of het naar `51.158.210.162` wijst wereldwijd

2. **What's My DNS**: https://www.whatsmydns.net/#A/nlappstore.nl
   - Check DNS propagation wereldwijd

3. **MXToolbox**: https://mxtoolbox.com/DNSLookup.aspx
   - Voer `nlappstore.nl` in
   - Check A record

### Test DNS Lokaal

```bash
# Test met verschillende DNS servers
dig @8.8.8.8 nlappstore.nl
dig @1.1.1.1 nlappstore.nl
dig @208.67.222.222 nlappstore.nl

# Test met nslookup
nslookup nlappstore.nl 8.8.8.8
nslookup www.nlappstore.nl 8.8.8.8
nslookup api.nlappstore.nl 8.8.8.8
```

**Verwacht resultaat**: Alle moeten naar `51.158.210.162` wijzen

## Na DNS Configuratie

1. **Wacht 5-60 minuten** voor DNS propagation
2. **Test DNS** met bovenstaande commando's
3. **Check certificaat status**:
   ```bash
   kubectl get certificate -n nl-appstore-registry
   kubectl describe certificate app-tls-secret -n nl-appstore-registry
   ```
4. **Certificaat wordt automatisch aangemaakt** (kan 2-5 minuten duren)
5. **HTTPS werkt automatisch** na certificaat

## Troubleshooting

### DNS wijst niet naar juiste IP

**Check**:
- Zijn de A records correct geconfigureerd?
- Is de TTL laag genoeg (300 seconden)?
- Heeft je DNS provider de records opgeslagen?

**Oplossing**:
- Verifieer records bij je DNS provider
- Wacht langer voor propagation (kan tot 48 uur duren, meestal 5-60 minuten)
- Test met verschillende DNS servers

### DNS wijst naar verkeerd IP

**Probleem**: DNS wijst naar ander IP dan `51.158.210.162`

**Oplossing**:
- Update A records naar `51.158.210.162`
- Wacht op propagation

### Certificaat wordt niet aangemaakt na DNS fix

**Oplossing**:
```bash
# Delete en recreate certificate
kubectl delete certificate app-tls-secret -n nl-appstore-registry
kubectl delete secret app-tls-secret -n nl-appstore-registry

# Cert-manager maakt automatisch nieuwe aan
kubectl get certificate -n nl-appstore-registry -w
```

---

## Snelle Check

```bash
# Test of DNS werkt
dig +short nlappstore.nl @8.8.8.8

# Als dit 51.158.210.162 teruggeeft → DNS werkt!
# Als dit leeg is → DNS nog niet geconfigureerd of niet gepropageerd
```

---

**Configureer DNS en test met bovenstaande commando's. Laat weten wanneer DNS werkt!** 🌐

