import { defineEventHandler, createError } from 'h3'
import { getDummyBySlug } from '~/stores/_dummyApps'

export default defineEventHandler((event) => {
  const slug = event.context.params?.slug as string
  const found = getDummyBySlug(slug)
  if (!found) {
    // 404 zodat fetchOne fallback kan doen
    throw createError({ statusCode: 404, statusMessage: 'Not found' })
  }
  return found
})
