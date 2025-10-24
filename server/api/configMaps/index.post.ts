
type ConfigMapPayload = {
  name: string
  data: Record<string, string>
  namespace: string
}

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig(event)
  const base = config.externalApiBase || 'http://localhost:8080'
  const auth = config.externalApiAuth

  // read the incoming JSON from the client
  const body = (await readBody(event)) as ConfigMapPayload
  // minimal validation
  for (const key of ['name', 'data', 'namespace'] as const) {
    if (!body?.[key] || (key !== 'data' && typeof body[key] !== 'string')) {
      throw createError({ statusCode: 400, statusMessage: `Missing or invalid field: ${key}` })
    }
  }

  try {
    // forward to the external API
    const res = await $fetch<unknown>(`${base}/configmap/${body.namespace}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(auth ? { Authorization: auth } : {}),
      },
      body: {
        name: body.name!,
        data: body.data!,
      }
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