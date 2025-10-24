import { defineEventHandler, readBody, createError } from 'h3'

type ClusterPayload = {
  name: string
  server: string
  caPEM: string
  bearerToken: string
  defaultNamespace?: string
  domain?: string
  Certificate?: string
  privateKey?: string
}

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event)
  const base = config.externalApiBase || 'http://localhost:8080'
  const auth = config.externalApiAuth

  // read the incoming JSON from the client
  const body = (await readBody(event)) as Partial<ClusterPayload>

  // minimal validation
  for (const key of ['name', 'server', 'caPEM', 'bearerToken'] as const) {
    if (!body?.[key] || typeof body[key] !== 'string') {
      throw createError({ statusCode: 400, statusMessage: `Missing or invalid field: ${key}` })
    }
  }

  try {
    // forward to the external API
    const res = await $fetch<unknown>(`${base}/clusters`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(auth ? { Authorization: auth } : {}),
      },
      body: {
        name: body.name!,
        server: body.server!,
        caPEM: body.caPEM!,
        bearerToken: body.bearerToken!,
        defaultNamespace: body.defaultNamespace ?? 'testing',
        Certificate: body.Certificate,
        privateKey: body.privateKey,
        domain: body.domain ?? '',
      },
      // optional: simple retry on network hiccups
      retry: 1,
    })

    return res
  } catch (err: any) {
    // normalize errors from the upstream API
    throw createError({
      statusCode: err?.statusCode || err?.response?.status || 502,
      statusMessage: err?.statusMessage || 'Upstream API error',
      data: err?.data || err?.response?._data || null,
    })
  }
})