<!-- /pages/apps/[slug].vue -->
<template>
  <div>
    <!-- Breadcrumbs -->
    <div class="surface-0 border-bottom-1 surface-border">
      <div class="container p-3">
        <Breadcrumb :model="bcItems">
          <template #item="{ item }">
            <NuxtLink :to="item.to" v-if="item.to" class="no-underline">{{ item.label }}</NuxtLink>
            <span v-else>{{ item.label }}</span>
          </template>
        </Breadcrumb>
      </div>
    </div>

    <!-- HERO -->
    <section class="surface-0">
      <div class="container p-4">
        <div class="grid gap-4 align-items-center">
          <div class="row-12 md:col-8">
            <div class="flex align-items-center gap-3">
              <img src="~/assets/images/model.png" alt="" width="56" height="56" class="border-circle" />
              <div>
                <h1 class="m-0 text-900">{{ app.name }}</h1>
                <p class="m-0 text-700">{{ app.subtitle }}</p>
              </div>
            </div>

            <!-- in de HERO, onder de tags -->
            <div class="mt-2">
              <Tag v-if="latestAudit" :value="auditStatusLabel(latestAudit.status)"
                :severity="auditStatusSeverity(latestAudit.status)" />
              <small v-if="latestAudit" class="ml-2 text-600">Laatste audit: {{ prettyDate(latestAudit.date) }}</small>
            </div>

          </div>

          <!-- CTA / Pricing -->
          <div class="row-12 md:row-4 w-full">
            <Card>
              <template #title>Acties</template>
              <template #content>
                <div class="flex align-items-center justify-content-between">
                  <div class="text-700">Gebruikstype</div>
                  <b class="text-900">{{ priceType(app.pricing) }}</b>
                </div>
                <div class="flex align-items-center justify-content-between mt-2">
                  <div class="text-700">Indicatie</div>
                  <b class="text-900">{{ priceLabel(app.pricing) }}</b>
                </div>

                <Divider />

                <!-- ORG: deploy naar eigen omgeving -->
                <div class="flex flex-column gap-2 mb-3" v-if="userStore.hasRole('developer') || userStore.hasRole('administrator')">
                  <NuxtLink :to="deployLink" class="no-underline">
                    <Button label="Deploy in jouw organisatie" icon="pi pi-cloud-upload" />
                  </NuxtLink>
                  <small class="text-600">Uitrollen in je eigen tenant of on-prem, volgens jouw policies.</small>
                </div>
              </template>

              <template #footer>
                <small class="text-600">Laatste update: {{ prettyDate(app.updatedAt) }}</small>
                <div class="mt-1">
                  <Tag v-if="app.lifecycle" :value="lifecycleLabel(app.lifecycle)" />
                </div>
              </template>
            </Card>
          </div>

        </div>
      </div>
    </section>

    <!-- BODY -->
    <section class="surface-0">
      <div class="container p-4 pt-0">
        <div class="grid">
          <!-- LEFT column -->
          <div class="col-12 lg:col-8">
            <!-- Screenshots / media -->
            <Galleria v-if="screens.length" :value="screens" :numVisible="4" :showThumbnails="true"
              :showIndicators="false" containerStyle="max-width: 100%;" class="mb-4">
              <template #item="slotProps">
                <img :src="slotProps.item" alt="Screenshot" style="width: 100%; height: 420px; object-fit: cover;" />
              </template>
              <template #thumbnail="slotProps">
                <img :src="slotProps.item" alt="Screenshot" style="width: 100px; height: 70px; object-fit: cover;" />
              </template>
            </Galleria>

            <!-- Description -->
            <Panel header="Beschrijving" class="mb-4">
              <div class="text-700" v-html="renderedDescription"></div>
            </Panel>

            <!-- Reviews -->
            <Panel header="Audits & naleving" class="mb-4">
              <div class="grid gap-4">
                <!-- Samenvatting links -->
                <div class="row-12 md:row-5 w-full">
                  <Card>
                    <template #title>Samenvatting naleving</template>
                    <template #content>
                      <div class="flex align-items-center justify-content-between mb-2">
                        <span class="text-700">Compliance-score</span>
                        <b>{{ complianceScore ?? '—' }}</b>
                      </div>
                      <ProgressBar :value="complianceScore || 0" :showValue="false" style="height: 10px" />
                      <div class="mt-3">
                        <div class="text-700">Laatste audit</div>
                        <div class="flex align-items-center gap-2 mt-1">
                          <Tag v-if="latestAudit" :value="auditStatusLabel(latestAudit.status)"
                            :severity="auditStatusSeverity(latestAudit.status)" />
                          <small v-if="latestAudit" class="text-600">{{ prettyDate(latestAudit.date) }}</small>
                          <small v-else class="text-600">Nog geen audits.</small>
                        </div>
                        <div class="mt-2 flex gap-2 flex-wrap">
                          <Tag v-for="d in (latestAudit?.domains || [])" :key="d" :value="domainLabel(d)" />
                        </div>
                      </div>
                    </template>
                  </Card>
                </div>

                <!-- Overzicht audits rechts -->
                <div class="row-12 md:row-7">
                  <div class="flex align-items-center justify-content-between mb-2">
                    <div class="flex gap-2">
                      <Dropdown v-model="filterSeverity" :options="severityOptions" placeholder="Filter ernst"
                        class="w-12rem" />
                      <Dropdown v-model="filterDomain" :options="domainOptions" placeholder="Domein" class="w-12rem" />
                    </div>
                    <div class="flex align-items-center gap-2">
                      <Checkbox v-model="onlyIssues" :binary="true" inputId="onlyIssues" />
                      <label for="onlyIssues">Alleen issues</label>
                    </div>
                  </div>

                  <div class="flex flex-column gap-3">
                    <Card v-for="a in filteredAudits" :key="a.id">
                      <template #title>
                        <div class="flex align-items-center gap-2">
                          <Tag :value="auditTypeLabel(a.type)" />
                          <span class="font-medium">{{ a.title }}</span>
                        </div>
                      </template>
                      <template #subtitle>
                        <div class="flex align-items-center gap-2">
                          <Tag :value="auditStatusLabel(a.status)" :severity="auditStatusSeverity(a.status)" />
                          <small class="text-600">{{ prettyDate(a.date) }}</small>
                          <small class="text-600">• Auditor: {{ a.auditor.name }}<span v-if="a.auditor.org"> ({{
                            a.auditor.org }})</span></small>
                          <small class="text-600">• Score: <b>{{ a.score }}</b></small>
                        </div>
                      </template>
                      <template #content>
                        <div class="flex gap-2 flex-wrap mb-2">
                          <Tag v-for="d in a.domains" :key="d" :value="domainLabel(d)" />
                        </div>

                        <DataTable :value="filteredFindings(a)" size="small" responsiveLayout="scroll">
                          <Column field="criterionKey" header="Criterion" style="width: 10rem" />
                          <Column field="title" header="Bevinding" />
                          <Column header="Ernst" style="width: 8rem">
                            <template #body="{ data }">
                              <Tag :value="severityLabel(data.severity)" :severity="severityTag(data.severity)" />
                            </template>
                          </Column>
                          <Column header="Status" style="width: 8rem">
                            <template #body="{ data }">
                              <Tag :value="findingStatusLabel(data.status)"
                                :severity="findingStatusSeverity(data.status)" />
                            </template>
                          </Column>
                          <Column header="Bewijs" style="width: 8rem">
                            <template #body="{ data }">
                              <a v-if="data.evidenceUrl" :href="data.evidenceUrl" target="_blank" rel="noopener"
                                class="no-underline">
                                <Button icon="pi pi-paperclip" text rounded />
                              </a>
                              <span v-else class="text-600">—</span>
                            </template>
                          </Column>
                        </DataTable>

                        <div v-if="a.notes" class="mt-2 text-700">
                          <i class="pi pi-info-circle mr-2" aria-hidden="true"></i>{{ a.notes }}
                        </div>
                      </template>
                    </Card>

                    <Message v-if="!filteredAudits.length" severity="info" :closable="false">
                      Geen audits gevonden met de gekozen filters.
                    </Message>
                  </div>
                </div>
              </div>
            </Panel>

            <!-- Changelog / Versies -->
            <Panel header="Versies & wijzigingen">
              <DataTable :value="app.versions || []" :rows="5" :paginator="true">
                <Column field="version" header="Versie" style="width: 10rem" />
                <Column field="date" header="Datum" style="width: 12rem">
                  <template #body="{ data }"><small>{{ prettyDate(data.date) }}</small></template>
                </Column>
                <Column header="Wijzigingen">
                  <template #body="{ data }"><span class="text-700">{{ data.notes }}</span></template>
                </Column>
              </DataTable>
            </Panel>
          </div>

          <!-- RIGHT column -->
          <Panel header="Toelatingseisen" class="mb-4">
            <ul class="list-none p-0 m-0">
              <li v-for="it in criteriaItems" :key="it.k"
                class="flex align-items-center justify-content-between py-2 border-bottom-1 surface-border">
                <span class="text-700">{{ it.label }}</span>
                <i :class="it.met ? 'pi pi-check-circle text-green-500' : 'pi pi-times-circle text-red-400'"
                  aria-hidden="true"></i>
              </li>
            </ul>
            <Divider />
            <div class="text-700 text-sm">Bewijsstukken: <NuxtLink to="/criteria">bekijk beleid</NuxtLink>
            </div>
          </Panel>
        </div>
      </div>
    </section>

    <!-- Related apps -->
    <section class="surface-0">
      <div class="container p-4 pt-0">
        <h2 class="m-0 text-900 mb-3">Gerelateerde apps</h2>
        <div class="grid">
          <div v-for="r in related" :key="r.slug" class="col-12 md:col-6 lg:col-3">
            <Card>
              <template #title>
                <div class="flex align-items-center gap-2">
                  <img :src="r.logo || '/logo.svg'" alt="" width="28" height="28" class="border-circle" />
                  <span>{{ r.name }}</span>
                </div>
              </template>
              <template #content>
                <p class="m-0 text-700">{{ r.summary }}</p>
              </template>
              <template #footer>
                <NuxtLink :to="`/apps/${r.slug}`"><Button label="Bekijk" icon="pi pi-eye" /></NuxtLink>
              </template>
            </Card>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useSeoMeta } from '#imports'
import { useAppCatalog } from '~/stores/appCatalog'
import { useUserStore } from '~/stores/user'
const userStore = useUserStore()

type Pricing =
  | { type: 'free' }
  | { type: 'one_time', price: number, currency: string }
  | { type: 'subscription', price: number, currency: string, interval: 'month' | 'year' }
  | { type: 'usage', unit: string }

const route = useRoute()
const slug = route.params.slug as string
const catalog = useAppCatalog()                       // ⬅️

const { data: appData } = await useAsyncData(`app:${slug}`, () => catalog.fetchOne(slug))
const app = computed<any>(() => appData.value)

useSeoMeta({
  title: () => app.value ? `${app.value.name} — IPO appstore` : 'IPO appstore',
  description: () => app.value?.subtitle || 'ML-app of bouwsteen uit de catalogus.'
})

const bcItems = computed(() => ([
  { label: 'Home', to: '/' },
  { label: 'App-store', to: '/apps' },
  { label: app.value?.name || '—' }
]))

function lifecycleLabel(s?: string) {
  if (!s) return '—'
  const map: Record<string, string> = {
    submitted: 'Ingediend',
    under_review: 'In review',
    changes_requested: 'Wijzigingen gevraagd',
    rejected: 'Afgewezen',
    published: 'Gepubliceerd'
  }
  return map[s] || s
}

const screens = computed<string[]>(() => app.value?.media?.screenshots || [])
const deployLink = computed(() => ({ path: '/deploy/new', query: { slug } }))

const requiredCriteria = [
  { k: 'open_source', label: 'Open-source' },
  { k: 'modular', label: 'Modulair' },
  { k: 'dpia', label: 'DPIA' },
  { k: 'license_clarity', label: 'Licentiestructuur' },
  { k: 'a11y', label: 'Toegankelijkheid' },
  { k: 'ai_risk_cat', label: 'AI-risicocategorie' },
  { k: 'functional_desc', label: 'Functionele beschrijving' },
  { k: 'self_hostable', label: 'Lokaal te draaien' },
  { k: 'tech_explainer', label: 'Technische uitleg' },
]

const extraLabelKeys = [
  { k: 'external_support', label: 'Externe ondersteuning (SLA)' },
  { k: 'human_rights_assessment', label: 'FRIA / IAMA' },
  { k: 'gov_built', label: 'Ontwikkeld door de overheid' },
  { k: 'user_guide', label: 'Handleiding' },
  { k: 'open_inference_api', label: 'Open Inference API' },
]

const criteriaItems = computed(() =>
  requiredCriteria.map(it => ({ ...it, met: !!app.value?.criteria?.[it.k]?.met }))
)
const extraActiveLabels = computed(() =>
  extraLabelKeys.filter(it => !!app.value?.labels?.[it.k]?.has)
)
const hasExtraLabels = computed(() => extraActiveLabels.value.length > 0)

function priceType(p: Pricing | undefined) {
  if (!p) return 'Onbekend'
  if (p.type === 'free') return 'Gratis'
  if (p.type === 'one_time') return 'Eenmalig'
  if (p.type === 'subscription') return 'Abonnement'
  if (p.type === 'usage') return 'Gebruik'
}
function priceLabel(p: Pricing | undefined) {
  if (!p) return '—'
  if (p.type === 'free') return '€ 0'
  if (p.type === 'one_time') return `${p.currency} ${p.price.toFixed(2)}`
  if (p.type === 'subscription') return `${p.currency} ${p.price}/${p.interval === 'month' ? 'mnd' : 'jr'}`
  if (p.type === 'usage') return `Per ${p.unit}`
}
function initials(name: string) { return (name || 'U').split(' ').map(x => x[0]).slice(0, 2).join('').toUpperCase() }
function prettyDate(d?: string | Date) { return d ? new Date(d).toLocaleDateString() : '—' }
function stageLabel(s: string) { return s === 'productie' ? 'Productie' : s === 'pilots' ? 'Pilots' : s === 'testfase' ? 'Testfase' : s }
function orgLabel(t: string) { return t === 'overheidsontwikkeld' ? 'Overheidsontwikkeld' : 'Commercieel' }

const related = computed(() => (app.value?.related || []).slice(0, 4))
const publisherInitials = computed(() => initials(app.value?.org?.name || 'ORG'))
const renderedDescription = computed(() => {
  const md = app.value?.descriptionMd || ''
  return md.replace(/\n{2,}/g, '</p><p>').replace(/\n/g, '<br/>').replace(/^/, '<p>').concat('</p>')
})

import { useAuditStore } from '~/stores/audits'

const audits = useAuditStore()
onMounted(() => audits.initDemo())

const appAudits = computed(() => audits.forApp(slug))
const latestAudit = computed(() => audits.latestForApp(slug))
const complianceScore = computed<number | null>(() => audits.scoreForApp(slug))

/* Filters */
const severityOptions = [
  { label: 'Alle ernstniveaus', value: null },
  { label: 'Info', value: 'info' },
  { label: 'Low', value: 'low' },
  { label: 'Medium', value: 'medium' },
  { label: 'High', value: 'high' },
  { label: 'Critical', value: 'critical' }
]
const domainOptions = [
  { label: 'Alle domeinen', value: null },
  { label: 'Privacy', value: 'privacy' },
  { label: 'Security', value: 'security' },
  { label: 'Ethiek', value: 'ethics' },
  { label: 'Toegankelijkheid', value: 'accessibility' },
  { label: 'Licenties', value: 'licensing' },
  { label: 'AI-risico', value: 'ai_risk' }
]
const filterSeverity = ref<string | null>(null)
const filterDomain = ref<string | null>(null)
const onlyIssues = ref(false)

const filteredAudits = computed(() => {
  return appAudits.value.filter(a => {
    // audit-level filters (alleen domein check; severity filter is op finding-niveau)
    return !filterDomain.value || a.domains.includes(filterDomain.value as any)
  })
})

function filteredFindings(a: any) {
  return (a.findings || []).filter((f: any) => {
    const okDomain = !filterDomain.value || a.domains.includes(filterDomain.value as any)
    const okSev = !filterSeverity.value || f.severity === filterSeverity.value
    const okIssue = !onlyIssues.value || f.status !== 'ok'
    return okDomain && okSev && okIssue
  })
}

/* Labels & visuals */
function auditTypeLabel(t: string) {
  return t === 'initial' ? 'Initiële audit'
    : t === 'periodic' ? 'Periodiek'
      : t === 'incident' ? 'Incident'
        : t === 'retest' ? 'Hertest'
          : t
}
function auditStatusLabel(s: string) {
  return s === 'passed' ? 'Goedgekeurd'
    : s === 'passed_with_notes' ? 'Goedgekeurd (met notities)'
      : s === 'attention' ? 'Aandachtspunt'
        : s === 'failed' ? 'Afgekeurd'
          : s === 'pending' ? 'In behandeling'
            : s
}
function auditStatusSeverity(s: string) {
  return s === 'passed' ? 'success'
    : s === 'passed_with_notes' ? 'info'
      : s === 'attention' ? 'warning'
        : s === 'failed' ? 'danger'
          : 'secondary'
}
function domainLabel(d: string) {
  return d === 'privacy' ? 'Privacy'
    : d === 'security' ? 'Security'
      : d === 'ethics' ? 'Ethiek'
        : d === 'accessibility' ? 'Toegankelijkheid'
          : d === 'licensing' ? 'Licenties'
            : d === 'ai_risk' ? 'AI-risico'
              : d
}
function severityLabel(s: string) {
  return s === 'info' ? 'Info'
    : s === 'low' ? 'Low'
      : s === 'medium' ? 'Medium'
        : s === 'high' ? 'High'
          : s === 'critical' ? 'Critical'
            : s
}
function severityTag(s: string) {
  return s === 'critical' ? 'danger'
    : s === 'high' ? 'warning'
      : s === 'medium' ? 'info'
        : s === 'low' ? 'secondary'
          : 'secondary'
}
function findingStatusLabel(s: string) {
  return s === 'ok' ? 'OK'
    : s === 'warning' ? 'Let op'
      : s === 'fail' ? 'Niet ok'
        : s
}
function findingStatusSeverity(s: string) {
  return s === 'ok' ? 'success'
    : s === 'warning' ? 'warning'
      : s === 'fail' ? 'danger'
        : 'secondary'
}

</script>

<style scoped>
.container {
  max-width: 1100px;
  margin: 0 auto;
}

h1 {
  font-weight: 800;
}

/* optioneel */
:deep(.p-card) {
  border-left: 3px solid var(--surface-border);
}
</style>
