<!-- /components/AppCard.vue -->
<template>
    <Card class="app-card" :class="{ 'is-clickable': !!to }" role="article" @click="to && nav(to)"
        @keyup.enter="to && nav(to)" tabindex="0">
        <!-- Titel -->
        <template #title>
            <div class="flex align-items-center gap-2">
                <img v-if="logoSrc" src="~/assets/images/model.png" :alt="`${app.name} logo`" width="28" height="28"
                    class="border-circle" loading="lazy" />
                <component :is="to ? NuxtLink : 'span'" :to="to" class="heading-link" aria-label="Open app">
                    <span class="heading-dark-sm">{{ app.name }}</span>
                </component>
            </div>
        </template>

        <!-- Subtitel -->
        <template #subtitle>
            <p class="m-0 text-700 line-clamp-3">{{ app.summary }}</p>
        </template>

        <!-- Content -->
        <template #content>
            <div class="flex gap-2 flex-wrap">
                <Tag :value="stageLabel" :severity="stageSeverity" />
                <Tag :value="app.orgType === 'overheidsontwikkeld' ? 'Overheidsontwikkeld' : 'Commercieel'"
                    severity="info" />
                <Tag v-for="c in (app.categories || []).slice(0, maxCategories)" :key="c" :value="c"
                    class="surface-50" />
            </div>
            <slot name="extra" />
        </template>

        <!-- Footer -->
        <template #footer>
            <div class="flex justify-content-between align-items-center w-full">
                <small class="text-600" v-if="updatedAt">
                    Laatst bijgewerkt: {{ prettyDate(updatedAt) }}
                </small>
                <div class="ml-auto">
                    <NuxtLink v-if="to" :to="to">
                        <Button :label="ctaLabel" severity="secondary" :icon="ctaIcon" />
                    </NuxtLink>
                    <slot name="actions" />
                </div>
            </div>
        </template>
    </Card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NuxtLink } from '#components'

type Stage = 'productie' | 'pilots' | 'testfase'
type Org = 'overheidsontwikkeld' | 'commercieel'

interface AppLike {
    slug: string
    name: string
    summary?: string
    stage?: Stage
    orgType?: Org
    categories?: string[]
    media?: { logoUrl?: string }
    updatedAt?: string | Date
    logo?: string // fallback veld dat je elders gebruikt
}

const props = withDefaults(defineProps<{
    app: AppLike
    to?: string | null
    ctaLabel?: string
    ctaIcon?: string
    maxCategories?: number
    defaultLogo?: string
}>(), {
    to: null,
    ctaLabel: 'Bekijk',
    ctaIcon: 'pi pi-eye',
    maxCategories: 3,
    defaultLogo: '~/assets/images/model.png'
})

function nav(to: string) { navigateTo(to) }

const logoSrc = computed(() =>
    props.app.media?.logoUrl || props.defaultLogo
)

const stageLabel = computed(() => {
    const s = props.app.stage
    return s === 'productie' ? 'Productie' : s === 'pilots' ? 'Pilots' : s === 'testfase' ? 'Testfase' : '—'
})

const stageSeverity = computed(() => {
    const s = props.app.stage
    if (s === 'productie') return 'success'
    if (s === 'pilots') return 'warning'
    if (s === 'testfase') return 'secondary'
    return 'secondary'
})

const updatedAt = computed(() => props.app.updatedAt)

function prettyDate(d?: string | Date) {
    return d ? new Date(d).toLocaleDateString() : ''
}
</script>

<style scoped>
/* klikbaar gedrag */
.is-clickable {
    cursor: pointer;
}

.is-clickable :deep(.p-card-body) {
    user-select: none;
}

/* titel link zonder blauw/underline */
.heading-link {
    text-decoration: none;
    color: inherit;
}

/* multi-line clamp */
.line-clamp-3 {
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

/* optioneel compactere card */
.app-card :deep(.p-card-title),
.app-card :deep(.p-card-subtitle) {
    margin-bottom: .5rem;
}
</style>
