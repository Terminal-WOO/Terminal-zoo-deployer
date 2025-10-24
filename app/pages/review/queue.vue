<template>
  <section class="container p-4">
    <h1 class="m-0 text-900">Mijn audits</h1>
    <p class="text-700 mt-2">Uw lijst aan toegewezen audits van modellen.</p>

    <DataTable :value="pendingApps" class="mt-3" :rows="10" :paginator="true" responsiveLayout="scroll">
      <Column header="Model">
        <template #body="{data}">
          <div class="font-medium">{{ data.name }}</div>
          <small class="text-600">{{ data.summary }}</small>
        </template>
      </Column>
      <Column field="status" header="Status" style="width:12rem">
        <template #body="{data}"><Tag :value="statusLabel(data.approved)" :severity="statusTag(data.approved)" /></template>
      </Column>
      <Column header="Acties" style="width:16rem">
        <template #body="{data}">
          <Button icon="pi pi-eye" label="Open" @click="openApprovalModal(data)" />
        </template>
      </Column>
    </DataTable>
  </section>

  <Dialog v-model:visible="approvalModalOpen" header="Applicatie goedkeuren" :modal="true" :draggable="false" :closable="true" :style="{width:'50rem'}">
    <div class="p-fluid">
      <div class="field">
        <label for="app-name" class="block mb-2">Naam</label>
        <InputText id="app-name" v-model="editingApp.name" class="w-full" disabled />
      </div>
      <div class="field">
        <label for="app-summary" class="block mb-2">Samenvatting</label>
        <InputTextarea id="app-summary" v-model="editingApp.summary" class="w-full" rows="3" disabled />
      </div>
      <div class="field">
        <label for="app-status" class="block mb-2">Huidige status</label>
        <Tag :value="statusLabel(editingApp.approved)" :severity="statusTag(editingApp.approved)" />
      </div>
      <!-- Review the labels of the application and approve them or not put 2 fields next to eachother -->
      <div>
        <label class="block mb-2">Criteria</label>
        <div class="grid">
          <div class="col-12">
            <div class="field flex items-center gap-2">
              <Checkbox v-model="editingApp.criteria.a11y.met" :inputId="'a11y'" binary />
              <label for="a11y" class="m-0">Voldoet aan toegankelijkheidseisen</label>
            </div>
            <div class="field flex items-center gap-2">
              <Checkbox v-model="editingApp.criteria.dpia.met" :inputId="'dpia'" binary />
              <label for="dpia" class="m-0">DPIA uitgevoerd</label>
            </div>
            <div class="field flex items-center gap-2">
              <Checkbox v-model="editingApp.criteria.open_source.met" :inputId="'open_source'" binary />
              <label for="open_source">Is open source</label>
            </div>
            <div class="field flex items-center gap-2">
              <Checkbox v-model="editingApp.criteria.modular.met" :inputId="'modular'" binary />
              <label for="modular">Is modulair</label>
            </div>
            <div class="field flex items-center gap-2">
              <Checkbox v-model="editingApp.criteria.license_clarity.met" :inputId="'license_clarity'" binary />
              <label for="license_clarity">Licentie is duidelijk</label>
            </div>
            <div class="field flex items-center gap-2">
              <Checkbox v-model="editingApp.criteria.ai_risk_cat.met" :inputId="'ai_risk_cat'" binary />
              <label for="ai_risk_cat">AI risicocategorie is bepaald</label>
            </div>
            <div class="field flex items-center gap-2">
              <Checkbox v-model="editingApp.criteria.functional_desc.met" :inputId="'functional_desc'" binary />
              <label for="functional_desc">Functionele beschrijving is aanwezig</label>
            </div>
            <div class="field flex items-center gap-2">
              <Checkbox v-model="editingApp.criteria.self_hostable.met" :inputId="'self_hostable'" binary />
              <label for="self_hostable">Kan zelf gehost worden</label>
            </div>
            <div class="field flex items-center gap-2">
              <Checkbox v-model="editingApp.criteria.tech_explainer.met" :inputId="'tech_explainer'" binary />
              <label for="tech_explainer">Technische uitleg is aanwezig</label>
            </div>
          </div>
        </div>
      </div>
      <div class="field flex gap-2">
        <Button type="button" label="Publiceren" icon="pi pi-check" class="p-button-success" @click="editingApp.approved='published'; updateAppConfig()" />
        <Button type="button" label="Wijzigingen vragen" icon="pi pi-pencil" class="p-button-warning" @click="editingApp.approved='changes_requested'; updateAppConfig()" />
        <Button type="button" label="Afkeuren" icon="pi pi-times" class="p-button-danger" @click="editingApp.approved='rejected'; updateAppConfig()" />
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })
import { computed } from 'vue'
import { useUserStore } from '~/stores/user'
import { useSubmissionsStore } from '~/stores/submissions'
import { useAppCatalog } from '~/stores/appCatalog'
import type { AppModel } from '~/types/types'

const user = useUserStore()
// eenvoudige role-gate (UI): echte bescherming kan ook met route-middleware
if (!user.hasRole('reviewer') && !user.hasRole('administrator')) {
  navigateTo('/'); // of toon Message
}

const catalog = useAppCatalog()

onMounted(async () => { await catalog.fetchUserApps() })

const pendingApps = computed(() => catalog.userApps.filter(a => ['submitted','under_review','changes_requested'].includes(a.approved)))
const statusLabel = (s:string)=>({submitted:'Ingediend',under_review:'In review',changes_requested:'Wijzigingen gevraagd'})[s]||s
const statusTag = (s:string)=> s==='changes_requested'?'warning':'info'

const approvalModalOpen = ref(false)
const editingApp = ref<AppModel>({} as AppModel)

function openApprovalModal(a: AppModel) {
  editingApp.value = a
  approvalModalOpen.value = true
}

async function updateAppConfig() {
  let applications = JSON.parse(JSON.stringify(catalog.userApps)) as AppModel[]
  const index = applications.findIndex(a => a.slug === editingApp.value.slug)
  if (index !== -1) {
    applications[index] = editingApp.value
    await $fetch('/api/apps/userApps', {
        method: 'POST',
        body: applications,
      })
  }

  await catalog.fetchUserApps()
  approvalModalOpen.value = false
  editingApp.value = {} as AppModel
}

const subs = useSubmissionsStore(); subs.initDemo()
</script>

<style scoped>
.container{ max-width:1100px; margin:0 auto; }
</style>
