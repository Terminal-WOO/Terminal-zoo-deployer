import { defineEventHandler, readBody } from 'h3'
export default defineEventHandler(async (event) => {
  const { code } = await readBody(event)
  if (!code) throw createError({ statusCode:400, statusMessage:'Code ontbreekt' })
  if (code !== '') throw createError({ statusCode:400, statusMessage:'Code ongeldig' })
  return { ok: true }
})
