<!-- /pages/community/index.vue -->
<template>
  <div class="community-page">
    <!-- Hero -->
    <section class="surface-0 border-bottom-1 surface-border">
      <div class="container p-4">
        <h1 class="m-0 text-900">Community & monitoring</h1>
        <p class="mt-2 text-700">
          Deel ervaringen, beoordeel apps op toelatingseisen (DPIA, toegankelijkheid, licenties) en volg
          meldingen over privacy, security, techniek en ethiek.
        </p>

        <!-- Search + filters + CTA -->
        <div class="grid mt-3">
          <div class="col-12 md:col-6">
            <div class="p-inputgroup">
              <InputText v-model="q" placeholder="Zoek discussies, apps of onderwerpen…" @keyup.enter="doSearch" />
              <Button icon="pi pi-search" label="Zoek" @click="doSearch" />
            </div>
          </div>
          <div class="col-6 md:col-3">
            <Dropdown v-model="filterCat" :options="categoryOptions" optionLabel="label" optionValue="value" placeholder="Categorie" class="w-full" />
          </div>
          <div class="col-6 md:col-3 flex md:justify-content-end">
            <Button label="Nieuwe bijdrage" icon="pi pi-plus" @click="dialogNew = true" severity="secondary" />
          </div>
        </div>

        <div class="mt-2 flex gap-2 flex-wrap">
          <Tag v-for="c in quickCats" :key="c" :value="c" @click="toggleQuick(c)" class="cursor-pointer" />
        </div>
      </div>
    </section>

    <!-- Tabs -->
    <section class="surface-0">
      <div class="container p-4 pt-3">
        <TabView>
          <!-- TAB 1: Discussies -->
          <TabPanel header="Discussies">
            <DataTable
              :value="filteredThreads"
              dataKey="id"
              :paginator="true"
              :rows="8"
              responsiveLayout="scroll"
              :rowHover="true"
              class="mt-2"
              :loading="loading"
              currentPageReportTemplate="{first}–{last} van {totalRecords}"
              :rowsPerPageOptions="[8,16,24]"
            >
              <Column field="title" header="Onderwerp" :sortable="true">
                <template #body="{ data }">
                  <div class="flex flex-column">
                    <div class="flex align-items-center gap-2">
                      <NuxtLink :to="`/community/thread/${data.id}`" class="no-underline text-900 hover:text-primary font-medium">{{ data.title }}</NuxtLink>
                      <Tag :value="data.category" :severity="tagSeverity(data.category)" />
                    </div>
                    <small class="text-600">Aangemaakt door {{ data.author }} • {{ fromNow(data.updatedAt) }}</small>
                  </div>
                </template>
              </Column>
              <Column header="Status" :sortable="true" field="status" style="width: 12rem">
                <template #body="{ data }">
                  <Tag :value="data.status" :severity="statusSeverity(data.status)" />
                </template>
              </Column>
              <Column field="replies" header="Reacties" :sortable="true" style="width: 8rem" />
              <Column field="app" header="Gerelateerd aan" style="width: 14rem">
                <template #body="{ data }">
                  <Chip :label="data.app" icon="pi pi-th-large" />
                </template>
              </Column>
              <Column header="Acties" style="width: 10rem">
                <template #body>
                  <div class="flex gap-2">
                    <Button icon="pi pi-eye" text rounded/>
                    <Button icon="pi pi-flag" text rounded/>
                  </div>
                </template>
              </Column>
            </DataTable>
          </TabPanel>

          <!-- TAB 4: Security advisories -->
          <TabPanel header="Security advisories">
            <div class="grid mt-2">
              <div class="col-12 md:col-3">
                <Panel header="Filter">
                  <div class="field">
                    <label>Ernst</label>
                    <Dropdown v-model="severity" :options="severityOptions" optionLabel="label" optionValue="value" class="w-full" placeholder="Alle" />
                  </div>
                  <div class="field">
                    <label>Component</label>
                    <Dropdown v-model="component" :options="componentOptions" class="w-full" placeholder="Alle" />
                  </div>
                  <Button label="Wissen" text icon="pi pi-times" @click="resetAdvisoryFilters" />
                </Panel>
              </div>
              <div class="col-12 md:col-9">
                <DataTable :value="filteredAdvisories" :paginator="true" :rows="6" responsiveLayout="scroll">
                  <Column field="id" header="ID" style="width: 9rem" />
                  <Column field="title" header="Titel" />
                  <Column field="component" header="Component" style="width: 12rem" />
                  <Column field="severity" header="Ernst" style="width: 10rem">
                    <template #body="{ data }">
                      <Tag :value="data.severity" :severity="severityColor(data.severity)" />
                    </template>
                  </Column>
                  <Column field="published" header="Geplaatst" style="width: 10rem">
                    <template #body="{ data }">
                      <small>{{ fromNow(data.published) }}</small>
                    </template>
                  </Column>
                  <Column header="Acties" style="width: 10rem">
                    <template #body>
                      <div class="flex gap-2">
                        <Button icon="pi pi-download" text rounded />
                        <Button icon="pi pi-external-link" text rounded />
                      </div>
                    </template>
                  </Column>
                </DataTable>
              </div>
            </div>
          </TabPanel>
        </TabView>

        <!-- Policy -->
        <Divider class="mt-5" />
        <Panel header="Moderatie- en publicatiebeleid" class="mt-3">
          <ul class="m-0 pl-3">
            <li class="mb-2">Respecteer privacy- en beveiligingsrichtlijnen; deel geen vertrouwelijke data.</li>
            <li class="mb-2">Onderbouw claims met bronnen; geef aan of iets een mening of een feit is.</li>
            <li class="mb-2">Kwetsbaarheden melden via het tabblad “Security advisories” of responsible disclosure.</li>
            <li>Moderators kunnen bijdragen labelen of sluiten bij onjuistheden of overtredingen.</li>
          </ul>
        </Panel>
      </div>
    </section>

    <!-- Dialoog: nieuwe bijdrage -->
    <Dialog v-model:visible="dialogNew" modal header="Nieuwe bijdrage" class="w-11 md:w-6">
      <div class="field">
        <label>Titel</label>
        <InputText v-model="newThread.title" placeholder="Korte, duidelijke titel" />
      </div>
      <div class="field">
        <label>Categorie</label>
        <Dropdown v-model="newThread.category" :options="categoryOptions" optionLabel="label" optionValue="value" class="w-full" placeholder="Kies categorie" />
      </div>
      <div class="field">
        <label>Gerelateerd aan (app)</label>
        <Dropdown v-model="newThread.app" :options="apps" class="w-full" placeholder="Kies app (optioneel)" />
      </div>
      <div class="field">
        <label>Bericht</label>
        <Textarea v-model="newThread.content" rows="6" autoResize placeholder="Omschrijf vraag, bevinding of voorstel…" />
      </div>
      <div class="flex justify-content-end gap-2 mt-2">
        <Button label="Annuleren" text @click="dialogNew = false" />
        <Button label="Plaatsen" icon="pi pi-send" @click="createThread" />
      </div>
      <Message v-if="newMsg" :severity="newOk ? 'success' : 'warn'" class="mt-3" :closable="false">{{ newMsg }}</Message>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useCommunityStore } from '~/stores/community'    // ⬅️

const store = useCommunityStore()
onMounted(() => store.init())     

// loading ref voor je Discussies-tabel (je template gebruikt :loading="loading")
const loading = ref(false)

/* UI filters */
const q = ref('')
const filterCat = ref<string | null>(null)
const quickCats = ['Privacy','Security','Techniek','Ethiek']
const categoryOptions = quickCats.map(c => ({ label:c, value:c }))
const toggleQuick = (c:string) => { filterCat.value = filterCat.value === c ? null : c }

/* Discussies */
const filteredThreads = computed(() =>
  store.threads.filter(t => {
    const matchQ = !q.value || [t.title, t.author, t.app].join(' ').toLowerCase().includes(q.value.toLowerCase())
    const matchC = !filterCat.value || t.category === filterCat.value
    return matchQ && matchC
  })
)


// zoek is reactief via q/filterCat → dus no-op is genoeg
const doSearch = () => {
  // eventueel: active tab naar "Discussies" zetten of analytics
}

/* Advisories filters (optioneel lokaal of via store getters) */
const severityOptions = [{label:'Alle', value:null}, {label:'Low', value:'Low'}, {label:'Medium', value:'Medium'}, {label:'High', value:'High'}, {label:'Critical', value:'Critical'}]
const componentOptions = computed(() => Array.from(new Set(store.advisories.map(a=>a.component))))
const severity = ref<string|null>(null)
const component = ref<string|null>(null)
const filteredAdvisories = computed(() =>
  store.advisories.filter(a => (!severity.value || a.severity===severity.value) && (!component.value || a.component===component.value))
)
const resetAdvisoryFilters = () => { severity.value = null; component.value = null }

/* Helpers */
function tagSeverity(c:string){ if(c==='Security')return 'danger'; if(c==='Privacy')return 'info'; if(c==='Techniek')return 'success'; if(c==='Ethiek')return 'warning'; }
function statusSeverity(s:string){ if(s==='Open')return 'info'; if(s==='Gesloten')return 'secondary'; if(s==='Opgelost')return 'success'; }
function severityColor(s:string){ if(s==='Critical')return 'danger'; if(s==='High')return 'warning'; if(s==='Medium')return 'info'; return 'secondary' }

/* Nieuwe bijdrage dialog (op je bestaande template bindings) */
const dialogNew = ref(false)
const apps = ['GovChat NL','Geo Buddy','Parlement Agent','Libre Infer','Overig']
const newThread = ref({ title:'', category:null as null|string, app:null as null|string, content:'' })
const newMsg = ref<string|null>(null); const newOk = ref(false)
function createThread() {
  if(!newThread.value.title || !newThread.value.category){ newOk.value=false; newMsg.value='Titel en categorie zijn verplicht.'; return }
  store.addThread({ title:newThread.value.title, author:'Jij', category:newThread.value.category as any, app:newThread.value.app || '—' })
  dialogNew.value=false; newOk.value=true; newMsg.value=null; newThread.value={ title:'', category:null, app:null, content:'' }
}

/* tiny util */
function fromNow(iso:string){ try { return new Date(iso).toLocaleString() } catch { return iso } }
</script>

<style scoped>
.container { max-width: 1100px; margin: 0 auto; }
.community-page :deep(.p-datatable .p-datatable-tbody > tr > td) { vertical-align: top; }
.cursor-pointer { cursor: pointer; }

.field {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}
</style>
