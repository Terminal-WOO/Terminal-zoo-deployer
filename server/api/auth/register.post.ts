import { defineEventHandler, readBody } from 'h3'
export default defineEventHandler(async (event) => {
  const { name, email, org, password } = await readBody(event)
  if (!name || !email || !password) throw createError({ statusCode:400, statusMessage:'Ontbrekende velden' })
  return { ok: true, verifyCode: '' } // stuur ‘e-mail’
})
