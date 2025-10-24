<!-- /pages/community/thread/[slug].vue -->
<template>
  <div class="thread-page">
    <!-- Breadcrumbs / Header -->
    <section class="surface-0 border-bottom-1 surface-border">
      <div class="container p-3">
        <Breadcrumb :model="bcItems">
          <template #item="{ item }">
            <NuxtLink v-if="item.to" :to="item.to" class="no-underline">{{ item.label }}</NuxtLink>
            <span v-else>{{ item.label }}</span>
          </template>
        </Breadcrumb>
      </div>
    </section>

    <!-- Thread hero -->
    <section class="surface-0">
      <div class="container p-4">
        <div v-if="thread" class="grid">
          <div class="col-12 lg:col-8">
            <div class="flex align-items-start justify-content-between gap-2">
              <div>
                <h1 class="m-0 text-900">{{ thread.title }}</h1>
                <div class="mt-2 flex align-items-center gap-2 flex-wrap">
                  <Tag :value="thread.category" :severity="tagSeverity(thread.category)" />
                  <Tag :value="thread.status" :severity="statusSeverity(thread.status)" />
                  <Chip :label="thread.app" icon="pi pi-th-large" />
                  <span class="text-600">
                    Door <b>{{ thread.author }}</b> • {{ prettyDateTime(thread.updatedAt) }}
                  </span>
                </div>
              </div>
              <div class="flex gap-2">
                <NuxtLink to="/community" class="no-underline">
                  <Button label="Terug" icon="pi pi-arrow-left" text />
                </NuxtLink>
                <Button :label="subscribed ? 'Gevolgd' : 'Volgen'" :icon="subscribed ? 'pi pi-check' : 'pi pi-bell'" @click="toggleSub" />
              </div>
            </div>

            <!-- Posts -->
            <Panel header="Discussie" class="mt-3">
              <div class="flex flex-column gap-3">
                <div v-for="p in posts" :key="p.id" class="post border-1 surface-border border-round p-3">
                  <div class="flex align-items-center gap-2">
                    <Avatar :label="p.initials" />
                    <div class="flex flex-column">
                      <span class="font-medium">{{ p.user }}</span>
                      <small class="text-600">{{ p.org || '—' }}</small>
                    </div>
                    <span class="ml-auto text-600"><small>{{ prettyDateTime(p.date) }}</small></span>
                  </div>
                  <p class="m-0 mt-2 text-700">{{ p.content }}</p>
                  <div v-if="p.labels?.length" class="mt-2 flex gap-2 flex-wrap">
                    <Tag v-for="l in p.labels" :key="l" :value="l" />
                  </div>
                </div>
              </div>
            </Panel>

            <!-- New reply -->
            <Panel header="Reageer" class="mt-3">
              <div class="field">
                <label class="block mb-2">Je reactie</label>
                <Textarea v-model="replyText" rows="5" autoResize placeholder="Schrijf je antwoord…" />
              </div>
              <div class="field">
                <label class="block mb-2">Labels (optioneel)</label>
                <Chips v-model="replyLabels" separator="," placeholder="Bijv. DPIA, Open-source" />
              </div>
              <div class="flex justify-content-end gap-2">
                <Button label="Plaatsen" icon="pi pi-send" :disabled="!replyText.trim()" @click="submitReply" />
              </div>
              <Message v-if="replyMsg" :severity="replyOk ? 'success' : 'warn'" :closable="false" class="mt-3">
                {{ replyMsg }}
              </Message>
            </Panel>
          </div>

          <!-- Sidebar -->
          <div class="col-12 lg:col-4">
            <Panel header="Thread-info">
              <ul class="list-none p-0 m-0">
                <li class="flex align-items-center justify-content-between py-2 border-bottom-1 surface-border">
                  <span class="text-700">Status</span>
                  <Tag :value="thread.status" :severity="statusSeverity(thread.status)" />
                </li>
                <li class="flex align-items-center justify-content-between py-2 border-bottom-1 surface-border">
                  <span class="text-700">Reacties</span>
                  <b>{{ thread.replies }}</b>
                </li>
                <li class="flex align-items-center justify-content-between py-2 border-bottom-1 surface-border">
                  <span class="text-700">Categorie</span>
                  <Tag :value="thread.category" :severity="tagSeverity(thread.category)" />
                </li>
                <li class="flex align-items-center justify-content-between py-2">
                  <span class="text-700">App</span>
                  <Chip :label="thread.app" icon="pi pi-th-large" />
                </li>
              </ul>
              <Divider />
              <div class="flex gap-2">
                <Button label="Markeer" icon="pi pi-flag" outlined />
                <Button label="Delen" icon="pi pi-share-alt" text @click="copyLink" />
              </div>
              <Message v-if="copied" severity="success" :closable="false" class="mt-2">Link gekopieerd</Message>
            </Panel>

            <Panel header="Gerelateerde threads" class="mt-3">
              <div class="flex flex-column gap-2">
                <div v-for="rt in related" :key="rt.id" class="border-1 surface-border p-2 border-round">
                  <NuxtLink :to="`/community/thread/${rt.id}`" class="no-underline font-medium text-900 hover:text-primary">
                    {{ rt.title }}
                  </NuxtLink>
                  <div class="mt-1 flex align-items-center gap-2">
                    <Tag :value="rt.category" :severity="tagSeverity(rt.category)" />
                    <small class="text-600">{{ prettyDate(rt.updatedAt) }}</small>
                  </div>
                </div>
                <div v-if="!related.length" class="text-600">Geen gerelateerde items.</div>
              </div>
            </Panel>
          </div>
        </div>

        <!-- Not found -->
        <div v-else>
          <Message severity="warn" :closable="false">Thread niet gevonden.</Message>
          <NuxtLink to="/community" class="inline-block mt-2"><Button label="Terug naar community" icon="pi pi-arrow-left" text /></NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useSeoMeta } from '#imports'
import { useCommunityStore } from '~/stores/community'

const route = useRoute()
const id = route.params.slug as string
const store = useCommunityStore()

onMounted(() => store.init()) // seed dummy data indien leeg

/* Thread ophalen uit store */
const thread = computed(() => store.threads.find(t => t.id === id))

/* SEO */
useSeoMeta({
  title: () => thread.value ? `${thread.value.title} — Community` : 'Community thread',
  description: () => thread.value ? `${thread.value.category} • ${thread.value.app}` : 'Community discussie'
})

/* Breadcrumbs */
const bcItems = computed(() => ([
  { label: 'Home', to: '/' },
  { label: 'Community', to: '/community' },
  { label: thread.value?.title || 'Thread' }
]))

/* Dummy posts (lokaal voor weergave; replies tellen via store) */
type Post = { id: string; user: string; initials: string; org?: string; date: string; content: string; labels?: string[] }
const posts = ref<Post[]>([])
function seedPosts() {
  if (!thread.value) { posts.value = []; return }
  // eenvoudige seed op basis van id, met fallback content
  const base = [
    { user: 'Anouk Visser', org: 'Gemeente Delft', content: 'We hebben de DPIA afgerond; template was duidelijk. Let op paragraaf 3.2.' },
    { user: 'Ruben Jansen', org: 'Provincie Utrecht', content: 'Rate limiting in n8n opgelost met een reverse proxy en burst limit.' },
    { user: 'Milan de Boer', org: 'Rijksoverheid', content: 'TLS-config is nu A+ in test; deel straks de Ansible role.' }
  ]
  posts.value = base.slice(0, 2 + ((id.charCodeAt(0) + id.length) % 2)).map((b, i) => ({
    id: `p${i}`,
    user: b.user,
    initials: initials(b.user),
    org: b.org,
    date: new Date(Date.now() - (i+1) * 3600_000).toISOString(),
    content: b.content,
    labels: i === 0 ? ['DPIA', 'Open-source'] : (i === 1 ? ['Techniek'] : [])
  }))
}
seedPosts()

/* Nieuwe reply */
const replyText = ref('')
const replyLabels = ref<string[]>([])
const replyMsg = ref<string | null>(null)
const replyOk = ref(false)

function submitReply() {
  replyMsg.value = null; replyOk.value = false
  if (!thread.value) return
  const p: Post = {
    id: 'p' + Math.random().toString(36).slice(2),
    user: 'Jij',
    initials: 'JIJ',
    org: '—',
    date: new Date().toISOString(),
    content: replyText.value.trim(),
    labels: [...replyLabels.value]
  }
  if (!p.content) { replyMsg.value = 'Bericht is leeg.'; return }
  posts.value.push(p)
  store.addReply(thread.value.id) // verhoog teller in store
  replyOk.value = true
  replyMsg.value = 'Reactie geplaatst.'
  replyText.value = ''
  replyLabels.value = []
}

/* Sidebar: related threads */
const related = computed(() => {
  if (!thread.value) return []
  // zelfde categorie, andere id; meest recent eerst
  return store.threads
    .filter(t => t.category === thread.value!.category && t.id !== thread.value!.id)
    .sort((a,b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
    .slice(0, 5)
})

/* Subscriben + copy link */
const subscribed = ref(false)
const toggleSub = () => { subscribed.value = !subscribed.value }
const copied = ref(false)
async function copyLink() {
  try {
    await navigator.clipboard.writeText(location.href)
    copied.value = true
    setTimeout(() => (copied.value = false), 1200)
  } catch { /* ignore */ }
}

/* Helpers */
function initials(name?: string) {
  const s = (name || '').trim()
  return s ? s.split(/\s+/).map(p => p[0]).slice(0,2).join('').toUpperCase() : 'U'
}
function tagSeverity(c: string) {
  if (c === 'Security') return 'danger'
  if (c === 'Privacy') return 'info'
  if (c === 'Techniek') return 'success'
  if (c === 'Ethiek') return 'warning'
  return 'secondary'
}
function statusSeverity(s: string) {
  if (s === 'Open') return 'info'
  if (s === 'Gesloten') return 'secondary'
  if (s === 'Opgelost') return 'success'
  return 'secondary'
}
function prettyDate(d?: string | Date) { return d ? new Date(d).toLocaleDateString() : '—' }
function prettyDateTime(d?: string | Date) { return d ? new Date(d).toLocaleString() : '—' }
</script>

<style scoped>
.container { max-width: 1100px; margin: 0 auto; }
.post { background: var(--surface-card); }
.field {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}
</style>
