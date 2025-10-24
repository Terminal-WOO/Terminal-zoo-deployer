
<template>
  <div class="surface-0 min-h-screen flex flex-column">
    <AppHeader :sticky="true" :elevated="true" :showSearch="true" />
    <main class="flex-1">
      <!-- Hero-secties kunnen full-width zijn; gebruik binnenpagina containers -->
      <slot />
    </main>
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.min-h-screen { min-height: 100vh; }
</style>