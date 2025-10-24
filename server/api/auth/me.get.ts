import { defineEventHandler, getCookie, createError } from 'h3'
export default defineEventHandler((event) => {
  const token = getCookie(event, 'auth')
  if (!token) throw createError({ statusCode: 401, statusMessage: 'Not authenticated' })
  return { token, user: { id: 'u1', name: 'Demo', email: 'demo@example.com', roles: ['member'] } }
})
