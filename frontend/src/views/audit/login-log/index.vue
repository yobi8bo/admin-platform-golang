<template>
  <main class="page">
    <PageHeader :title="t('page.loginLog.title')" :description="t('page.loginLog.description')" />
    <DataTable
      row-key="id"
      :title="t('page.loginLog.tableTitle')"
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
import { loginLogs } from '@/api/audit'
import { formatDateTime } from '@/utils/datetime'

const { t } = useI18n()
const rows = ref([])
const loading = ref(false)
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = computed(() => [
  { title: t('field.username'), dataIndex: 'username', width: 160 },
  { title: t('field.ip'), dataIndex: 'ip', width: 150 },
  { title: t('field.status'), key: 'status', width: 120 },
  { title: t('field.message'), dataIndex: 'message' },
  { title: t('field.createdAt'), dataIndex: 'createdAt', width: 180 }
])

async function load() {
  loading.value = true
  try {
    const data = await loginLogs({ page: pagination.current, pageSize: pagination.pageSize })
    rows.value = (data.list || []).map((item) => ({
      ...item,
      message: translateLoginMessage(item.message),
      createdAt: formatDateTime(item.createdAt)
    }))
    pagination.total = data.total
  } finally {
    loading.value = false
  }
}

function translateLoginMessage(message) {
  const messages = {
    'login success': '登录成功',
    'invalid username or password': '用户名或密码错误',
    'invalid request': '无效请求',
    'token issue failed': '令牌签发失败'
  }
  return messages[message] || message
}

function onTableChange(p) {
  pagination.current = p.current
  pagination.pageSize = p.pageSize
  load()
}

onMounted(load)
</script>
