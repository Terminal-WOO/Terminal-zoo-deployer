import { defineStore } from 'pinia'
import { useDeployments } from '~/composables/API/useDeployments';
import { useClusters } from '~/composables/API/useClusters';
import { useConfigSecrets } from '~/composables/API/useConfigSecrets';
import type { ManualDeploymentConfig } from '~/types/types'

// export type DeployStatus = 'draft'|'queued'|'deploying'|'running'|'failed'|'canceled'|'decommissioned'

export interface Deployment {
  id: string;
  creationTimestamp: string;
  status: string;
  config: ManualDeploymentConfig;
}

export interface ClusterDeployment {
  name: string;
  server: string;
  k8sVersion: string;
  defaultNamespace: string;
}

export const useDeploymentsStore = defineStore('deployments', {
  state: () => ({
    activeClusterConfig: null as ClusterDeployment | null,
    usercluster: 'demo',
    userNamespace: 'testing',
    deployments: [] as Deployment[],
    clusterDeployments: [] as ClusterDeployment[],
    configMaps: [] as any[],
    secrets: [] as any[],
  }),
  actions: {
    async loadDeployments() {
      try {
        const deployments: any[] = [];
        console.log('clusters', this.clusterDeployments);
        // fetch default cluster deployments
        const clusterDeployments = await useDeployments().getDeployments(this.userNamespace);
        deployments.push(clusterDeployments);

        // fetch deployments from all configured clusters
        for (const cluster of this.clusterDeployments) {
          console.log(`Loading deployments for cluster: ${cluster.name}`);
          const clusterDeployments = await useDeployments().getDeployments(this.userNamespace, cluster.name);
          // append to deployments array
          console.log(`Deployments from cluster ${cluster.name}:`, clusterDeployments);
          deployments.push(clusterDeployments);
        }
        console.log('Fetched deployments:', deployments);
        this.deployments = [];
        for (const dep of deployments) {
          for (const item of dep.items) {
            if (item.metadata.uid in this.deployments.map(d => d.id)) {
              console.log(`Skipping duplicate deployment with UID: ${item.metadata.uid}`);
              continue; // skip duplicates
            }
            this.deployments.push(
              {
                id: item.metadata.uid,
                creationTimestamp: item.metadata.creationTimestamp,
                status: item.status?.conditions?.find((c: any) => c.type === 'Available')?.status === 'True' ? 'running' : 'pending',
                config: {
                  deploymentName: item.metadata.name,
                  namespace: item.metadata.namespace,
                  containerImage: item.spec.template.spec.containers[0].image,
                  replicas: item.spec.replicas,
                  resources: {
                    cpuRequests: item.spec.template.spec.containers[0].resources?.requests?.cpu || "unknown",
                    cpuLimits: item.spec.template.spec.containers[0].resources?.limits?.cpu || "unknown",
                    memoryRequests: item.spec.template.spec.containers[0].resources?.requests?.memory || "unknown",
                    memoryLimits: item.spec.template.spec.containers[0].resources?.limits?.memory || "unknown",
                  },
                  ports: item.spec.template.spec.containers[0].ports.map((p: any) => ({
                    ContainerPort: p.containerPort,
                    Protocol: p.protocol
                  })),
                }
              }
            );
          }
        }

        console.log('Processed deployments:', this.deployments);
      } catch (error) {
        console.error('Failed to load deployments:', error);
      }
    },
    async loadConfigMaps() {
      try {
        const configMaps = await useConfigSecrets().readConfigMaps(this.userNamespace);
        this.configMaps = configMaps || [];
      } catch (error) {
        console.error('Failed to load config maps:', error);
      }
    },
    async loadSecrets() {
      try {
        const secrets = await useConfigSecrets().readSecrets(this.userNamespace);
        this.secrets = secrets || [];
      } catch (error) {
        console.error('Failed to load secrets:', error);
      }
    },
    currentActiveClusterConfig() {
      return this.activeClusterConfig;
    },
    setActiveClusterConfig(config: ClusterDeployment) {
      this.activeClusterConfig = config;
    },
    async loadClusterDeployments() {
      try {
        const clusters = await useClusters().getClusters();
        console.log('Fetched clusters:', clusters);

        this.clusterDeployments = clusters?.map((cl: any) => ({
          name: cl.name,
          server: cl.server,
          k8sVersion: cl.k8sVersion || 'unknown',
          defaultNamespace: cl?.defaultNamespace || 'default',
        })) || [];
      } catch (error) {
        console.error('Failed to load cluster deployments:', error);
      }
    }
  },
});