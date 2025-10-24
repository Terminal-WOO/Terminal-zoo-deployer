import { defineEventHandler, deleteCookie } from 'h3'
export default defineEventHandler((event) => {
  deleteCookie(event, 'auth', { path: '/' })
  return { ok: true }
})
