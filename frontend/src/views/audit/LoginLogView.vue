<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('audit.loginLog') }}</h1>
        <p class="page-subtitle">{{ t('audit.loginLogDesc') }}</p>
      </div>
    </div>
    <div class="panel">
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="pagination" @change="onPageChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'success' ? 'green' : 'red'">{{ statusLabel(record.status) }}</a-tag>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { getLoginLogs } from '../../api/audit'

const { t } = useI18n()
const loading = ref(false)
const rows = ref([])
const query = reactive({ page: 1, pageSize: 10, total: 0 })
const columns = computed(() => [
  { title: t('common.id'), dataIndex: 'id', width: 80 },
  { title: t('audit.username'), dataIndex: 'username' },
  { title: 'IP', dataIndex: 'ip' },
  { title: t('common.status'), key: 'status', width: 120 },
  { title: t('common.message'), dataIndex: 'message' },
  { title: t('common.createdAt'), dataIndex: 'createdAt', width: 180 }
])
const pagination = computed(() => ({ current: query.page, pageSize: query.pageSize, total: query.total, showSizeChanger: true }))

onMounted(load)

async function load() {
  loading.value = true
  try {
    const data = await getLoginLogs(query)
    rows.value = data.list
    query.total = data.total
  } finally {
    loading.value = false
  }
}

function onPageChange(page) {
  query.page = page.current
  query.pageSize = page.pageSize
  load()
}

function statusLabel(status) {
  if (status === 'success') return t('common.success')
  if (status === 'failed') return t('common.failed')
  return status || '-'
}
</script>
