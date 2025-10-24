<!-- /pages/apps/index.vue -->
<template>
  <div class="p-4">
    <!-- Search -->
    <div class="flex gap-3 align-items-center">
      <InputText
        v-model="q"
        placeholder="Zoek applicatie…"
        class="w-full"
        @keyup.enter="doSearch"
        aria-label="Zoek applicatie"
      />
      <Button icon="pi pi-search" @click="doSearch" />
    </div>

    <!-- Tabs -->
    <TabMenu class="mt-3" :model="tabItems" />

    <!-- Filters -->
    <div class="grid mt-3">
      <div class="col-12 md:col-4">
        <Dropdown
          v-model="cat"
          :options="catOptions"
          optionLabel="label" optionValue="value"
          class="w-full"
          placeholder="Categorie"
        />
      </div>
      <div class="col-6 md:col-4">
        <Dropdown
          v-model="stage"
          :options="stageOptions"
          optionLabel="label" optionValue="value"
          class="w-full"
          placeholder="Stage"
        />
      </div>
      <div class="col-6 md:col-4">
        <Dropdown
          v-model="org"
          :options="orgOptions"
          optionLabel="label" optionValue="value"
          class="w-full"
          placeholder="Type"
        />
      </div>
      <div class="col-12 flex gap-2 justify-content-end">
        <Button label="Wissen" icon="pi pi-times" text @click="resetFilters" />
      </div>
    </div>

    <!-- Result grid -->
    <div class="grid mb-4">
      <div v-for="a in visibleApps" :key="a.slug" class="col-12 md:col-6 lg:col-4">
        <AppCard :app="a" :to="`/apps/${a.slug}`" />
      </div>

      <div v-if="!loading && !visibleApps.length" class="col-12">
        <Message severity="info" :closable="false" v-if="q && q !== ''">
          Geen resultaten voor “{{ q }}” in tab “{{ tabLabel }}”.
        </Message>
        <Message severity="info" :closable="false" v-if="!q || q === ''">
          Geen resultaten.
        </Message>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useAsyncData, useRoute, useRouter } from '#app'
import { useAppCatalog } from '~/stores/appCatalog'
import type { AppModel } from '~/types/types'

/* route sync */
const route = useRoute()
const router = useRouter()

/* state */
const catalog = useAppCatalog()
const q = ref<string>((route.query.q as string) || '')
const tab = ref<'alle'|'geinstalleerd'|'testfase'|'overheidsontwikkeld'|'commercieel'>(
  (route.query.tab as any) || 'alle'
)
const cat   = ref<string|null>((route.query.cat as string) || null)
const stage = ref<'productie'|'pilots'|'testfase'|null>((route.query.stage as any) || null)
const org   = ref<'overheidsontwikkeld'|'commercieel'|null>((route.query.org as any) || null)

/* hydrate list */
const { pending } = await useAsyncData('apps:list', () => catalog.fetchAll())
const loading = computed(() => pending.value || catalog.loading)

/* tabs */
const tabItems = [
  { label: 'Alle apps',            command: () => (tab.value = 'alle') },
  { label: 'Geïnstalleerd',        command: () => (tab.value = 'geinstalleerd') },
  { label: 'Testfase',             command: () => (tab.value = 'testfase') },
  { label: 'Overheidsontwikkeld',  command: () => (tab.value = 'overheidsontwikkeld') },
  { label: 'Commercieel',          command: () => (tab.value = 'commercieel') }
]
const tabLabel = computed(() => {
  const map: Record<typeof tab.value,string> = {
    alle:'Alle apps', geinstalleerd:'Geïnstalleerd', testfase:'Testfase',
    overheidsontwikkeld:'Overheidsontwikkeld', commercieel:'Commercieel'
  }
  return map[tab.value]
})

/* filter options (categorieën uit data) */
const catOptions = computed(() => {
  const set = new Set<string>()
  for (const a of catalog.apps) (a.categories || []).forEach(c => set.add(c))
  return [{ label: 'Alle categorieën', value: null as any }].concat(
    Array.from(set).sort().map(c => ({ label: c, value: c }))
  )
})

const stageOptions = [
  { label:'Alle stages', value: null },
  { label:'Productie',   value: 'productie' },
  { label:'Pilots',      value: 'pilots' },
  { label:'Testfase',    value: 'testfase' }
]
const orgOptions = [
  { label:'Alle types',           value: null },
  { label:'Overheidsontwikkeld',  value: 'overheidsontwikkeld' },
  { label:'Commercieel',          value: 'commercieel' }
]

/* tab filter mapping */
const byTab: Record<typeof tab.value, (a: AppModel) => boolean> = {
  alle: () => true,
  geinstalleerd: (a) => a.stage === 'productie',
  testfase: (a) => a.stage === 'pilots' || a.stage === 'testfase',
  overheidsontwikkeld: (a) => a.orgType === 'overheidsontwikkeld',
  commercieel: (a) => a.orgType === 'commercieel'
}

/* results */
const visibleApps = computed(() => {
  return catalog.apps.filter(a => {
    // tab
    if (!byTab[tab.value](a)) return false
    // zoekterm
    if (q.value && q.value !== '') {
      const ql = q.value.toLowerCase()
      const inName = a.name.toLowerCase().includes(ql)
      const inDesc = (a.summary || '').toLowerCase().includes(ql)
      if (!inName && !inDesc) return false
    }
    // categorie
    if (cat.value && cat.value !== null) {
      if (!a.categories || !a.categories.includes(cat.value)) return false
    }
    // stage
    if (stage.value && stage.value !== null) {
      if (a.stage !== stage.value) return false
    }
    // orgType
    if (org.value && org.value !== null) {
      if (a.orgType !== org.value) return false
    }
    return true
  })
})

/* actions */
function doSearch() { /* client-side filter; serversearch kan later */ }
function resetFilters() {
  q.value = ''
  tab.value = 'alle'
  cat.value = null
  stage.value = null
  org.value = null
}

/* query sync (→ URL) */
watch([q, tab, cat, stage, org], () => {
  const query: Record<string, any> = {
    q: q.value || undefined,
    tab: tab.value !== 'alle' ? tab.value : undefined,
    cat: cat.value || undefined,
    stage: stage.value || undefined,
    org: org.value || undefined
  }
  router.replace({ query })
})

/* query sync (← URL) – voor back/forward & deep links */
watch(() => route.query, (qr) => {
  q.value     = (qr.q as string) || ''
  tab.value   = (qr.tab as any) || 'alle'
  cat.value   = (qr.cat as any) || null
  stage.value = (qr.stage as any) || null
  org.value   = (qr.org as any) || null
})
</script>

<style scoped>
/* geen extra styles nodig; PrimeFlex regelt grid */
</style>
