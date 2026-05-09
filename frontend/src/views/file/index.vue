<template>
  <main class="page">
    <PageHeader :title="t('page.file.title')" :description="t('page.file.description')">
      <template #actions>
        <a-upload :show-upload-list="false" :custom-request="customUpload">
          <a-button type="primary" v-permission="'file:upload'">{{ t('common.upload') }}</a-button>
        </a-upload>
      </template>
    </PageHeader>

    <SearchForm :model="query" :title="t('page.file.searchTitle')">
      <a-form-item :label="t('field.keyword')">
        <a-input v-model:value="query.keyword" :placeholder="t('page.file.keywordPlaceholder')" allow-clear />
      </a-form-item>
      <template #actions>
        <a-button type="primary" @click="search">{{ t('common.query') }}</a-button>
        <a-button @click="reset">{{ t('common.reset') }}</a-button>
      </template>
    </SearchForm>

    <DataTable
      row-key="id"
      :title="t('page.file.tableTitle')"
      :columns="columns"
      :data-source="rows"
      :pagination="pagination"
      :loading="loading"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'size'">
          {{ formatSize(record.size) }}
        </template>
        <template v-else-if="column.key === 'actions'">
          <a-space>
            <a-button type="link" size="small" @click="downloadFile(record)">{{ t('common.download') }}</a-button>
            <ConfirmModal :title="t('page.file.deleteConfirm')" @confirm="removeFile(record)">
              <a-button type="link" size="small" danger v-permission="'file:delete'">{{ t('common.delete') }}</a-button>
            </ConfirmModal>
          </a-space>
        </template>
      </template>
    </DataTable>
  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import ConfirmModal from '@/components/common/ConfirmModal.vue'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import SearchForm from '@/components/common/SearchForm.vue'
import { deleteFile, fileDownloadUrl, fileList, uploadFile } from '@/api/file'
import { formatDateTime } from '@/utils/datetime'

const { t } = useI18n()
const loading = ref(false)
const rows = ref([])
const query = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const columns = computed(() => [
  { title: t('field.originalName'), dataIndex: 'originalName', width: 260 },
  { title: t('field.contentType'), dataIndex: 'contentType', width: 180 },
  { title: t('field.size'), key: 'size', width: 120 },
  { title: t('field.createdAt'), dataIndex: 'createdAt', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 140 }
])

async function customUpload({ file, onSuccess, onError }) {
  try {
    await uploadFile(file)
    message.success(t('page.file.uploadSuccess'))
    await load()
    onSuccess()
  } catch (error) {
    onError(error)
  }
}

async function load() {
  loading.value = true
  try {
    const data = await fileList({
      page: pagination.current,
      pageSize: pagination.pageSize,
      keyword: query.keyword
    })
    rows.value = (data.list || []).map((item) => ({
      ...item,
      createdAt: formatDateTime(item.createdAt)
    }))
    pagination.total = data.total
  } finally {
    loading.value = false
  }
}

function search() {
  pagination.current = 1
  load()
}

function reset() {
  query.keyword = ''
  search()
}

function onTableChange(p) {
  pagination.current = p.current
  pagination.pageSize = p.pageSize
  load()
}

async function downloadFile(record) {
  const data = await fileDownloadUrl(record.id)
  const link = document.createElement('a')
  link.href = data.url
  link.download = record.originalName || 'download'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

async function removeFile(record) {
  await deleteFile(record.id)
  message.success(t('message.deleteSuccess'))
  if (rows.value.length === 1 && pagination.current > 1) {
    pagination.current -= 1
  }
  await load()
}

function formatSize(size) {
  const value = Number(size || 0)
  if (value < 1024) return `${value} B`
  if (value < 1024 * 1024) return `${(value / 1024).toFixed(1)} KB`
  if (value < 1024 * 1024 * 1024) return `${(value / 1024 / 1024).toFixed(1)} MB`
  return `${(value / 1024 / 1024 / 1024).toFixed(1)} GB`
}

onMounted(load)
</script>
