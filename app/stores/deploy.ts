import { defineStore } from 'pinia'

export type DeployStatus = 'draft'|'queued'|'deploying'|'running'|'failed'|'canceled'|'decommissioned'

export interface DeploySpec {
  slug?: string | null
  llm: string | null
  vectordb: string | null
  extras: string[]
  provider: 'azure'|'hetzner'|'onprem'|null
  region: string | null
  runtime: 'k8s'|'docker'|null
  resources: { cpu: number; ramGiB: number; gpu?: 'none'|'a10'|'t4'|'l4' }
  env: { tenantId?: string; domain?: string }
  secrets: { apiKeys: Record<string,string> }
  notes: string
}

export interface Deployment {
  id: string
  spec: DeploySpec
  status: DeployStatus
  createdAt: string
  updatedAt: string
  logs?: string[]
}

const defaultSpec = (): DeploySpec => ({
  slug: null,
  llm: null,
  vectordb: null,
  extras: [],
  provider: null,
  region: null,
  runtime: 'k8s',
  resources: { cpu: 2, ramGiB: 8, gpu: 'none' },
  env: {},
  secrets: { apiKeys: {} },
  notes: ''
})

export const useDeployStore = defineStore('deploy', {
  state: () => ({
    draft: defaultSpec(),
    deployments: [] as Deployment[]
  }),
  getters: {
    isValid(s): boolean {
      const m = s.draft
      return !!(m.llm && m.vectordb && m.provider && (m.provider!=='onprem' ? m.region : true) && m.runtime)
    },
    summary(s) {
      const m = s.draft
      return [
        `App: ${m.slug || '—'}`,
        `LLM: ${m.llm || '—'}`,
        `Vector DB: ${m.vectordb || '—'}`,
        `Extras: ${m.extras.join(', ') || '—'}`,
        `Provider: ${m.provider || '—'} ${m.region ? '('+m.region+')' : ''}`,
        `Runtime: ${m.runtime || '—'}`,
        `Resources: ${m.resources.cpu} CPU / ${m.resources.ramGiB}GiB RAM / GPU ${m.resources.gpu || 'none'}`
      ].join('\n')
    },
    count: (s) => s.deployments.length,
    recent: (s) => [...s.deployments].sort((a,b)=>new Date(b.createdAt).getTime()-new Date(a.createdAt).getTime()),
    byId: (s) => (id: string) => s.deployments.find(d => d.id === id)
  },
  actions: {
    async fetchAll() {
      const arr = await $fetch<Deployment[]>('/api/deploy').catch(() => null)
      if (arr) this.deployments = arr
    },
    async createDeployment() {
      try {
        const dep = await $fetch<Deployment>('/api/deploy', { method:'POST', body:{ spec: this.draft } })
        this.deployments.unshift(dep); return dep
      } catch {
        // fallback (local) — desnoods laat je dit weg als je altijd API wilt
        const id = 'dep_' + Math.random().toString(36).slice(2)
        const now = new Date().toISOString()
        const dep: Deployment = { id, spec: JSON.parse(JSON.stringify(this.draft)), status:'running', createdAt: now, updatedAt: now, logs: [] }
        this.deployments.unshift(dep); return dep
      }
    },
    async redeployFrom(id: string) {
      const dep = await $fetch<Deployment>(`/api/deploy/${id}/redeploy`, { method:'POST' }).catch(() => null)
      if (dep) this.deployments.unshift(dep)
      return dep
    },
    async cancel(id: string) {
      const dep = await $fetch<Deployment>(`/api/deploy/${id}`, { method:'PATCH', body:{ action:'cancel' } }).catch(() => null)
      if (dep) {
        const i = this.deployments.findIndex(d=>d.id===id)
        if (i>=0) this.deployments[i] = dep
      }
    },

    initDraft(partial?: Partial<DeploySpec>) {
      this.draft = { ...defaultSpec(), ...(partial || {}) }
    },
    set<K extends keyof DeploySpec>(key: K, val: DeploySpec[K]) { (this.draft[key] as any) = val },
    setResource(partial: Partial<DeploySpec['resources']>) { this.draft.resources = { ...this.draft.resources, ...partial } },
    addExtra(x: string) { if (!this.draft.extras.includes(x)) this.draft.extras.push(x) },
    removeExtra(x: string) { this.draft.extras = this.draft.extras.filter(e => e !== x) },
    setApiKey(service: string, key: string) { this.draft.secrets.apiKeys[service] = key },

    clearDraft() { this.draft = defaultSpec() },

  

    updateStatus(id: string, status: DeployStatus) {
      const d = this.deployments.find(x => x.id === id); if (!d) return
      d.status = status; d.updatedAt = new Date().toISOString()
    },

    decommission(id: string) { this.updateStatus(id, 'decommissioned') },

    remove(id: string) { this.deployments = this.deployments.filter(d => d.id !== id) },

    
  }
})
