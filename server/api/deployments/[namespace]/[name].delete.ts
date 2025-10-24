// server/api/deployments/[namespace]/[name].delete.ts
import { defineEventHandler, getRouterParam, createError } from 'h3'

export default defineEventHandler(async (event) => {
  const cfg = useRuntimeConfig(event)
  const base = cfg.externalApiBase
  const auth = cfg.externalApiAuth

  const namespace = getRouterParam(event, 'namespace')
  const name = getRouterParam(event, 'name')
  if (!namespace || !name) {
    throw createError({ statusCode: 400, statusMessage: 'Missing namespace or name' })
  }

  try {
    const res = await $fetch.raw(`${base}/deployments/${encodeURIComponent(namespace)}/${encodeURIComponent(name)}`, {
      method: 'DELETE',
      headers: { ...(auth ? { Authorization: auth } : {}) },
    })

    if (res.status === 204) {
      return { message: 'Deployment deleted successfully' }
    }
    // If upstream returns JSON body on delete, forward it:
    try {
      const data = await res.json()
      return data
    } catch {
      throw createError({ statusCode: res.status, statusMessage: `Upstream returned ${res.status}` })
    }
  } catch (err: any) {
    throw createError({
      statusCode: err?.statusCode || err?.response?.status || 502,
      statusMessage: err?.statusMessage || 'Upstream API error (delete deployment)',
      data: err?.data || err?.response?._data || null,
    })
  }
})
