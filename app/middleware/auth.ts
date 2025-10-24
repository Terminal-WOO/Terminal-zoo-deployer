export default defineNuxtRouteMiddleware((to) => {
  const user = useUserStore()
  if (!user.token) {
    return navigateTo({ path: '/auth/login', query: { next: to.fullPath } })
  }
})