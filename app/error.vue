

<template>
  <div class="surface-0 min-h-screen flex align-items-center justify-content-center p-3">
    <Card class="w-full md:w-6">
      <template #title>Er ging iets mis</template>
      <template #content>
        <p class="m-0 text-700">
          {{ (error?.statusCode || 'Onbekende fout') + ': ' + (error?.message || 'Er is een fout opgetreden.') }}
        </p>
      </template>
      <template #footer>
        <Button label="Terug naar start" icon="pi pi-home" @click="handleError" />
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
const error = useError()
const handleError = () => clearError({ redirect: '/' })
</script>

<style scoped>
.min-h-screen { min-height: 100vh; }
</style>