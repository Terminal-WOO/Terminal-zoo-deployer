<template>
  <footer class="surface-0 border-top-1 surface-border">
    <div class="container p-4">
      <div class="grid">
        <div class="col-12 md:col-6">
          <div class="flex align-items-center gap-2 mb-2">
            <BrandLogo size="md" class="mb-3" />
          </div>
          <p class="m-0 text-700">
            Transparante AI-oplossingen met betrouwbare criteria, community-monitoring en bouwstenen
            voor snelle en veilige implementaties.
          </p>
          <p class="m-0 text-700 mt-2">
            Deze website is ontwikkeld door Clappform B.V.
          </p>
        </div>

        <div class="col-6 md:col-3">
          <h4 class="mt-0">Product</h4>
          <ul class="list-none p-0 m-0">
            <li><NuxtLink to="/apps" class="footer-link">App-store</NuxtLink></li>
            <li><NuxtLink to="/criteria" class="footer-link">Toelatingseisen</NuxtLink></li>
            <li><NuxtLink to="/community" class="footer-link">Community</NuxtLink></li>
            <li><NuxtLink to="/deploy/new" class="footer-link">Nieuwe deploy</NuxtLink></li>
          </ul>
        </div>

        <div class="col-6 md:col-3">
          <h4 class="mt-0">Juridisch</h4>
          <ul class="list-none p-0 m-0">
            <li><NuxtLink to="/legal/terms" class="footer-link">Voorwaarden</NuxtLink></li>
            <li><NuxtLink to="/legal/privacy" class="footer-link">Privacy</NuxtLink></li>
            <li><NuxtLink to="/legal/cookies" class="footer-link">Cookies</NuxtLink></li>
          </ul>
        </div>
      </div>

      <Divider class="mt-4" />

      <div class="text-600 text-sm">© {{ new Date().getFullYear() }} {{ brand }}. Alle rechten voorbehouden.</div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import BrandLogo from '~/components/BrandLogo.vue'
const props = withDefaults(defineProps<{ brand?: string }>(), { brand: 'Clappform B.V' })
const brand = props.brand
</script>

<style scoped>
.container { max-width: 1100px; margin: 0 auto; }
.text-xl { font-size: 1.25rem; }
.footer-link {
  display: inline-block; padding: .25rem 0; text-decoration: none; color: var(--text-color);
}
.footer-link:hover { color: var(--primary-color); }
</style>
