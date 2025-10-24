export default defineNuxtRouteMiddleware((_to) => {
  const user = useUserStore()
  if (user.token) { return navigateTo('/dashboard') }
})
