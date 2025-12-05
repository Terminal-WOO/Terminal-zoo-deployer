# Kubeconfig Setup - Scaleway Kubernetes

Deze guide legt uit hoe je kubeconfig configureert om toegang te krijgen tot je Scaleway Kubernetes cluster.

## Overzicht

Kubeconfig is een configuratiebestand dat kubectl vertelt hoe het verbinding moet maken met je Kubernetes cluster.

## Methode 1: Via Scaleway CLI (Aanbevolen)

### Stap 1: Installeer Scaleway CLI (Als Je Die Nog Niet Hebt)

```bash
# Mac
brew install scw

# Linux
curl -o /usr/local/bin/scw -L "https://github.com/scaleway/scaleway-cli/releases/latest/download/scw-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
chmod +x /usr/local/bin/scw

# Windows
# Download van: https://github.com/scaleway/scaleway-cli/releases
```

### Stap 2: Authenticate Met Scaleway

```bash
scw init
# Volg de instructies:
# - Enter your access key: <jouw-access-key>
# - Enter your secret key: <jouw-secret-key>
# - Select default region: nl-ams (of andere)
# - Select default zone: nl-ams-1 (of andere)
```

### Stap 3: Installeer Kubeconfig

```bash
# Vind je cluster ID eerst
scw k8s cluster list

# Installeer kubeconfig voor je cluster
scw k8s kubeconfig install <cluster-id>

# Of als je maar 1 cluster hebt:
scw k8s kubeconfig install
```

### Stap 4: Verifieer Connectie

```bash
# Test connectie
kubectl get nodes

# Je zou moeten zien:
# NAME                    STATUS   ROLES    AGE   VERSION
# scw-xxxxx               Ready    <none>   5m    v1.31.x
```

---

## Methode 2: Via Scaleway Console (Handmatig)

### Stap 1: Download Kubeconfig

1. **Ga naar Scaleway Console**
   - https://console.scaleway.com

2. **Ga naar Kubernetes → Kapsule**

3. **Klik op je cluster**

4. **Klik op "Download kubeconfig"** of **"Kubeconfig"** tab

5. **Download het bestand** (meestal `kubeconfig.yaml`)

### Stap 2: Configureer Kubeconfig

#### Optie A: Vervang Standaard Kubeconfig

```bash
# Backup huidige kubeconfig (als je die hebt)
mv ~/.kube/config ~/.kube/config.backup

# Kopieer nieuwe kubeconfig
mkdir -p ~/.kube
cp ~/Downloads/kubeconfig.yaml ~/.kube/config

# Set juiste permissions
chmod 600 ~/.kube/config
```

#### Optie B: Gebruik Specifieke Context (Aanbevolen)

```bash
# Kopieer kubeconfig naar specifieke locatie
mkdir -p ~/.kube
cp ~/Downloads/kubeconfig.yaml ~/.kube/scaleway-config

# Exporteer KUBECONFIG environment variable
export KUBECONFIG=~/.kube/scaleway-config

# Of voeg toe aan ~/.bashrc of ~/.zshrc
echo 'export KUBECONFIG=~/.kube/scaleway-config' >> ~/.zshrc
source ~/.zshrc
```

#### Optie C: Merge Met Bestaande Config

```bash
# Als je al andere clusters hebt geconfigureerd
export KUBECONFIG=~/.kube/config:~/.kube/scaleway-config
kubectl config view --flatten > ~/.kube/config_merged
mv ~/.kube/config_merged ~/.kube/config
```

### Stap 3: Verifieer Connectie

```bash
# Test connectie
kubectl get nodes

# Check huidige context
kubectl config current-context

# List alle contexts
kubectl config get-contexts
```

---

## Methode 3: Via kubectl Direct (Als Je Cluster Details Hebt)

Als je cluster details hebt (API server URL, certificaat, token):

```bash
# Set cluster
kubectl config set-cluster scaleway-cluster \
  --server=https://<api-server-url> \
  --certificate-authority=<ca-cert-file>

# Set credentials
kubectl config set-credentials scaleway-user \
  --token=<bearer-token>

# Set context
kubectl config set-context scaleway-context \
  --cluster=scaleway-cluster \
  --user=scaleway-user

# Use context
kubectl config use-context scaleway-context
```

---

## Troubleshooting

### Probleem: "Unable to connect to the server"

**Oplossing 1**: Check kubeconfig locatie
```bash
# Check waar kubectl kubeconfig zoekt
kubectl config view

# Check KUBECONFIG environment variable
echo $KUBECONFIG

# Check of bestand bestaat
ls -la ~/.kube/config
```

**Oplossing 2**: Verifieer kubeconfig inhoud
```bash
# Check kubeconfig
cat ~/.kube/config

# Zorg dat het geldige YAML is
# Check of server URL correct is
# Check of certificaten aanwezig zijn
```

**Oplossing 3**: Check netwerk connectie
```bash
# Test of je cluster bereikbaar is
curl -k https://<api-server-url>/version

# Check firewall rules
```

### Probleem: "certificate signed by unknown authority"

**Oplossing**:
```bash
# Download kubeconfig opnieuw van Scaleway Console
# Of gebruik:
scw k8s kubeconfig install <cluster-id> --overwrite
```

### Probleem: "context not found"

**Oplossing**:
```bash
# List alle contexts
kubectl config get-contexts

# Set juiste context
kubectl config use-context <context-name>

# Of check kubeconfig
kubectl config view
```

### Probleem: "permission denied"

**Oplossing**:
```bash
# Check file permissions
ls -la ~/.kube/config

# Set juiste permissions
chmod 600 ~/.kube/config
```

---

## Verificatie Checklist

Na kubeconfig setup:

- [ ] `kubectl get nodes` werkt
- [ ] Je ziet je cluster nodes
- [ ] `kubectl config current-context` toont juiste context
- [ ] Je kunt namespaces zien: `kubectl get namespaces`

---

## Meerdere Clusters Beheren

Als je meerdere clusters hebt:

### Switch Tussen Clusters

```bash
# List alle contexts
kubectl config get-contexts

# Switch context
kubectl config use-context <context-name>

# Check huidige context
kubectl config current-context
```

### Gebruik Aliases

Voeg toe aan `~/.bashrc` of `~/.zshrc`:

```bash
# Alias voor Scaleway cluster
alias k8s-scaleway='kubectl config use-context scaleway-context'

# Alias voor andere cluster
alias k8s-local='kubectl config use-context docker-desktop'
```

---

## Kubeconfig Locatie

### Standaard Locaties

- **Mac/Linux**: `~/.kube/config`
- **Windows**: `%USERPROFILE%\.kube\config`

### Environment Variable

```bash
# Set specifieke kubeconfig
export KUBECONFIG=~/.kube/scaleway-config

# Of voor meerdere configs (worden gemerged)
export KUBECONFIG=~/.kube/config:~/.kube/scaleway-config
```

---

## Quick Start Commando's

```bash
# 1. Installeer Scaleway CLI (als nodig)
brew install scw  # Mac
# Of download van GitHub voor andere OS

# 2. Authenticate
scw init

# 3. Installeer kubeconfig
scw k8s kubeconfig install <cluster-id>

# 4. Verifieer
kubectl get nodes

# 5. Check context
kubectl config current-context
```

---

## Voor GitHub Actions

Als je kubeconfig nodig hebt voor GitHub Actions (meestal niet nodig met Scaleway CLI):

1. **Download kubeconfig** van Scaleway Console
2. **Base64 encode**:
   ```bash
   cat ~/.kube/config | base64 -w 0
   ```
3. **Voeg toe als GitHub Secret**: `KUBECONFIG`
4. **Workflow gebruikt dit automatisch** (zie `.github/workflows/deploy.yml`)

**Maar**: De workflow gebruikt Scaleway CLI, dus dit is meestal niet nodig!

---

## Volgende Stappen

Na kubeconfig setup:

1. ✅ **Verifieer connectie**: `kubectl get nodes`
2. ✅ **Check namespaces**: `kubectl get namespaces`
3. ✅ **Ga verder met deployment**: `make k8s-apply`

---

**Tip**: Gebruik Scaleway CLI (`scw k8s kubeconfig install`) - dit is de makkelijkste methode! 🚀

