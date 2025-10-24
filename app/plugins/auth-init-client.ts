export default defineNuxtPlugin(async () => {
  const user = useUserStore()
  await user.fetchSession()
})