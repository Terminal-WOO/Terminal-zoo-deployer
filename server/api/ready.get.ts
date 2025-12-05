import { defineEventHandler } from 'h3'

export default defineEventHandler(() => {
  // Nuxt app is ready if it can respond
  return { status: 'ready' }
})

