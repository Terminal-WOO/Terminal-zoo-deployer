import { defineEventHandler, createError } from 'h3'
import { getDeployment, createDeployment, addLog, updateStatus } from '../../../../server/utils/deployDb'

export default defineEventHandler((event) => {
  const id = event.context.params!.id
  const src = getDeployment(id)
  if (!src) throw createError({ statusCode: 404, statusMessage: 'Source deployment not found' })
  const dep = createDeployment(src.spec)
  addLog(dep.id, `Redeploy van ${id}`)
  updateStatus(dep.id, 'deploying'); updateStatus(dep.id, 'running')
  return dep
})
