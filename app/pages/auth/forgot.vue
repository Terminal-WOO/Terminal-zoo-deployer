<template>
  <div>
    <Card>
      <template #title>Wachtwoord resetten</template>
      <template #content>
        <Message v-if="msg" :severity="ok ? 'success' : 'info'" :closable="false">{{ msg }}</Message>

        <div class="field">
          <label for="email">E-mail</label>
          <InputText id="email" v-model="email" type="email" placeholder="naam@organisatie.nl" />
        </div>

        <Button class="w-full" label="Stuur reset-link" icon="pi pi-envelope" :loading="loading" @click="send" />
        <Divider />
        <NuxtLink to="/auth/login">Terug naar inloggen</NuxtLink>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'auth', middleware: ['guest'] })

import { ref } from 'vue'
const email = ref(''); const msg = ref(''); const ok = ref(false); const loading = ref(false)

const send = async () => {
  msg.value = ''; ok.value=false; loading.value=true
  try {
    const res = await $fetch('/api/auth/forgot', { method:'POST', body:{ email: email.value } })
    ok.value = true
    msg.value = 'Als je e-mail bestaat, ontvang je zo een link. Voor demo: reset token = ' + res.resetToken
  } catch (e:any) {
    msg.value = e?.data?.message || 'Versturen mislukt'
  } finally { loading.value=false }
}
</script>

<style scoped>
.field { margin-bottom: 1rem; }
</style>
