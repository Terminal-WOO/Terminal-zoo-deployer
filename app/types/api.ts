import type { Port, Resources } from "./types"

export interface DeploymentRequest {
  deploymentName: string
  namespace: string
  image: string
  replicas: number
  ports: Port[]
  resources: Resources
}