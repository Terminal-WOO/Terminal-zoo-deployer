# Troubleshooting: Leeg Cluster Dashboard

Als je een cluster ziet in de overview maar niets ziet als je erop klikt, volg deze stappen.

## Mogelijke Oorzaken

1. **Cluster heeft geen node pools** (meest waarschijnlijk)
2. **Cluster is nog aan het laden**
3. **UI bug / cache probleem**
4. **Cluster is in fout status**
5. **Permissions probleem**

---

## Stap 1: Check Cluster Status

### In Scaleway Console

1. **Ga naar Kubernetes → Kapsule**
2. **Kijk naar je cluster in de lijst**
   - Check de **Status** kolom
   - Mogelijke statussen:
     - ✅ **Ready** - Cluster is klaar
     - 🟡 **Creating** - Cluster wordt nog aangemaakt
     - 🔴 **Error** - Er is een fout
     - 🟡 **Updating** - Cluster wordt geüpdatet

3. **Klik op je cluster naam**
4. **Check de Overview pagina**
   - Zie je een error message?
   - Zie je "No node pools"?
   - Zie je een loading spinner?

---

## Stap 2: Check Node Pools

### Via Scaleway Console

1. **In je cluster, klik op "Node Pools" tab**
   - Zie je een lege lijst?
   - Zie je "Create your first node pool"?
   - Zie je een error?

2. **Als er geen node pools zijn:**
   - Dit is waarschijnlijk het probleem!
   - Je moet een node pool aanmaken

### Via kubectl (Als je toegang hebt)

```bash
# Check of je toegang hebt tot het cluster
kubectl get nodes

# Als je "No resources found" ziet:
# - Cluster heeft geen nodes
# - Je moet een node pool aanmaken
```

---

## Stap 3: Oplossingen

### Oplossing 1: Node Pool Aanmaken

Als je cluster geen node pools heeft:

1. **In je cluster dashboard**
2. **Klik op "Node Pools" tab**
3. **Klik "Create Node Pool"** of **"Add Node Pool"**

4. **Configureer node pool:**
   - **Pool Name**: `default-pool`
   - **Node Type**: `DEV1-M` (4 vCPU, 8GB RAM)
   - **Node Count**: `1`
   - **Auto-scaling**: 
     - Min: `1`
     - Max: `2`
   - **Spot Instances**: ❌ Disabled

5. **Klik "Create"**
6. **Wacht 5-10 minuten** tot node is aangemaakt

### Oplossing 2: Refresh / Cache Leegmaken

1. **Hard refresh pagina**
   - Windows/Linux: `Ctrl + F5`
   - Mac: `Cmd + Shift + R`

2. **Log uit en weer in**
   - Scaleway Console → Profiel → Logout
   - Log weer in

3. **Probeer andere browser**
   - Chrome, Firefox, Safari

### Oplossing 3: Check Cluster Details

1. **In cluster overview, check:**
   - **Cluster ID** (is deze zichtbaar?)
   - **Region** (welke region?)
   - **Kubernetes Version**
   - **Status**

2. **Check Events/Logs**
   - Kijk naar "Events" of "Activity" tab
   - Zie je errors?

---

## Stap 4: Verifieer Cluster Exists

### Via Scaleway CLI

```bash
# List alle clusters
scw k8s cluster list

# Check specifiek cluster details
scw k8s cluster get <cluster-id>

# Check node pools
scw k8s pool list <cluster-id>
```

### Via kubectl

```bash
# Probeer connectie
kubectl cluster-info

# Check nodes
kubectl get nodes

# Check namespaces
kubectl get namespaces
```

---

## Veelvoorkomende Scenario's

### Scenario 1: Cluster Zonder Node Pools

**Symptomen**:
- Cluster zichtbaar in overview
- Leeg dashboard als je klikt
- "No node pools" message

**Oplossing**:
- Maak een node pool aan (zie Oplossing 1 hierboven)

### Scenario 2: Cluster Nog Aan Het Laden

**Symptomen**:
- Status is "Creating" of "Updating"
- Loading spinner zichtbaar
- Geen node pools zichtbaar

**Oplossing**:
- Wacht 5-10 minuten
- Refresh pagina
- Check status opnieuw

### Scenario 3: Cluster In Error Status

**Symptomen**:
- Status is "Error"
- Error message zichtbaar
- Rode indicator

**Oplossing**:
1. Check error message
2. Meestal: verwijder cluster en maak nieuwe aan
3. Of contact Scaleway support

### Scenario 4: Oude/Lege Cluster

**Symptomen**:
- Cluster bestaat al lang
- Geen node pools
- Geen workloads

**Oplossing**:
- Verwijder oude cluster
- Maak nieuwe cluster aan met node pool

---

## Stap 5: Nieuwe Cluster Aanmaken (Als Nodig)

Als je cluster niet werkt, maak een nieuwe aan:

### Via Scaleway Console

1. **Verwijder oude cluster eerst** (optioneel)
   - Kubernetes → Kapsule → Selecteer cluster → Delete

2. **Maak nieuwe cluster aan**
   - Kubernetes → Kapsule → Create Cluster
   - Volg: **[K8S_CLUSTER_SETUP.md](K8S_CLUSTER_SETUP.md)**

3. **Belangrijk**: Zorg dat je een **node pool aanmaakt** tijdens cluster creation!

---

## Quick Checklist

- [ ] Check cluster status (Ready/Creating/Error?)
- [ ] Check node pools tab (zijn er node pools?)
- [ ] Hard refresh pagina (Ctrl+F5 / Cmd+Shift+R)
- [ ] Probeer andere browser
- [ ] Check via kubectl: `kubectl get nodes`
- [ ] Maak node pool aan als die ontbreekt

---

## Wat Te Doen Nu?

### Optie A: Cluster Heeft Geen Node Pools

1. **Klik op "Node Pools" tab** in je cluster
2. **Klik "Create Node Pool"**
3. **Configureer**:
   - Type: DEV1-M
   - Count: 1
   - Auto-scaling: 1-2
4. **Create** en wacht 5-10 minuten

### Optie B: Cluster Werkt Niet

1. **Noteer Cluster ID** (voor referentie)
2. **Verwijder oude cluster** (als je die niet nodig hebt)
3. **Maak nieuwe cluster aan** met node pool direct

### Optie C: Cluster Is Nog Aan Het Laden

1. **Wacht 5-10 minuten**
2. **Refresh pagina**
3. **Check status opnieuw**

---

## Hulp Nodig?

Als niets werkt:

1. **Check Scaleway Status Page**: https://status.scaleway.com/
2. **Contact Scaleway Support**: Via console of support@scaleway.com
3. **Check Documentatie**: https://www.scaleway.com/en/docs/kubernetes-kapsule/

---

## Volgende Stappen

Nadat je cluster werkt met node pools:

1. ✅ **Noteer Cluster ID** (nodig voor GitHub Secrets)
2. ✅ **Verifieer nodes draaien**: `kubectl get nodes`
3. ✅ **Ga verder met**: Container Registry setup (Stap 3 in QUICK_START.md)

---

**Tip**: Meestal is het probleem dat er geen node pools zijn. Maak er een aan en het cluster zou moeten werken! 🚀


