// /stores/criteria.ts
import { defineStore } from 'pinia'

export type CriteriaKey =
  | 'open_source' | 'modular' | 'dpia' | 'license_clarity' | 'a11y' | 'ai_risk_cat'
  | 'functional_desc' | 'self_hostable' | 'tech_explainer'

export type Status = 'ok'|'missing'|'na'
export type Domain = 'Privacy'|'Security'|'Techniek'|'Ethiek'|'Juridisch'

export interface Criterion {
  k: CriteriaKey
  label: string
  domain: Domain
  desc: string
}
export interface CriteriaRecord {
  checks: Record<CriteriaKey, Status>
  evidence: Record<CriteriaKey, { url: string; note: string }>
  updatedAt?: string
}

const ALL: Criterion[] = [
  { k:'open_source', label:'Open-source', domain:'Juridisch', desc:'OSS licentie of gemotiveerde uitzondering.' },
  { k:'modular', label:'Modulair', domain:'Techniek', desc:'Vervangbare bouwstenen, geen lock-in.' },
  { k:'dpia', label:'DPIA', domain:'Privacy', desc:'Geldige DPIA aanwezig.' },
  { k:'license_clarity', label:'Licentiestructuur', domain:'Juridisch', desc:'Licenties per component/dataset helder.' },
  { k:'a11y', label:'Toegankelijkheid (WCAG)', domain:'Ethiek', desc:'Minimaal WCAG 2.1 AA, verslag beschikbaar.' },
  { k:'ai_risk_cat', label:'AI-risicocategorie', domain:'Ethiek', desc:'Geklasseerd + maatregelen.' },
  { k:'functional_desc', label:'Functionele beschrijving', domain:'Techniek', desc:'Doel, beperkingen, in-/output.' },
  { k:'self_hostable', label:'Lokaal te draaien', domain:'Techniek', desc:'Self-host optie of onderbouwing.' },
  { k:'tech_explainer', label:'Technische uitleg', domain:'Techniek', desc:'Architectuur, dataflows, logging.' }
]

const defaultChecks = (): CriteriaRecord['checks'] => ({
  open_source:'missing', modular:'missing', dpia:'missing', license_clarity:'missing',
  a11y:'missing', ai_risk_cat:'missing', functional_desc:'missing', self_hostable:'missing', tech_explainer:'missing'
})
const defaultEvidence = (): CriteriaRecord['evidence'] =>
  Object.fromEntries(ALL.map(c => [c.k, { url:'', note:'' }])) as CriteriaRecord['evidence']

export const useCriteriaStore = defineStore('criteria', {
  state: () => ({
    all: ALL as Criterion[],
    records: {} as Record<string, CriteriaRecord>, // 🔑 per appSlug
    activeStep: 0
  }),
  getters: {
    record: (s) => (appId: string): CriteriaRecord => {
      return s.records[appId] ?? { checks: defaultChecks(), evidence: defaultEvidence() }
    },
    progressPct: (s) => (appId: string) => {
      const r = s.records[appId] ?? { checks: defaultChecks(), evidence: defaultEvidence() }
      const ok = Object.values(r.checks).filter(v => v === 'ok').length
      return Math.round((ok / s.all.length) * 100)
    },
    byDomain: (s) => (d: Domain) => s.all.filter(c => c.domain === d),
    evidenceLabel: (s) => (appId: string, k: CriteriaKey) => {
      const e = (s.records[appId]?.evidence ?? defaultEvidence())[k]
      if (e.url) return 'URL toegevoegd'
      if (e.note) return 'Notitie'
      return 'Geen bewijs'
    },
    evidenceSeverity: (s) => (appId: string, k: CriteriaKey) => {
      const e = (s.records[appId]?.evidence ?? defaultEvidence())[k]
      if (e.url) return 'success'
      if (e.note) return 'info'
      return 'secondary'
    }
  },
  actions: {
    ensure(appId: string) {
      if (!this.records[appId]) this.records[appId] = {
        checks: defaultChecks(),
        evidence: defaultEvidence(),
        updatedAt: new Date().toISOString()
      }
    },
    setCheck(appId: string, k: CriteriaKey, v: Status) {
      this.ensure(appId)
      this.records[appId].checks[k] = v
      this.records[appId].updatedAt = new Date().toISOString()
    },
    setEvidence(appId: string, k: CriteriaKey, field: 'url'|'note', value: string) {
      this.ensure(appId)
      this.records[appId].evidence[k][field] = value
      this.records[appId].updatedAt = new Date().toISOString()
    },
    nextStep() { this.activeStep = Math.min(this.activeStep+1, 3) },
    prevStep() { this.activeStep = Math.max(this.activeStep-1, 0) },
    reset(appId: string) {
      this.records[appId] = { checks: defaultChecks(), evidence: defaultEvidence(), updatedAt: new Date().toISOString() }
      this.activeStep = 0
    }
  }
})
