import { defineEventHandler, createError } from 'h3'

export default defineEventHandler(async (event) => {
  const cfg = useRuntimeConfig(event)
  const base = cfg.externalApiBase
  const auth = cfg.externalApiAuth

  try {
    const res = await $fetch(`${base}/clusters`, {
      headers: { ...(auth ? { Authorization: auth } : {}) },
    })

    console.log(res)

    return res
  } catch (err: any) {
    throw createError({
      statusCode: err?.statusCode || err?.response?.status || 502,
      statusMessage: err?.statusMessage || 'Upstream API error (get clusters)',
      data: err?.data || err?.response?._data || null,
    })
  }
})