// script
import { defineStore } from 'pinia'
import type { User } from '~/types/types'

export const useUserStore = defineStore('user', {
  state: () => ({
    session: {} as User,
    roles: [] as string[],
    token: '' as string
  }),
  getters: {
    isAuthenticated: (state) => {
      return !!(state.token && state.session && state.session.id)
    },
  },
  actions: {
    fetchSession() { return this.session },
    hasRole(role: string): boolean { return this.roles.includes(role) },
    hasAnyRole(): boolean { return this.roles.length > 0 },

    login(userData: User) {
      this.session = userData
      this.roles = userData.roles
      this.token = userData.token || ''
    },
    logout() {
      this.session = {} as User
      this.roles = []
      this.token = ''
    }
  }
})