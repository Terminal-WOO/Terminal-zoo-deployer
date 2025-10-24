<!-- /components/HomeHero.vue -->
<template>
    <section class="surface-0 hero">
        <div class="container p-4">
            <div class="grid">
                <!-- LEFT: headline, copy, zoek, features -->
                <div class="col-12 lg:col-7">
                    <!-- bovenste pill in eigen thema -->
                    <div class="flex align-items-center gap-2 mb-3">
                        <Tag value="Voor publieke organisaties" severity="info" rounded />
                        <NuxtLink to="/criteria"
                            class="no-underline text-700 hover:text-900 flex align-items-center gap-2">
                            Transparante toelatingseisen
                            <i class="pi pi-arrow-right" aria-hidden="true"></i>
                        </NuxtLink>
                    </div>

                    <!-- grote titel in jouw toon -->
                    <h1 class="m-0 text-900 hero-title">
                        Deploy betrouwbare AI-applicaties
                    </h1>

                    <!-- subcopy in thema -->
                    <p class="mt-3 text-700 text-lg">
                        Combineer LLM, vectordatabase en agents. Beoordeel op duidelijke criteria
                        zoals <b>DPIA</b>, <b>toegankelijkheid</b> en <b>licenties</b>, en rol uit in jouw eigen
                        omgeving.
                    </p>

                    <!-- zoek -->
                    <form class="mt-3 p-inputgroup hero-search" role="search" aria-label="Zoek in catalogus"
                        @submit.prevent="doSearch">
                        <InputText v-model="q" placeholder="Zoek in catalogus — bijv. 'RAG burgerbrieven'"
                            aria-label="Zoekterm" @keyup.enter="doSearch" />
                        <Button icon="pi pi-search" label="Zoek" type="submit" />
                    </form>

                    <!-- quick links naar categorieën (optioneel, in thema) -->
                    <div class="mt-3 flex flex-wrap gap-2">
                    <Button text size="small" label="Chatbots" @click="goCat('Chatbot')" />
                    <Button text size="small" label="RAG" @click="goCat('RAG')" />
                    <Button text size="small" label="Agents" @click="goCat('Agents')" />
                    <Button text size="small" label="Nederlands" @click="goCat('Nederlands')" />

                    </div>

                    <!-- features 3 kolommen, afgestemd op jouw proposities -->
                    <div class="grid mt-5">
                        <div class="col-12 md:col-4">
                            <div class="flex align-items-start gap-3">
                                <span class="hero-icon"><i class="pi pi-list-check" aria-hidden="true"></i></span>
                                <div>
                                    <div class="font-bold">Toelatingseisen</div>
                                    <small class="text-600">DPIA, toegankelijkheid, AI-risico en licenties helder in
                                        beeld.</small>
                                </div>
                            </div>
                        </div>
                        <div class="col-12 md:col-4">
                            <div class="flex align-items-start gap-3">
                                <span class="hero-icon"><i class="pi pi-cloud" aria-hidden="true"></i></span>
                                <div>
                                    <div class="font-bold">Cloud-agnostisch</div>
                                    <small class="text-600">Azure, Hetzner of eigen (on-prem) omgeving.</small>
                                </div>
                            </div>
                        </div>
                        <div class="col-12 md:col-4">
                            <div class="flex align-items-start gap-3">
                                <span class="hero-icon"><i class="pi pi-chart-line" aria-hidden="true"></i></span>
                                <div>
                                    <div class="font-bold">Community & monitoring</div>
                                    <small class="text-600">Feedback, audits en gebruiks-/kostendashboard.</small>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- RIGHT: één afbeelding -->
                <div class="col-12 lg:col-5 flex align-items-center justify-content-center">
                    <div class="img-box rounded">
                        <img src="@/assets/images/ml-appstore.png" alt="Demo afbeelding 1"
                            loading="lazy" style="width:100%;height:auto;display:block;border-radius:12px;" />
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAppCatalog } from '~/stores/appCatalog'   

const router = useRouter()
const catalog = useAppCatalog()    
const q = ref('')

const doSearch = () => {
    catalog.q = q.value  
    router.push({ path: '/apps', query: q.value ? { q: q.value } : undefined })
}
const goCat = (cat: string) => {
  catalog.selectedCats = [cat]                      
  router.push({ path: '/apps', query: { cat } })
}
</script>

<style scoped>
.container {
    max-width: 1100px;
    margin: 0 auto;
}

/* Hero typografie in lijn met je thema */
.hero-title {
    font-weight: 800;
    line-height: 1.15;
    letter-spacing: -0.01em;
    font-size: clamp(2rem, 3.5vw + 1rem, 3rem);
}

/* Zoekbalk iets royaler, in PrimeVue-stijl */
.hero-search .p-inputtext {
    padding: 1rem;
}

.hero-search .p-button {
    padding: 1rem 1.25rem;
}

/* Feature-iconen in themakleuren */
.hero-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: 10px;
    background: var(--surface-100);
    color: var(--primary-color);
}

/* Kaarten rechts (alleen visuals, kleuren volgen het thema) */
.hero-card :deep(.p-card) {
    height: 100%;
}

.hero-card {
    border-radius: 12px;
    overflow: hidden;
}
</style>
