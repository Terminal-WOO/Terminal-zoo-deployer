# Scaleway CLI - Namespace Commando's

Handige commando's om namespaces te vinden en beheren met Scaleway CLI.

## Container Registry Namespaces

### List Alle Registry Namespaces

```bash
# List alle container registry namespaces
scw registry namespace list

# Met meer details
scw registry namespace list -o json

# Filter op region
scw registry namespace list region=nl-ams
```

### Get Specifieke Namespace Details

```bash
# Get details van een namespace
scw registry namespace get <namespace-id>

# Of met naam
scw registry namespace get name=<namespace-name>
```

### Find Namespace ID

```bash
# List namespaces en zoek naar je namespace
scw registry namespace list | grep nl-appstore-registry

# Of gebruik JSON output voor betere filtering
scw registry namespace list -o json | jq '.[] | select(.name=="nl-appstore-registry")'
```

---

## Kubernetes Namespaces

### List Kubernetes Namespaces in Cluster

```bash
# Via kubectl (na kubeconfig setup)
kubectl get namespaces

# Of specifiek
kubectl get namespace nl-appstore-registry
```

### Via Scaleway CLI (Cluster Info)

```bash
# List clusters
scw k8s cluster list

# Get cluster details (bevat namespace info)
scw k8s cluster get <cluster-id>
```

---

## Handige Commando's

### Find Registry Namespace ID

```bash
# Quick find
scw registry namespace list | grep -i appstore

# Of met JSON
scw registry namespace list -o json | jq -r '.[] | "\(.name) - \(.id)"'
```

### Check Namespace Exists

```bash
# Check of namespace bestaat
scw registry namespace list | grep nl-appstore-registry

# Of
scw registry namespace get name=nl-appstore-registry
```

### Get Namespace URL

```bash
# Registry URL is meestal:
# rg.<region>.scw.cloud/<namespace-name>

# Voor nl-appstore-registry in Amsterdam:
# rg.nl-ams.scw.cloud/nl-appstore-registry
```

---

## Voor Je Project

### Check Je Registry Namespace

```bash
# List alle namespaces
scw registry namespace list

# Zoek naar nl-appstore-registry
scw registry namespace list | grep nl-appstore-registry

# Get details
scw registry namespace get name=nl-appstore-registry
```

### Verifieer Namespace Naam

```bash
# Check of namespace bestaat
if scw registry namespace list | grep -q "nl-appstore-registry"; then
  echo "Namespace gevonden!"
else
  echo "Namespace niet gevonden - maak aan via console of CLI"
fi
```

---

## Troubleshooting

### "Namespace not found"

**Oplossing**:
1. Check of je de juiste region gebruikt:
   ```bash
   scw registry namespace list region=nl-ams
   ```

2. Check of namespace bestaat in Scaleway Console

3. Maak namespace aan als die niet bestaat:
   ```bash
   scw registry namespace create name=nl-appstore-registry region=nl-ams
   ```

### "Permission denied"

**Oplossing**:
1. Check of je ingelogd bent:
   ```bash
   scw config get access_key
   ```

2. Re-authenticate indien nodig:
   ```bash
   scw init
   ```

---

## Quick Reference

```bash
# List namespaces
scw registry namespace list

# Get namespace details
scw registry namespace get name=nl-appstore-registry

# Create namespace (als nodig)
scw registry namespace create name=nl-appstore-registry region=nl-ams

# Delete namespace (voorzichtig!)
scw registry namespace delete name=nl-appstore-registry
```

---

**Tip**: Gebruik `scw registry namespace list` om alle beschikbare namespaces te zien! 📦

