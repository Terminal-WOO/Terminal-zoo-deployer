import { defineStore } from 'pinia'
import { $fetch } from 'ofetch'
import { DUMMY_APPS, getDummyBySlug } from './_dummyApps'
import type { AppModel } from '~/types/types'

export const useAppCatalog = defineStore('catalog', {
  state: () => ({
    apps: [] as AppModel[],
    userApps: [] as AppModel[],
    loading: false,
    error: '',
    q: '' as string,
    selectedCats: [] as string[],
    sort: 'name' as 'name'|'updated',
    page: 1,
    pageSize: 12
  }),
  getters: {
    filtered(state): AppModel[] {
      let list = [...state.apps]
      if (state.q) {
        const q = state.q.toLowerCase()
        list = list.filter(a => [a.name, a.subtitle, a.summary].join(' ').toLowerCase().includes(q))
      }
      if (state.selectedCats.length) {
        list = list.filter(a => state.selectedCats.every(c => a.categories.includes(c)))
      }
      if (state.sort === 'name') list.sort((a,b)=>a.name.localeCompare(b.name))
      if (state.sort === 'updated') list.sort((a,b)=>new Date(b.updatedAt||0).getTime()-new Date(a.updatedAt||0).getTime())
      return list
    },
    paged(): AppModel[] {
      const start = (this.page-1)*this.pageSize
      return this.filtered.slice(start, start + this.pageSize)
    },
    bySlug: (s) => (slug: string) => s.apps.find(a => a.slug === slug)
  },
  actions: {
    async fetchAll(params?: { q?: string; cat?: string; includeUnpublished?: boolean }) {
      this.loading = true; this.error = ''
      try {
        const api = await $fetch<AppModel[]>('/api/apps', { params }).catch(() => null)
        let list = api && api.length ? api : [...DUMMY_APPS]
        if (!params?.includeUnpublished) list = list.filter(a => a.approved === 'published')
        this.apps = list
      } catch (e:any) {
        console.error('Error fetching apps:', e)
        this.error = e?.message || 'Kon apps niet laden'
        this.apps = [...DUMMY_APPS].filter(a => a.approved === 'published')
      } finally { this.loading = false }
    },
    async fetchUserApps(params?: { q?: string; cat?: string; includeUnpublished?: boolean }) {
      this.loading = true; this.error = ''
      try {
        const api = await $fetch<AppModel[]>('/api/apps/userApps', { params }).catch(() => null)
        this.userApps = api && api.length ? api : []
      } catch (e:any) {
        console.error('Error fetching user apps:', e)
        this.error = e?.message || 'Kon user apps niet laden'
        this.userApps = []
      } finally { this.loading = false }
    },
    async deleteApp() {
      await this.fetchUserApps()
      await this.fetchAll()
    },
    async fetchOne(slug: string) {
      let app = this.bySlug(slug)
      if (app) return app
      try {
        app = await $fetch<AppModel>(`/api/apps/${slug}`).catch(() => null as any)
      } catch { /* noop */ }
      if (!app) app = getDummyBySlug(slug) || { ...DUMMY_APPS[0], slug }
      this.apps = uniqueBySlug([app, ...this.apps])
      return app
    },
  }
})

function uniqueBySlug(list: AppModel[]) {
  const map = new Map<string,AppModel>(); list.forEach(a=>map.set(a.slug,a)); return [...map.values()]
}
