
<template>
  <div class="surface-0 min-h-screen flex flex-column">
    <AppHeader :sticky="true" :elevated="true" :showSearch="true" />
    <main class="flex-1">
      <div class="container p-4">
        <slot />
      </div>
    </main>
    <AppFooter />
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.container { max-width: 1100px; margin: 0 auto; }
.min-h-screen { min-height: 100vh; }
</style>