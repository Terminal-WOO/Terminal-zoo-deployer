<template>
  <div>
    <Card>
      <template #title>Account aanmaken</template>
      <template #content>
        <Message v-if="msg" :severity="ok ? 'success' : 'warn'" :closable="false">{{ msg }}</Message>

        <div class="field">
          <label for="name">Naam</label>
          <InputText id="name" v-model="name" placeholder="Voornaam Achternaam" />
        </div>

        <div class="field">
          <label for="email">E-mail</label>
          <InputText id="email" v-model="email" type="email" placeholder="naam@organisatie.nl" />
        </div>

        <div class="field">
          <label for="org">Organisatie</label>
          <InputText id="org" v-model="org" placeholder="Bijv. Gemeente Delft" />
        </div>

        <div class="field">
          <label for="pwd">Wachtwoord</label>
          <Password id="pwd" v-model="password" :feedback="true" toggleMask />
        </div>

        <Button class="mt-2 w-full" label="Account aanmaken" icon="pi pi-user-plus" :loading="loading" @click="doRegister" />
        <Divider />
        <div class="text-center">
          <small class="text-700">Heb je al een account?</small>
          <NuxtLink to="/auth/login" class="ml-2">Inloggen</NuxtLink>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'auth', middleware: ['guest'] })

import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const name = ref(''); const email = ref(''); const org = ref(''); const password = ref('')
const msg = ref(''); const ok = ref(false)

const doRegister = async () => {
  msg.value = ''; ok.value = false; loading.value = true
  try {
    await $fetch('/api/auth/register', { method: 'POST', body: { name: name.value, email: email.value, org: org.value, password: password.value } })
    ok.value = true
    msg.value = 'Registratie gelukt. We hebben je een verificatie-e-mail gestuurd.'
    setTimeout(() => { router.push('/auth/login') }, 1200)
  } catch (e:any) {
    msg.value = e?.data?.message || e?.message || 'Registreren mislukt'
  } finally { loading.value = false }
}
</script>

<style scoped>
.field { margin-bottom: 1rem; }
</style>
