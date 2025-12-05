# Scaleway UI Troubleshooting - Lege Resource Pagina's

Als je resources ziet in de overview maar niets ziet als je erop klikt, volg deze stappen.

## Veelvoorkomende Problemen

1. **UI bug / cache probleem**
2. **Resources zijn leeg (geen node pools, geen instances, etc.)**
3. **Permissions probleem**
4. **Browser probleem**
5. **Resource is nog aan het laden**

---

## Stap 1: Quick Fixes (Probeer Dit Eerst)

### Fix 1: Hard Refresh

1. **Hard refresh de pagina**
   - Windows/Linux: `Ctrl + F5` of `Ctrl + Shift + R`
   - Mac: `Cmd + Shift + R`

2. **Clear browser cache**
   - Chrome: Settings → Privacy → Clear browsing data → Cached images
   - Firefox: Settings → Privacy → Clear Data → Cached Web Content
   - Safari: Develop → Empty Caches

### Fix 2: Andere Browser

1. **Probeer een andere browser**
   - Chrome, Firefox, Safari, Edge
   - Soms lost dit UI bugs op

### Fix 3: Log Uit en Weer In

1. **Log uit van Scaleway Console**
2. **Log weer in**
3. **Probeer opnieuw**

### Fix 4: Incognito/Private Mode

1. **Open browser in incognito/private mode**
2. **Log in op Scaleway Console**
3. **Check of resources nu zichtbaar zijn**

---

## Stap 2: Check Welke Resources Je Ziet

### In Resources Overview

Noteer wat je ziet:
- [ ] **Kubernetes Clusters** (hoeveel?)
- [ ] **Instances** (hoeveel?)
- [ ] **Container Registry** (hoeveel?)
- [ ] **Load Balancers** (hoeveel?)
- [ ] **Volumes** (hoeveel?)

### Check Elke Resource Type

#### Kubernetes Clusters

1. **Klik op een cluster**
2. **Wat zie je?**
   - Lege pagina?
   - "No node pools" message?
   - Loading spinner?
   - Error message?

**Als leeg:**
- Cluster heeft waarschijnlijk geen node pools
- Zie: **[TROUBLESHOOTING_CLUSTER.md](TROUBLESHOOTING_CLUSTER.md)**

#### Instances

1. **Klik op een instance**
2. **Wat zie je?**
   - Instance details?
   - Lege pagina?
   - Error?

**Als leeg:**
- Instance bestaat mogelijk niet meer
- Of er is een UI bug

#### Container Registry

1. **Klik op een registry**
2. **Wat zie je?**
   - Registry details?
   - Images?
   - Lege pagina?

**Als leeg:**
- Registry heeft mogelijk geen images
- Of er is een UI bug

---

## Stap 3: Via Scaleway CLI (Alternatief)

Als de UI niet werkt, gebruik de CLI:

### Installeer Scaleway CLI

```bash
# Mac
brew install scw

# Linux
curl -o /usr/local/bin/scw -L "https://github.com/scaleway/scaleway-cli/releases/latest/download/scw-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
chmod +x /usr/local/bin/scw

# Windows
# Download van: https://github.com/scaleway/scaleway-cli/releases
```

### Authenticate

```bash
scw init
# Volg de instructies
```

### Check Resources

```bash
# Kubernetes clusters
scw k8s cluster list

# Instances
scw instance server list

# Container registries
scw registry namespace list

# Volumes
scw instance volume list
```

---

## Stap 4: Specifieke Oplossingen

### Probleem: Kubernetes Cluster Leeg

**Symptomen**: Cluster zichtbaar, maar leeg dashboard

**Oplossing**:
1. Check "Node Pools" tab
2. Als leeg: maak node pool aan
3. Zie: **[TROUBLESHOOTING_CLUSTER.md](TROUBLESHOOTING_CLUSTER.md)**

### Probleem: Instance Leeg

**Symptomen**: Instance zichtbaar, maar geen details

**Oplossing**:
1. Check of instance nog bestaat: `scw instance server list`
2. Probeer direct URL: `https://console.scaleway.com/instance/servers/<instance-id>`
3. Check of instance niet is verwijderd

### Probleem: Container Registry Leeg

**Symptomen**: Registry zichtbaar, maar geen images

**Oplossing**:
1. Dit is normaal als je nog geen images hebt gepusht
2. Registry is leeg maar bestaat wel
3. Je kunt deze gebruiken voor je deployment

---

## Stap 5: Direct URL's Proberen

Soms werkt direct navigatie beter:

### Kubernetes Cluster

```
https://console.scaleway.com/k8s/clusters/<cluster-id>
```

### Instance

```
https://console.scaleway.com/instance/servers/<instance-id>
```

### Container Registry

```
https://console.scaleway.com/registry/namespaces/<namespace-id>
```

**Hoe vind je de ID?**
- In de overview lijst, hover over de resource naam
- Of check de URL als je erop klikt
- Of gebruik CLI: `scw <resource-type> list`

---

## Stap 6: Check Browser Console

Als je technisch bent:

1. **Open browser developer tools**
   - Chrome/Firefox: `F12`
   - Mac: `Cmd + Option + I`

2. **Check Console tab**
   - Zie je JavaScript errors?
   - Zie je network errors?

3. **Check Network tab**
   - Zie je failed requests?
   - Zie je 404 errors?

**Als je errors ziet:**
- Noteer de error message
- Probeer hard refresh
- Of contact Scaleway support met error details

---

## Stap 7: Check Scaleway Status

Soms is het een Scaleway probleem:

1. **Check status page**: https://status.scaleway.com/
2. **Check voor outages**
3. **Wacht even en probeer opnieuw**

---

## Praktische Aanpak Nu

### Optie A: Probeer Quick Fixes

1. ✅ Hard refresh (`Ctrl + F5`)
2. ✅ Probeer andere browser
3. ✅ Log uit en weer in
4. ✅ Probeer incognito mode

### Optie B: Gebruik CLI

Als UI niet werkt:

```bash
# Installeer Scaleway CLI
scw init

# Check resources
scw k8s cluster list
scw instance server list
scw registry namespace list
```

### Optie C: Direct URL's

1. Noteer resource ID's uit overview
2. Gebruik direct URL's (zie Stap 5 hierboven)

---

## Wat Te Doen Voor Deployment

Voor automatische deployment heb je nodig:

1. **Kubernetes Cluster ID**
   - Via CLI: `scw k8s cluster list`
   - Of in console URL: `.../clusters/<cluster-id>`

2. **Container Registry Namespace**
   - Via CLI: `scw registry namespace list`
   - Of in console: Registry → Namespace naam

3. **API Keys**
   - Scaleway Console → IAM → API Keys

**Je kunt deployment setup doen zonder de UI!**

---

## Hulp Nodig?

Als niets werkt:

1. **Gebruik Scaleway CLI** (werkt meestal beter dan UI)
2. **Contact Scaleway Support**: support@scaleway.com
3. **Check Community**: https://community.scaleway.com/

---

## Quick Checklist

- [ ] Hard refresh geprobeerd (`Ctrl + F5`)
- [ ] Andere browser geprobeerd
- [ ] Uitgelogd en weer ingelogd
- [ ] Incognito mode geprobeerd
- [ ] Scaleway CLI geïnstalleerd en geprobeerd
- [ ] Direct URL's geprobeerd
- [ ] Browser console gecheckt voor errors
- [ ] Scaleway status page gecheckt

---

**Tip**: Als de UI niet werkt, gebruik Scaleway CLI - dat werkt meestal betrouwbaarder! 🚀


