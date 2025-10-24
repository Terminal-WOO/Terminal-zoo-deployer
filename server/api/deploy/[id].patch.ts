import { defineEventHandler, readBody, createError } from 'h3'
import { getDeployment, updateStatus, addLog, type DeployStatus } from '../../../server/utils/deployDb'

export default defineEventHandler( async (event) => {
  const id = event.context.params!.id
  const d = getDeployment(id)
  if (!d) throw createError({ statusCode: 404, statusMessage: 'Not found' })

  const body = await readBody<{ action?: 'cancel'|'decommission'; status?: DeployStatus }>(event)
  if (body?.action === 'cancel') { updateStatus(id, 'canceled'); addLog(id, 'Deployment geannuleerd'); return getDeployment(id) }
  if (body?.action === 'decommission') { updateStatus(id, 'decommissioned'); addLog(id, 'Deployment uit dienst'); return getDeployment(id) }
  if (body?.status) { updateStatus(id, body.status); addLog(id, `Status → ${body.status}`); return getDeployment(id) }

  throw createError({ statusCode: 400, statusMessage: 'No action/status' })
})
