import { defineEventHandler, createError } from 'h3'
import { getDeployment } from '../../../server/utils/deployDb'

export default defineEventHandler((event) => {
  const id = event.context.params!.id
  const d = getDeployment(id)
  if (!d) throw createError({ statusCode: 404, statusMessage: 'Not found' })
  return d
})
