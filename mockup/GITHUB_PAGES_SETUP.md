# GitHub Pages Setup

De mockup is geconfigureerd voor deployment naar GitHub Pages.

## Automatische Setup

De GitHub Actions workflow (`github-pages.yml`) is al geconfigureerd. Je hoeft alleen GitHub Pages in te schakelen in de repository settings.

## Stappen om GitHub Pages te activeren:

1. **Ga naar Repository Settings**
   - Open je GitHub repository
   - Klik op **Settings** (bovenaan de repository)

2. **Ga naar Pages**
   - In het linker menu, scroll naar **Pages**
   - Of ga direct naar: `https://github.com/Terminal-WOO/Terminal-zoo-deployer/settings/pages`

3. **Configureer Source**
   - **Source**: Selecteer **GitHub Actions**
   - Klik **Save**

4. **Wacht op eerste deployment**
   - Ga naar **Actions** tab
   - Je zou een workflow run moeten zien: **"Deploy Mockup to GitHub Pages"**
   - Wacht tot deze voltooid is (groen ✅)

5. **Bekijk je site**
   - Na succesvolle deployment, ga naar:
   - `https://terminal-woo.github.io/Terminal-zoo-deployer/`
   - Of check de URL in Settings → Pages → "Your site is live at..."

## Automatische Updates

Na elke push naar `main` branch met wijzigingen in de `mockup/` directory, wordt de site automatisch bijgewerkt.

## Troubleshooting

### Workflow draait niet
- Check of GitHub Pages is ingeschakeld (Settings → Pages → Source = GitHub Actions)
- Check of de workflow file bestaat: `.github/workflows/github-pages.yml`

### Site is niet bereikbaar
- Wacht 1-2 minuten na eerste deployment
- Check Actions tab voor errors
- Verifieer dat de workflow succesvol is voltooid

### Bestanden worden niet gevonden
- Verifieer dat alle bestanden in `mockup/` directory staan
- Check of de paths in `index.html` correct zijn (relatief, niet absoluut)

## Handmatige Deployment (optioneel)

Als je handmatig wilt deployen:

```bash
# Clone repository
git clone https://github.com/Terminal-WOO/Terminal-zoo-deployer.git
cd Terminal-zoo-deployer

# Maak wijzigingen in mockup/
# ...

# Commit en push
git add mockup/
git commit -m "Update mockup"
git push origin main

# GitHub Actions zal automatisch deployen
```

---

**Status**: ✅ Workflow geconfigureerd  
**URL**: Wordt beschikbaar na activatie van GitHub Pages  
**Laatste update**: 2025-01-XX

