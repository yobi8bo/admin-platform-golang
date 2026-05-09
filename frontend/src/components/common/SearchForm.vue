<template>
  <section class="search-form" :aria-label="title">
    <div v-if="title || $slots.actions" class="search-form__header">
      <strong v-if="title">{{ title }}</strong>
      <div v-if="$slots.actions" class="search-form__actions">
        <slot name="actions" />
      </div>
    </div>
    <a-form class="search-form__body" layout="inline" :model="model">
      <slot />
    </a-form>
  </section>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: ''
  },
  model: {
    type: Object,
    default: () => ({})
  }
})
</script>
