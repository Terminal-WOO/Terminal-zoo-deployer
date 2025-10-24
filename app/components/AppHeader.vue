<template>
  <header
    class=""
    :class="[{ 'header-sticky': true }]"
    role="banner"
  >
    <!-- Desktop (MegaMenu) -->
    <div class="container pt-2 desktop-only">
      <MegaMenu :model="megaModel" orientation="horizontal">
        <!-- Brand (links) -->
        <template #start>
          <BrandLogo size="sm" @click="goHome" class="cursor-pointer mr-3" />
        </template>

        <!-- Einde (rechts): zoek + dashboard -->
        <template #end>
          <div class="flex align-items-center gap-2">
            <div v-if="showSearch" class="p-inputgroup">
              <InputText
                v-model="q"
                placeholder="Zoek in catalogus"
                aria-label="Zoekterm"
                @keyup.enter="doSearch"
              />
              <Button icon="pi pi-search" @click="doSearch" aria-label="Zoeken" />
            </div>
            <NuxtLink to="/dashboard" class="no-underline">
              <Button label="Dashboard" icon="pi pi-user" size="small" />
            </NuxtLink>
          </div>
        </template>
      </MegaMenu>
    </div>

    <!-- Mobiel (Sidebar blijft) -->
    <div class="container p-2 mobile-only">
      <div class="flex align-items-center justify-content-between">
        <div @click="goHome" aria-label="Ga naar startpagina">
          <BrandLogo size="md" />
        </div>

        <div class="flex align-items-center gap-2">
          <Button v-if="showSearch" icon="pi pi-search" text rounded aria-label="Zoeken" @click="mobileOpen = true" />
          <Button icon="pi pi-bars" text rounded aria-label="Menu" @click="toggleMobile" />
        </div>
      </div>

      <Sidebar v-model:visible="mobileOpen" position="right" :dismissable="true" aria-label="Mobiele navigatie">
        <div class="flex flex-column gap-3">
          <div v-if="showSearch" class="p-inputgroup">
            <InputText v-model="q" placeholder="Zoeken…" aria-label="Zoekterm" @keyup.enter="doSearch" />
            <Button icon="pi pi-search" @click="doSearch" aria-label="Zoeken" />
          </div>

          <Divider />

          <nav class="flex flex-column">
            <NuxtLink
              v-for="it in flatLinks"
              :key="it.to"
              :to="it.to"
              class="mobile-link"
              :class="{ 'is-active': isActive(it.to) }"
              @click="mobileOpen=false"
            >
              <i :class="it.icon" class="mr-2" />
              {{ it.label }}
            </NuxtLink>
          </nav>

          <Divider />

          <NuxtLink to="/dashboard" class="no-underline">
            <Button label="Dashboard" icon="pi pi-user" class="w-full" />
          </NuxtLink>
        </div>
      </Sidebar>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAppCatalog } from '~/stores/appCatalog'

const props = withDefaults(defineProps<{
  brand?: string
  sticky?: boolean
  elevated?: boolean
  showSearch?: boolean
}>(), { brand: 'IPO appstore', sticky: true, elevated: true, showSearch: true })

const emit = defineEmits<{ (e:'search', q:string): void }>()

const router = useRouter()
const route = useRoute()
const catalog = useAppCatalog()

const q = ref('')
const mobileOpen = ref(false)

/** MegaMenu model
 * Alleen het EERSTE item (App-store) heeft kolommen voor de mega dropdown.
 * De rest zijn directe links zonder submenu.
 */
import { useUserStore } from '~/stores/user'
const userStore = useUserStore()
const megaModel = computed(() => {
  const items = [
    {
      label: 'App-store',
      icon: 'pi pi-th-large',
      items: [
        [
          {
            label: 'Verkennen',
            items: [
              { label: 'Alle apps', icon: 'pi pi-list', command: () => nav('/apps') },
              { label: 'Chatbots', icon: 'pi pi-comments', command: () => nav('/apps?cat=Chatbot') },
              { label: 'RAG', icon: 'pi pi-book', command: () => nav('/apps?cat=RAG') },
              { label: 'Agents', icon: 'pi pi-bolt', command: () => nav('/apps?cat=Agents') },
              { label: 'Nederlands', icon: 'pi pi-flag', command: () => nav('/apps?cat=Nederlands') }
            ]
          }
        ],
        [
          {
            label: 'Ontwikkeling door',
            items: [
              { label: 'Overheidsontwikkeld', icon: 'pi pi-shield', command: () => nav('/apps?org=overheidsontwikkeld') },
              { label: 'Commercieel', icon: 'pi pi-briefcase', command: () => nav('/apps?org=commercieel') }
            ]
          }
        ],
        [
          {
            label: 'Fase',
            items: [
              { label: 'Pilot', icon: 'pi pi-cog', command: () => nav('/apps?stage=pilots') },
              { label: 'Test', icon: 'pi pi-check', command: () => nav('/apps?stage=test') },
              { label: 'Productie', icon: 'pi pi-rocket', command: () => nav('/apps?stage=productie') }
            ]
          }
        ]
      ]
    },
    // Directe links (geen dropdown)
    { label: 'Review', icon: 'pi pi-inbox', command: () => nav('/review/queue'), role: 'reviewer' },
    { label: 'Community', icon: 'pi pi-users', command: () => nav('/community') },
    { label: 'Toelatingseisen', icon: 'pi pi-check-circle', command: () => nav('/criteria') }
  ]
  // Only show 'Review' if user has role 'reviewer' or 'admin'
  return items.filter(item => {
    if (item.label === 'Review') {
      return userStore.roles.includes('reviewer') || userStore.roles.includes('administrator')
    }
    return true
  })
})


/** Vlakke links lijst voor mobiel */
const flatLinks = computed(() => {
  const links = [
    { label:'App-store', to:'/apps', icon:'pi pi-th-large' },
    { label:'Mijn audits', to:'/review/queue', icon:'pi pi-inbox', role: 'reviewer' },
    { label:'Dashboard', to:'/dashboard', icon:'pi pi-chart-line' },
    { label:'Community', to:'/community', icon:'pi pi-users' },
    { label:'Toelatingseisen', to:'/criteria', icon:'pi pi-check-circle' }
  ]
  // Only show 'Mijn audits' if user has role 'reviewer' or 'admin'
  return links.filter(link => {
    if (link.label === 'Mijn audits') {
      return userStore.roles.includes('reviewer') || userStore.roles.includes('administrator')
    }
    return true
  })
})

/** Actieve route check (voor mobiel lijstje) */
const isActive = (to: string) => route.path.startsWith(to)

/** Zoeken */
const doSearch = async () => {
  emit('search', q.value)
  catalog.q = q.value
  await router.push({ path: '/apps', query: q.value ? { q: q.value } : undefined })
  mobileOpen.value = false
}

/** Helpers */
const toggleMobile = () => (mobileOpen.value = !mobileOpen.value)
const goHome = () => router.push('/')
const nav = (to: string) => router.push(to)
</script>

<style scoped>
.container { margin: 0 auto; }

/* Desktop vs mobiel */
.desktop-only { display: block; }
.mobile-only { display: none; }
@media (max-width: 991px) {
  .desktop-only { display: none; }
  .mobile-only { display: block; }
}

/* Sticky */
.header-sticky { position: sticky; top: 0; z-index: 1000; }

/* Mobiele links */
.mobile-link {
  display: flex; align-items: center;
  padding: .75rem .5rem; border-radius: .375rem;
  text-decoration: none; color: var(--text-color);
}
.mobile-link.is-active {
  background: var(--surface-100);
  font-weight: 600;
}
.text-lg { font-size: 1.125rem; }
</style>
