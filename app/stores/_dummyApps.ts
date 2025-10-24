// /stores/_dummyApps.ts
import type {
  AppModel,
} from "~/types/types";

const base: AppModel = {
  slug: 'govchat-nl',
  name: 'GovChat NL',
  subtitle: 'Veilige chat-assistent met RAG op beleid & wetgeving.',
  summary: 'Chatbot met bronverwijzing en logging.',
  stage: 'productie',
  orgType: 'overheidsontwikkeld',
  imageName: 'ollama/ollama:latest',
  categories: ['RAG', 'Chatbot', 'Nederlands'],
  media: {
    logoUrl: '~/assets/images/logo.svg',
    screenshots: [
      'https://picsum.photos/id/1020/1200/600',
      'https://picsum.photos/id/1018/1200/600',
      'https://picsum.photos/id/1039/1200/600'
    ]
  },
  approved: 'published',
  pricing: { type: 'subscription', price: 199, currency: 'EUR', interval: 'month' },
  labels: { external_support: { has: true }, human_rights_assessment: { has: true }, gov_built: { has: true }, user_guide: { has: true }, open_inference_api: { has: false } },
  criteria: {
    open_source: { met: true }, modular: { met: true }, dpia: { met: true }, license_clarity: { met: true },
    a11y: { met: true }, ai_risk_cat: { met: true }, functional_desc: { met: true }, self_hostable: { met: true }, tech_explainer: { met: true }
  },
  org: { name: 'Rijksoverheid', url: 'https://www.rijksoverheid.nl' },
  versions: [
    { version: '1.4.0', date: new Date().toISOString(), notes: 'Nieuwe bronverwijzing & auditlog.' },
    { version: '1.3.2', date: new Date(Date.now() - 7 * 864e5).toISOString(), notes: 'Bugfixes en prestatieverbeteringen.' }
  ],
  reviews: [
    { id: 'r1', user: 'Eva Janssen', stars: 5, date: new Date().toISOString(), comment: 'Sterk NL en duidelijke bronverwijzing.', labels: ['DPIA OK'] },
    { id: 'r2', user: 'Tom de Vries', stars: 4, date: new Date().toISOString(), comment: 'Werkt goed; index kan sneller.', labels: ['Open-source'] }
  ],
  related: [
    { slug: 'geo-buddy', name: 'Geo Buddy', summary: 'Ruimtelijke AI voor kaarten', },
    { slug: 'parlement-agent', name: 'Parlement Agent', summary: 'Samenvat en vergelijkt Kamerstukken' },
    { slug: 'wl-toegankelijkheid-database', name: 'WL Toegankelijkheid', summary: 'Database voor wetgeving en beleid' },
    { slug: 'libre-infer', name: 'Libre Infer', summary: 'Open-source inference API (LiteLLM)' },
    { slug: 'docs-rag', name: 'Docs RAG', summary: 'Document AI voor beleid' }
  ],
  updatedAt: new Date().toISOString(),

}

const terminalWebapp = {
  slug: 'terminal-webapp',
  name: 'WOO Toegangkelijkheid Webapplicatie',
  subtitle: 'Web applicatie voor verbinding met WOO API.',
  summary: 'Webapplicatie die verbinding maakt met de WOO API om toegankelijke rapporten op te zoeken en om gebruik te kunnen maken van AI om rapporten samen te vatten en te bevragen op inhoud.',
  stage: 'productie',
  orgType: 'commercieel',
  imageName: 'clappform/terminal_zoo_webapp:1.0.3',
  categories: ['RAG', 'Chatbot', 'Nederlands'],
  media: {
    logoUrl: '~/assets/images/logo_terminal.svg',
    screenshots: [
      'https://clappformimages.blob.core.windows.net/social-rewards/hackathon/woo_webapp 1.png',
      'https://clappformimages.blob.core.windows.net/social-rewards/hackathon/woo_webapp 2.png'
    ]
  },
  approved: 'published',
  pricing: { type: 'subscription', price: 199, currency: 'EUR', interval: 'month' },
  labels: { external_support: { has: true }, human_rights_assessment: { has: true }, gov_built: { has: true }, user_guide: { has: true }, open_inference_api: { has: false } },
  criteria: {
    open_source: { met: true }, modular: { met: true }, dpia: { met: true }, license_clarity: { met: true },
    a11y: { met: true }, ai_risk_cat: { met: true }, functional_desc: { met: true }, self_hostable: { met: true }, tech_explainer: { met: true }
  },
  org: { name: 'Clappform', url: 'https://www.clappform.com' },
  versions: [
    { version: '1.4.0', date: new Date().toISOString(), notes: 'Nieuwe bronverwijzing & auditlog.' },
    { version: '1.3.2', date: new Date(Date.now() - 7 * 864e5).toISOString(), notes: 'Bugfixes en prestatieverbeteringen.' }
  ],
  reviews: [
    { id: 'r1', user: 'Eva Janssen', stars: 5, date: new Date().toISOString(), comment: 'Sterk NL en duidelijke bronverwijzing.', labels: ['DPIA OK'] },
    { id: 'r2', user: 'Tom de Vries', stars: 4, date: new Date().toISOString(), comment: 'Werkt goed; index kan sneller.', labels: ['Open-source'] }
  ],
  related: [
    { slug: 'terminal-api', name: 'Terminal API', summary: 'API voor het opzoeken van documenten' }
  ],
  updatedAt: new Date().toISOString(),
} as AppModel

const terminalapi = {
  slug: 'terminal-api',
  name: 'WOO Toegangkelijkheid API',
  subtitle: 'API voor verbinding met elastic search en AI services.',
  summary: 'API die verbinding maakt met elastic search om toegankelijke rapporten op te zoeken en om gebruik te kunnen maken van AI om rapporten samen te vatten en te bevragen op inhoud.',
  stage: 'productie',
  orgType: 'commercieel',
  imageName: 'clappform/terminal_zoo_api:1.0.8',
  categories: ['RAG', 'Chatbot', 'Nederlands'],
  media: {
    logoUrl: '~/assets/images/logo_terminal.svg',
    screenshots: [
      'https://clappformimages.blob.core.windows.net/social-rewards/hackathon/woo_webapp 1.png',
      'https://clappformimages.blob.core.windows.net/social-rewards/hackathon/woo_webapp 2.png'
    ]
  },
  approved: 'published',
  pricing: { type: 'subscription', price: 199, currency: 'EUR', interval: 'month' },
  labels: { external_support: { has: true }, human_rights_assessment: { has: true }, gov_built: { has: true }, user_guide: { has: true }, open_inference_api: { has: false } },
  criteria: {
    open_source: { met: true }, modular: { met: true }, dpia: { met: true }, license_clarity: { met: true },
    a11y: { met: true }, ai_risk_cat: { met: true }, functional_desc: { met: true }, self_hostable: { met: true }, tech_explainer: { met: true }
  },
  org: { name: 'Clappform', url: 'https://www.clappform.com' },
  versions: [
    { version: '1.4.0', date: new Date().toISOString(), notes: 'Nieuwe audit log.' },
    { version: '1.3.2', date: new Date(Date.now() - 7 * 864e5).toISOString(), notes: 'Bugfixes en prestatieverbeteringen.' }
  ],
  reviews: [
    { id: 'r1', user: 'Eva Janssen', stars: 5, date: new Date().toISOString(), comment: 'Sterk NL en duidelijke bronverwijzing.', labels: ['DPIA OK'] },
    { id: 'r2', user: 'Tom de Vries', stars: 4, date: new Date().toISOString(), comment: 'Werkt goed; index kan sneller.', labels: ['Open-source'] }
  ],
  related: [
    { slug: 'terminal-webapp', name: 'Terminal Webapp', summary: 'Webapp voor het opzoeken van documenten' }
  ],
  updatedAt: new Date().toISOString(),
} as AppModel

function clone(slug: string, idx: number): AppModel {
  return {
    ...base,
    slug,
    name: slug.split('-').map(s => s[0].toUpperCase() + s.slice(1)).join(' '),
    categories: idx % 2 ? ['RAG', 'Docs'] : ['Agent'],
    updatedAt: new Date(Date.now() - idx * 864e5).toISOString()
  }
}

export const DUMMY_APPS: AppModel[] = [
  terminalWebapp,
  terminalapi,
  {
    ...terminalWebapp,
    slug: 'wl-toegankelijkheid-database',
    name: 'WL Toegankelijkheid Database',
    subtitle: 'Database voor wetgeving en beleid',
    summary: 'Database voor wetgeving en beleid met toegankelijke AI-interface.',
  },
  base,
  clone('geo-buddy', 1),
  clone('parlement-agent', 2),
  clone('libre-infer', 3),
  clone('docs-rag', 4),
  clone('audit-bot', 5),
 
  
]

export function getDummyBySlug(slug: string): AppModel | undefined {
  return DUMMY_APPS.find(a => a.slug === slug)
}
