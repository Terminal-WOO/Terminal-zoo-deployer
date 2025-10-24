<template>
  <!-- linkerkolom wordt door /layouts/auth.vue gerenderd -->
  <section>
    <h2 class="heading-dark m-0 text-lg">Login to Dashboard</h2>
    <p class="mt-1 text-700">Vul onderstaande velden in om verder te gaan.</p>

    <form @submit.prevent="doLogin" novalidate>
      <div class="field">
        <label for="email">Email</label>
        <span class="p-input-icon-left w-full">
          <InputText
            id="email"
            v-model.trim="email"
            type="email"
            placeholder="Enter email address"
            autocomplete="email"
            class="w-full"
            @blur="touchedEmail=true"
            :aria-invalid="touchedEmail && !emailOk"
          />
        </span>
        <small v-if="touchedEmail && !emailOk" class="err">Vul een geldig e-mailadres in.</small>
      </div>

      <div class="field">
        <label for="password">Password</label>
        <span class="p-input-icon-left w-full">
          <InputText
            id="password"
            type="password"
            v-model="password"
            :feedback="false"
            placeholder="Enter Password"
            autocomplete="current-password"
            class="w-full"
            @blur="touchedPw=true"
            :aria-invalid="touchedPw && !pwOk"
          />
        </span>
        <small v-if="touchedPw && !pwOk" class="err">Wachtwoord is verplicht.</small>
      </div>

      <div class="flex justify-content-between align-items-center mb-2">
        <NuxtLink to="/auth/forgot" class="link">Forgot Password?</NuxtLink>
      </div>

      <Button
        type="submit"
        class="w-full"
        :loading="loading"
        :disabled="!canSubmit || loading"
        label="Login"
      />
      <Message v-if="error" severity="error" :closable="false" class="mt-2">{{ error }}</Message>
    </form>
  </section>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'auth', middleware: ['guest'] })

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '~/stores/user'

import type { User } from '~/types/types'

const router = useRouter()
const user = useUserStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const touchedEmail = ref(false)
const touchedPw = ref(false)
const emailOk = computed(() => /^\S+@\S+\.\S+$/.test(email.value))
const pwOk = computed(() => password.value.length > 0)
const canSubmit = computed(() => emailOk.value && pwOk.value)

const doLogin = async () => {
  error.value = ''
  touchedEmail.value = true
  touchedPw.value = true
  if (!canSubmit.value) return

  loading.value = true
  try {
    const res = await $fetch('/api/auth/login', {
      method: 'POST',
      body: { email: email.value, password: password.value }
    })

    user.login(res.user as User)
    user.token = res.token
    router.push('/dashboard')
  } catch (e:any) {
    error.value = e?.data?.message || e?.message || 'Inloggen mislukt'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* typografie */
:root { --ai-deep-green: #0E2D19; }
.heading-dark { color: var(--ai-deep-green); font-weight: 800; }
.text-lg { font-size: 1.125rem; }

.field { margin-bottom: 1rem; }
.err { color: var(--red-500); display: inline-block; margin-top: .25rem; }

/* socials */
.grid-soc { display: grid; grid-template-columns: 1fr 1fr; gap: .5rem; margin: 1rem 0 .5rem; }

/* divider “of” */
.div-or { position: relative; text-align: center; margin: .75rem 0 1rem; }
.div-or::before,
.div-or::after {
  content: ""; position: absolute; top: 50%; width: 40%; height: 1px;
  background: var(--surface-300);
}
.div-or::before { left: 0; }
.div-or::after { right: 0; }
.div-or > span { display: inline-block; padding: 0 .5rem; color: var(--text-color-secondary); }

/* links */
.link { text-decoration: underline; }
.w-full { width: 100%; }
</style>
