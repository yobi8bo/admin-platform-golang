<template>
  <a-tag :color="config.color" class="status-tag">
    <a-badge :status="config.badge" />
    {{ config.label }}
  </a-tag>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  status: {
    type: [String, Number, Boolean],
    default: ''
  }
})

const { t, te } = useI18n()

const statusMap = {
  enabled: { color: 'success', badge: 'success' },
  disabled: { color: 'default', badge: 'default' },
  success: { color: 'success', badge: 'success' },
  failed: { color: 'error', badge: 'error' },
  error: { color: 'error', badge: 'error' },
  warning: { color: 'warning', badge: 'warning' },
  pending: { color: 'processing', badge: 'processing' },
  processing: { color: 'processing', badge: 'processing' },
  200: { color: 'success', badge: 'success' }
}

const config = computed(() => {
  const key = String(props.status)
  const base = statusMap[key] || { color: 'default', badge: 'default' }
  const labelKey = `status.${key}`
  return {
    ...base,
    label: te(labelKey) ? t(labelKey) : key || t('status.unknown')
  }
})
</script>
