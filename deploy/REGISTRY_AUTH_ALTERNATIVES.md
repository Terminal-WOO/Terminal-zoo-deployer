# Container Registry Authenticatie - Alternatieven

Als je geen aparte Registry Secret Key kunt aanmaken, zijn er alternatieven.

## Optie 1: Gebruik API Secret Key (Meest Eenvoudig)

**Scaleway Container Registry accepteert meestal API keys voor authenticatie.**

### Setup

In GitHub Secrets, gebruik dezelfde Secret Key voor beide:

```yaml
# GitHub Secrets:
SCALEWAY_SECRET_KEY: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx  # Van API Keys
SCR_SECRET_KEY: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx  # ZELFDE waarde!
```

**Dit werkt in de meeste gevallen!**

### Verificatie

Test of dit werkt:

```bash
# Login met API keys
docker login rg.nl-ams.scw.cloud -u nologin -p <jouw-api-secret-key>

# Als dit werkt, kun je dezelfde key gebruiken voor SCR_SECRET_KEY
```

---

## Optie 2: Via Scaleway CLI

Als je Scaleway CLI hebt geconfigureerd:

### Check Registry Toegang

```bash
# Test of je registry toegang hebt
scw registry namespace list

# Als dit werkt, gebruik je API Secret Key
```

### Gebruik CLI voor Docker Login

```bash
# Configureer Docker login via CLI
scw registry login

# Dit configureert automatisch Docker met je API keys
```

---

## Optie 3: Check Registry Settings

### Via Scaleway Console

1. **Ga naar Container Registry**
2. **Klik op je registry**
3. **Check "Settings" of "Access" tab**
4. **Zie je authenticatie opties?**
   - API Key authentication
   - Service Account
   - etc.

### Via API

Sommige registries gebruiken API keys standaard:

```bash
# Test met API keys
curl -u nologin:<api-secret-key> https://rg.nl-ams.scw.cloud/v2/
```

---

## Voor GitHub Actions Workflow

### Als Je Geen Registry Secret Key Hebt

Update `.github/workflows/deploy.yml` om dezelfde secret te gebruiken:

```yaml
# In workflow, gebruik dezelfde secret voor beide
- name: Login to Scaleway Container Registry
  uses: docker/login-action@v3
  with:
    registry: ${{ env.REGISTRY }}
    username: ${{ secrets.SCR_USERNAME }}
    password: ${{ secrets.SCALEWAY_SECRET_KEY }}  # Gebruik API Secret Key
```

Of gebruik een aparte secret maar met dezelfde waarde:

```yaml
# In GitHub Secrets, zet beide naar dezelfde waarde:
SCALEWAY_SECRET_KEY: <jouw-api-secret-key>
SCR_SECRET_KEY: <jouw-api-secret-key>  # Zelfde waarde
```

---

## Verificatie: Werkt Het?

### Test Lokaal

```bash
# Test Docker login met API Secret Key
docker login rg.nl-ams.scw.cloud \
  -u nologin \
  -p <jouw-api-secret-key>

# Als dit werkt, kun je dezelfde key gebruiken!
```

### Test in GitHub Actions

Na het pushen naar GitHub:

1. **Check workflow logs**
2. **Zoek naar "Login to Scaleway Container Registry"**
3. **Als het werkt**: ✅ Je kunt dezelfde key gebruiken
4. **Als het faalt**: Probeer alternatieven hieronder

---

## Troubleshooting

### Probleem: Docker Login Fails

**Fout**: `unauthorized: authentication required`

**Oplossing 1**: Check of je de juiste key gebruikt
```bash
# Verifieer API keys werken
scw config get secret_key
```

**Oplossing 2**: Check registry URL
```bash
# Verifieer registry URL
scw registry namespace list
# Noteer de juiste registry URL
```

**Oplossing 3**: Gebruik expliciete authenticatie
```bash
# Probeer met expliciete credentials
docker login rg.nl-ams.scw.cloud \
  -u nologin \
  -p $(scw config get secret_key)
```

### Probleem: Registry Accepteert Geen API Key

**Als API keys niet werken**:

1. **Check Scaleway Documentatie**: https://www.scaleway.com/en/docs/containers/registry/
2. **Contact Scaleway Support**: Mogelijk is er een specifieke setup nodig
3. **Check Registry Type**: Sommige registry types hebben andere authenticatie

---

## Aanbevolen Aanpak

### Stap 1: Probeer API Secret Key Eerst

```yaml
# In GitHub Secrets:
SCR_SECRET_KEY: <jouw-api-secret-key>  # Zelfde als SCALEWAY_SECRET_KEY
```

### Stap 2: Test Lokaal

```bash
docker login rg.nl-ams.scw.cloud -u nologin -p <api-secret-key>
```

### Stap 3: Als Het Werkt

✅ Gebruik dezelfde key voor beide secrets in GitHub

### Stap 4: Als Het Niet Werkt

- Check Scaleway documentatie
- Contact Scaleway support
- Check of er een andere authenticatie methode is

---

## Samenvatting

**Kort antwoord**: 

✅ **Ja, je kunt meestal je API Secret Key gebruiken voor Container Registry!**

In GitHub Secrets:
- `SCALEWAY_SECRET_KEY`: Je API Secret Key
- `SCR_SECRET_KEY`: **Zelfde waarde** (je API Secret Key)

**Dit werkt in de meeste gevallen!**

---

**Tip**: Test eerst lokaal met `docker login` om te verifiëren dat je API Secret Key werkt voor registry toegang! 🐳


