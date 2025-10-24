// /data/apps/dummy.ts

export type User = {
  id: string
  name: string
  email: string
  roles: string[]
  organization?: string
  org_type?: 'overheidsontwikkeld' | 'commercieel'
  website?: string
  token?: string
}

export type CriteriaKey =
  | 'open_source' | 'modular' | 'dpia' | 'license_clarity' | 'a11y' | 'ai_risk_cat'
  | 'functional_desc' | 'self_hostable' | 'tech_explainer'

export type ExtraLabelKey =
  | 'external_support' | 'human_rights_assessment' | 'gov_built'
  | 'user_guide' | 'open_inference_api'

export type ApprovalStatus = 'submitted' | 'under_review' | 'changes_requested' | 'rejected' | 'published'

export type Pricing =
  | { type: 'free' }
  | { type: 'one_time', price: number, currency: string }
  | { type: 'subscription', price: number, currency: string, interval: 'month' | 'year' }
  | { type: 'usage', unit: string }

export interface AppModel {
  slug: string
  name: string
  subtitle: string
  summary?: string
  stage: 'productie' | 'pilots' | 'testfase'
  orgType: 'overheidsontwikkeld' | 'commercieel'
  categories: string[]
  media?: { logoUrl?: string; screenshots?: string[] }
  descriptionMd?: string
  pricing?: Pricing
  criteria?: Record<CriteriaKey, { met: boolean; evidenceUrl?: string }>
  labels?: Record<ExtraLabelKey, { has: boolean; evidenceUrl?: string }>
  org?: { name: string; url?: string }
  versions?: { version: string; date: string; notes: string }[]
  reviews?: { id: string; user: string; stars: number; date: string; comment: string; labels: string[] }[]
  related?: { slug: string; name: string; summary: string; logo?: string }[]
  updatedAt?: string
  approved: ApprovalStatus
  imageName?: string
}

export const DUMMY_APP: AppModel = {
  slug: 'govchat-nl',
  name: 'GovChat NL',
  subtitle: 'Veilige chat-assistent met RAG op beleid & wetgeving.',
  summary: 'Chatbot voor NL-overheid met bronverwijzing en logging.',
  stage: 'productie',
  orgType: 'overheidsontwikkeld',
  categories: ['RAG', 'Chatbot', 'Nederlands'],

  media: {
    logoUrl: '~/assets/images/logo.svg',
    screenshots: [
      'https://picsum.photos/id/1020/1200/600',
      'https://picsum.photos/id/1018/1200/600',
      'https://picsum.photos/id/1039/1200/600'
    ]
  },

  descriptionMd:
    `GovChat NL helpt medewerkers bij het opstellen en toetsen van brieven en antwoorden.
De app gebruikt retrieval-augmented generation (RAG) met bronnen uit overheid.nl en lokale beleidsdatabanken.

**Kernwaarden:** transparantie, bronverwijzing en naleving van DPIA/richtlijnen.`,

  pricing: { type: 'subscription', price: 199, currency: 'EUR', interval: 'month' },

  criteria: {
    open_source: { met: true },
    modular: { met: true },
    dpia: { met: true, evidenceUrl: '/docs/dpia.pdf' },
    license_clarity: { met: true },
    a11y: { met: true },
    ai_risk_cat: { met: true },
    functional_desc: { met: true },
    self_hostable: { met: true },
    tech_explainer: { met: true }
  },

  labels: {
    external_support: { has: true },
    human_rights_assessment: { has: true },
    gov_built: { has: true },
    user_guide: { has: true },
    open_inference_api: { has: false }
  },

  org: { name: 'Rijksoverheid', url: 'https://www.rijksoverheid.nl' },

  versions: [
    { version: '1.4.0', date: new Date().toISOString(), notes: 'Nieuwe bronverwijzing & auditlog.' },
    { version: '1.3.2', date: new Date(Date.now() - 7 * 864e5).toISOString(), notes: 'Bugfixes en prestatieverbeteringen.' }
  ],

  reviews: [
    { id: 'r1', user: 'Eva Janssen', stars: 5, date: new Date().toISOString(), comment: 'Sterk in Nederlands met duidelijke bronverwijzing.', labels: ['DPIA OK'] },
    { id: 'r2', user: 'Tom de Vries', stars: 4, date: new Date().toISOString(), comment: 'Werkt goed, graag snellere zoekindex.', labels: ['Open-source'] }
  ],

  related: [
    { slug: 'geo-buddy', name: 'Geo Buddy', summary: 'Ruimtelijke AI voor kaarten' },
    { slug: 'parlement-agent', name: 'Parlement Agent', summary: 'Samenvat en vergelijkt Kamerstukken' },
    { slug: 'libre-infer', name: 'Libre Infer', summary: 'Open-source inference API (LiteLLM)' },
    { slug: 'docs-rag', name: 'Docs RAG', summary: 'Document AI voor beleid' }
  ],

  updatedAt: new Date().toISOString()
}

export interface ClusterConfig {
  name: string
  auth: "kubeconfig" | "inCluster"| "serviceAccountToken"
  kubeconfig?: string // only for auth = kubeconfig
  server?: string // only for auth = serviceAccountToken
  bearerToken?: string // only for auth = serviceAccountToken
  caPEM?: string // only for auth = serviceAccountToken
  domain?: string // optional, for reference
  privateKey?: string // optional, for reference
  Certificate?: string // optional, for reference
}

export interface ManualDeploymentConfig {
  deploymentName: string
  namespace: string
  containerImage: string
  replicas: number
  ports: Port[]
  resources: Resources
  clusterName?: string // optional, for selecting a specific cluster
}

export interface Resources {
  cpuRequests: string;
  cpuLimits: string;
  memoryRequests: string;
  memoryLimits: string;
}

export interface Port {
  ContainerPort: number;
}