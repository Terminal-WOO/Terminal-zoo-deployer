<!-- /pages/deploy/new.vue -->
<template>
  <div>
    <main class="surface-0">
      <div class="container p-4">
        <h1 class="m-0">Deployment</h1>
        <p class="mt-2 text-700">Kies één eenvoudige optie om te deployen.</p>

        <!-- Keuze -->
        <Panel header="Kies deployment-type" class="mb-3">
          <div class="grid">
            <div class="col-12 md:col-4">
              <Card :class="['h-full', 'clickable-card', { 'border-1 border-primary': mode === 'self' }]" role="button"
                tabindex="0" :aria-pressed="mode === 'self'" @click="selectMode('self')"
                @keydown.enter.prevent="selectMode('self')" @keydown.space.prevent="selectMode('self')">
                <template #title>Zelf deployen (YAML)</template>
                <template #subtitle>Download een YAML-configuratiebestand</template>
                <template #content>
                  <p class="m-0 text-700">Geschikt voor eigen Kubernetes of CI/CD.</p>
                </template>
                <template #footer>
                  <div class="flex align-items-center gap-2">
                    <RadioButton inputId="mode-self" name="mode" value="self" v-model="mode" />
                    <label for="mode-self">Selecteer</label>
                  </div>
                </template>
              </Card>
            </div>

            <div class="col-12 md:col-4">
              <Card
                :class="['h-full', 'clickable-card', { 'border-1 border-primary': mode === 'azure' }, { 'card-disabled': !azurePossible }]"
                role="button" tabindex="0" :aria-pressed="mode === 'azure'" @click="selectMode('azure')"
                @keydown.enter.prevent="selectMode('azure')" @keydown.space.prevent="selectMode('azure')">
                <template #title>Deploy op custom cluster</template>
                <template #subtitle>Eenvoudige managed uitrol</template>
                <template #content>
                  <p class="m-0 text-700">Registreer deze deployment als custom cluster.</p>
                  <Dropdown
                    v-model="deployments.activeClusterConfig"
                    :options="deployments.clusterDeployments"
                    optionLabel="name"
                    placeholder="Selecteer cluster"
                  />
                </template>
                <template #footer>
                  <div class="flex align-items-center gap-2">
                    <RadioButton inputId="mode-azure" name="mode" value="azure" v-model="mode" />
                    <label for="mode-azure">Selecteer</label>
                  </div>
                </template>
              </Card>
            </div>

            <div class="col-12 md:col-4">
              <Card :class="['h-full', 'clickable-card', { 'border-1 border-primary': mode === 'hosted' }]"
                role="button" tabindex="0" :aria-pressed="mode === 'hosted'" @click="selectMode('hosted')"
                @keydown.enter.prevent="selectMode('hosted')" @keydown.space.prevent="selectMode('hosted')">
                <template #title>Gehost door een partij</template>
                <template #subtitle>Kies hostingprovider</template>
                <template #content>
                  <Dropdown v-model="hoster" :options="hosterOptions" optionLabel="label" optionValue="value"
                    placeholder="Selecteer partij" class="w-full" :disabled="mode !== 'hosted'" />
                  
                  <!-- Database selection for terminal-api -->
                   <div v-if="route?.query?.slug === 'terminal-webapp' && mode === 'hosted' && hoster" class="mt-3">
                    <label for="database-select" class="block mb-2 text-sm font-medium">Selecteer API</label>
                    <Dropdown 
                      id="database-select"
                      v-model="selectedAPIs" 
                      :options="apiOptions" 
                      optionLabel="label" 
                      optionValue="value"
                      placeholder="Selecteer API" 
                      class="w-full" />
                  </div>
                  <div v-if="route?.query?.slug === 'terminal-api' && mode === 'hosted' && hoster" class="mt-3">
                    <label for="database-select" class="block mb-2 text-sm font-medium">Selecteer database</label>
                    <Dropdown 
                      id="database-select"
                      v-model="selectedDatabase" 
                      :options="databaseOptions" 
                      optionLabel="label" 
                      optionValue="value"
                      placeholder="Selecteer database" 
                      class="w-full" />
                  </div>
                </template>
                <template #footer>
                  <div class="flex align-items-center gap-2">
                    <RadioButton inputId="mode-hosted" name="mode" value="hosted" v-model="mode" />
                    <label for="mode-hosted">Selecteer</label>
                  </div>
                </template>
              </Card>
            </div>
          </div>
        </Panel>

        <!-- Basisgegevens -->
        <Panel header="Details" class="mb-3">
          <div class="grid">
            <div class="col-12 md:col-6">
              <div class="field">
                <label for="dep-name">Naam deployment</label>
                <InputText id="dep-name" v-model="meta.name" placeholder="bijv. burgerbrieven-assistant"
                  class="w-full" />
              </div>
            </div>
            <div class="col-12 md:col-6">
              <div class="field">
                <label for="dep-app">App</label>
                <Dropdown id="dep-app" v-model="selectedApp" :options="appChoices" optionLabel="label"
                  optionValue="value" placeholder="Kies een app (optioneel)" class="w-full" :disabled="true" />
              </div>
            </div>
          </div>
          <div class="field">
            <label for="dep-notes">Omschrijving (optioneel)</label>
            <Textarea id="dep-notes" v-model="meta.description" rows="3" autoResize placeholder="Korte beschrijving"
              class="w-full" />
          </div>
        </Panel>

        <Panel header="Configuratie">
          <div class="grid">
            <div class="col-12 md:col-6">
              <div class="field">
                <label for="dep-replicas">CPU cores aanvraag</label>
                <InputText id="dep-replicas" v-model="manualConfig.resources.cpuRequests" :min="1" class="w-full" />
              </div>
            </div>
            <div class="col-12 md:col-6">
              <div class="field">
                <label for="dep-replicas">CPU cores limiet</label>
                <InputText id="dep-replicas" v-model="manualConfig.resources.cpuLimits" :min="1" class="w-full" />
              </div>
            </div>
            <div class="col-12 md:col-6">
              <div class="field">
                <label for="dep-replicas">RAM aanvraag (MiB)</label>
                <InputText id="dep-replicas" v-model="manualConfig.resources.memoryRequests" :min="1" class="w-full" />
              </div>
            </div>
            <div class="col-12 md:col-6">
              <div class="field">
                <label for="dep-replicas">RAM limiet (MiB)</label>
                <InputText id="dep-replicas" v-model="manualConfig.resources.memoryLimits" :min="1" class="w-full" />
              </div>
            </div>
            <div class="col-12 md:col-6">
              <div class="field">
                <label for="dep-replicas">Aantal replicas</label>
                <InputNumber id="dep-replicas" v-model="manualConfig.replicas" :min="1" class="w-full" />
              </div>
            </div>
            <div class="col-12 md:col-6">
              <div class="field">
                <label for="dep-ports">Poorten</label>
                <div v-for="(port, index) in manualConfig.ports" :key="index"
                  class="flex align-items-center gap-2 mb-2">
                  <InputNumber :id="`dep-port-${index}`" v-model="port.ContainerPort" placeholder="bijv. 80"
                    class="w-full" type="number" />
                  <Button icon="pi pi-times" class="p-button-text p-button-danger p-button-sm" type="button"
                    @click="manualConfig.ports.splice(index, 1)" :disabled="manualConfig.ports.length === 1"
                    aria-label="Verwijder poort" />
                </div>
                <Button type="button" icon="pi pi-plus" class="p-button-text p-button-sm mb-2"
                  @click="manualConfig.ports.push({ ContainerPort: 80 })" />
              </div>
            </div>
          </div>
        </Panel>

        <!-- Acties -->
        <div class="flex align-items-center justify-content-between mt-2">
          <div class="flex gap-2">
            <Button v-if="mode === 'self'" label="Download YAML" icon="pi pi-download" :disabled="!canDownloadYaml"
              @click="downloadYaml" />
            <small v-if="mode === 'self' && !yamlDownloaded" class="text-600 align-self-center">Download eerst de YAML
              om
              te
              kunnen deployen.</small>
          </div>

          <div class="flex align-items-center gap-2">
            <Message v-if="errorMsg" severity="error" :closable="false">{{ errorMsg }}</Message>
            <Button label="Deploy" icon="pi pi-cloud-upload" :disabled="!canSave" :loading="busy"
              @click="saveDeployment" />
          </div>
        </div>

        <!-- Success -->
        <Dialog v-model:visible="savedDialog" modal header="Opgeslagen" :closable="false" class="w-11 md:w-5">
          <div class="flex flex-column align-items-center text-center">
            <i class="pi pi-check-circle text-4xl text-green-500 mb-2" aria-hidden="true"></i>
            <h2 class="m-0">Deployment opgeslagen</h2>
            <p class="mt-3 text-700">ID: <b>{{ savedId }}</b></p>
            <div class="flex gap-2 mt-2">
              <NuxtLink :to="`/dashboard?deployment=${savedId}`"><Button label="Naar dashboard"
                  icon="pi pi-external-link" />
              </NuxtLink>
            </div>
          </div>
        </Dialog>

        <!-- Add a cost window to the right when hosted is selected -->
        <Dialog v-if="mode === 'hosted' && !acceptedCosts" position="right" :visible="true" modal :closable="false"
          header="Kostenindicatie" class="w-11 md:w-5">
          <div class="flex flex-column align-items-center text-center">
            <i class="pi pi-info-circle text-4xl text-blue-500 mb-2" aria-hidden="true"></i>
            <h2 class="m-0">Kostenindicatie</h2>
            <p class="mt-3 text-700">De kosten voor gehoste oplossingen variëren afhankelijk van de gekozen provider en
              configuratie.</p>

            <!-- Price indication -->
            <p class="mt-3 text-700">Voor uw geselecteerde hostingprovider en configuratie kunnen de kosten als volgt
            </p>
            <ul class="mt-2 text-left">
              <li>Basis hosting: €20 - €100 per maand</li>
              <li>Extra opslag: €0,10 per GB per maand</li>
              <li>Dataverkeer: €0,05 per GB</li>
              <li>Ondersteuning en SLA's: variabel, afhankelijk van het serviceniveau</li>
            </ul>
            <p class="mt-3 text-700">Deze prijzen zijn indicatief en kunnen variëren op basis van specifieke
              vereisten en gebruik.</p>
            <p class="mt-3 text-700">Neem contact op met de hostingprovider voor een gedetailleerde offerte op maat van
              uw behoeften.</p>

            <div class="flex gap-2 mt-2">
              <Button label="Ik begrijp de kosten" icon="pi pi-check" @click="
                acceptedCosts = true
                " />
            </div>
          </div>
        </Dialog>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { useRoute } from 'vue-router'
  import { useDeployStore } from '~/stores/deploy'
  import { useDeployments } from '~/composables/API/useDeployments'
  import { useDeploymentsStore } from '~/stores/deployments'
  import { useUserStore } from '~/stores/user'
  import type { ManualDeploymentConfig, Resources, Port } from '~/types/types'

  type Mode = 'self' | 'azure' | 'hosted'

  const route = useRoute()
  const deploy = useDeployStore()
  const deployments = useDeploymentsStore()
  const userStore = useUserStore()

  /* eenvoudige state */
  const mode = ref<Mode>('self')
  const hoster = ref<string | null>(null)
  const hosterOptions = [{ label: 'Clappform', value: 'Clappform' }]

  // Database selection for terminal-api
  const selectedDatabase = ref<string | null>(null)
  const databaseOptions = [
    { label: 'Terminal ZOO Database ', value: 'terminal-zoo-database' }
  ]
  const selectedAPIs = ref<string | null>(null)
  const apiOptions = [
    { label: 'Terminal Webapp API', value: 'terminal-webapp-api' }
  ]

  const meta = ref({ name: '', description: '' })
  const selectedApp = ref<string | null>(null)

  const busy = ref(false)
  const errorMsg = ref<string | null>(null)
  const savedDialog = ref(false)
  const savedId = ref<string | null>(null)

  const acceptedCosts = ref(false)

  /* optioneel: prefill vanuit query ?slug=... */
  const appChoices = ref([
    { label: '— geen —', value: null },
    { label: 'GovChat NL', value: 'govchat-nl', imageName: 'ollama/ollama:latest' },
    { label: 'Geo Buddy', value: 'geo-buddy', imageName: 'ollama/ollama:latest' },
    { label: 'Parlement Agent', value: 'parlement-agent', imageName: 'ollama/ollama:latest' },
    { label: 'Terminal API', value: 'terminal-api', imageName: 'terminal/api:latest' },
  ])
  onMounted(() => {
    const slug = typeof route.query.slug === 'string' ? route.query.slug : null
    if (slug && !selectedApp.value) {
      const hit = appChoices.value.find(a => a.value === slug)
      if (hit) selectedApp.value = hit.value
    }
  })

  const azurePossible = computed(() => deployments.clusterDeployments.length > 0)

  // Create through the API
  const manualConfig = ref<ManualDeploymentConfig>({
    deploymentName: "my-deployment",
    namespace: 'testing',
    containerImage: 'nginx:latest',
    replicas: 1,
    ports: [{ ContainerPort: 80 }],
    resources: {
      cpuLimits: '250m',   // in cores
      cpuRequests: '100m',
      memoryLimits: '256Mi', // in bytes
      memoryRequests: '128Mi',
    }
  })


  /* YAML generator (minimaal) + download-tracking */
  const lastDownloadedYaml = ref<string | null>(null)

  function buildYaml() {
    const lines = [
      'apiVersion: v1',
      'kind: AICoDeployment',
      'metadata:',
      `  name: ${meta.value.name}`,
      'spec:',
      `  mode: ${mode.value}`,
      `  app: ${selectedApp.value || 'null'}`,
      '  notes: |',
      `    ${meta.value.description || ''}`.replace(/\n/g, '\n    ')
    ]
    return lines.join('\n') + '\n'
  }

  const yamlDownloaded = computed(() => lastDownloadedYaml.value === buildYaml())

  function downloadYaml() {
    const yaml = buildYaml()
    const blob = new Blob([yaml], { type: 'text/yaml;charset=utf-8' })
    const a = document.createElement('a')
    a.href = URL.createObjectURL(blob)
    a.download = `${meta.value.name.replace(/\s+/g, '-').toLowerCase()}.yaml`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(a.href)
    // Markeer dat exact deze YAML is gedownload
    lastDownloadedYaml.value = yaml
  }

  /* validaties */
  const canDownloadYaml = computed(() => !!meta.value.name)
  const canSave = computed(() => {
    if (!meta.value.name) return false
    if (mode.value === 'hosted')
      if (!hoster.value) return false
    // Check the manual config values with the composable function 'validateResources' from api/deployments
    if (!useDeployments().validateResources(manualConfig.value.resources)) return false
    if (manualConfig.value.replicas < 1) return false
    for (const port of manualConfig.value.ports) {
      if (!port.ContainerPort || port.ContainerPort < 1 || port.ContainerPort > 65535) return false
    }
    // Vereis dat de huidige YAML is gedownload vóór deploy in 'self' modus
    if (mode.value === 'self' && !yamlDownloaded.value) return false
    return true
  })

  /* Klikbaar maken van gehele kaarten */
  function selectMode(next: Mode) {
    mode.value = next

    if (mode.value != next)
      acceptedCosts.value = next === 'hosted' ? false : acceptedCosts.value
    
    // Als je wisselt van modus of potentieel de inhoud wijzigt, dwing opnieuw downloaden af
    if (next === 'self') {
      // reset alleen als de inhoud straks niet meer overeenkomt
      // (het computed checkt al op exacte match, dus hier niets nodig)
    } else {
      // in andere modi speelt YAML niet, maar we resetten om verwarring te voorkomen
      lastDownloadedYaml.value = null
    }
  }

  /* Store save (generiek) */
  async function saveDeployment() {
    let deployment = null
    if (!canSave.value) return
    busy.value = true
    errorMsg.value = null
    try {
      // const id = 'dep_' + Math.random().toString(36).slice(2, 10)
      // const record = {
      //   id,
      //   createdAt: new Date().toISOString(),
      //   status: 'queued',
      //   spec: {
      //     slug: selectedApp.value,
      //     mode: mode.value,                 // 'self' | 'azure' | 'hosted'
      //     provider: mode.value === 'azure' ? 'azure' : (mode.value === 'hosted' ? hoster.value : 'self'),
      //     notes: meta.value.description || '',
      //     resources: { cpu: 2, ramGiB: 4, gpu: 'none' }
      //   },
      //   meta: { name: meta.value.name }
      // }

      if (mode.value === 'hosted') {
        manualConfig.value.deploymentName = meta.value.name.toLowerCase().replace(/\s+/g, '-')
        manualConfig.value.containerImage = appChoices.value.find(a => a.value === selectedApp.value)?.imageName || 'nginx:latest'
        deployment = await useDeployments().createDeployment(manualConfig.value)
        if (typeof deployment !== 'object') throw new Error(`Deployment API gaf fout: ${deployment.status} ${deployment.statusText}`)
      } else if (mode.value === 'azure') {
        if (!deployments.activeClusterConfig) throw new Error('Geen actieve cluster-configuratie gevonden.')
        // Vul de minimale velden in

        deployment = await useDeployments().createDeployment({
          deploymentName: meta.value.name.toLowerCase().replace(/\s+/g, '-'),
          namespace: userStore.activeNamespace || 'testing',
          containerImage: manualConfig.value.containerImage || 'nginx:latest',
          replicas: manualConfig.value.replicas,
          ports: manualConfig.value.ports as Port[],
          resources: manualConfig.value.resources,
          clusterName: deployments.activeClusterConfig.name,
        })
        if (typeof deployment !== 'object') throw new Error(`Deployment API gaf fout: ${deployment.status} ${deployment.statusText}`)
      } else {
        // self-hosted
        const yaml = buildYaml()
        deployment = { metadata: { uid: 'self-' + Math.random().toString(36).slice(2, 10) } } // dummy ID
        // In echte wereld zou je hier de YAML opslaan of iets dergelijks
      }

      // // Probeer nette store API, val terug op simpele push
      // const anyDeploy = deploy as any
      // if (typeof anyDeploy.add === 'function') {
      //   await anyDeploy.add(record)
      // } else if (Array.isArray(anyDeploy.deployments)) {
      //   anyDeploy.deployments.unshift(record)
      // } else {
      //   anyDeploy.recent = [record, ...(Array.isArray(anyDeploy.recent) ? anyDeploy.recent : [])]
      // }

      savedId.value = deployment.metadata.uid
      savedDialog.value = true
    } catch (e: any) {
      errorMsg.value = e?.message || 'Opslaan mislukt.'
    } finally {
      busy.value = false
    }
  }

  function resetForm() {
    savedDialog.value = false
    meta.value = { name: '', description: '' }
    selectedApp.value = null
    mode.value = 'self'
    hoster.value = null
    selectedDatabase.value = null
    selectedAPIs.value = null
    savedId.value = null
    errorMsg.value = null
    lastDownloadedYaml.value = null
  }
</script>

<style scoped>
  .container {
    max-width: 900px;
    margin: 0 auto;
  }

  .border-primary {
    border-color: var(--primary-color) !important;
  }

  .field {
    margin-bottom: 1rem;
  }

  .h-full {
    height: 100%;
  }

  /* Hele kaart klikbaar + focus-stijlen */
  .clickable-card {
    cursor: pointer;
  }

  .clickable-card:focus {
    outline: 2px solid var(--primary-color);
    outline-offset: 2px;
  }

  .clickable-card:focus-within {
    outline: 2px solid var(--primary-color);
    outline-offset: 2px;
  }

  .card-disabled {
    opacity: 0.6;
    /* faded look */
    pointer-events: none;
    /* prevent clicks inside */
    filter: grayscale(100%);
    /* optional grayscale */
    position: relative;
  }

  .card-disabled::after {
    content: "Configureer een custom cluster.";
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.5);
    font-weight: bold;
    color: #070707;
  }
</style>