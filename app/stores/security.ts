import { defineStore } from 'pinia'

export type Severity = 'Low'|'Medium'|'High'|'Critical'
export interface Advisory { id:string; title:string; severity:Severity; component:string; published:string }

export const useSecurityStore = defineStore('security', {
  state: () => ({
    advisories: [] as Advisory[],
    filters: { severity: null as Severity|null, component: null as string|null }
  }),
  getters: {
    filtered(s): Advisory[] {
      return s.advisories.filter(a =>
        (!s.filters.severity || a.severity === s.filters.severity) &&
        (!s.filters.component || a.component === s.filters.component)
      )
    },
    severities: (s) => Array.from(new Set(s.advisories.map(a=>a.severity)))
  },
  actions: {
    init() { if (!this.advisories.length) this.seed() },
    seed() {
      const ago = (n:number,d:'h'|'d') => { const t=new Date(); d==='h'?t.setHours(t.getHours()-n):t.setDate(t.getDate()-n); return t.toISOString() }
      this.advisories = [
        { id:'ADV-2025-004', title:'Pinecone API key leak in CI logs', severity:'Critical', component:'Pinecone', published: ago(3,'h') },
        { id:'ADV-2025-003', title:'LiteLLM prompt logs onversleuteld op disk', severity:'Low', component:'LiteLLM', published: ago(12,'d') }
      ]
    },
    setFilter(partial: Partial<typeof this.filters>) { this.filters = { ...this.filters, ...partial } },
    clearFilters() { this.filters = { severity:null, component:null } }
  }
})
