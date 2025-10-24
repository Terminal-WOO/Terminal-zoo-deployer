
export function useConfigSecrets() {
  const createConfigMap = async (payload: { name: string; data: Record<string, string>; namespace: string }) => {
    return await $fetch('/api/configMaps', { method: 'POST', body: payload })
  }

  const readConfigMaps = async (namespace: string) => {
    return await $fetch(`/api/configMaps/${encodeURIComponent(namespace)}`, { method: 'GET' })
  }

  const createSecret = async (payload: { name: string; data: Record<string, string>; namespace: string }) => {
    return await $fetch('/api/secrets', { method: 'POST', body: payload })
  }

  const readSecrets = async (namespace: string) => {
    return await $fetch(`/api/configMaps/${encodeURIComponent(namespace)}`, { method: 'GET' })
  }

  return { createConfigMap, readConfigMaps, createSecret, readSecrets }
}
