import { defineEventHandler, readBody, createError } from 'h3'
import { validateResources } from "../../utils/validation"
import type { ManualDeploymentConfig } from '~/types/types'
import type { DeploymentRequest } from '~/types/api'

export default defineEventHandler(async (event) => {
  const cfg = useRuntimeConfig(event)
  const base = cfg.externalApiBase
  const auth = cfg.externalApiAuth

  const payload = (await readBody(event)) as ManualDeploymentConfig

  const selectedCluster = payload?.clusterName || null

  // Basic validation
  const required = ['deploymentName','namespace','containerImage'] as const
  for (const k of required) {
    if (!payload?.[k] || typeof payload[k] !== 'string') {
      throw createError({ statusCode: 400, statusMessage: `Missing or invalid field: ${k}` })
    }
  }
  if (typeof payload.replicas !== 'number' || payload.replicas < 0) {
    throw createError({ statusCode: 400, statusMessage: `Invalid replicas` })
  }
  if (!Array.isArray(payload.ports)) {
    throw createError({ statusCode: 400, statusMessage: `ports must be an array` })
  }
  if (!validateResources(payload.resources)) {
    throw createError({ statusCode: 400, statusMessage: 'Invalid resource format' })
  }

  const body: DeploymentRequest = {
    deploymentName: payload.deploymentName,
    namespace: payload.namespace,
    image: payload.containerImage,
    replicas: payload.replicas,
    ports: payload.ports,
    resources: payload.resources,
  }

  console.log('Creating deployment with payload:', body, 'on cluster:', selectedCluster || 'default')

  try {
    const res = await $fetch(`${base}/deployments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(auth ? { Authorization: auth } : {}),
        ...(selectedCluster ? { 'cluster-name': selectedCluster } : {}),
      },
      body,
    })

    console.log('Deployment created successfully:', res)
    return res
  } catch (err: any) {
    console.error('Error creating deployment:', err)
    throw createError({
      statusCode: err?.statusCode || err?.response?.status || 502,
      statusMessage: err?.statusMessage || 'Upstream API error (create deployment)',
      data: err?.data || err?.response?._data || null,
    })
  }
})
