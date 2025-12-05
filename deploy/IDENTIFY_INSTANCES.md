# Instances Identificeren en Verwijderen

Deze guide helpt je om te bepalen welke instance je moet verwijderen.

## Stap 1: Check Beide Instances

### In Scaleway Console

1. **Ga naar Instances Dashboard**
   - Scaleway Console → **Instances** → **Compute**
   - Of direct: https://console.scaleway.com/instance/servers

2. **Noteer Details van Beide Instances**

Voor elke instance, noteer:
- **Naam** (bijv. "my-cluster-node-1")
- **Type** (DEV1-S, DEV1-M, etc.)
- **Status** (Running, Stopped, etc.)
- **Created Date** (wanneer aangemaakt)
- **Tags/Labels** (als aanwezig)
- **IP Address**

---

## Stap 2: Identificeer Welke Verwijderen

### Checklist: Welke Instance Verwijderen?

#### ✅ Verwijder Deze Instance Als:

- [ ] **Type is DEV1-S** (oude, kleinere instance)
- [ ] **Status is Stopped** (niet actief gebruikt)
- [ ] **Oudste creation date** (eerder aangemaakt)
- [ ] **Geen belangrijke naam** (bijv. "test", "old", "backup")
- [ ] **Geen workloads draaien** (check via kubectl als Kubernetes node)

#### ❌ Behoud Deze Instance Als:

- [ ] **Type is DEV1-M** (nieuwe, grotere instance)
- [ ] **Status is Running** (actief gebruikt)
- [ ] **Nieuwste creation date** (recent aangemaakt)
- [ ] **Belangrijke naam** (bijv. "production", "main-cluster")
- [ ] **Workloads draaien** (check via kubectl)

---

## Stap 3: Verifieer Via kubectl (Als Kubernetes Nodes)

Als beide instances Kubernetes nodes zijn:

### Check Welke Node Actief Is

```bash
# List alle nodes
kubectl get nodes -o wide

# Je ziet:
# NAME                    STATUS   ROLES    AGE   VERSION
# scw-dev1-s-xxxxx        Ready    <none>   30d   v1.31.x
# scw-dev1-m-xxxxx        Ready    <none>   2d    v1.31.x
```

### Check Welke Pods Waar Draaien

```bash
# Check pod distributie
kubectl get pods -o wide -A

# Zoek naar pods op DEV1-S node
kubectl get pods -o wide -A | grep scw-dev1-s
```

**Als er pods op DEV1-S draaien:**
1. Wacht tot ze naar DEV1-M migreren (of forceer met drain)
2. Dan kun je DEV1-S veilig verwijderen

---

## Stap 4: Veilige Verwijdering

### Optie A: Instance is Standalone (Niet in Kubernetes)

1. **Stop Instance Eerst** (veiliger)
   - Klik op instance → **Stop**
   - Wacht tot status "Stopped" is

2. **Verwijder Instance**
   - Klik op instance → **Delete**
   - Bevestig verwijdering

3. **Verwijder Volumes** (als je die niet meer nodig hebt)
   - Ga naar **Volumes** tab
   - Verwijder attached volumes

### Optie B: Instance is Kubernetes Node

#### Methode 1: Via Scaleway Console (Aanbevolen)

1. **Ga naar Kubernetes Cluster**
   - Scaleway Console → **Kubernetes** → **Kapsule**
   - Klik op je cluster

2. **Ga naar Node Pools**
   - Klik **Node Pools** tab
   - Zoek pool met DEV1-S nodes

3. **Verwijder Node Pool**
   - Klik op **3 puntjes** → **Delete**
   - Bevestig verwijdering

#### Methode 2: Via kubectl (Als je toegang hebt)

```bash
# 1. Drain node (verplaats pods naar andere nodes)
kubectl drain <dev1-s-node-name> --ignore-daemonsets --delete-emptydir-data

# 2. Verwijder node
kubectl delete node <dev1-s-node-name>

# 3. Verwijder via Scaleway Console (als node nog bestaat)
```

---

## Stap 5: Verificatie

### Na Verwijdering

1. **Check Instances Dashboard**
   - Je zou nu alleen DEV1-M moeten zien
   - Of alleen de instance die je wilt behouden

2. **Check Kubernetes** (als van toepassing)
   ```bash
   kubectl get nodes
   # Je zou alleen DEV1-M node moeten zien
   ```

3. **Check Pods Draaien Nog**
   ```bash
   kubectl get pods -n terminal-zoo
   # Alle pods moeten nog draaien (op DEV1-M)
   ```

4. **Check Kosten** (na 24-48 uur)
   - Scaleway Console → **Billing** → **Usage**
   - Kosten zouden moeten dalen met ~€15/maand

---

## Veelvoorkomende Scenario's

### Scenario 1: Twee Instances, Beide Running

**Situatie**: 
- Instance A: DEV1-S, Running, 30 dagen oud
- Instance B: DEV1-M, Running, 2 dagen oud

**Actie**:
1. ✅ Verwijder Instance A (DEV1-S)
2. ✅ Behoud Instance B (DEV1-M)

### Scenario 2: Eén Running, Eén Stopped

**Situatie**:
- Instance A: DEV1-S, Stopped, 30 dagen oud
- Instance B: DEV1-M, Running, 2 dagen oud

**Actie**:
1. ✅ Verwijder Instance A (DEV1-S, stopped)
2. ✅ Behoud Instance B (DEV1-M, running)

### Scenario 3: Beide Zelfde Type

**Situatie**:
- Instance A: DEV1-M, Running, 30 dagen oud
- Instance B: DEV1-M, Running, 2 dagen oud

**Actie**:
1. Check welke instance workloads heeft
2. Check welke instance je nieuw hebt aangemaakt
3. Verwijder de oude (meestal Instance A)

### Scenario 4: Kubernetes Cluster Nodes

**Situatie**:
- Node Pool A: DEV1-S nodes (oude pool)
- Node Pool B: DEV1-M nodes (nieuwe pool)

**Actie**:
1. ✅ Verifieer alle pods draaien op DEV1-M pool
2. ✅ Verwijder DEV1-S node pool
3. ✅ Behoud DEV1-M node pool

---

## Veiligheid Checklist

Voordat je verwijdert:

- [ ] **Backup gemaakt** (als er belangrijke data op staat)
- [ ] **Workloads gecontroleerd** (check of er pods/apps draaien)
- [ ] **Juiste instance geïdentificeerd** (DEV1-S, oudste, stopped)
- [ ] **Nieuwe instance werkt** (DEV1-M draait en werkt)
- [ ] **Geen dependencies** (geen andere services afhankelijk van deze instance)

---

## Troubleshooting

### Kan Instance Niet Verwijderen

**Probleem**: "Cannot delete: instance is in use"

**Oplossing**:
1. Stop instance eerst
2. Check of er attached volumes zijn
3. Verwijder volumes eerst (als je die niet nodig hebt)
4. Probeer opnieuw

### Instance Blijft Hangen

**Probleem**: Instance blijft in "Deleting" status

**Oplossing**:
1. Wacht 5-10 minuten (kan even duren)
2. Refresh de pagina
3. Als het langer duurt, contact Scaleway support

### Pods Draaien Nog Op DEV1-S

**Probleem**: Na verwijdering draaien pods nog op DEV1-S node

**Oplossing**:
1. Dit zou niet moeten gebeuren als je correct hebt gedrained
2. Check: `kubectl get pods -o wide`
3. Als pods nog op DEV1-S staan, forceer migratie:
   ```bash
   kubectl delete pod <pod-name> -n <namespace>
   ```

---

## Quick Reference

### Identificeer Instance Type
```bash
# Via Scaleway Console
Instances → Compute → Check "Type" kolom

# Via CLI
scw instance server list
```

### Check Kubernetes Nodes
```bash
kubectl get nodes -o wide
```

### Check Pods op Nodes
```bash
kubectl get pods -o wide -A
```

---

## Hulp Nodig?

Als je nog steeds niet zeker weet welke instance te verwijderen:

1. **Noteer beide instance namen/types**
2. **Check creation dates**
3. **Check welke actief gebruikt wordt**
4. **Verwijder de oude DEV1-S** (meestal de veilige keuze)

**Veilige regel**: Verwijder altijd de DEV1-S instance als je een DEV1-M hebt!

---

**Tip**: Maak een screenshot van beide instances voordat je verwijdert, voor referentie! 📸


