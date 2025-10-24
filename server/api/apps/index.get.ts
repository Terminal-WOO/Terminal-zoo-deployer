import { defineEventHandler, getQuery } from 'h3'
import fs from 'fs/promises'
import { DUMMY_APPS } from '~/stores/_dummyApps'
import { AppModel } from '~/types/types'


export default defineEventHandler(async (event) => {
  const q = getQuery(event).q?.toString().toLowerCase() || ''
  const cat = getQuery(event).cat?.toString()

  let list: AppModel[] = [...DUMMY_APPS]
  const storage = useStorage('assets:server')

  let userApps = await storage.getItem<AppModel[]>('userApps.json') || []
  if (q) list = list.filter(a => [a.name, a.subtitle, a.summary].join(' ').toLowerCase().includes(q))
  if (cat) list = list.filter(a => a.categories.includes(cat))

  list = [...userApps, ...list] // User apps first
  return list
})
