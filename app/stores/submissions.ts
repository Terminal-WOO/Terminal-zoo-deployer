import { defineStore } from 'pinia'
import type { AppModel } from './_dummyApps'

export type SubmissionStatus = 'draft'|'submitted'|'under_review'|'changes_requested'|'rejected'|'verified'|'published'

export interface Submission {
  id: string
  app: Partial<AppModel> & { slug: string, name: string }
  status: SubmissionStatus
  createdAt: string
  updatedAt: string
  criteria: {
    checks: Record<string,'ok'|'missing'|'na'>
    evidence: Record<string,{ url:string; note:string }>
  }
  reviewerNotes?: string
}

const defaultCriteria = () => ({ checks:{}, evidence:{} })

export const useSubmissionsStore = defineStore('submissions', {
  state: () => ({ list: [] as Submission[] }),
  getters: {
    recent: s => [...s.list].sort((a,b)=>new Date(b.updatedAt).getTime()-new Date(a.updatedAt).getTime())
  },
  actions: {
    initDemo() {
      if (this.list.length) return
      this.list.push({
        id:'sub_1',
        app: { slug:'geo-buddy', name:'Geo Buddy' },
        status:'under_review',
        createdAt:new Date(Date.now()-2*864e5).toISOString(),
        updatedAt:new Date(Date.now()-1*864e5).toISOString(),
        criteria: defaultCriteria(),
        reviewerNotes:'Graag licentiematrix aanvullen.'
      })
    },
    createDraft(from?: { slug?: string; name?: string }) {
      const id = 'sub_' + Math.random().toString(36).slice(2)
      const now = new Date().toISOString()
      const sub: Submission = {
        id,
        app: { slug: from?.slug || 'nieuw-model', name: from?.name || 'Nieuw model' },
        status:'draft',
        createdAt: now, updatedAt: now,
        criteria: defaultCriteria()
      }
      this.list.unshift(sub)
      return sub
    },
    submit(id: string) {
      const s = this.list.find(x=>x.id===id); if(!s) return
      s.status = 'submitted'; s.updatedAt = new Date().toISOString()
    },
    setStatus(id: string, status: SubmissionStatus, notes?: string) {
      const s = this.list.find(x=>x.id===id); if(!s) return
      s.status = status; s.reviewerNotes = notes || s.reviewerNotes
      s.updatedAt = new Date().toISOString()
    }
  }
})
