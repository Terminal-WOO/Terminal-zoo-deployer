<!-- /pages/index.vue -->
<template>
  <div>
    <!-- HERO -->
    <div class="home-hero-outer">
      <HomeHero class="hero-space-top" />
    </div>

    <!-- TRUST/CRITERIA STRIP -->
    <section class="surface-50 border-y-1 surface-border">
      <div class="container mx-auto px-3 py-4">
        <div class="flex flex-wrap gap-3 align-items-center justify-content-between">
          <div v-for="c in criteria" :key="c.k" class="flex align-items-center gap-2">
            <i :class="c.icon" aria-hidden="true"></i>
            <span class="text-800">{{ c.lbl }}</span>
          </div>
          <NuxtLink to="/criteria" class="no-underline ml-auto">
            <Button label="Bekijk alle toelatingseisen" class="p-button-text" icon="pi pi-list-check" />
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- CATEGORIEËN -->
    <section class="surface-0">
      <div class="container mx-auto px-3 py-6">
        <div class="flex align-items-end justify-content-between">
          <h2 class="m-0 heading-dark">Categorieën</h2>
          <NuxtLink to="/apps" class="no-underline">
            <Button class="p-button-text" icon="pi pi-arrow-right" label="Alles bekijken" />
          </NuxtLink>
        </div>

        <div class="grid mt-3">
          <div v-for="cat in categories" :key="cat.key" class="col-12 md:col-6 lg:col-4 card-eq">
            <Card class="card-fill hover:surface-100 cursor-pointer" @click="goToCategory(cat.key)" role="link">
              <template #title>
                <div class="flex align-items-center gap-2">
                  <i :class="cat.icon" aria-hidden="true"></i>
                  <span class="heading-dark-sm">{{ cat.label }}</span>
                </div>
              </template>
              <template #content>
                <p class="m-0 text-700 line-clamp-2">
                  Zoek snel binnen {{ cat.label.toLowerCase() }} en vergelijk op labels.
                </p>
              </template>
            </Card>
          </div>
        </div>
      </div>
    </section>

    <!-- UITGELICHT -->
    <section class="surface-0">
      <div class="container mx-auto px-3 pb-1">
        <h2 class="heading-dark">Uitgelicht</h2>
      </div>

      <div class="container mx-auto px-3 pb-6">
        <div class="grid">
          <div v-for="a in featured" :key="a.slug" class="col-12 md:col-6 lg:col-3 card-eq">
            <AppCard :app="a" :to="`/apps/${a.slug}`" />
          </div>
        </div>
      </div>
    </section>

    <!-- ZO WERKT HET -->
    <section class="surface-0">
      <div class="container mx-auto px-3 pb-6">
        <h2 class="heading-dark">Zo werkt het</h2>
        <div class="grid mt-2">
          <div class="col-12 md:col-3 card-eq">
            <Card class="card-fill">
              <template #title><span class="heading-dark-sm">1. Selecteer</span></template>
              <template #content>
                <p class="m-0 text-700">Kies een app of combineer bouwstenen.</p>
                <div class="mt-2"><Tag value="LLM" /><Tag value="Vectordb" class="ml-2" /><Tag value="Agents" class="ml-2" /></div>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3 card-eq">
            <Card class="card-fill">
              <template #title><span class="heading-dark-sm">2. Configureer</span></template>
              <template #content>
                <p class="m-0 text-700">Stel LLM (GPT/Mistral/Gemini), Qdrant/Pinecone/Azure AI Search en extra’s in.</p>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3 card-eq">
            <Card class="card-fill">
              <template #title><span class="heading-dark-sm">3. Deploy</span></template>
              <template #content>
                <p class="m-0 text-700">Rol uit naar jouw cloud of on-prem — cloud-agnostisch.</p>
              </template>
            </Card>
          </div>
          <div class="col-12 md:col-3 card-eq">
            <Card class="card-fill">
              <template #title><span class="heading-dark-sm">4. Monitor</span></template>
              <template #content>
                <p class="m-0 text-700">Beheer labels, audits en community-feedback.</p>
              </template>
            </Card>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="surface-50 border-top-1 surface-border">
      <div class="container mx-auto px-3 py-6">
        <div class="grid align-items-center">
          <div class="col-12 md:col-8">
            <h2 class="m-0 heading-dark">Klaar om te beginnen?</h2>
            <p class="mt-2 text-700">
              Bekijk de catalogus en ontdek betrouwbare AI-toepassingen met transparante toelatingseisen.
            </p>
          </div>
          <div class="col-12 md:col-4 flex md:justify-content-end gap-2">
            <NuxtLink to="/apps"><Button label="Bekijk catalogus" icon="pi pi-arrow-right" /></NuxtLink>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useHead } from '#imports'

definePageMeta({ layout: 'landing' })

useHead({
  title: 'App Store — Selecteer, Configureer, Deploy',
  meta: [
    { name: 'description', content: 'Vind, beoordeel en deploy AI/ML-applicaties met bouwstenen en duidelijke toelatingseisen & labels.' },
    { property: 'og:title', content: 'App Store' },
    { property: 'og:description', content: 'Selecteer, configureer en deploy betrouwbare AI-toepassingen met transparante criteria.' },
    { property: 'og:type', content: 'website' }
  ]
})

const router = useRouter()

/* Categorieën */
const categories = [
    { key: 'Agent',     label: 'Agent',        icon: 'pi pi-user' },
    { key: 'Chatbot',   label: 'Chatbot',      icon: 'pi pi-comments' },
    { key: 'Docs',      label: 'Docs',         icon: 'pi pi-book' },
    { key: 'Nederlands',label: 'Nederlands',   icon: 'pi pi-globe' },
    { key: 'RAG',       label: 'RAG',          icon: 'pi pi-database' }
]

function goToCategory(key: string) {
  // nieuwe query key "category" + backward-compat "cat"
  router.push({ path: '/apps', query: { cat: key } })
}

/* Uitgelicht (dummy) */
type FeaturedApp = {
  slug: string; name: string; summary: string; tags: string[];
  stage: 'productie'|'pilots'|'testfase'; orgType: 'overheidsontwikkeld'|'commercieel'; logo?: string
}
const featured = ref<FeaturedApp[]>([
  { slug: 'govchat-nl',      name: 'GovChat NL',      summary: 'Veilige chat-assistent met RAG op beleid & wetgeving.', tags: ['RAG','Chatbot','Nederlands'], stage: 'productie', orgType: 'overheidsontwikkeld' },
  { slug: 'geo-buddy',       name: 'Geo Buddy',       summary: 'Ruimtelijke AI voor kaarten en plankaarten-analyse.',   tags: ['Geodata','Vision'],          stage: 'pilots',    orgType: 'commercieel' },
  { slug: 'parlement-agent', name: 'Parlement Agent', summary: 'Agent die Kamerstukken samenvat en vergelijkt.',        tags: ['Agents','Parlementaire documenten'], stage: 'testfase', orgType: 'overheidsontwikkeld' },
  { slug: 'wl-toegankelijkheid-database',     name: 'WL Toegankelijkheid Database',     summary: 'Database voor wetgeving en beleid',               tags: ['Open-source','Inference API'], stage: 'productie', orgType: 'commercieel' },
  { slug: 'libre-infer',     name: 'Libre Infer',     summary: 'Open-source inference API met LiteLLM.',               tags: ['Open-source','Inference API'], stage: 'productie', orgType: 'commercieel' }
])

/* Trust strip */
const criteria = [
  { k: 'open_source', lbl: 'Open-source',        icon: 'pi pi-github' },
  { k: 'dpia',        lbl: 'DPIA OK',            icon: 'pi pi-check-circle' },
  { k: 'a11y',        lbl: 'Toegankelijk',       icon: 'pi pi-unlock' },
  { k: 'ai_risk',     lbl: 'AI-risicocat.',      icon: 'pi pi-exclamation-triangle' },
  { k: 'self_host',   lbl: 'Lokaal te draaien',  icon: 'pi pi-home' }
]

/* Helper voor stage tags */
const stageMap: Record<FeaturedApp['stage'], string> = {
  productie: 'Productie', pilots: 'Pilots', testfase: 'Testfase'
}
function stageTag(s: FeaturedApp['stage']) {
  if (s === 'productie') return 'success'
  if (s === 'pilots') return 'warning'
  return 'secondary'
}

/* JSON-LD (ingesloten) */
const ldJson = computed(() => ({
  '@context': 'https://schema.org',
  '@type': 'WebSite',
  name: 'App Store',
  url: 'https://example.org',
  potentialAction: {
    '@type': 'SearchAction',
    target: 'https://example.org/apps?q={search_term_string}',
    'query-input': 'required name=search_term_string'
  }
}))
useHead({
  script: [{ type: 'application/ld+json', children: JSON.stringify(ldJson.value) }]
})
</script>

<style>
/* ——— VISUEEL THEMA ——— */
:root { /* super-donkergroen die bij je preset past */
  --ai-deep-green: #0E2D19;
}
</style>

<style scoped>
.container { max-width: 1100px; }

/* Hero spacing tweaks */
.hero-space-top { margin-top: 1.25rem; padding-top: .5rem; }
.home-hero-outer :deep(.p-button) { margin-right: .5rem; } /* ruimte tussen hero-link & illustratie */
.home-hero-outer :deep(img) { margin-left: .5rem; }        /* beetje marge richting afbeelding rechts */

/* geen gitzwart voor headings */
.heading-dark { color: var(--ai-deep-green); font-weight: 800; }
.heading-dark-sm { color: var(--ai-deep-green); font-weight: 700; }

/* ——— GELIJKE HOOGTE CARDS ——— */
.card-eq { display: flex; }
.card-eq > :deep(.p-card) { width: 100%; }

.card-fill :deep(.p-card) {
  height: 100%;
  display: flex;
  flex-direction: column;
}
.card-fill :deep(.p-card-body) {
  display: flex;
  flex-direction: column;
  flex: 1;
}
.card-fill :deep(.p-card-content) { flex: 1; }

/* compacter titels/subtitles */
.card-fill :deep(.p-card-title),
.card-fill :deep(.p-card-subtitle) { margin-bottom: .5rem; }

/* multi-line clamp voor gelijke hoogte teksten */
.line-clamp-2 { display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.line-clamp-3 { display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }
</style>
