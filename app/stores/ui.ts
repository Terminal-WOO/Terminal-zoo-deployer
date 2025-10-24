import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', {
  state: () => ({
    dark: false as boolean,
    locale: 'nl' as 'nl'|'en',
    sidebarOpen: false as boolean,
    loading: 0 as number,
    toastQueue: [] as { id:string; severity?:'success'|'info'|'warn'|'error'; summary:string; detail?:string }[]
  }),
  getters: {
    isLoading: (s) => s.loading > 0
  },
  actions: {
    init() { /* persist handled by plugin */ },
    setDark(v: boolean) { this.dark = v },
    toggleDark() { this.dark = !this.dark },
    setLocale(l: 'nl'|'en') { this.locale = l },
    openSidebar() { this.sidebarOpen = true },
    closeSidebar() { this.sidebarOpen = false },
    withLoading<T>(p: Promise<T>) {
      this.loading++; return p.finally(() => { this.loading = Math.max(0, this.loading - 1) })
    },
    pushToast(t: Omit<this['toastQueue'][number],'id'>) {
      const id = Math.random().toString(36).slice(2); this.toastQueue.push({ id, ...t }); return id
    },
    removeToast(id: string) { this.toastQueue = this.toastQueue.filter(t => t.id !== id) }
  }
})
