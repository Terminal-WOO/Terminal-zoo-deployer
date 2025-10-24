// Simple in-memory DB for deployments (dev/mock)
export type DeployStatus = 'queued'|'deploying'|'running'|'failed'|'canceled'|'decommissioned'
export interface DeploySpec {
  slug?: string | null
  llm: string | null
  vectordb: string | null
  extras: string[]
  provider: 'azure'|'hetzner'|'onprem'|null
  region: string | null
  runtime: 'k8s'|'docker'|null
  resources: { cpu: number; ramGiB: number; gpu?: 'none'|'a10'|'t4'|'l4' }
  env: { tenantId?: string; domain?: string }
  secrets: { apiKeys: Record<string,string> }
  notes: string
}
export interface Deployment {
  id: string
  spec: DeploySpec
  status: DeployStatus
  createdAt: string
  updatedAt: string
  logs: string[]
}

const db = { deployments: [] as Deployment[] }

export function listDeployments(): Deployment[] {
  return [...db.deployments].sort((a,b)=>new Date(b.createdAt).getTime()-new Date(a.createdAt).getTime())
}
export function getDeployment(id: string): Deployment | undefined {
  return db.deployments.find(d => d.id === id)
}
export function createDeployment(spec: DeploySpec): Deployment {
  const id = 'dep_' + Math.random().toString(36).slice(2)
  const now = new Date().toISOString()
  const dep: Deployment = { id, spec, status:'queued', createdAt: now, updatedAt: now, logs: [] }
  db.deployments.unshift(dep)
  return dep
}
export function updateStatus(id: string, status: DeployStatus) {
  const d = getDeployment(id); if (!d) return
  d.status = status; d.updatedAt = new Date().toISOString()
}
export function addLog(id: string, line: string) {
  const d = getDeployment(id); if (!d) return
  d.logs.push(`[${new Date().toISOString()}] ${line}`)
  d.updatedAt = new Date().toISOString()
}
