# Secret Keys Uitleg - Welke Heb Je Nodig?

Er zijn verschillende soorten secret keys in Scaleway. Hier is het verschil.

## Soorten Secret Keys

### 1. Scaleway API Keys (Al Gemaaakt?)
- **Waar voor**: Algemene Scaleway API toegang
- **Gebruikt voor**: Kubernetes cluster management, CLI toegang
- **Waar te vinden**: Scaleway Console → IAM → API Keys
- **Heeft 2 delen**:
  - **Access Key**: Begint met `SCW...`
  - **Secret Key**: UUID formaat

### 2. Container Registry Secret Keys (Anders!)
- **Waar voor**: Alleen voor Container Registry toegang
- **Gebruikt voor**: Docker image pull/push
- **Waar te vinden**: Scaleway Console → Container Registry → Je Registry → Secrets tab
- **Formaat**: UUID formaat (anders dan API Secret Key)

---

## Verschil Tussen Beide

| Type | Waar Voor | Waar Te Vinden | Voorbeeld |
|------|-----------|----------------|-----------|
| **API Keys** | Algemene Scaleway toegang | IAM → API Keys | `SCWXXX...` + UUID |
| **Registry Secret Key** | Container Registry alleen | Registry → Secrets | UUID (anders!) |

**Belangrijk**: Dit zijn **verschillende keys**! Je hebt beide nodig.

---

## Check Welke Je Al Hebt

### Check 1: API Keys (Al Gemaaakt?)

1. **Ga naar Scaleway Console**
2. **IAM** → **API Keys**
3. **Zie je keys?**
   - ✅ Ja → Je hebt API keys
   - ❌ Nee → Maak nieuwe aan (zie hieronder)

### Check 2: Registry Secret Key (Nog Nodig?)

1. **Ga naar Container Registry**
2. **Klik op je registry**
3. **Ga naar "Secrets" tab**
4. **Zie je secret keys?**
   - ✅ Ja → Kopieer deze (dit is anders dan API Secret Key!)
   - ❌ Nee → Maak nieuwe aan (zie hieronder)

---

## Welke Heb Je Nodig Voor Deployment?

Voor automatische deployment heb je **beide** nodig:

### GitHub Secrets Die Je Nodig Hebt:

1. **SCALEWAY_ACCESS_KEY** 
   - Van: API Keys (Access Key deel)
   - Voor: Scaleway CLI toegang

2. **SCALEWAY_SECRET_KEY**
   - Van: API Keys (Secret Key deel)
   - Voor: Scaleway CLI authenticatie

3. **SCALEWAY_K8S_CLUSTER_ID**
   - Van: Kubernetes cluster
   - Voor: Cluster identificatie

4. **SCR_NAMESPACE**
   - Van: Container Registry
   - Voor: Registry namespace naam

5. **SCR_USERNAME**
   - Meestal: `nologin` (standaard)
   - Voor: Registry login

6. **SCR_SECRET_KEY** ⚠️ **DIT IS ANDERS!**
   - Van: Container Registry → Secrets tab
   - Voor: Docker registry pull/push
   - **Dit is NIET hetzelfde als SCALEWAY_SECRET_KEY!**

---

## Stap-voor-Stap: Registry Secret Key Aanmaken

### Optie 1: Via Registry Secrets Tab (Als Beschikbaar)

1. **Ga naar Container Registry**
   - Scaleway Console → **Container Registry**

2. **Klik op je registry**

3. **Ga naar "Secrets" tab**

4. **Klik "Generate new secret key"**

5. **Geef een naam**: `github-actions` (of andere naam)

6. **Klik "Generate"**

7. **BELANGRIJK**: Kopieer de secret key direct!
   - ⚠️ Deze wordt maar 1x getoond
   - Dit is een **andere key** dan je API Secret Key

8. **Bewaar deze** - dit is je `SCR_SECRET_KEY` voor GitHub Secrets

### Optie 2: Gebruik API Secret Key (Als Registry Secrets Niet Beschikbaar)

**Als je geen aparte Registry Secret Key kunt aanmaken**, gebruik dan je **API Secret Key** voor beide:

```yaml
# In GitHub Secrets:
SCALEWAY_SECRET_KEY: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx  # Van API Keys
SCR_SECRET_KEY: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx  # ZELFDE als SCALEWAY_SECRET_KEY
```

**Dit werkt vaak ook!** Scaleway Container Registry accepteert meestal API keys voor authenticatie.

### Optie 3: Via Scaleway CLI

Als de UI niet werkt, probeer via CLI:

```bash
# Login met API keys
scw init

# Test registry toegang
scw registry namespace list

# Als dit werkt, gebruik je API Secret Key voor SCR_SECRET_KEY
```

### Als Je Nog Geen Registry Hebt

1. **Maak eerst een registry aan**:
   - Container Registry → **Create Registry**
   - Kies **Amsterdam** region
   - Geef een naam

2. **Dan volg stappen hierboven** om secret key aan te maken

---

## Verificatie: Welke Keys Heb Je?

### Checklist

- [ ] **API Access Key** (begint met `SCW...`)
  - Voor: `SCALEWAY_ACCESS_KEY` in GitHub
  
- [ ] **API Secret Key** (UUID formaat)
  - Voor: `SCALEWAY_SECRET_KEY` in GitHub
  
- [ ] **Registry Secret Key** (UUID formaat, ANDERS!)
  - Voor: `SCR_SECRET_KEY` in GitHub
  - ⚠️ Dit is een **andere** key dan API Secret Key!

---

## Veelvoorkomende Verwarring

### ❌ Fout: Zelfde Key Gebruiken

```yaml
# FOUT - gebruik niet dezelfde key voor beide!
SCALEWAY_SECRET_KEY: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
SCR_SECRET_KEY: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx  # ZELFDE KEY!
```

### ✅ Correct: Verschillende Keys

```yaml
# CORRECT - gebruik verschillende keys!
SCALEWAY_SECRET_KEY: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx  # Van API Keys
SCR_SECRET_KEY: yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy  # Van Registry Secrets
```

---

## Quick Check: Heb Je Beide?

### Via Scaleway Console

1. **Check API Keys**:
   - IAM → API Keys
   - Zie je een key? → ✅ Je hebt API keys
   - Kopieer Access Key en Secret Key

2. **Check Registry Secret Key**:
   - Container Registry → Je Registry → Secrets tab
   - Zie je een secret key? → ✅ Je hebt registry key
   - Kopieer deze key (dit is anders!)

### Via CLI

```bash
# Check API keys (via config)
scw config get access_key
scw config get secret_key

# Check registry (je moet ingelogd zijn)
scw registry namespace list
# Als dit werkt, heb je API keys goed geconfigureerd

# Voor registry secret key, moet je deze handmatig checken in console
# CLI kan deze niet tonen (security)
```

---

## Voor GitHub Secrets

Je hebt **6 verschillende secrets** nodig:

| GitHub Secret | Waar Van | Type |
|--------------|----------|------|
| `SCALEWAY_ACCESS_KEY` | API Keys | Access Key (`SCW...`) |
| `SCALEWAY_SECRET_KEY` | API Keys | Secret Key (UUID) |
| `SCALEWAY_K8S_CLUSTER_ID` | Kubernetes | Cluster ID (UUID) |
| `SCR_NAMESPACE` | Container Registry | Namespace naam (string) |
| `SCR_USERNAME` | Standaard | `nologin` (string) |
| `SCR_SECRET_KEY` | Registry Secrets OF API Keys | Secret Key (UUID, kan hetzelfde zijn als API Secret Key) |

---

## Samenvatting

**Antwoord op je vraag**: 

✅ **JA, de Registry Secret Key is ANDERS** dan je API Secret Key!

- **API Secret Key**: Voor algemene Scaleway toegang (van IAM → API Keys)
- **Registry Secret Key**: Alleen voor Container Registry (van Registry → Secrets)

Je hebt **beide** nodig voor automatische deployment!

---

## Volgende Stappen

1. ✅ **Check of je API keys hebt** (IAM → API Keys)
2. ✅ **Check of je Registry Secret Key hebt** (Registry → Secrets)
3. ✅ **Als je Registry Secret Key mist**: Maak nieuwe aan (zie hierboven)
4. ✅ **Voeg beide toe aan GitHub Secrets** (verschillende keys!)

---

**Tip**: Maak een notitie van welke key voor wat is, om verwarring te voorkomen! 📝

