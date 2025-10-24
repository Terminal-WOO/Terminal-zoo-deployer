import { defineEventHandler, getQuery } from 'h3'
import fs from 'fs/promises'
import { AppModel } from '~/types/types'


export default defineEventHandler(async (event) => {
  const q = getQuery(event).q?.toString().toLowerCase() || ''
  const cat = getQuery(event).cat?.toString()

  // Always read fresh from disk
  const storage = useStorage('assets:server')
  const userApps = await storage.getItem<AppModel[]>('userApps.json') || []
  let list: AppModel[] = [...userApps]

  if (q) list = list.filter(a => [a.name, a.subtitle, a.summary].join(' ').toLowerCase().includes(q))
  if (cat) list = list.filter(a => a.categories.includes(cat))

  return list
})
