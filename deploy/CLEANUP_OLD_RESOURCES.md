# Oude Resources Verwijderen - Kostenbesparing

Deze guide helpt je oude DEV1-S instances en andere ongebruikte resources te verwijderen om kosten te besparen.

> **💡 Eerst**: Als je meerdere instances ziet en niet weet welke te verwijderen, zie **[IDENTIFY_INSTANCES.md](IDENTIFY_INSTANCES.md)** voor hulp bij identificatie.

## Wat Verwijderen?

### 1. Oude Kubernetes Nodes (DEV1-S)
Als je een oude DEV1-S node hebt in je Kubernetes cluster

### 2. Standalone Instances (DEV1-S)
Als je een standalone DEV1-S instance hebt (niet in Kubernetes)

### 3. Oude Clusters
Als je een oud Kubernetes cluster hebt die je niet meer gebruikt

---

## Methode 1: Kubernetes Node Verwijderen

### Via Scaleway Console

1. **Ga naar je Kubernetes Cluster**
   - Scaleway Console → **Kubernetes** → **Kapsule**
   - Klik op je cluster naam

2. **Ga naar Node Pools**
   - Klik op **Node Pools** tab
   - Je ziet alle node pools

3. **Verwijder Node Pool met DEV1-S**
   - Zoek de pool met DEV1-S nodes
   - Klik op de **3 puntjes** (menu) → **Delete**
   - Bevestig de verwijdering

**⚠️ Let op**: 
- Als dit je enige node pool is, moet je eerst een nieuwe DEV1-M pool aanmaken
- Of verwijder het hele cluster als je het niet meer nodig hebt

### Via kubectl (Als je nog toegang hebt)

```bash
# Check nodes
kubectl get nodes

# Drain node (verplaats pods naar andere nodes)
kubectl drain <node-name> --ignore-daemonsets --delete-emptydir-data

# Verwijder node
kubectl delete node <node-name>
```

---

## Methode 2: Standalone Instance Verwijderen

Als je een standalone DEV1-S instance hebt (niet in Kubernetes):

### Via Scaleway Console

1. **Ga naar Instances**
   - Scaleway Console → **Instances** → **Compute**
   - Of direct: https://console.scaleway.com/instance/servers

2. **Zoek je DEV1-S Instance**
   - Filter op naam of type
   - Zoek naar instances met type "DEV1-S"

3. **Verwijder Instance**
   - Klik op je instance
   - Klik op **Delete** (of **More** → **Delete**)
   - Bevestig verwijdering
   - Optioneel: verwijder ook de attached volumes

**⚠️ Let op**: 
- Backup belangrijke data eerst!
- Verwijder ook attached volumes als je die niet meer nodig hebt

### Via Scaleway CLI

```bash
# List instances
scw instance server list

# Stop instance eerst (optioneel)
scw instance server stop <instance-id>

# Verwijder instance
scw instance server delete <instance-id>

# Verwijder volumes (als je die niet meer nodig hebt)
scw instance volume delete <volume-id>
```

---

## Methode 3: Hele Kubernetes Cluster Verwijderen

Als je een oud cluster hebt die je niet meer gebruikt:

### Via Scaleway Console

1. **Ga naar Kubernetes**
   - Scaleway Console → **Kubernetes** → **Kapsule**

2. **Selecteer Cluster**
   - Klik op je oude cluster

3. **Verwijder Cluster**
   - Klik op **Delete** (meestal rechtsboven)
   - Typ de cluster naam ter bevestiging
   - Klik **Delete Cluster**

**⚠️ Let op**: 
- Dit verwijdert ALLES: nodes, workloads, data
- Backup belangrijke data eerst!
- Dit kan niet ongedaan gemaakt worden

### Via Scaleway CLI

```bash
# List clusters
scw k8s cluster list

# Verwijder cluster
scw k8s cluster delete <cluster-id>
```

---

## Methode 4: Bulk Cleanup - Alle Ongebruikte Resources

### Check Alle Resources

1. **Instances**
   - Scaleway Console → **Instances** → **Compute**
   - Check welke instances je niet meer gebruikt

2. **Kubernetes Clusters**
   - Scaleway Console → **Kubernetes** → **Kapsule**
   - Check welke clusters je niet meer gebruikt

3. **Volumes**
   - Scaleway Console → **Instances** → **Volumes**
   - Check welke volumes niet meer attached zijn

4. **Container Registry**
   - Scaleway Console → **Container Registry**
   - Verwijder oude images om storage kosten te besparen

5. **Load Balancers**
   - Scaleway Console → **Load Balancers**
   - Check welke je niet meer gebruikt (~€10/maand per LB)

6. **IP Addresses**
   - Scaleway Console → **IP Addresses**
   - Verwijder ongebruikte IPs (~€0.01/maand per IP)

---

## Kosten Check

### Voor Verwijdering

```bash
# Check huidige kosten
# Scaleway Console → Billing → Usage
```

### Na Verwijdering

Na het verwijderen van DEV1-S instance:
- **Besparing**: ~€15/maand (als je DEV1-S had)
- **Nieuwe kosten**: ~€30/maand (voor DEV1-M cluster)

**Netto verschil**: +€15/maand, maar je krijgt:
- ✅ Meer resources (4 vCPU / 8GB vs 2 vCPU / 4GB)
- ✅ Betere performance
- ✅ Meer headroom voor groei

---

## Veilige Verwijdering Checklist

Voordat je iets verwijdert:

- [ ] **Backup belangrijke data**
- [ ] **Check of resource nog gebruikt wordt**
- [ ] **Verifieer dat er geen dependencies zijn**
- [ ] **Noteer wat je verwijdert** (voor referentie)
- [ ] **Check kosten na verwijdering**

---

## Troubleshooting

### Kan Node Pool Niet Verwijderen

**Probleem**: "Cannot delete node pool: nodes are in use"

**Oplossing**:
1. Verwijder eerst alle workloads van die nodes
2. Of maak een nieuwe pool aan en migreer workloads
3. Dan kun je de oude pool verwijderen

### Instance Staat Op "Deleting"

**Probleem**: Instance blijft hangen in "Deleting" status

**Oplossing**:
1. Wacht 5-10 minuten (kan even duren)
2. Check of er attached volumes zijn die eerst verwijderd moeten worden
3. Contact Scaleway support als het langer duurt

### Cluster Kan Niet Verwijderd Worden

**Probleem**: "Cluster has active resources"

**Oplossing**:
1. Verwijder eerst alle workloads: `kubectl delete all --all -n <namespace>`
2. Verwijder namespaces: `kubectl delete namespace <namespace>`
3. Dan kun je het cluster verwijderen

---

## Quick Commands

### Check Resources

```bash
# Kubernetes clusters
scw k8s cluster list

# Instances
scw instance server list

# Volumes
scw instance volume list

# Load balancers
scw lb lb list
```

### Verwijder Resources

```bash
# Verwijder instance
scw instance server delete <instance-id>

# Verwijder volume
scw instance volume delete <volume-id>

# Verwijder cluster
scw k8s cluster delete <cluster-id>
```

---

## Na Verwijdering

### Verifieer Kosten

1. Ga naar **Scaleway Console** → **Billing** → **Usage**
2. Check of kosten zijn gedaald
3. Wacht 24-48 uur voor accurate kosten weergave

### Setup Nieuwe Cluster

Als je een nieuwe DEV1-M cluster nodig hebt:
- Volg: **[K8S_CLUSTER_SETUP.md](K8S_CLUSTER_SETUP.md)**

---

## Hulp Nodig?

- **Scaleway Support**: Via console of support@scaleway.com
- **Documentatie**: https://www.scaleway.com/en/docs/
- **Community**: https://community.scaleway.com/

---

**Tip**: Verwijder regelmatig ongebruikte resources om kosten te besparen! 💰

