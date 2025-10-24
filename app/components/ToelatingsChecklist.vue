

<template>
  <Panel header="Toelatingseisen">
    <div class="grid">
      <div class="col-12 md:col-6" v-for="item in [
        {k:'open_source', lbl:'Open-source app'},
        {k:'modular', lbl:'Modulaire opzet'},
        {k:'dpia', lbl:'DPIA goedgekeurd'},
        {k:'license_clarity', lbl:'Licentiestructuur inzichtelijk'},
        {k:'a11y', lbl:'Digitale toegankelijkheid'},
        {k:'ai_risk_cat', lbl:'AI risicocategorie bepaald'},
        {k:'functional_desc', lbl:'Functionele beschrijving'},
        {k:'self_hostable', lbl:'Lokaal draaien mogelijk'},
        {k:'tech_explainer', lbl:'Technische uitleg beschikbaar'},
      ]" :key="item.k">
        <div class="flex items-center gap-2">
          <Checkbox :binary="true" :modelValue="modelValue?.[item.k]?.met" @update:modelValue="toggle(item.k)" />
          <span>{{ item.lbl }}</span>
        </div>
      </div>
    </div>
  </Panel>

  <Panel header="Extra certificaten en labels" class="mt-3">
    <div class="grid">
      <div class="col-12 md:col-6" v-for="item in [
        {k:'external_support', lbl:'Externe ondersteuning (SLA)'},
        {k:'human_rights_assessment', lbl:'FRIA/IAMA mensenrechtentoets'},
        {k:'gov_built', lbl:'Ontwikkeld door de overheid'},
        {k:'user_guide', lbl:'Gebruikershandleiding'},
        {k:'open_inference_api', lbl:'Open-source Inference API'},
      ]" :key="item.k">
        <div class="flex items-center gap-2">
          <InputSwitch :modelValue="modelValue?.[item.k]?.has" @update:modelValue="v => emit('update:modelValue',{...modelValue,[item.k]:{...(modelValue[item.k]||{}),has:v}})" />
          <span>{{ item.lbl }}</span>
        </div>
      </div>
    </div>
  </Panel>
</template>

<script setup lang="ts">
defineProps<{ modelValue: any }>()
const emit = defineEmits(['update:modelValue'])
const toggle = (k:string) => emit('update:modelValue', { ...$props.modelValue, [k]: { ...( $props.modelValue[k]||{} ), met: !($props.modelValue[k]?.met) }})
</script>