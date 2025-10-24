import { defineEventHandler, getRouterParam, createError } from 'h3'

export default defineEventHandler(async (event) => {
  const cfg = useRuntimeConfig(event)
  const base = cfg.externalApiBase
  const auth = cfg.externalApiAuth

  const clusterName = getRouterParam(event, 'clusterName')
  if (!clusterName) {
    throw createError({ statusCode: 400, statusMessage: 'Missing clusterName' })
  }

  try {
    const res = await $fetch(`${base}/clusters/${encodeURIComponent(clusterName)}`, {
      headers: { ...(auth ? { Authorization: auth } : {}) },
    })

    return res
  } catch (err: any) {
    throw createError({
      statusCode: err?.statusCode || err?.response?.status || 502,
      statusMessage: err?.statusMessage || 'Upstream API error (get clusters)',
      data: err?.data || err?.response?._data || null,
    })
  }
})