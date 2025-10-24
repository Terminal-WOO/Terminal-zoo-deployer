import { defineEventHandler } from 'h3'
import { listDeployments } from '../../utils/deployDb'

export default defineEventHandler(() => listDeployments())
