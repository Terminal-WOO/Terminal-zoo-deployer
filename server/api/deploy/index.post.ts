// /server/api/deploy/index.post.ts
import { defineEventHandler, readBody, createError } from 'h3'
import { createDeployment, updateStatus, addLog } from '../../utils/deployDb' // ⬅️ let op: ../../ i.p.v. ../../../

export default defineEventHandler(async (event) => {
  const body = await readBody<any>(event)

  // ✅ accepteer beide vormen: { spec: {...} } of direct {...}
  const spec = body?.spec ?? body
  if (!spec || typeof spec !== 'object') {
    throw createError({ statusCode: 400, statusMessage: 'Missing spec' })
  }

  // nette validatie met concrete melding
  const required = ['llm','vectordb','provider','runtime'] as const
  const missing = required.filter(k => !spec[k])
  if (missing.length) {
    throw createError({ statusCode: 422, statusMessage: 'Incomplete spec: ' + missing.join(', ') })
  }

  const dep = createDeployment(spec)
  addLog(dep.id, `Deployment aangemaakt voor ${spec.slug || 'onbekend'}`)
  updateStatus(dep.id, 'deploying'); addLog(dep.id, 'Container build gestart')
  updateStatus(dep.id, 'running');   addLog(dep.id, 'Dienst actief')

  return dep
})
