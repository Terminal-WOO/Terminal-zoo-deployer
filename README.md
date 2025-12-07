# App Store - AI Application Deployment Platform

Een moderne applicatie voor het deployen en beheren van AI-toepassingen, gebouwd met Nuxt.js (frontend) en Go (backend).

## Overzicht

De applicatie bestaat uit:
- **Frontend**: Nuxt.js applicatie met moderne UI
- **Backend**: Go Kubernetes API server voor deployment management
- **Infrastructuur**: Kubernetes cluster op Scaleway met container registry

## Features

- Transparante AI-oplossingen met betrouwbare criteria
- Community-monitoring en bouwstenen
- Snelle en veilige implementaties
- Kubernetes-native deployment

## 📚 Platform Engineering Documentatie

Dit platform is gebouwd volgens de principes uit "Effective Platform Engineering" met **10 modules** verdeeld over **3 fasen**.

### 🚀 Quick Start

- **[Platform Engineering Overview](docs/PLATFORM_ENGINEERING_OVERVIEW.md)** - Hoofdnavigatie met alle links
- **[Executive Summary](docs/PLATFORM_ENGINEERING_SUMMARY.md)** - Snelle samenvatting
- **[Platform Engineering Modules](deploy/PLATFORM_ENGINEERING_MODULES.md)** - Volledige module beschrijvingen

### 📖 Belangrijkste Documentatie

**Foundation**:
- [Platform Visie](docs/platform-vision.md) - Visie, missie en doelen
- [Platform Roadmap](docs/platform-roadmap.md) - Strategische roadmap
- [Platform Domains](docs/platform-domains.md) - Product domains
- [Architecture Decisions](docs/architecture-decisions/README.md) - ADRs

**Building**:
- [Governance Framework](docs/governance-framework.md) - Governance en compliance
- [Observability Guide](docs/observability-guide.md) - Observability platform
- [Infrastructure Guide](docs/infrastructure-guide.md) - Infrastructure as code

**Scaling**:
- [Scaling Architecture](docs/scaling-architecture.md) - Scaling strategieën
- [Product Evolution](docs/product-evolution.md) - Platform evolutie
- [Developer Experience](docs/developer-experience.md) - DevEx improvements

## Setup

Make sure to install dependencies:

```bash
# npm
npm install

# pnpm
pnpm install

# yarn
yarn install

# bun
bun install
```

## Development Server

Start the development server on `http://localhost:3000`:

```bash
# npm
npm run dev

# pnpm
pnpm dev

# yarn
yarn dev

# bun
bun run dev
```

## Production

Build the application for production:

```bash
# npm
npm run build

# pnpm
pnpm build

# yarn
yarn build

# bun
bun run build
```

Locally preview production build:

```bash
# npm
npm run preview

# pnpm
pnpm preview

# yarn
yarn preview

# bun
bun run preview
```

Check out the [deployment documentation](https://nuxt.com/docs/getting-started/deployment) for more information.
