

<template>
  <div class="surface-0 min-h-screen flex flex-column">
    <AppHeader :sticky="true" :elevated="true" :showSearch="false" />

    <div class="container p-0 flex-1">
      <div class="grid m-0">
        <!-- Content -->
        <section class="">
          <slot />
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.container {  margin: 0 auto; }
.min-h-screen { min-height: 100vh; }
@media (min-width: 992px) {
  aside { position: sticky; top: 64px; /* onder de header */ }
}
</style>
