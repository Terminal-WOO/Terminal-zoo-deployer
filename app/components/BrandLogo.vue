<template>
  <span class="brand" :class="sizeClass" aria-label="Logo">
    <img :src="resolvedLogo" alt="" aria-hidden="true" />
  </span>
</template>

<script setup lang="ts">
// Default logo comes from /assets (build-time import)
import defaultLogo from '~/assets/images/logo.svg'

const props = withDefaults(defineProps<{
  size?: 'xs'|'sm'|'md'|'lg'|'xl'
  variant?: 'default'
}>(), { size: 'md', variant: 'default' })

const resolvedLogo = computed(() => defaultLogo)
const sizeClass = computed(() => `brand--${props.size}`)
</script>

<style scoped>
.brand { display: inline-flex; align-items: center; }
.brand img { display: block; width: 100%; height: auto; }
.brand--xs { width: 60px; }
.brand--sm { width: 100px; }
.brand--md { width: 120px; }
.brand--lg { width: 140px; }
.brand--xl { width: 160px; }
</style>