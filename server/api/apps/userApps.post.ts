import fs from 'fs/promises'

export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  
  const storage = useStorage('assets:server')
  await storage.setItem('userApps.json', body)
  return { success: true }
})