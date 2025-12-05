# Scaleway CLI Configuratie - Zone/Region Aanpassen

Als je resources niet ziet, kan het zijn dat je default zone/region niet klopt. Hier is hoe je dit aanpast.

## Probleem

Je resources staan in een andere zone/region dan je default zone, waardoor je ze niet ziet.

## Oplossing: Default Zone Aanpassen

### Methode 1: Via `scw config`

```bash
# Check huidige configuratie
scw config get

# Je ziet waarschijnlijk:
# default_zone: nl-ams-1
# default_region: nl-ams

# Pas default zone aan
scw config set default_zone=nl-ams-1
scw config set default_region=nl-ams

# Of voor andere zones:
# nl-ams-1 (Amsterdam)
# fr-par-1 (Paris)
# fr-par-2 (Paris)
# pl-waw-1 (Warsaw)
```

### Methode 2: Via Configuratie Bestand

Het configuratie bestand staat op:
- **Mac/Linux**: `~/.config/scw/config.yaml`
- **Windows**: `%USERPROFILE%\.config\scw\config.yaml`

Bewerk het bestand:

```yaml
# ~/.config/scw/config.yaml
access_key: SCWXXXXXXXXXXXXXXXXX
secret_key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
default_organization_id: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
default_project_id: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
default_zone: nl-ams-1          # Pas dit aan
default_region: nl-ams          # Pas dit aan
```

### Methode 3: Via `scw init` Opnieuw

```bash
# Herinitialiseer configuratie
scw init

# Kies de juiste zone/region wanneer gevraagd
# Meestal: nl-ams-1 (Amsterdam)
```

---

## Beschikbare Zones/Regions

### Amsterdam (Aanbevolen voor Nederland)
- **Zone**: `nl-ams-1`
- **Region**: `nl-ams`
- **Meestal goedkoper**

### Paris
- **Zone**: `fr-par-1` of `fr-par-2`
- **Region**: `fr-par`

### Warsaw
- **Zone**: `pl-waw-1`
- **Region**: `pl-waw`

---

## Check Welke Zone Je Resources Hebben

### Via Scaleway Console

1. Ga naar je resource (bijv. Kubernetes cluster)
2. Check de **Region** kolom
3. Noteer de region (bijv. "Amsterdam", "Paris")

### Via CLI

```bash
# Kubernetes clusters (toont region)
scw k8s cluster list

# Instances (toont zone)
scw instance server list

# Container registries (toont region)
scw registry namespace list
```

---

## Stap-voor-Stap: Zone Aanpassen

### Stap 1: Check Huidige Config

```bash
scw config get
```

Noteer:
- `default_zone`: ?
- `default_region`: ?

### Stap 2: Check Waar Je Resources Zijn

```bash
# Probeer resources te zien (zonder zone flag)
scw k8s cluster list

# Als leeg, probeer met expliciete zone
scw k8s cluster list zone=nl-ams-1
scw k8s cluster list zone=fr-par-1
```

### Stap 3: Pas Default Zone Aan

```bash
# Voor Amsterdam (meestal gebruikt)
scw config set default_zone=nl-ams-1
scw config set default_region=nl-ams

# Verifieer
scw config get
```

### Stap 4: Test

```bash
# Check of je resources nu ziet
scw k8s cluster list
scw instance server list
scw registry namespace list
```

---

## Voor Kubernetes Clusters

Kubernetes clusters hebben een **region** (niet zone):

```bash
# List clusters in specifieke region
scw k8s cluster list region=nl-ams
scw k8s cluster list region=fr-par

# Of zonder flag (gebruikt default region)
scw k8s cluster list
```

---

## Voor Container Registry

Container registries hebben ook een **region**:

```bash
# List registries in specifieke region
scw registry namespace list region=nl-ams
scw registry namespace list region=fr-par
```

---

## Voor Instances

Instances hebben een **zone**:

```bash
# List instances in specifieke zone
scw instance server list zone=nl-ams-1
scw instance server list zone=fr-par-1
```

---

## Permanente Oplossing

### Option 1: Set Default Zone

```bash
# Set voor alle commando's
scw config set default_zone=nl-ams-1
scw config set default_region=nl-ams
```

### Option 2: Gebruik Zone Flag in Commando's

```bash
# Gebruik altijd expliciete zone
scw k8s cluster list region=nl-ams
scw instance server list zone=nl-ams-1
```

### Option 3: Environment Variables

```bash
# Set environment variables
export SCW_DEFAULT_ZONE=nl-ams-1
export SCW_DEFAULT_REGION=nl-ams

# Of voeg toe aan ~/.bashrc of ~/.zshrc
echo 'export SCW_DEFAULT_ZONE=nl-ams-1' >> ~/.zshrc
echo 'export SCW_DEFAULT_REGION=nl-ams' >> ~/.zshrc
source ~/.zshrc
```

---

## Troubleshooting

### Probleem: Resources Nog Steeds Niet Zichtbaar

**Oplossing**:
1. Check of je de juiste zone hebt:
   ```bash
   scw config get
   ```

2. Probeer expliciete zone:
   ```bash
   scw k8s cluster list region=nl-ams
   scw k8s cluster list region=fr-par
   ```

3. Check alle regions:
   ```bash
   # Probeer alle beschikbare regions
   for region in nl-ams fr-par pl-waw; do
     echo "=== $region ==="
     scw k8s cluster list region=$region
   done
   ```

### Probleem: Config Bestand Niet Gevonden

**Oplossing**:
```bash
# Maak configuratie opnieuw
scw init

# Kies de juiste zone/region
```

### Probleem: Meerdere Projects

Als je meerdere projects hebt:

```bash
# Check huidige project
scw config get default_project_id

# List alle projects
scw account project list

# Switch project
scw config set default_project_id=<project-id>
```

---

## Quick Fix Commando

Voor Amsterdam (meest gebruikt):

```bash
# Set default zone en region
scw config set default_zone=nl-ams-1
scw config set default_region=nl-ams

# Verifieer
scw config get

# Test
scw k8s cluster list
scw instance server list
scw registry namespace list
```

---

## Configuratie Bestand Locatie

### Mac/Linux
```bash
# Config bestand
~/.config/scw/config.yaml

# Bekijk configuratie
cat ~/.config/scw/config.yaml

# Bewerk configuratie
nano ~/.config/scw/config.yaml
# of
vim ~/.config/scw/config.yaml
```

### Windows
```
%USERPROFILE%\.config\scw\config.yaml
```

---

## Voorbeeld Configuratie Bestand

```yaml
# ~/.config/scw/config.yaml
access_key: SCWXXXXXXXXXXXXXXXXX
secret_key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
default_organization_id: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
default_project_id: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
default_zone: nl-ams-1
default_region: nl-ams
api_url: https://api.scaleway.com
insecure: false
```

---

## Verificatie

Na het aanpassen:

```bash
# Check configuratie
scw config get

# Test resources
scw k8s cluster list
scw instance server list
scw registry namespace list

# Je zou nu je resources moeten zien!
```

---

**Tip**: Meestal is `nl-ams-1` (Amsterdam) de juiste zone voor Nederlandse gebruikers! 🇳🇱


