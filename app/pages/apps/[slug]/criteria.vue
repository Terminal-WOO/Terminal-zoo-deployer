<!-- /pages/apps/[slug]/criteria.vue -->
<template>
  <div class="container p-4">
    <div class="flex align-items-center gap-3" v-if="app">
        <img :src="app?.media?.logoUrl || '/logo.svg'" alt="" width="56" height="56" class="border-circle" />
      <div>
        <h1 class="m-0 text-900">Beoordeling — {{ app?.name || slug }}</h1>
        <small class="text-700">Voortgang: {{ progress }}%</small>
      </div>
      <span class="ml-auto">
        <NuxtLink :to="`/apps/${slug}`"><Button text icon="pi pi-arrow-left" label="Terug naar app" /></NuxtLink>
      </span>
    </div>

    <div class="grid mt-3">
      <div class="col-12 md:col-8">
        <Steps :model="steps" :activeIndex="activeStep" />
      </div>
      <div class="col-12 md:col-4">
        <ProgressBar :value="progress" />
      </div>
    </div>

    <div class="grid mt-3">
      <div v-for="c in stepCriteria" :key="c.k" class="col-12 md:col-6">
        <Card>
          <template #title>{{ c.label }}</template>
          <template #subtitle>{{ c.domain }}</template>
          <template #content>
            <p class="m-0 text-700">{{ c.desc }}</p>

            <div class="mt-3 flex align-items-center gap-3">
              <Dropdown
                :modelValue="rec.checks[c.k]"
                :options="statusOpts"
                optionLabel="label" optionValue="value"
                class="w-12rem"
                @update:modelValue="v => setCheck(c.k, v)"
              />
              <Tag :value="evidenceLabel(c.k)" :severity="evidenceSeverity(c.k)" />
            </div>

            <div class="field mt-3">
              <label class="block mb-1">Bewijs (URL)</label>
              <InputText
                :modelValue="rec.evidence[c.k].url"
                @update:modelValue="v => setEvidence(c.k, 'url', String(v))"
                placeholder="https://…"
              />
            </div>
            <div class="field">
              <label class="block mb-1">Opmerking</label>
              <Textarea
                :modelValue="rec.evidence[c.k].note"
                @update:modelValue="v => setEvidence(c.k, 'note', String(v))"
                rows="3" autoResize
              />
            </div>
          </template>
        </Card>
      </div>
    </div>

    <div class="flex justify-content-between mt-2">
      <Button label="Vorige" icon="pi pi-arrow-left" :disabled="activeStep===0" @click="prevStep" />
      <div class="flex gap-2">
        <Button label="Reset" icon="pi pi-refresh" outlined @click="resetAll" />
        <Button label="Volgende" iconPos="right" icon="pi pi-arrow-right" :disabled="activeStep===steps.length-1" @click="nextStep" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })

import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAppCatalog } from '~/stores/appCatalog'
import { useCriteriaStore, type CriteriaKey, type Status } from '~/stores/criteria'

const route = useRoute()
const slug = route.params.slug as string

const catalog = useAppCatalog()
const crit = useCriteriaStore()

onMounted(async () => {
  await catalog.fetchAll().catch(()=>{})
  crit.ensure(slug) // maak app-record aan
})

const app = computed(() => catalog.bySlug(slug))
const rec = computed(() => crit.record(slug))

const steps = [{ label:'Juridisch' },{ label:'Privacy' },{ label:'Techniek' },{ label:'Ethiek' }]
const stepMap = { 0:'Juridisch', 1:'Privacy', 2:'Techniek', 3:'Ethiek' } as const
const activeStep = computed({ get:()=>crit.activeStep, set:v=>crit.activeStep=v })
const stepCriteria = computed(() => crit.all.filter(c => c.domain === stepMap[activeStep.value]))

const statusOpts = [
  { label:'Voldaan', value:'ok' },
  { label:'Ontbreekt', value:'missing' },
  { label:'N.v.t.', value:'na' }
]

function setCheck(k: CriteriaKey, v: Status) { crit.setCheck(slug, k, v) }
function setEvidence(k: CriteriaKey, field: 'url'|'note', value: string) { crit.setEvidence(slug, k, field, value) }

const evidenceLabel = (k: CriteriaKey) => crit.evidenceLabel(slug, k)
const evidenceSeverity = (k: CriteriaKey) => crit.evidenceSeverity(slug, k)
const progress = computed(() => crit.progressPct(slug))

const nextStep = () => crit.nextStep()
const prevStep = () => crit.prevStep()
const resetAll = () => crit.reset(slug)
</script>

<style scoped>
.container { max-width: 1100px; margin: 0 auto; }
.field { margin-bottom: 1rem; }
</style>
