// composables/useDeployments.ts
import type { ManualDeploymentConfig } from '~/types/types'

export function useDeployments() {
  const createDeployment = async (payload: ManualDeploymentConfig) => {
    return await $fetch('/api/deployments', { method: 'POST', body: payload })
  }

  const getDeployments = async (namespace: string, clusterName?: string) => {
    return await $fetch(`/api/deployments/${encodeURIComponent(namespace)}` + (clusterName ? `?clusterName=${encodeURIComponent(clusterName)}` : ''))
  }

  const deleteDeployment = async (namespace: string, name: string) => {
    return await $fetch(`/api/deployments/${encodeURIComponent(namespace)}/${encodeURIComponent(name)}`, {
      method: 'DELETE',
    })
  }

  const validateResources = (resources: any): boolean => {
    if (typeof resources !== 'object' || resources === null) return false
    const keys = ['cpuRequests', 'cpuLimits', 'memoryRequests', 'memoryLimits']
    for (const key of keys) {
      if (typeof resources[key] !== 'string') return false
    }
    return true
  }

  return { createDeployment, getDeployments, deleteDeployment, validateResources }
}