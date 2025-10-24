import { defineEventHandler, readBody, createError } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event)
  const base = config.externalApiBase || 'http://localhost:8080'
  const auth = config.externalApiAuth
  const namespace = event.context.params?.id
  if (!namespace || typeof namespace !== 'string') {
    throw createError({ statusCode: 400, statusMessage: `Missing or invalid namespace parameter` })
  }

  console.log("Fetching secrets for namespace:", namespace);

  try {
    // forward to the external API
    const res = await $fetch<unknown>(`${base}/secrets/${encodeURIComponent(namespace)}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        ...(auth ? { Authorization: auth } : {}),
      },
    })

    console.log("Received response:", res);

    return res
  }
  catch (err: any) {
    // normalize errors from the upstream API
    console.error("Error fetching secrets:", err);
    throw createError({
      statusCode: err?.statusCode || err?.response?.status || 502,
      statusMessage: err?.statusMessage || 'Upstream API error',
      data: err?.data || err?.response?._data || null,
    })
  }
})