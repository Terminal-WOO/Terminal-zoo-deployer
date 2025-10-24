import { defineEventHandler, readBody, setCookie, H3Event } from 'h3'

// For demo: load from Nitro assets (works in serverless bundles)
async function loadUsers() {
  const storage = useStorage('assets:server') // reads /server/assets
  const users = await storage.getItem<any[]>('users.json') // already parsed
  return Array.isArray(users) ? users : []
}


function setAuthCookie(event: H3Event, token: string, remember = true) {
  setCookie(event, 'auth', token, {
    httpOnly: true, sameSite: 'lax', path: '/',
    maxAge: remember ? 60 * 60 * 24 * 30 : undefined
  })
}

export default defineEventHandler(async (event) => {
  const { email, password, remember } = await readBody(event)
  if (!email || !password || String(password).length < 4) {
    throw createError({ statusCode: 400, statusMessage: 'Ongeldige inloggegevens' })
  }
  const token = 'tkn_' + Math.random().toString(36).slice(2)
  setAuthCookie(event, token, remember !== false)

  // Load user data from a database in a real application
  // Here we just simulate with dummy data
  const users = await loadUsers()
  const user = users.find((u: any) => u.email === email)
  if (!user) {
    throw createError({ statusCode: 401, statusMessage: 'Gebruiker niet gevonden' })
  }

  if (user.password !== password) {
    throw createError({ statusCode: 401, statusMessage: 'Onjuist wachtwoord' })
  }

  return { token, user: { id: user.id, name: user.name, email: user.email, roles: user.roles, organization: user.organization, website: user.website, token: token, org_type: user.org_type } }
})