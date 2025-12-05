# Kubernetes Cluster Setup - Kostenbewust

Deze guide helpt je een minimalistisch en goedkoop Kubernetes cluster op te zetten op Scaleway.

## Kosten Overzicht

### Minimale Setup (Beschikbaar)
- **Node Type**: DEV1-M (4 vCPU, 8GB RAM)
- **Aantal Nodes**: 1
- **Geschatte Kosten**: ~€30/maand
- **Beschikbaarheid**: Meestal beschikbaar in alle regio's

### Waarom DEV1-M Goed Is
- ✅ Meer headroom voor toekomstige groei
- ✅ Betere performance voor meerdere pods
- ✅ Ruimte voor monitoring/logging tools
- ✅ Nog steeds kostenefficiënt (~€30/maand)

### Alternatieven (indien beschikbaar)
- **DEV1-S** (2 vCPU, 4GB RAM): ~€15/maand (niet altijd beschikbaar)
- **DEV1-L** (8 vCPU, 16GB RAM): ~€60/maand (niet nodig voor start)

## Stap-voor-Stap: Cluster Aanmaken

### Stap 1: Ga naar Kubernetes Kapsule

1. Log in op https://console.scaleway.com
2. Ga naar **Kubernetes** → **Kapsule**
3. Klik op **Create Cluster**

### Stap 2: Cluster Basis Configuratie

Vul het formulier in:

**Cluster Name**: 
```
terminal-zoo-cluster
```
(of een andere naam die je wilt)

**Region**: 
```
Amsterdam (nl-ams)
```
(meestal goedkoper dan andere regio's)

**Kubernetes Version**: 
```
1.31.x (latest stable)
```
(gebruik de laatste stabiele versie)

**CNI Plugin**: 
```
Cilium
```
(standaard, goed voor kostenbewuste setup)

**Ingress Controller**: 
```
NGINX
```
(standaard, nodig voor ingress)

Klik **Next**

### Stap 3: Node Pool Configuratie (BELANGRIJK!)

Dit is waar je kosten bespaart. Configureer als volgt:

#### Pool 1: Default Pool (Minimaal)

**Pool Name**: 
```
default-pool
```

**Node Type**: 
```
DEV1-M
```
- 4 vCPU
- 8GB RAM
- **~€30/maand per node**
- ✅ Meestal beschikbaar als minimum

**Node Count**: 
```
1
```
(Start met 1 node, je kunt later opschalen)

**Auto-scaling**: 
- ✅ **Enable Auto-scaling**
- **Min Nodes**: `1`
- **Max Nodes**: `2`
(Dit voorkomt onverwachte kosten door automatisch opschalen)

**Spot Instances**: 
```
❌ Disabled
```
(Voor productie, gebruik geen spot instances - die kunnen worden gestopt)

**Container Runtime**: 
```
containerd
```
(standaard, goedkoper dan Docker)

Klik **Next**

### Stap 4: Review en Aanmaken

1. **Review je configuratie**:
   - ✅ 1x DEV1-M node
   - ✅ Auto-scaling: 1-2 nodes
   - ✅ Totale kosten: ~€30-60/maand (afhankelijk van scaling)

2. **Accept Terms** (vink aan)

3. Klik **Create Cluster**

### Stap 5: Wacht op Cluster Creation

- Dit duurt ongeveer **5-10 minuten**
- Je ziet de status: "Creating..." → "Ready"
- **Noteer de Cluster ID** die wordt getoond (je hebt deze nodig voor GitHub Secrets)

---

## Cluster ID Vinden

Na het aanmaken:

1. Klik op je cluster naam
2. Op de **Overview** pagina, zoek naar **Cluster ID**
3. Kopieer de UUID (bijv. `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`)
4. **Bewaar deze** - je hebt hem nodig voor GitHub Secrets

**Of via CLI:**
```bash
scw k8s cluster list
# Zoek je cluster en kopieer de ID
```

---

## Kosten Optimalisatie Tips

### 1. Start Klein
- Begin met **1 node** (DEV1-S)
- Je applicatie past hier ruimschoots op
- Schaal alleen op als je echt meer resources nodig hebt

### 2. Auto-Scaling Configuratie
```
Min: 1 node  (altijd minimaal 1 voor beschikbaarheid)
Max: 2 nodes (voorkomt onverwachte kosten)
```

### 3. Monitor Resource Gebruik
Na deployment, check regelmatig:
```bash
kubectl top nodes
kubectl top pods -n terminal-zoo
```

Als gebruik consistent laag is (<30%), kun je zelfs overwegen:
- Kleinere node type (als beschikbaar)
- Of blijf bij 1 node

### 4. Gebruik Resource Limits
De applicatie is al geconfigureerd met lage resource limits:
- Frontend: 100m CPU / 128Mi RAM request
- Backend: 50m CPU / 64Mi RAM request

Dit past ruimschoots op een DEV1-S node.

### 5. Schakel Onnodige Services Uit
- ❌ Monitoring tools (indien niet nodig)
- ❌ Extra ingress controllers
- ❌ Log aggregation (gebruik kubectl logs)

---

## Kosten Breakdown

### Maandelijkse Kosten (Met DEV1-M)

| Component | Kosten |
|-----------|--------|
| Kubernetes Cluster (1x DEV1-M) | ~€30 |
| Container Registry (10GB) | ~€1-2 |
| Load Balancer (optioneel) | ~€10 |
| **Totaal (zonder LB)** | **~€31-32** |
| **Totaal (met LB)** | **~€41-42** |

### Wat Past op DEV1-M?

Met 1x DEV1-M node (4 vCPU, 8GB RAM):
- ✅ Je frontend pod (100m CPU / 128Mi RAM)
- ✅ Je backend pod (50m CPU / 64Mi RAM)
- ✅ System overhead (~500m CPU / 1GB RAM)
- ✅ Ruimte voor scaling (nog ~3.3 CPU / 6.8GB RAM over)
- ✅ Ruimte voor monitoring/logging tools
- ✅ Ruimte voor meerdere replica's indien nodig

**Conclusie**: Je hebt veel ruimte! Perfect voor groei.

---

## Verificatie

Na het aanmaken van je cluster:

### 1. Download kubeconfig
```bash
# Via Scaleway Console:
# Kubernetes → Je cluster → Download kubeconfig

# Of via CLI:
scw k8s kubeconfig install <cluster-id>
```

### 2. Test Connectie
```bash
kubectl get nodes

# Je zou moeten zien:
# NAME                    STATUS   ROLES    AGE   VERSION
# scw-xxxxx               Ready    <none>   5m    v1.31.x
```

### 3. Check Node Resources
```bash
kubectl describe node

# Check:
# - CPU: 2 cores beschikbaar
# - Memory: ~4GB beschikbaar
```

---

## Troubleshooting

### Cluster Creation Fails
- ✅ Check of je voldoende quota hebt in Scaleway
- ✅ Verifieer je payment method is geconfigureerd
- ✅ Probeer een andere region (Amsterdam is meestal beschikbaar)

### Node Start Niet
- ✅ Check node status in Scaleway Console
- ✅ Verifieer auto-scaling settings
- ✅ Check of er voldoende resources zijn in de region

### Hoge Kosten
- ✅ Check aantal nodes: `kubectl get nodes`
- ✅ Verifieer auto-scaling max is ingesteld op 2
- ✅ Monitor resource usage regelmatig
- ✅ Overweeg kleinere node type als gebruik laag is

---

## Volgende Stappen

Na het aanmaken van je cluster:

1. ✅ **Noteer Cluster ID** (nodig voor GitHub Secrets)
2. ✅ **Download kubeconfig** (voor lokale testing)
3. ✅ **Verifieer connectie**: `kubectl get nodes`
4. ➡️ **Ga verder met**: Container Registry setup (Stap 3)

---

## FAQ

**Q: Kan ik later opschalen?**
A: Ja! Je kunt altijd nodes toevoegen of een grotere node type kiezen. Start klein en schaal op indien nodig.

**Q: Wat als 1 node niet genoeg is?**
A: Auto-scaling voegt automatisch een 2e node toe als resources opraken. Je kunt ook handmatig opschalen.

**Q: Kan ik kosten nog verder verlagen?**
A: 
- Gebruik spot instances (niet aanbevolen voor productie)
- Schakel monitoring/logging uit
- Gebruik geen Load Balancer (gebruik NodePort)
- Overweeg serverless voor minder kritieke componenten

**Q: Wat als ik per ongeluk te veel nodes heb?**
A: Auto-scaling schaalt automatisch af naar 1 node als resources niet meer nodig zijn. Je kunt ook handmatig nodes verwijderen.

---

**Gefeliciteerd! 🎉** Je hebt nu een kostenefficiënt Kubernetes cluster!

**Geschatte maandelijkse kosten**: ~€31-32 (zonder Load Balancer)

**Voordeel van DEV1-M**: Meer headroom voor toekomstige groei zonder direct op te hoeven schalen!

