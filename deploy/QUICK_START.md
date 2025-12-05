# Quick Start: Automatische Deployment Setup

Volg deze stappen om automatische deployment naar Scaleway te configureren.

> **⚠️ UI Problemen?** Als je resources ziet maar niets ziet als je erop klikt, zie **[SCALEWAY_UI_TROUBLESHOOTING.md](SCALEWAY_UI_TROUBLESHOOTING.md)** voor hulp.

> **💡 Zone Probleem?** Als je resources niet ziet via CLI, check je default zone: **[SCALEWAY_CLI_CONFIG.md](SCALEWAY_CLI_CONFIG.md)**

## Stap 1: Scaleway API Keys Aanmaken (5 minuten)

### 1.1 Log in op Scaleway Console
1. Ga naar https://console.scaleway.com
2. Log in met je account

### 1.2 Maak API Keys aan
1. Klik op je **profiel** (rechtsboven) → **API Keys**
2. Klik op **Generate API Key**
3. Geef een naam: `GitHub Actions Deployment`
4. Klik **Generate**
5. **BELANGRIJK**: Kopieer beide keys direct:
   - **Access Key**: Begint met `SCW...` (bijv. `SCWXXXXXXXXXXXXXXXXX`)
   - **Secret Key**: UUID formaat (bijv. `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`)
   
   ⚠️ **De Secret Key wordt maar 1x getoond!** Kopieer deze nu.

6. Bewaar deze keys tijdelijk in een tekstbestand (we voegen ze straks toe aan GitHub)

---

## Stap 2: Kubernetes Cluster ID Vinden (10-15 minuten)

> **💡 Tip**: Als je een oude DEV1-S instance hebt die je niet meer nodig hebt, verwijder die eerst om kosten te besparen. Zie **[CLEANUP_OLD_RESOURCES.md](CLEANUP_OLD_RESOURCES.md)** voor instructies.

### 2.1 Maak een Kubernetes Cluster aan (als je die nog niet hebt)

**📖 Voor gedetailleerde instructies**: Zie **[K8S_CLUSTER_SETUP.md](K8S_CLUSTER_SETUP.md)**

**⚠️ Probleem met leeg cluster?** Zie **[TROUBLESHOOTING_CLUSTER.md](TROUBLESHOOTING_CLUSTER.md)**

**Kort overzicht:**
1. Ga naar **Kubernetes** → **Kapsule** → **Create Cluster**
2. **Cluster Name**: `terminal-zoo-cluster`
3. **Region**: `Amsterdam (nl-ams)`
4. **Kubernetes Version**: `1.31.x` (latest)
5. **Node Pool Configuratie** (BELANGRIJK - maak deze aan!):
   - **Node Type**: `DEV1-M` (4 vCPU, 8GB RAM) - **~€30/maand**
     - ⚠️ Dit is meestal het minimum beschikbare type
   - **Node Count**: `1` (start klein!)
   - **Auto-scaling**: Min `1`, Max `2` (voorkomt hoge kosten)
   - **Spot Instances**: ❌ Disabled (voor productie)
6. Klik **Create Cluster** en wacht 5-10 minuten

**💰 Kosten**: ~€30/maand voor 1 node (nog steeds kostenefficiënt!)
**✅ Voordeel**: Meer headroom voor groei zonder direct op te hoeven schalen

**💡 Belangrijk**: Zorg dat je een **node pool aanmaakt** tijdens cluster creation! Als je cluster leeg lijkt, heeft het waarschijnlijk geen node pools.

### 2.2 Kopieer Cluster ID
1. Na het aanmaken, klik op je cluster naam
2. Op de **Overview** pagina, zoek naar **Cluster ID**
3. Kopieer de UUID (bijv. `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`)
4. Bewaar deze tijdelijk

**Of via CLI:**
```bash
scw k8s cluster list
# Zoek je cluster en kopieer de ID
```

---

## Stap 3: Container Registry Info Verzamelen (3 minuten)

### 3.1 Ga naar Container Registry
1. In Scaleway Console, ga naar **Container Registry**
2. Als je nog geen registry hebt, klik **Create Registry**
3. Kies **Amsterdam** region
4. Noteer je **Registry URL** (bijv. `rg.nl-ams.scw.cloud`)
5. Noteer je **Namespace** naam (bijv. `nl-appstore-registry`)

**Of via CLI:**
```bash
# List alle namespaces
scw registry namespace list

# Zoek naar je namespace
scw registry namespace list | grep nl-appstore-registry
```

> **📖 Voor meer CLI commando's**: Zie **[SCW_NAMESPACE_COMMANDS.md](SCW_NAMESPACE_COMMANDS.md)**

### 3.2 Genereer Registry Secret Key
1. Ga naar je registry → **Secrets** tab
2. Klik **Generate new secret key**
3. Geef een naam: `github-actions`
4. Klik **Generate**
5. **Kopieer de secret key** (wordt maar 1x getoond)
6. Bewaar deze tijdelijk

> **💡 Belangrijk**: De Registry Secret Key is **ANDERS** dan je API Secret Key! 
> - API Secret Key: Voor Scaleway toegang (van IAM → API Keys)
> - Registry Secret Key: Alleen voor Container Registry (van Registry → Secrets)
> 
> Zie **[SECRET_KEYS_EXPLAINED.md](SECRET_KEYS_EXPLAINED.md)** voor uitleg over het verschil.

---

## Stap 3.5: Kubeconfig Configureren (5 minuten)

> **💡 Belangrijk**: Voordat je kubectl commando's kunt gebruiken, moet je kubeconfig configureren!

### Via Scaleway CLI (Aanbevolen)

```bash
# 1. Installeer Scaleway CLI (als je die nog niet hebt)
brew install scw  # Mac
# Of download van: https://github.com/scaleway/scaleway-cli/releases

# 2. Authenticate
scw init
# Volg instructies, gebruik je API keys

# 3. Installeer kubeconfig
scw k8s kubeconfig install <cluster-id>
# Of als je maar 1 cluster hebt:
scw k8s kubeconfig install

# 4. Verifieer
kubectl get nodes
# Je zou je nodes moeten zien!
```

### Via Scaleway Console (Handmatig)

1. **Ga naar Kubernetes → Kapsule → Je cluster**
2. **Klik "Download kubeconfig"**
3. **Kopieer naar**: `~/.kube/config`
4. **Set permissions**: `chmod 600 ~/.kube/config`
5. **Test**: `kubectl get nodes`

> **📖 Voor uitgebreide instructies**: Zie **[KUBECONFIG_SETUP.md](KUBECONFIG_SETUP.md)**

---

## Stap 4: GitHub Secrets Configureren (5 minuten)

### 4.1 Ga naar GitHub Repository Settings
1. Open je GitHub repository
2. Klik op **Settings** (bovenaan de repository)
3. In het linker menu, klik **Secrets and variables** → **Actions**

### 4.2 Voeg Secrets Toe

Klik voor elk secret op **New repository secret** en voeg toe:

#### Secret 1: SCALEWAY_ACCESS_KEY
- **Name**: `SCALEWAY_ACCESS_KEY`
- **Value**: Je Access Key (begint met `SCW...`)
- Klik **Add secret**

#### Secret 2: SCALEWAY_SECRET_KEY
- **Name**: `SCALEWAY_SECRET_KEY`
- **Value**: Je Secret Key (UUID formaat)
- Klik **Add secret**

#### Secret 3: SCALEWAY_K8S_CLUSTER_ID
- **Name**: `SCALEWAY_K8S_CLUSTER_ID`
- **Value**: Je Kubernetes Cluster ID (UUID)
- Klik **Add secret**

#### Secret 4: SCR_NAMESPACE
- **Name**: `SCR_NAMESPACE`
- **Value**: Je Container Registry namespace naam (bijv. `my-registry`)
- Klik **Add secret**

#### Secret 5: SCR_USERNAME
- **Name**: `SCR_USERNAME`
- **Value**: `nologin` (dit is standaard voor Scaleway)
- Klik **Add secret**

#### Secret 6: SCR_SECRET_KEY
- **Name**: `SCR_SECRET_KEY`
- **Value**: Je Container Registry secret key
- Klik **Add secret**

### 4.3 Verifieer Secrets
Je zou nu 6 secrets moeten hebben:
- ✅ SCALEWAY_ACCESS_KEY
- ✅ SCALEWAY_SECRET_KEY
- ✅ SCALEWAY_K8S_CLUSTER_ID
- ✅ SCR_NAMESPACE
- ✅ SCR_USERNAME
- ✅ SCR_SECRET_KEY

---

## Stap 5: Kubernetes Manifests Aanpassen (5 minuten)

### 5.1 Update Image URLs

Open `k8s/frontend-deployment.yaml` en zoek naar:
```yaml
image: rg.nl-ams.scw.cloud/your-registry/ai-co:latest
```

Vervang `your-registry` met je eigen registry namespace:
```yaml
image: rg.nl-ams.scw.cloud/<jouw-registry-namespace>/ai-co:latest
```

Doe hetzelfde in `k8s/backend-deployment.yaml`:
```yaml
image: rg.nl-ams.scw.cloud/<jouw-registry-namespace>/ai-co:latest
```

### 5.2 Update Ingress (Optioneel)

Als je een custom domain hebt, open `k8s/ingress.yaml` en vervang:
```yaml
- host: your-domain.com
```
met je eigen domain.

---

## Stap 6: Eerste Deployment Testen (10 minuten)

### 6.1 Commit en Push
```bash
# Check wat er gewijzigd is
git status

# Voeg wijzigingen toe
git add k8s/

# Commit
git commit -m "Configure automatic deployment to Scaleway"

# Push naar main branch
git push origin main
```

### 6.2 Monitor GitHub Actions
1. Ga naar je GitHub repository
2. Klik op **Actions** tab
3. Je zou een workflow run moeten zien: **"Build and Deploy"**
4. Klik erop om de progress te volgen

### 6.3 Check Workflow Status

De workflow heeft 3 jobs:
1. **Build Frontend** - Bouwt Nuxt image
2. **Build Backend** - Bouwt Go image  
3. **Deploy to Kubernetes** - Deployt naar cluster

Wacht tot alle 3 groen zijn (✅).

---

## Stap 7: Verifieer Deployment (5 minuten)

### 7.1 Check Kubernetes Pods
```bash
# Verifieer kubectl werkt
kubectl get nodes

# Check pods
kubectl get pods -n terminal-zoo

# Je zou moeten zien:
# nuxt-frontend-xxxxx   1/1   Running
# go-backend-xxxxx       1/1   Running
```

### 7.2 Check Services
```bash
kubectl get services -n terminal-zoo

# Je zou moeten zien:
# nuxt-frontend   ClusterIP   10.x.x.x   3000/TCP
# go-backend      ClusterIP   10.x.x.x   8080/TCP
```

### 7.3 Test Health Endpoints
```bash
# Port forward naar frontend
kubectl port-forward -n terminal-zoo svc/nuxt-frontend 3000:3000

# In een andere terminal, test:
curl http://localhost:3000/api/health
# Verwacht: {"status":"healthy"}

# Test backend
kubectl port-forward -n terminal-zoo svc/go-backend 8080:8080

# In een andere terminal:
curl http://localhost:8080/health
# Verwacht: {"status":"healthy"}
```

---

## Troubleshooting

### Workflow Fails: "authentication failed"
- ✅ Check of `SCALEWAY_ACCESS_KEY` en `SCALEWAY_SECRET_KEY` correct zijn
- ✅ Verifieer de keys niet zijn verlopen
- ✅ Check of de keys de juiste permissions hebben

### Workflow Fails: "cluster not found"
- ✅ Verifieer `SCALEWAY_K8S_CLUSTER_ID` is correct
- ✅ Check of het cluster actief is in Scaleway Console

### Workflow Fails: "ImagePullBackOff"
- ✅ Verifieer `SCR_SECRET_KEY` is correct
- ✅ Check of image pull secret bestaat:
  ```bash
  kubectl get secret scr-secret -n terminal-zoo
  ```
- ✅ Maak secret aan indien nodig:
  ```bash
  kubectl create secret docker-registry scr-secret \
    --docker-server=rg.nl-ams.scw.cloud \
    --docker-username=nologin \
    --docker-password=<jouw-registry-secret-key> \
    -n terminal-zoo
  ```

### Pods Starten Niet
```bash
# Check pod logs
kubectl logs -l app=nuxt-frontend -n terminal-zoo
kubectl logs -l app=go-backend -n terminal-zoo

# Check pod events
kubectl describe pod -l app=nuxt-frontend -n terminal-zoo
```

---

## Volgende Stappen

Na succesvolle deployment:

1. **Monitor eerste paar deployments** - Check of alles goed gaat
2. **Setup notifications** - Ontvang emails bij failed deployments
3. **Review kosten** - Check Scaleway Console → Billing
4. **Document custom configs** - Noteer eventuele aanpassingen

---

## Hulp Nodig?

- Check `deploy/AUTOMATIC_DEPLOYMENT_SETUP.md` voor uitgebreide documentatie
- Check `deploy/SETUP_CHECKLIST.md` voor een complete checklist
- GitHub Actions logs: Repository → Actions → Klik op failed workflow

---

**Gefeliciteerd! 🎉** Je hebt nu automatische deployment naar Scaleway geconfigureerd!

Bij elke push naar `main` branch wordt je applicatie automatisch gebouwd en gedeployed.

