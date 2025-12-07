# Platform Engineering Mockup

Een visuele mockup van het Platform Engineering platform, gebaseerd op alle documentatie.

## Overzicht

Deze mockup visualiseert alle 10 Platform Engineering modules, verdeeld over 3 fasen:
- **Fase 1: Foundation** (3 modules) ✅
- **Fase 2: Building** (5 modules) ✅ (3 voltooid, 2 in ontwikkeling)
- **Fase 3: Scaling** (2 modules) ✅

## Features

- **Overzicht**: Visuele weergave van alle fasen en modules
- **Module Details**: Gedetailleerde informatie per module
- **Success Metrics**: DORA metrics en platform metrics
- **Documentatie Links**: Links naar GitHub waar de documentatie staat (opent in nieuwe tab)

## Gebruik

### Lokaal Openen

Open `index.html` in je browser:

```bash
# Via Python
python3 -m http.server 8000

# Via Node.js
npx serve .

# Via PHP
php -S localhost:8000
```

Ga dan naar `http://localhost:8000`

### Direct Openen

Je kunt `index.html` ook direct openen in je browser door dubbel te klikken op het bestand.

### Documentatie Links

De documentatie links openen automatisch naar GitHub:
- Links verwijzen naar `https://github.com/Terminal-WOO/Terminal-zoo-deployer`
- Bestanden worden geopend in een nieuwe tab
- Werkt zowel lokaal als online

## Structuur

```
mockup/
├── index.html          # Hoofdpagina
├── css/
│   └── style.css       # Styling
├── js/
│   └── app.js          # Interactiviteit
└── README.md            # Deze file
```

## Secties

### Hero Section
- Titel en subtitel
- Statistieken (10 modules, 8 voltooid, 3 fasen, 23 documenten)

### Overview Section
- Visuele weergave van alle 3 fasen
- Status badges per fase
- Module lijsten per fase

### Modules Section
- Alle 10 modules in een grid
- Status indicators (voltooid/in ontwikkeling)
- Features per module

### Metrics Section
- DORA metrics (Deployment Frequency, Lead Time, Change Failure Rate, MTTR)
- Platform metrics (Developer Satisfaction, Cognitive Load)

### Documentation Section
- Links naar alle belangrijke documentatie
- Directe toegang tot guides en frameworks

## Customization

### Kleuren Aanpassen

Pas de CSS variabelen aan in `css/style.css`:

```css
:root {
    --primary: #3b82f6;
    --secondary: #8b5cf6;
    --success: #10b981;
    /* ... */
}
```

### Modules Toevoegen/Wijzigen

Bewerk de `modules` array in `js/app.js`:

```javascript
const modules = [
    {
        id: 'X.X',
        title: 'Module Titel',
        phase: 'Foundation',
        description: 'Beschrijving',
        status: 'completed',
        features: ['Feature 1', 'Feature 2']
    },
    // ...
];
```

## Browser Support

- Chrome/Edge (laatste 2 versies)
- Firefox (laatste 2 versies)
- Safari (laatste 2 versies)

## Referenties

- [Platform Engineering Overview](../docs/PLATFORM_ENGINEERING_OVERVIEW.md)
- [Platform Engineering Modules](../deploy/PLATFORM_ENGINEERING_MODULES.md)
- [Executive Summary](../docs/PLATFORM_ENGINEERING_SUMMARY.md)

---

**Status**: Actief  
**Laatste update**: 2025-01-XX

