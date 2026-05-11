<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('dashboard.title') }}</h1>
        <p class="page-subtitle">{{ t('dashboard.desc') }}</p>
      </div>
      <a-button @click="refresh">{{ t('dashboard.refresh') }}</a-button>
    </div>
    <div class="stat-grid">
      <div class="stat-card" v-for="item in stats" :key="item.label">
        <div class="stat-label">{{ item.label }}</div>
        <div class="stat-value">{{ item.value }}</div>
      </div>
    </div>
    <div class="panel">
      <a-timeline>
        <a-timeline-item color="green">{{ t('dashboard.authReady') }}</a-timeline-item>
        <a-timeline-item color="green">{{ t('dashboard.systemReady') }}</a-timeline-item>
        <a-timeline-item color="green">{{ t('dashboard.fileReady') }}</a-timeline-item>
        <a-timeline-item>{{ t('dashboard.auditReady') }}</a-timeline-item>
      </a-timeline>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../../stores/auth'

const { t } = useI18n()
const auth = useAuthStore()
const stats = computed(() => [
  { label: t('dashboard.menuCount'), value: countMenus(auth.menus) },
  { label: t('dashboard.permissionCount'), value: auth.permissions.length },
  { label: t('dashboard.roleCount'), value: auth.profile?.roles?.length || 0 },
  { label: t('dashboard.currentUser'), value: auth.profile?.username || '-' }
])

function countMenus(items) {
  return items.reduce((sum, item) => sum + 1 + countMenus(item.children || []), 0)
}

function refresh() {
  auth.bootstrap()
}
</script>
