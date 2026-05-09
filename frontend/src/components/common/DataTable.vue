<template>
  <section class="data-table" :class="rootClass" :style="rootStyle">
    <div v-if="title || $slots.actions" class="data-table__toolbar">
      <div>
        <h3 v-if="title">{{ title }}</h3>
        <p v-if="description">{{ description }}</p>
      </div>
      <div v-if="$slots.actions" class="data-table__actions">
        <slot name="actions" />
      </div>
    </div>
    <a-table
      v-bind="tableAttrs"
      :columns="columns"
      :data-source="dataSource"
      :loading="loading"
      :pagination="pagination"
      :scroll="scroll"
    >
      <template #bodyCell="slotProps">
        <slot name="bodyCell" v-bind="slotProps" />
      </template>
      <template #emptyText>
        <a-empty :description="emptyText" />
      </template>
    </a-table>
  </section>
</template>

<script setup>
import { computed, useAttrs } from 'vue'

defineOptions({ inheritAttrs: false })

const attrs = useAttrs()

const rootClass = computed(() => attrs.class)
const rootStyle = computed(() => attrs.style)
const tableAttrs = computed(() => {
  const { class: _class, style: _style, ...rest } = attrs
  return rest
})

defineProps({
  title: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  },
  columns: {
    type: Array,
    required: true
  },
  dataSource: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  pagination: {
    type: [Object, Boolean],
    default: false
  },
  scroll: {
    type: Object,
    default: () => ({ x: 'max-content' })
  },
  emptyText: {
    type: String,
    default: ''
  }
})
</script>
