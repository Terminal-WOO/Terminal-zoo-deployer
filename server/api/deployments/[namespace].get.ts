import { defineEventHandler, getRouterParam, createError } from 'h3'

export default defineEventHandler(async (event) => {
  const cfg = useRuntimeConfig(event)
  const base = cfg.externalApiBase
  const auth = cfg.externalApiAuth

  const namespace = getRouterParam(event, 'namespace')
  if (!namespace) {
    throw createError({ statusCode: 400, statusMessage: 'Missing namespace' })
  }

  try {
    const res = await $fetch(`${base}/deployments/${encodeURIComponent(namespace)}`, {
      headers: { ...(auth ? { Authorization: auth } : {}) },
    })

    return res
  } catch (err: any) {
    throw createError({
      statusCode: err?.statusCode || err?.response?.status || 502,
      statusMessage: err?.statusMessage || 'Upstream API error (get deployments)',
      data: err?.data || err?.response?._data || null,
    })
  }
})