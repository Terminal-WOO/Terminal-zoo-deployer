<template>
  <div>
    <Card>
      <template #title>E-mail verifiëren</template>
      <template #content>
        <Message v-if="msg" :severity="ok ? 'success' : 'warn'" :closable="false">{{ msg }}</Message>

        <p class="text-700">Code: <b>{{ code }}</b></p>
        <Button label="Verifieer" icon="pi pi-check-circle" :loading="loading" @click="verify" />
        <Divider />
        <NuxtLink to="/auth/login">Verder naar inloggen</NuxtLink>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'auth', middleware: ['guest'] })

import { ref } from 'vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const code = route.params.code as string
const msg = ref(''); const ok = ref(false); const loading = ref(false)

const verify = async () => {
  msg.value=''; ok.value=false; loading.value=true
  try {
    await $fetch('/api/auth/verify', { method:'POST', body:{ code } })
    ok.value = true; msg.value = 'E-mail geverifieerd.'
  } catch (e:any) {
    msg.value = e?.data?.message || 'Verificatie mislukt'
  } finally { loading.value=false }
}
</script>

<style scoped>
</style>
