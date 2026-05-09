<template>
  <a-card class="metric-card" :bordered="false">
    <div class="metric-card__header">
      <span>{{ title }}</span>
      <slot name="icon" />
    </div>
    <a-statistic :value="value" :suffix="suffix" />
    <p v-if="description">{{ description }}</p>
    <a-progress v-if="progress !== null" :percent="progress" :show-info="false" size="small" />
  </a-card>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  suffix: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  },
  progress: {
    type: Number,
    default: null
  }
})
</script>
