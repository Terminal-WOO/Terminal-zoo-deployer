import { defineStore } from 'pinia'

export interface Thread {
  id: string
  title: string
  author: string
  category: 'Privacy'|'Security'|'Techniek'|'Ethiek'
  status: 'Open'|'Gesloten'|'Opgelost'
  replies: number
  app: string
  updatedAt: string
}
export interface Review {
  id: string; user: string; org: string; app: string; stars: number; comment: string; labels: string[]; date: string
}
export interface CriteriaEvent {
  id: string; title: string; domain: string; labels: string[]; description: string; date: string; icon?: string; color?: string
}
export interface Advisory {
  id:string; title:string; severity:'Low'|'Medium'|'High'|'Critical'; component:string; published:string
}

export const useCommunityStore = defineStore('community', {
  state: () => ({
    threads: [] as Thread[],
    reviews: [] as Review[],
    criteriaEvents: [] as CriteriaEvent[],
    leaderboard: [] as { id:string; name:string; initials:string; org:string; score:number; badges:string[] }[],
    advisories: [] as Advisory[]
  }),
  getters: {
    byCategory: (s) => (c: Thread['category']) => s.threads.filter(t => t.category === c),
    advisoryFilters: (s) => ({
      severities: Array.from(new Set(s.advisories.map(a=>a.severity)))
    })
  },
  actions: {
    init() {
      if (!this.threads.length) this.seed()
    },
    seed() {
      const ago = (n:number, unit:'h'|'d'|'w') => {
        const d = new Date()
        if (unit==='h') d.setHours(d.getHours()-n)
        if (unit==='d') d.setDate(d.getDate()-n)
        if (unit==='w') d.setDate(d.getDate()-(n*7))
        return d.toISOString()
      }
      this.threads = [
        { id:'t1', title:'DPIA-resultaten voor GovChat NL', author:'Anouk', category:'Privacy', status:'Open', replies:8, app:'GovChat NL', updatedAt: ago(2,'d') },
        { id:'t2', title:'API rate limiting n8n agent', author:'Ruben', category:'Techniek', status:'Opgelost', replies:5, app:'Libre Infer', updatedAt: ago(4,'h') },
        { id:'t3', title:'TLS-config Hetzner ingress', author:'Milan', category:'Security', status:'Open', replies:12, app:'Geo Buddy', updatedAt: ago(1,'d') }
      ]
      this.reviews = [
        { id:'r1', user:'Eva Janssen', org:'Gemeente Delft', app:'GovChat NL', stars:5, comment:'Sterk NL en goede bronvermelding.', labels:['DPIA OK','A11y'], date:new Date().toISOString() },
        { id:'r2', user:'Tom de Vries', org:'Provincie Utrecht', app:'Geo Buddy', stars:4, comment:'Kaarten goed; performance kan beter.', labels:['Open-source'], date:new Date().toISOString() }
      ]
      this.criteriaEvents = [
        { id:'c1', title:'Nieuwe paragraaf mensenrechtentoets', domain:'Ethiek', labels:['FRIA/IAMA'], description:'Toevoeging met voorbeelden.', date: ago(3,'d'), icon:'pi pi-book', color:'var(--primary-color)' }
      ]
      this.leaderboard = [
        { id:'u1', name:'Anouk Visser', initials:'AV', org:'Gemeente Delft', score:1280, badges:['Privacy','DPIA OK','Reviewer'] }
      ]
      this.advisories = [
        { id:'ADV-2025-004', title:'Pinecone API key leak in CI logs', severity:'Critical', component:'Pinecone', published: ago(3,'h') },
        { id:'ADV-2025-002', title:'n8n webhook exposure zonder IP-allowlist', severity:'High', component:'n8n', published: ago(6,'d') }
      ]
    },
    addThread(t: Omit<Thread,'id'|'updatedAt'|'replies'|'status'> & Partial<Pick<Thread,'status'|'replies'>>) {
      const id = 't' + Math.random().toString(36).slice(2)
      this.threads.unshift({ id, replies:0, status:'Open', updatedAt:new Date().toISOString(), ...t } as Thread)
      return id
    },
    addReply(threadId: string) {
      const t = this.threads.find(x=>x.id===threadId); if (t) { t.replies++; t.updatedAt = new Date().toISOString() }
    },
    addAdvisory(a: Advisory) { this.advisories.unshift(a) }
  }
})
