<template>
  <a-popconfirm
    :title="title"
    :ok-text="okText || t('common.confirm')"
    :cancel-text="cancelText || t('common.cancel')"
    @confirm="emit('confirm')"
  >
    <slot />
  </a-popconfirm>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

defineProps({
  title: {
    type: String,
    required: true
  },
  okText: {
    type: String,
    default: ''
  },
  cancelText: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['confirm'])
const { t } = useI18n()
</script>
