<!-- /pages/criteria/index.vue -->
<template>
  <div class="container p-4">
    <h1 class="m-0 text-900">Toelatingseisen</h1>
    <p class="mt-2 text-700">
      Dit is de <b>algemene set criteria</b>. Selecteer een app om een <b>beoordeling</b> te starten.
    </p>

    <div class="grid mt-3">
      <div class="col-12 md:col-8">
        <div class="p-inputgroup">
          <InputText v-model="q" placeholder="Zoek in eisen…" />
          <Button icon="pi pi-search" label="Zoek" />
        </div>
      </div>
    </div>

    <Panel header="Toelatingseisen-overzicht" class="mt-3">
      <DataTable :value="rows" :rows="10" :paginator="true" responsiveLayout="scroll">
        <Column field="domain" header="Domein" style="width: 12rem" sortable />
        <Column field="label" header="Eis" sortable />
        <Column header="Beschrijving">
          <template #body="{ data }"><span class="text-700">{{ data.desc }}</span></template>
        </Column>
      </DataTable>
    </Panel>

    <Panel header="Richtlijnen & templates" class="mt-3">
      <div class="grid">
        <div v-for="t in templates" :key="t.k" class="col-12 md:col-6 lg:col-4">
          <Card>
            <template #title>{{ t.title }}</template>
            <template #content><p class="m-0 text-700">{{ t.desc }}</p></template>
            <template #footer><a :href="t.href" target="_blank" rel="noopener"><Button label="Download" icon="pi pi-download" /></a></template>
          </Card>
        </div>
      </div>
    </Panel>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useCriteriaStore } from '~/stores/criteria'
import { useAppCatalog } from '~/stores/appCatalog'

const crit = useCriteriaStore()
const catalog = useAppCatalog()

onMounted(() => { catalog.fetchAll().catch(()=>{}) })

const q = ref('')
const rows = computed(() =>
  crit.all.filter(c =>
    !q.value || [c.label, c.desc, c.domain].join(' ').toLowerCase().includes(q.value.toLowerCase())
  )
)

const templates = [
  { k:'dpia', title:'DPIA template', desc:'Standaard DPIA-sjabloon met voorbeelden.', href:'#' },
  { k:'wcag', title:'WCAG checklijst', desc:'Checklist voor toegankelijkheid (AA).', href:'#' },
  { k:'license', title:'Licentie-matrix', desc:'OSS/closed-matrix en voorbeelden.', href:'#' },
  { k:'ai-risk', title:'AI-risico indeling', desc:'Formulier indeling & maatregelen.', href:'#' }
]
</script>

<style scoped>
.container { max-width: 1100px; margin: 0 auto; }
</style>
