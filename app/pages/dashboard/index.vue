<!-- /pages/dashboard.vue -->
<template>
  <div class="dashboard">
    <section class="surface-0">
      <div class="container p-0">
        <div class="grid no-gutter">
          <!-- MAIN CONTENT -->
          <main class="col-12 p-3">
            <!-- HERO -->
            <section :id="ids.hero" class="surface-0 border-bottom-1 surface-border pb-3">
              <div class="p-3">
                <div class="flex align-items-center gap-3">
                  <Avatar :label="initials(
                    userStore.session?.name || 'User'
                  )
                    " size="large" />
                  <div>
                    <h1 class="m-0 text-900">
                      Welkom,
                      {{
                        userStore.session?.name ||
                        "gebruiker"
                      }}
                    </h1>
                    <p class="m-0 text-700">
                      {{ orgLabel }} · {{ today }}
                    </p>
                  </div>
                </div>
              </div>
            </section>

            <!-- KPI's -->
            <section :id="ids.kpis" class="surface-0">
              <div class="p-3 pt-4">
                <div class="grid">
                  <div class="col-12 md:col-6 lg:col-3">
                    <Card class="kpi">
                      <template #content>
                        <div class="flex align-items-center justify-content-between">
                          <div>
                            <small class="text-600">Apps in
                              catalogus</small>
                            <div class="kpi-num">
                              {{ appsCount }}
                            </div>
                          </div>
                          <i class="pi pi-th-large text-2xl text-primary"></i>
                        </div>
                      </template>
                    </Card>
                  </div>
                  <div class="col-12 md:col-6 lg:col-3">
                    <Card class="kpi">
                      <template #content>
                        <div class="flex align-items-center justify-content-between">
                          <div>
                            <small class="text-600">Open
                              advisories</small>
                            <div class="kpi-num">
                              {{
                                advisoriesCount
                              }}
                            </div>
                          </div>
                          <i class="pi pi-exclamation-triangle text-2xl text-primary"></i>
                        </div>
                      </template>
                    </Card>
                  </div>
                  <div class="col-12 md:col-6 lg:col-3">
                    <Card class="kpi">
                      <template #content>
                        <div class="flex align-items-center justify-content-between">
                          <div>
                            <small class="text-600">Community
                              threads</small>
                            <div class="kpi-num">
                              {{ threadsCount }}
                            </div>
                          </div>
                          <i class="pi pi-comments text-2xl text-primary"></i>
                        </div>
                      </template>
                    </Card>
                  </div>
                  <div class="col-12 md:col-6 lg:col-3">
                    <Card class="kpi">
                      <template #content>
                        <div class="flex align-items-center justify-content-between">
                          <div>
                            <small class="text-600">Deployments</small>
                            <div class="kpi-num">
                              {{
                                deploymentsCount
                              }}
                            </div>
                          </div>
                          <i class="pi pi-cloud-upload text-2xl text-primary"></i>
                        </div>
                      </template>
                    </Card>
                  </div>
                </div>
              </div>
            </section>

            <!-- DEPLOYMENTS -->
            <section :id="ids.deploy" class="surface-0" v-if="userStore.hasRole('developer') || userStore.hasRole('administrator')">
              <div class="p-3 pt-0">
                <Panel header="Deployments (laatste 10)">
                  <DataTable :value="recentDeployments" :rows="6" :paginator="true" responsiveLayout="scroll"
                    size="small">
                    <Column field="id" header="ID" style="width: 11rem">
                      <template #body="{ data }"><code>{{
                        shortId(data.id)
                      }}</code></template>
                    </Column>
                    <Column header="Naam" style="min-width: 14rem">
                      <template #body="{ data }">
                        {{
                          data.config
                            .deploymentName || "—"
                        }}
                      </template>
                    </Column>
                    <Column header="App" style="min-width: 12rem">
                      <template #body="{ data }">{{
                        data.config.containerImage ||
                        "—"
                      }}</template>
                    </Column>
                    <Column header="Doel" style="min-width: 14rem">
                      <template #body="{ data }">
                        {{
                          data.config?.provider ||
                          "Clappform"
                        }}
                        <small v-if="data.config?.region">({{
                          data.config?.region ||
                          "EU West"
                        }})</small>
                      </template>
                    </Column>
                    <Column header="Resources" style="min-width: 12rem">
                      <template #body="{ data }">
                        {{
                          data.config?.resources
                            ?.cpuRequests || "—"
                        }}
                        Cores /
                        {{
                          data.config?.resources
                            ?.memoryRequests || "—"
                        }}
                      </template>
                    </Column>
                    <Column header="Status" style="width: 10rem">
                      <template #body="{ data }">
                        <Tag :value="data.status" :severity="statusTag(data.status)
                          " />
                      </template>
                    </Column>
                    <Column header="Geplaatst" style="width: 10rem">
                      <template #body="{ data }"><small>{{
                        prettyDate(
                          data.creationTimestamp
                        )
                          }}</small></template>
                    </Column>
                    <Column header="Acties" style="width: 14rem">
                      <template #body="{ data }">
                        <div class="flex gap-2">
                          <Button icon="pi pi-eye" text rounded @click="goToApp(data.id)" />
                          <Button icon="pi pi-refresh" text rounded @click="
                            quickRedeploy(
                              data.id
                            )
                            " />
                          <Button icon="pi pi-times" text rounded severity="danger" @click="
                            cancelDeployment(
                              data.config
                                .deploymentName
                            )
                            " />
                        </div>
                      </template>
                    </Column>
                  </DataTable>
                </Panel>
              </div>
            </section>

            <section :id="ids.configSecrets" class="surface-0" v-if="userStore.hasRole('developer') || userStore.hasRole('administrator')">
              <div class="p-3 pt-0">
                <Panel header="Config mappen & Secrets">
                  <div class="grid gap-4">
                    <div class="col-12">
                      <div v-if="configMaps.length === 0" class="p-3">
                        <h3>Config Mappen</h3>
                        <p class="m-0 text-600">
                          Je hebt nog geen config mappen
                          aangemaakt.
                        </p>
                        <Button icon="pi pi-plus" label="Voeg Config Map toe" @click="openConfigMapModal" />
                      </div>
                      <div v-else>
                        <DataTable :value="configMaps" :rows="6" :paginator="true" responsiveLayout="scroll">
                          <template #header>
                            <div class="flex flex-wrap items-center justify-between" style="justify-content: space-between;">
                              <span class="text-xl font-bold">Geconfigureerde
                                Config Mappen</span>
                              <Button icon="pi pi-plus" rounded @click="
                                openConfigMapModal()
                                " size="small" />
                            </div>
                          </template>
                          <Column header="Naam" style="min-width: 14rem">
                            <template #body="{ data }">
                              <span>{{ data.metadata.name }}</span>
                            </template>
                          </Column>
                          <Column header="Namespace" style="min-width: 14rem">
                            <template #body="{ data }">
                              <span>{{ data.metadata.namespace }}</span>
                            </template>
                          </Column>
                        </DataTable>
                      </div>
                    </div>
                    <Divider />
                    <div class="col-12">
                      <div v-if="deploymentsStore.secrets.length === 0" class="p-3">
                        <h3>Secrets</h3>
                        <p class="m-0 text-600">
                          Je hebt nog geen secrets
                          aangemaakt.
                        </p>
                        <Button icon="pi pi-plus" label="Voeg Secret toe" @click="openSecretModal" />
                      </div>
                      <div v-else>
                        <DataTable :value="secrets" :rows="6" :paginator="true" responsiveLayout="scroll">
                          <template #header>
                            <div class="flex flex-wrap items-center justify-between" style="justify-content: space-between;">
                              <span class="text-xl font-bold">Geconfigureerde
                                Secrets</span>
                              <Button icon="pi pi-plus" rounded @click="
                                openSecretModal()
                                " size="small" />
                            </div>
                          </template>
                          <Column header="Naam" style="min-width: 14rem">
                            <template #body="{ data }">
                              <span>{{ data.metadata.name }}</span>
                            </template>
                          </Column>
                          <Column header="Namespace" style="min-width: 14rem">
                            <template #body="{ data }">
                              <span>{{ data.metadata.namespace }}</span>
                            </template>
                          </Column>
                        </DataTable>
                      </div>
                    </div>
                  </div>
                </Panel>

              </div>
            </section>

            <!-- Cluster configuratie -->
            <section :id="ids.clusters" class="surface-0" v-if="userStore.hasRole('developer') || userStore.hasRole('administrator')">
              <div class="p-3 pt-0">
                <Panel header="Custom cluster Configuratie">
                  <div v-if="clusterConfigs.length === 0" class="p-3">
                    <p class="m-0 text-600">
                      Je hebt nog geen cluster geconfigureerd.
                    </p>
                    <Button type="button" label="Configureer cluster" icon="pi pi-plus" class="mt-2" @click="openClusterConfigModal(clusterConfig)" />
                  </div>
                  <div v-else>
                    <DataTable :value="clusterConfigs" :rows="6" :paginator="true" responsiveLayout="scroll"
                      size="small">
                      <template #header>
                        <div class="flex flex-wrap items-center justify-between" style="justify-content: space-between;">
                          <span class="text-xl font-bold">Geconfigureerde
                            Clusters</span>
                          <Button icon="pi pi-plus" rounded @click="
                            openClusterConfigModal(
                              clusterConfig
                            )
                            " size="small" />
                        </div>
                      </template>
                      <Column header="Naam" style="min-width: 14rem">
                        <template #body="{ data }">
                          <span>{{ data.name }}</span>
                        </template>
                      </Column>
                      <Column header="Geselecteerd" style="width: 10rem">
                        <template #body="{ data }">
                          <Tag v-if="deploymentsStore.activeClusterConfig?.name === data.name" value="Actief" label="Active" severity="success" />
                        </template>
                      </Column>
                      <Column header="API Server" style="min-width: 18rem">
                        <template #body="{ data }">
                          <span>{{ data.server }}</span>
                        </template>
                      </Column>
                      <Column header="K8s Versie" style="width: 10rem">
                        <template #body="{ data }">
                          <Tag :value="data.k8sVersion" severity="info" />
                        </template>
                      </Column>
                      <Column header="Standaard Namespace" style="width: 16rem">
                        <template #body="{ data }">
                          <span>{{ data.defaultNamespace }}</span>
                        </template>
                      </Column>
                      <Column header="Acties" style="width: 14rem">
                        <template #body="{ data }">
                          <div class="flex gap-2">
                            <!-- set cluster config active -->
                            <Button
                              v-if="deploymentsStore.activeClusterConfig?.name !== data.name"
                              icon="pi pi-check"
                              text
                              rounded
                              @click="
                                setActiveClusterConfig(data)
                              "
                            />
                          </div>
                        </template>
                      </Column>
                    </DataTable>
                  </div>
                </Panel>
              </div>


            </section>

            <!-- UPLOADED APPLICATIONS -->
            <section :id="ids.uploaded" class="surface-0" v-if="userStore.hasRole('developer') || userStore.hasRole('administrator')">
              <div class="p-3 pt-0">
                <Panel header="Uploaded Applications">
                  <div v-if="uploadedImages.length === 0" class="p-3">
                    <p class="m-0 text-600">
                      Je hebt nog geen applicaties
                      geüpload.
                    </p>
                    <Button type="button" label="Creëer applicatie" icon="pi pi-plus" class="mt-2" @click="openCreateModal()" />
                  </div>
                  <div v-else>
                    <DataTable :value="uploadedImages" :rows="6" :paginator="true" responsiveLayout="scroll"
                      size="small" >
                      <template #header>
                        <div class="flex flex-wrap items-center justify-between w-full" style="justify-content: space-between;">
                          <span class="text-xl font-bold">Geüploade
                            Applicaties</span>
                          <Button icon="pi pi-plus" rounded @click="
                            openCreateModal()
                            " size="small" />
                        </div>
                      </template>
                      <Column header="App" style="min-width: 14rem">
                        <template #body="{ data }">
                          <span>{{ data.name }}</span>
                        </template>
                      </Column>
                      <Column header="Fase" style="width: 10rem">
                        <template #body="{ data }">
                          <Tag :value="data.stage" severity="info" />
                        </template>
                      </Column>
                      <Column header="Categorieën" style="min-width: 14rem">
                        <template #body="{ data }">
                          <div class="flex flex-wrap gap-1">
                            <Tag v-for="cat in data.categories" :key="cat" :value="cat" severity="info" />
                          </div>
                        </template>
                      </Column>
                      <Column header="Review status" style="width: 10rem">
                        <template #body="{ data }">
                          <Tag :value="statusLabel(data.approved)" :severity="statusTagReview(data.approved)" />
                        </template>
                      </Column>
                      <Column header="Laatst bijgewerkt" style="width: 14rem">
                        <template #body="{ data }">
                          <small>{{
                            prettyDate(data.updatedAt)
                          }}</small>
                        </template>
                      </Column>
                      <Column header="Acties" style="width: 14rem">
                        <template #body="{ data }">
                          <div class="flex gap-2">
                            <Button icon="pi pi-pencil" text rounded @click="editApplication(data)" />
                            <Button icon="pi pi-times" text rounded severity="danger" @click="
                            deleteApplication(uploadedImages.indexOf(
                                  data
                                ));
                            " />
                          </div>
                        </template>
                      </Column>

                    </DataTable>
                  </div>
                </Panel>
              </div>
            </section>

            <!-- <section :id="ids.hosts" class="surface-0">
              <div class="p-3 pt-0">
                <Panel header="Host clusters">

                </Panel>
              </div>
            </section> -->

            <!-- UPDATES -->
            <section :id="ids.updates" class="surface-0">
              <div class="p-3 pt-0">
                <Panel header="Updates — Community & Security">
                  <DataTable :value="activityRows" :rows="8" :paginator="true" responsiveLayout="scroll" size="small">
                    <Column header="Type" style="width: 10rem">
                      <template #body="{ data }">
                        <div class="flex align-items-center gap-2">
                          <i :class="data.icon"></i>
                          <span>{{ data.type }}</span>
                        </div>
                      </template>
                    </Column>

                    <Column header="Titel" style="min-width: 18rem">
                      <template #body="{ data }">
                        <NuxtLink v-if="data.link" :to="data.link"
                          class="no-underline font-medium text-900 hover:text-primary">
                          {{ data.title }}
                        </NuxtLink>
                        <span v-else>{{
                          data.title
                        }}</span>
                      </template>
                    </Column>

                    <Column header="Label" style="width: 10rem">
                      <template #body="{ data }">
                        <Tag :value="data.label" :severity="data.labelSeverity
                          " />
                      </template>
                    </Column>

                    <Column header="Extra" style="min-width: 14rem">
                      <template #body="{ data }">
                        <span v-if="data.extra">{{
                          data.extra
                        }}</span>
                        <span v-else>—</span>
                      </template>
                    </Column>

                    <Column header="Datum" style="width: 10rem">
                      <template #body="{ data }">
                        <small>{{
                          prettyDate(data.date)
                        }}</small>
                      </template>
                    </Column>
                  </DataTable>
                </Panel>
              </div>
            </section>

            <!-- CREATE / EDIT CLUSTER CONFIG MODAL -->
            <Dialog header="Cluster Config" v-model:visible="clusterConfigModalOpen" :modal="true" :closable="false" :style="{ width: '600px' }">
              <template #header :style="">
                <h3 class="m-0">Cluster Configuratie</h3>
              </template>
              <div class="p-3">
                <div class="field">
                  <label for="cluster-name" class="block mb-2">Cluster Naam</label>
                  <InputText id="cluster-name" v-model="clusterConfig.name" class="w-full" />
                </div>
                <div class="field">
                  <label for="cluster-k8s-version" class="block mb-2">URL</label>
                  <InputText id="cluster-k8s-version" v-model="clusterConfig.domain" class="w-full" />
                </div>
                <div class="field">
                  <label for="cluster-namespace" class="block mb-2">TLS Key</label>
                  <InputText id="cluster-namespace" v-model="clusterConfig.privateKey" class="w-full" />
                </div>
                <div class="field">
                  <label for="cluster-k8s-version" class="block mb-2">TLS Certificaat</label>
                  <InputText id="cluster-k8s-version" v-model="clusterConfig.Certificate" class="w-full" />
                </div>
                <div class="field">
                  <label for="cluster-auth" class="block mb-2">Authenticatie Methode</label>
                  <Dropdown id="cluster-auth" v-model="clusterConfig.auth" optionLabel="label" optionValue="value" :options="[
                    { label: 'Kubeconfig', value: 'kubeconfig' },
                    { label: 'In-Cluster', value: 'inCluster' },
                    { label: 'Service Account Token', value: 'serviceAccountToken' },
                  ]" class="w-full" />
                </div>
                <div v-if="clusterConfig.auth === 'kubeconfig'" class="field">
                  <label for="cluster-kubeconfig" class="block mb-2">Kubeconfig (YAML)</label>
                  <Textarea id="cluster-kubeconfig" v-model="clusterConfig.kubeconfig" class="w-full" rows="6" />
                </div>
                <div v-if="clusterConfig.auth === 'serviceAccountToken'">
                  <div class="field">
                    <label for="cluster-server" class="block mb-2">API Server URL</label>
                    <InputText id="cluster-server" v-model="clusterConfig.server" class="w-full" />
                  </div>
                  <div class="field">
                    <label for="cluster-bearerToken" class="block mb-2">Bearer Token</label>
                    <InputText id="cluster-bearerToken" v-model="clusterConfig.bearerToken" class="w-full" />
                  </div>
                  <div class="field">
                    <label for="cluster-caPEM" class="block mb-2">CA Certificate (PEM)</label>
                    <Textarea id="cluster-caPEM" v-model="clusterConfig.caPEM" class="w-full" rows="4" />
                  </div>
                </div>

                <Accordion>
                  <AccordionTab header="Hoe maak ik een Service Account Token aan?">
                    <div class="p-3">
                      <p>Je kunt een Service Account Token aanmaken met behulp van de volgende stappen:</p>
                      <code-block language="bash">
                        <pre><code>
  # 1) Create a namespace (optional)
  kubectl create ns ext-api

  # 2) Create a ServiceAccount
  kubectl -n ext-api create serviceaccount api-client

  # 3) Bind minimal RBAC (example: scale deployments in any namespace)
  kubectl apply -f - <<'YAML'
  apiVersion: rbac.authorization.k8s.io/v1
  kind: ClusterRole
  metadata: { name: external-scaler }
  rules:
  - apiGroups: ["apps"]
    resources: ["deployments/scale"]
    verbs: ["get","update","patch"]
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get","list","watch"]
  YAML

  kubectl create clusterrolebinding external-scaler-binding \
    --clusterrole=external-scaler \
    --serviceaccount=ext-api:api-client

  # 4) Get the API server URL & CA
  SERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
  CADATA_B64=$(kubectl config view --raw --minify -o jsonpath='{.clusters[0].cluster.certificate-authority-data}')
  echo "$CADATA_B64" | base64 -d > ca.pem

  # 5) Get a token for the SA (short-lived by default)
  # Kubernetes v1.24+: use TokenRequest API via kubectl
  TOKEN=$(kubectl -n ext-api create token api-client --duration=3600s)
  echo "$TOKEN"
  </code></pre>
                        </code-block>
                    </div>
                  </AccordionTab>
                </Accordion>
                
              </div>
              <template #footer>
                <div class="flex justify-content-end gap-2">
                  <Button label="Annuleren" icon="pi pi-times" class="p-button-text" @click="clusterConfigModalOpen = false" />
                  <Button label="Opslaan" icon="pi pi-check"  @click="saveClusterConfig()" />
                </div>
              </template>
            </Dialog>

            <!-- CREATE / EDIT APPLICATION MODAL -->
            <Dialog header="Nieuwe applicatie" v-model:visible="createModalOpen" :modal="true" :closable="false"
              :style="{ width: '600px' }">
              <template #header>
                <h3 class="m-0 text-900">Nieuwe applicatie</h3>
              </template>
              <div class="p-3">
                <div class="field">
                  <label for="app-name" class="block mb-2">Naam</label>
                  <InputText id="app-name" v-model="appConfig.name" class="w-full" />
                </div>
                <div class="field">
                  <label for="app-subtitle" class="block mb-2">Subtitel</label>
                  <InputText id="app-subtitle" v-model="appConfig.subtitle" class="w-full" />
                </div>
                <div class="field">
                  <label for="app-summary" class="block mb-2">Korte omschrijving</label>
                  <Textarea id="app-summary" v-model="appConfig.summary" class="w-full" rows="2" />
                </div>
                <div class="field">
                  <label for="app-stage" class="block mb-2">Fase</label>
                  <Dropdown id="app-stage" v-model="appConfig.stage" :options="[
                    'productie',
                    'pilots',
                    'testfase',
                  ]" class="w-full" />
                </div>
                <div class="field">
                  <label for="app-orgtype" class="block mb-2">Type organisatie</label>  
                  <Dropdown id="app-orgtype" v-model="appConfig.orgType" :options="[
                    'overheidsontwikkeld',
                    'commercieel',
                  ]" class="w-full" :disabled="true" />
                </div>
                <div class="field">
                  <label for="app-categories" class="block mb-2">Categorieën</label>
                  <Chips id="app-categories" v-model="appConfig.categories" class="w-full" />
                </div>
                <div class="field">
                  <label for="app-logo" class="block mb-2">Logo URL</label>
                  <InputText id="app-logo" v-model="appConfig.media.logoUrl" class="w-full" />
                </div>
                <div class="field">
                  <label for="app-screenshots" class="block mb-2">Screenshots (komma-gescheiden
                    URLs)</label>
                  <InputText id="app-screenshots" v-model="appConfig.media.screenshots" class="w-full" />
                </div>
                <div class="field">
                  <label for="app-description" class="block mb-2">Beschrijving (Markdown)</label>
                  <Textarea id="app-description" v-model="appConfig.descriptionMd" class="w-full" rows="4" />
                </div>
                <div class="field">
                  <label for="app-pricing" class="block mb-2">Prijsmodel</label>
                  <Dropdown id="app-pricing" v-model="appConfig.pricing" :options="pricingOptions" optionLabel="label"
                    optionValue="value" class="w-full" />
                </div>
                <!-- Show fields based on pricing model -->
                <div v-if="typeof appConfig.pricing === 'object' &&
                  appConfig.pricing !== null &&
                  'type' in appConfig.pricing && appConfig.pricing.type === 'one_time'">
                  <div class="field">
                    <label for="app-price-onetime" class="block mb-2">Eenmalige prijs (EUR)</label>
                    <InputNumber id="app-price-onetime" v-model="appConfig.pricing.price" mode="currency" currency="EUR"
                      locale="nl-NL" class="w-full" />
                  </div>
                </div>
                <div v-if="
                  typeof appConfig.pricing === 'object' &&
                  appConfig.pricing !== null &&
                  'type' in appConfig.pricing &&
                  appConfig.pricing.type === 'subscription'
                ">
                  <div class="field">
                    <label for="app-price-subscription" class="block mb-2">Abonnementsprijs (EUR)</label>
                    <InputNumber id="app-price-subscription" v-model="appConfig.pricing.price" mode="currency"
                      currency="EUR" locale="nl-NL" class="w-full" />
                  </div>
                </div>
                <div v-if="typeof appConfig.pricing === 'object' &&
                  appConfig.pricing !== null &&
                  'type' in appConfig.pricing && appConfig.pricing.type === 'usage'">
                  <div class="field">
                    <label for="app-usage-unit" class="block mb-2">Gebruiksunit</label>
                    <InputText id="app-usage-unit" v-model="appConfig.pricing.unit" class="w-full" />
                  </div>
                </div>
                <div class="field">
                  <label for="app-orgname" class="block mb-2">Organisatienaam</label>
                  <InputText id="app-orgname" v-model="appConfig.org.name" class="w-full" />
                </div>
                <div class="field">
                  <label for="app-orgurl" class="block mb-2">Organisatie URL</label>
                  <InputText id="app-orgurl" v-model="appConfig.org.url" class="w-full" />
                </div>
                <div class="field">
                  <label for="app-repo" class="block mb-2">Repository URL</label>
                  <InputText id="app-repo" v-model="appConfig.imageName" class="w-full" />
                </div>
                <!-- Versions -->
                <label for="app-versions" class="block mb-2">Versies</label>
                <div class="field" v-for="version in appConfig.versions" :key="version.id">
                  <div class="p-3 border-1 surface-border border-round">
                    <div class="field">
                      <label for="version-number" class="block mb-2">Versienummer</label>
                      <InputText id="version-number" v-model="version.version" class="w-full" />
                    </div>
                    <div class="field">
                      <label for="version-date" class="block mb-2">Datum</label>
                      <Calendar id="version-date" v-model="version.date" dateFormat="dd-mm-yy" showIcon
                        class="w-full" />
                    </div>
                    <div class="field">
                      <label for="version-notes" class="block mb-2">Release notes (Markdown)</label>
                      <Textarea id="version-notes" v-model="version.notes" class="w-full" rows="3" />
                    </div>
                    <Button label="Verwijder versie" icon="pi pi-times" class="mt-2" @click="
                      appConfig.versions.splice(
                        appConfig.versions.indexOf(
                          version
                        ),
                        1
                      )
                      " />
                  </div>
                </div>
                <Button label="Voeg versie toe" icon="pi pi-plus" class="mt-2" @click="
                  appConfig.versions.push({
                    version: '',
                    date: new Date().toISOString(),
                    notes: '',
                  })
                  " />
              </div>
              <template #footer>
                <Button type="button" label="Annuleren" icon="pi pi-times" class="p-button-text" @click="createModalClose()" />
                <Button type="button" label="Opslaan" icon="pi pi-check" @click="createApplication()" />
              </template>
            </Dialog>

            <!-- CREATE Config Map -->
            <Dialog header="Config map configuratie" v-model:visible="configMapModalOpen" :modal="true" :closable="false"
            :style="{ width: '600px' }">
            <template #header>
              <h3 class="m-0 text-900">Config map aanmaken</h3>
            </template>
            <div class="p-3">
              <div class="field">
                <label for="configmap-name" class="block mb-2">Naam</label>
                <InputText id="configmap-name" v-model="configMapConfig.name" class="w-full" />
              </div>
              <div class="field">
                <label for="configmap-data" class="block mb-2">Data (key=value, één per regel)</label>
                <Textarea id="configmap-data" v-model="configMapConfig.data" class="w-full" rows="6" />
              </div>
            </div>
            <template #footer>
              <Button type="button" label="Annuleren" icon="pi pi-times" class="p-button-text" @click="configMapModalOpen = false" />
              <Button type="button" label="Opslaan" icon="pi pi-check" @click="createConfigMap()" />
            </template>
            </Dialog>

            <!-- CREATE SECRET -->
            <Dialog header="Secret configuratie" v-model:visible="secretModalOpen" :modal="true" :closable="false"
            :style="{ width: '600px' }">
            <template #header>
              <h3 class="m-0 text-900">Secret aanmaken</h3>
            </template>
            <div class="p-3">
              <div class="field">
                <label for="configmap-name" class="block mb-2">Naam</label>
                <InputText id="configmap-name" v-model="secretConfig.name" class="w-full" />
              </div>
              <div class="field">
                <label for="configmap-data" class="block mb-2">Data (key=value, één per regel)</label>
                <Textarea id="configmap-data" v-model="secretConfig.data" class="w-full" rows="6" />
              </div>
            </div>
            <template #footer>
              <Button type="button" label="Annuleren" icon="pi pi-times" class="p-button-text" @click="secretModalOpen = false" />
              <Button type="button" label="Opslaan" icon="pi pi-check" @click="createSecret()" />
            </template>
            </Dialog>

          </main>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
  definePageMeta({ middleware: ["auth"], layout: "dashboard" });

  import { computed, onMounted, ref } from "vue";
  import { useUserStore } from "~/stores/user";
  import { useAppCatalog } from "~/stores/appCatalog";
  import { useCriteriaStore } from "~/stores/criteria";
  import { useCommunityStore } from "~/stores/community";
  import { useSecurityStore } from "~/stores/security";
  import { useDeployStore, type DeployStatus } from "~/stores/deploy";

  import type {
    AppModel,
    Pricing,
    ApprovalStatus,
    CriteriaKey,
    ExtraLabelKey,
  } from "~/types/types";
  import { useDeploymentsStore } from "~/stores/deployments";
  import { useDeployments } from "~/composables/API/useDeployments";
  import { useClusters } from "~/composables/API/useClusters";
  import { useConfigSecrets } from "~/composables/API/useConfigSecrets";

  /* Stores */
  const userStore = useUserStore();
  const catalog = useAppCatalog();
  const criteria = useCriteriaStore();
  const community = useCommunityStore();
  const security = useSecurityStore();
  const deploy = useDeployStore();

  const deploymentsStore = useDeploymentsStore();
  await deploymentsStore.loadClusterDeployments(); // Fetch clusters on store init
  await deploymentsStore.loadDeployments(); // Fetch deployments on store init
  await deploymentsStore.loadConfigMaps(); // Fetch config maps on store init
  await deploymentsStore.loadSecrets(); // Fetch secrets on store init

  const deploymentsApi = useDeployments();
  const configSecretsApi = useConfigSecrets();

  /* Hydrate */
  onMounted(async () => {
    await catalog.fetchAll().catch(() => { });
    await catalog.fetchUserApps().catch(() => { });
    uploadedImages.value = catalog.userApps;
    community.init();
    security.init();
  });

  /* IDs for sections */
  const ids = {
    kpis: "dash-kpis",
    hero: "dash-hero",
    assess: "dash-assess",
    catalog: "dash-catalog",
    security: "dash-security",
    community: "dash-community",
    deploy: "dash-deploy",
    updates: "dash-updates",
    uploaded: "dash-uploaded",
    configSecrets: "dash-config-secrets",
  };

  /* clusterConfigs */
  import type { ClusterConfig } from "~/types/types";
  const clusterConfigs = ref(<ClusterDeployment[]>deploymentsStore.clusterDeployments || []);
  const clusterConfigModalOpen = ref(false)
  const clusterConfig = ref<ClusterConfig>({} as ClusterConfig);
  const openClusterConfigModal = (config: ClusterConfig) => {
    clusterConfig.value = { ...config };
    clusterConfigModalOpen.value = true;
  };

 const setActiveClusterConfig = async (config: ClusterDeployment) => {
    deploymentsStore.setActiveClusterConfig(config);
 }

  const saveClusterConfig = async () => {
    if (!clusterConfig.value.name) {
      alert("Vul een naam in voor de clusterconfiguratie.");
      return;
    }

    try {
      await useClusters().createCluster(clusterConfig.value);
      clusterConfigModalOpen.value = false;
    } catch (error) {
      console.error("Fout bij het opslaan van clusterconfiguratie:", error);
    }
  };


  /* Basics */
  const statusLabel = (s:string)=>({submitted:'Ingediend',under_review:'In review',changes_requested:'Wijzigingen gevraagd'})[s]||s
  const statusTagReview = (s:string)=> s==='changes_requested'?'warning':'info'
  const today = new Date().toLocaleDateString();
  const orgLabel = computed(() => userStore.session?.organization || "Zonder organisatie");
  const initials = (name: string) =>
    name
      .split(" ")
      .map((x) => x[0])
      .slice(0, 2)
      .join("")
      .toUpperCase();

  /* KPI's */
  const appsCount = computed(() => catalog.apps.length);
  const advisoriesCount = computed(() => security.advisories.length);
  const threadsCount = computed(() => community.threads.length);
  const deploymentsCount = computed(() => {
    const list = Array.isArray((deploymentsStore as any).deployments)
      ? (deploymentsStore as any).deployments
      : Array.isArray(deploy.deployments)
        ? deploy.deployments
        : [];
    return list.length;
  });

  /* Toelatingseisen voortgang */
  function hasMeaningfulRecord(slug: string) {
    const rec = criteria.record(slug);
    const anyCheck = Object.values(rec.checks || {}).some(
      (s: any) => s === "ok" || s === "na"
    );
    const anyEvidence = Object.values(rec.evidence || {}).some(
      (e: any) => !!e?.url || !!e?.note
    );
    return anyCheck || anyEvidence;
  }
  const assessed = computed(() =>
    catalog.apps
      .filter((a) => criteria && hasMeaningfulRecord(a.slug))
      .map((a) => ({
        slug: a.slug,
        name: a.name,
        logo: a.media?.logoUrl || "/logo.svg",
        progress: criteria.progressPct(a.slug),
      }))
  );

  // ConfigMaps & Secrets
  const configMapModalOpen = ref(false);
  const configMapConfig = ref({
    name: "",
    data: "",
  });

  const secretModalOpen = ref(false);
  const secretConfig = ref({
    name: "",
    data: "",
  });
  
  const openSecretModal = () => {
    secretModalOpen.value = true;
  };
  const openConfigMapModal = () => {
    configMapModalOpen.value = true;
  };

  const createSecret = async () => {
    if (!secretConfig.value.name) {
      alert("Vul een naam in voor het secret.");
      return;
    }

    // Parse data from textarea (key=value per line)
    const dataLines = secretConfig.value.data.split("\n");
    const data: Record<string, string> = {};
    dataLines.forEach((line) => {
      const [key, value] = line.split("=");
      if (key && value) {
        data[key.trim()] = value.trim();
      }
    });

    try {
      await configSecretsApi.createSecret({
        namespace: deploymentsStore.userNamespace,
        name: secretConfig.value.name,
        data,
      });
      secretModalOpen.value = false;
      // Refresh secrets in store
      await deploymentsStore.loadSecrets();
    } catch (error) {
      console.error("Fout bij het aanmaken van secret:", error);
    }
  };

  const createConfigMap = async () => {
    if (!configMapConfig.value.name) {
      alert("Vul een naam in voor de config map.");
      return;
    }

    // Parse data from textarea (key=value per line)
    const dataLines = configMapConfig.value.data.split("\n");
    const data: Record<string, string> = {};
    dataLines.forEach((line) => {
      const [key, value] = line.split("=");
      if (key && value) {
        data[key.trim()] = value.trim();
      }
    });

    try {
      await configSecretsApi.createConfigMap({
        namespace: deploymentsStore.userNamespace,
        name: configMapConfig.value.name,
        data,
      });
      configMapModalOpen.value = false;
      // Refresh config maps in store
      await deploymentsStore.loadConfigMaps();
    } catch (error) {
      console.error("Fout bij het aanmaken van config map:", error);
    }
  };



  const uploadedImages = ref<AppModel[]>([]); // User-uploaded apps

  function goToApp(deploymentId: string) {
    // Get the image name from the deployment
    const deployment = (deploymentsStore as any).deployments.find(
      (d: any) => d.id === deploymentId
    );

    if (!deployment) {
      alert("Deployment niet gevonden.");
      return;
    }

    // Get corresponding app from catalog
    const app = catalog.apps.find((a) => a.imageName === deployment.config.containerImage);
    if (app) {
      // Navigate to app detail page
      window.location.href = `/apps/${app.slug}`;
    } else {
      alert("App niet gevonden in de catalogus.");
    }
  }

  const sevTag = (s: string) =>
    s === "Critical"
      ? "danger"
      : s === "High"
        ? "warning"
        : s === "Medium"
          ? "info"
          : "secondary";

  const catTag = (c: string) =>
    c === "Security"
      ? "danger"
      : c === "Privacy"
        ? "info"
        : c === "Techniek"
          ? "success"
          : "warning";

  const activityRows = computed(() => {
    const advisories = security.advisories.map((a) => ({
      id: a.id,
      type: "Advisory",
      icon: "pi pi-exclamation-triangle",
      title: a.title,
      label: a.severity,
      labelSeverity: sevTag(a.severity),
      extra: a.component,
      date: a.published,
      link: "/security", // pas aan wanneer je detailroute hebt (bv. `/security/advisory/${a.id}`)
    }));

    const threads = community.threads.map((t) => ({
      id: t.id,
      type: "Community",
      icon: "pi pi-comments",
      title: t.title,
      label: t.category,
      labelSeverity: catTag(t.category),
      extra: `Door ${t.author} • ${t.replies} reacties`,
      date: t.updatedAt,
      link: `/community/thread/${t.id}`,
    }));

    const all = [...advisories, ...threads];
    all.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
    return all.slice(0, 10);
  });

  /* Uploaded apps */
  const createModalOpen = ref(false);
  const appConfig = ref<AppModel>({
    slug: "",
    name: "",
    subtitle: "",
    summary: "",
    stage: "testfase",
    orgType: userStore.session?.org_type || "overheidsontwikkeld",
    categories: [],
    approved: 'submitted' as ApprovalStatus,
    media: {
      logoUrl: "",
      screenshots: [],
    },
    descriptionMd: "",
    pricing: "free" as unknown as Pricing,
    criteria: {
      open_source: { met: false },
      modular: { met: false },
      dpia: { met: false },
      license_clarity: { met: false },
      a11y: { met: false },
      ai_risk_cat: { met: false },
      functional_desc: { met: false },
      self_hostable: { met: false },
      tech_explainer: { met: false },
    },
    labels: {
      external_support: { has: true },
      human_rights_assessment: { has: true },
      gov_built: { has: true },
      user_guide: { has: true },
      open_inference_api: { has: false },
    },
    org: { name: userStore.session?.organization || "", url: userStore.session?.website || "" },
    versions: [],
    reviews: [],
    related: [],
    updatedAt: new Date().toISOString(),
    imageName: "",
  });

  const pricingOptions = [
    { label: "Gratis", value: { type: "free" } },
    {
      label: "Eenmalig",
      value: { type: "one_time", price: 0, currency: "EUR" },
    },
    {
      label: "Abonnement (per maand)",
      value: {
        type: "subscription",
        price: 0,
        currency: "EUR",
        interval: "month",
      },
    },
    {
      label: "Abonnement (per jaar)",
      value: {
        type: "subscription",
        price: 0,
        currency: "EUR",
        interval: "year",
      },
    },
    { label: "Gebruik (usage)", value: { type: "usage", unit: "" } },
  ];

  const createModalClose = () => {
    createModalOpen.value = false;
    appConfig.value = {
      slug: "",
      name: "",
      subtitle: "",
      summary: "",
      stage: "testfase",
      orgType: userStore.session?.org_type || "overheidsontwikkeld",
      categories: [],
      media: {
        logoUrl: "",
        screenshots: [],
      },
      approved: 'submitted' as ApprovalStatus,
      descriptionMd: "",
      pricing: "free" as unknown as Pricing,
      criteria: {
        open_source: { met: false },
        modular: { met: false },
        dpia: { met: false },
        license_clarity: { met: false },
        a11y: { met: false },
        ai_risk_cat: { met: false },
        functional_desc: { met: false },
        self_hostable: { met: false },
        tech_explainer: { met: false },
      },
      labels: {
        external_support: { has: true },
        human_rights_assessment: { has: true },
        gov_built: { has: true },
        user_guide: { has: true },
        open_inference_api: { has: false },
      },
      org: { name: "Rijksoverheid", url: "https://www.rijksoverheid.nl" },
      versions: [],
      reviews: [],
      related: [],
      updatedAt: new Date().toISOString(),
      imageName: "",
    } as AppModel;
  };

  function openCreateModal() {
    modalMode.value = "create";
    createModalOpen.value = true;
  }

  function editApplication(app: AppModel) {
    appConfig.value = { ...app };
    modalMode.value = "edit";
    createModalOpen.value = true;
  }

  async function deleteApplication(index: number) {
    if (index !== -1) {
      uploadedImages.value.splice(index, 1);
      await $fetch('/api/apps/userApps', {
        method: 'POST',
        body: uploadedImages.value,
      })
    }
    catalog.deleteApp()
  }

  async function createApplication() {
    if (modalMode.value === "edit") {
      // Remove old version
      const index = uploadedImages.value.findIndex((a) => a.slug === appConfig.value.slug);
      if (index !== -1) {
        uploadedImages.value.splice(index, 1);
      }

      appConfig.value.approved = 'submitted' as ApprovalStatus; // Reset approval status on edit
      appConfig.value.updatedAt = new Date().toISOString(); // Update timestamp

      // Add updated version
      uploadedImages.value.push(appConfig.value);
      await $fetch('/api/apps/userApps', {
        method: 'POST',
        body: uploadedImages.value,
      })
      catalog.fetchUserApps(); // Refresh user apps in store
      createModalClose();
      return;
    }

    // Generate a slug based on the name
    const baseSlug = appConfig.value.name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, "-")
      .replace(/^-+|-+$/g, "");
    let slug = baseSlug;

    let counter = 1;
    while (uploadedImages.value.some((app) => app.slug === slug)) {
      slug = `${baseSlug}-${counter}`;
      counter++;
    }
    appConfig.value.slug = slug;

    // Initialize versions array
    uploadedImages.value.push(appConfig.value);
    await $fetch('/api/apps/userApps', {
      method: 'POST',
      body: uploadedImages.value,
    })

    catalog.fetchUserApps(); // Refresh user apps in store
    createModalClose();
  }

  const modalMode = ref<"create" | "edit">("create");

  const configMaps = computed(() => deploymentsStore.configMaps || []);
  const secrets = computed(() => deploymentsStore.secrets || []);

  /* Deployments */
  const recentDeployments = computed(() => {
    const list = Array.isArray((deploymentsStore as any).deployments)
      ? (deploymentsStore as any).deployments
      : Array.isArray(deploymentsStore.deployments)
        ? deploymentsStore.deployments
        : [];
    const sorted = [...list].sort(
      (a: any, b: any) =>
        new Date(b.createdAt || 0).getTime() -
        new Date(a.createdAt || 0).getTime()
    );
    return sorted.slice(0, 10);
  });
  const shortId = (id: string) => id.replace(/^dep_/, "").slice(0, 8);
  const statusTag = (s: DeployStatus) => {
    if (s === "running") return "success";
    if (s === "deploying" || s === "queued") return "info";
    if (s === "failed") return "danger";
    if (s === "canceled" || s === "decommissioned") return "secondary";
    return "secondary";
  };
  async function quickRedeploy(id: string) {
    await deploy.redeployFrom(id);
  }
  async function cancelDeployment(id: string) {
    await deploymentsApi.deleteDeployment(deploymentsStore.userNamespace, id);
    await deploymentsStore.loadDeployments();
  }

  /* Misc */
  function prettyDate(d?: string | Date) {
    return d ? new Date(d).toLocaleDateString() : "—";
  }
</script>

<style scoped>
  .container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .no-gutter {
    margin: 0;
  }

  .no-gutter>[class^="col-"] {
    padding: 0;
  }

  .kpi .kpi-num {
    font-size: 1.75rem;
    font-weight: 700;
    line-height: 1.1;
  }

  /* Sidebar visual tweaks */
  :deep(.p-panelmenu .p-panelmenu-header .p-panelmenu-header-content) {
    padding: 0.75rem 1rem;
  }

  :deep(.p-panelmenu .p-menuitem-text) {
    font-weight: 500;
  }

  :deep(.p-panelmenu .p-menuitem-link) {
    padding: 0.5rem 0.75rem;
  }

  pre {
  background: #0b1020;
  color: #e7eaf6;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
  font-family: ui-monospace, monospace;
  font-size: 0.9rem;
}
</style>
