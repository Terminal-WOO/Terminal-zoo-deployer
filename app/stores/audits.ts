// /stores/audits.ts
import { defineStore } from 'pinia'

export type AuditDomain = 'privacy'|'security'|'ethics'|'accessibility'|'licensing'|'ai_risk'
export type Severity = 'info'|'low'|'medium'|'high'|'critical'
export type FindingStatus = 'ok'|'warning'|'fail'
export type AuditType = 'initial'|'periodic'|'incident'|'retest'
export type AuditStatus = 'passed'|'passed_with_notes'|'attention'|'failed'|'pending'

export interface AuditFinding {
  id: string
  criterionKey: string              // bijv. 'dpia', 'a11y'
  title: string
  description: string
  severity: Severity
  status: FindingStatus
  evidenceUrl?: string
  remediation?: string
}

export interface Audit {
  id: string
  appSlug: string
  title: string
  type: AuditType
  date: string
  auditor: { name: string; org?: string }
  status: AuditStatus
  score: number                      // 0-100
  domains: AuditDomain[]
  findings: AuditFinding[]
  notes?: string
}

function daysAgo(n: number) {
  return new Date(Date.now() - n*864e5).toISOString()
}

export const useAuditStore = defineStore('audits', {
  state: () => ({
    audits: [] as Audit[],
    hydrated: false
  }),
  getters: {
    forApp: (s) => (slug: string) =>
      s.audits.filter(a => a.appSlug === slug)
              .sort((a,b)=>new Date(b.date).getTime()-new Date(a.date).getTime()),
    latestForApp: (s) => (slug: string) =>
      s.audits.filter(a => a.appSlug === slug)
              .sort((a,b)=>new Date(b.date).getTime()-new Date(a.date).getTime())[0] || null,
    scoreForApp: (s) => (slug: string) => {
      const items = s.audits.filter(a => a.appSlug === slug)
      if (!items.length) return null
      // gewogen: laatste audit telt dubbel
      const [latest, ...rest] = items.sort((a,b)=>new Date(b.date).getTime()-new Date(a.date).getTime())
      const sum = latest.score*2 + rest.reduce((acc,a)=>acc+a.score,0)
      return Math.round(sum / (rest.length + 2))
    }
  },
  actions: {
    initDemo() {
      if (this.hydrated) return
      this.audits = [
        // GovChat NL — initial + periodic + retest
        {
          id: 'aud_govchat_initial',
          appSlug: 'govchat-nl',
          title: 'Initiële peer-review & criteria check',
          type: 'initial',
          date: daysAgo(120),
          auditor: { name: 'Peer review team', org: 'Community' },
          status: 'passed_with_notes',
          score: 86,
          domains: ['privacy','licensing','ai_risk','accessibility','security'],
          findings: [
            { id:'f1', criterionKey:'dpia', title:'DPIA beschikbaar', description:'DPIA is aangeleverd en beoordeeld.', severity:'info', status:'ok', evidenceUrl:'/docs/dpia.pdf' },
            { id:'f2', criterionKey:'license_clarity', title:'Licenties helder', description:'Licentiematrix compleet.', severity:'low', status:'ok' },
            { id:'f3', criterionKey:'a11y', title:'Toegankelijkheid', description:'WCAG 2.1 AA grotendeels ok; contrast in 2 schermen aan te passen.', severity:'low', status:'warning', remediation:'Contrast fix bij volgende release.' },
            { id:'f4', criterionKey:'ai_risk_cat', title:'AI-risicocategorie onderbouwing', description:'Ingevuld conform handreiking.', severity:'info', status:'ok' }
          ],
          notes: 'Goedgekeurd met minor opmerkingen (a11y contrast).'
        },
        {
          id: 'aud_govchat_periodic',
          appSlug: 'govchat-nl',
          title: 'Periodieke audit Q2',
          type: 'periodic',
          date: daysAgo(45),
          auditor: { name: 'Ethische hackers', org: 'Testteam NL' },
          status: 'attention',
          score: 80,
          domains: ['security','privacy','ai_risk'],
          findings: [
            { id:'f5', criterionKey:'security_headers', title:'Security headers', description:'CSP mist report-uri.', severity:'medium', status:'warning', remediation:'CSP header uitbreiden en monitoren.' },
            { id:'f6', criterionKey:'logging', title:'Auditlog volledigheid', description:'Niet alle admin-acties worden gelogd.', severity:'medium', status:'warning', remediation:'Admin-events toevoegen aan auditlog.' }
          ],
          notes: 'Geen blocker; aanbevelingen opnemen in volgende sprint.'
        },
        {
          id: 'aud_govchat_retest',
          appSlug: 'govchat-nl',
          title: 'Hertest contrast + logging',
          type: 'retest',
          date: daysAgo(10),
          auditor: { name: 'Peer review team', org: 'Community' },
          status: 'passed',
          score: 92,
          domains: ['accessibility','security'],
          findings: [
            { id:'f7', criterionKey:'a11y', title:'Contrast hersteld', description:'Contrastproblemen opgelost in release 1.4.0.', severity:'info', status:'ok' },
            { id:'f8', criterionKey:'logging', title:'Auditlog uitgebreid', description:'Admin-events worden nu gelogd.', severity:'info', status:'ok' }
          ]
        },

        // Geo Buddy — initial
        {
          id: 'aud_geo_initial',
          appSlug: 'geo-buddy',
          title: 'Initiële criteria check',
          type: 'initial',
          date: daysAgo(60),
          auditor: { name: 'Peer review team', org:'Community' },
          status: 'attention',
          score: 74,
          domains: ['privacy','licensing','security'],
          findings: [
            { id:'g1', criterionKey:'dpia', title:'DPIA ontbreekt', description:'Nog niet aangeleverd.', severity:'high', status:'fail', remediation:'DPIA opstellen en aanleveren.' },
            { id:'g2', criterionKey:'license_clarity', title:'Licentie onduidelijk', description:'Onduidelijk gebruik van kaart-tiles.', severity:'medium', status:'warning' }
          ]
        },

        // Parlement Agent — initial
        {
          id:'aud_parl_initial',
          appSlug:'parlement-agent',
          title:'Initiële peer-review',
          type:'initial',
          date: daysAgo(30),
          auditor:{ name:'Peer review team' },
          status:'passed',
          score: 90,
          domains:['privacy','licensing','ai_risk'],
          findings:[
            { id:'p1', criterionKey:'dpia', title:'DPIA OK', description:'Goedgekeurd.', severity:'info', status:'ok' }
          ]
        }
      ]
      this.hydrated = true
    }
  }
})
