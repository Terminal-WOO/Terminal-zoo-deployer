export function useClusters() {

  const createCluster = async (payload: object) => {
    const body = {
      name: payload.name,
      server: payload.server,
      caPEM: payload.caPEM,                       // string in Go, OK as-is
      bearerToken: payload.bearerToken,
      defaultNamespace: payload.defaultNamespace,
      domain: payload.domain,
      privateKey: payload.privateKey,                      // omitempty on the server
      Certificate: payload.Certificate
    };


    return await $fetch('/api/clusters', { method: 'POST', body: body })
  }

  const getClusters = async () => {
    return await $fetch(`/api/clusters`)
  }

  const getCluster = async (clusterName: string) => {
    return await $fetch(`/api/clusters/${encodeURIComponent(clusterName)}`)
  }

  const deleteCluster = async (clusterName: string) => {
    return await $fetch(`/api/clusters/${encodeURIComponent(clusterName)}`, {
      method: 'DELETE',
    })
  }

  return { createCluster, getClusters, getCluster, deleteCluster }
}