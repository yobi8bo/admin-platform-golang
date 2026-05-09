<template>
  <main class="page">
    <PageHeader :title="t('page.operationLog.title')" :description="t('page.operationLog.description')" />
    <DataTable
      row-key="id"
      :title="t('page.operationLog.tableTitle')"
      :columns="columns"
      :data-source="rows"
      :pagination="pagination"
      :loading="loading"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <StatusTag :status="record.status" />
        </template>
      </template>
    </DataTable>
  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { operationLogs } from '@/api/audit'
import { formatDateTime } from '@/utils/datetime'

const { t } = useI18n()
const rows = ref([])
const loading = ref(false)
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = computed(() => [
  { title: t('field.username'), dataIndex: 'username', width: 160 },
  { title: t('field.createdAt'), dataIndex: 'createdAt', width: 180 },
  { title: t('field.operation'), dataIndex: 'operation' },
  { title: t('field.status'), key: 'status', width: 120 }
])

async function load() {
  loading.value = true
  try {
    const data = await operationLogs({ page: pagination.current, pageSize: pagination.pageSize })
    rows.value = (data.list || []).map((item) => ({
      ...item,
      createdAt: formatDateTime(item.createdAt)
    }))
    pagination.total = data.total
  } finally {
    loading.value = false
  }
}

function onTableChange(p) {
  pagination.current = p.current
  pagination.pageSize = p.pageSize
  load()
}

onMounted(load)
</script>
