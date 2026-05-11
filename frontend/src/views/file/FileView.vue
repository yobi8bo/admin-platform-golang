<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('file.title') }}</h1>
        <p class="page-subtitle">{{ t('file.desc') }}</p>
      </div>
      <a-upload v-permission="'file:upload'" :show-upload-list="false" :custom-request="customUpload">
        <a-button type="primary" :loading="uploading">{{ t('file.upload') }}</a-button>
      </a-upload>
    </div>
    <div class="panel">
      <div class="toolbar">
        <a-input-search v-model:value="query.keyword" :placeholder="t('file.keyword')" allow-clear @search="load" />
      </div>
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="pagination" @change="onPageChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'size'">{{ formatSize(record.size) }}</template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a @click="download(record.id)">{{ t('common.download') }}</a>
              <a-popconfirm :title="t('file.deleteConfirm')" @confirm="remove(record.id)">
                <a v-permission="'file:delete'" class="danger">{{ t('common.delete') }}</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import { deleteFile, getDownloadUrl, getFiles, uploadFile } from '../../api/file'

const { t } = useI18n()
const loading = ref(false)
const uploading = ref(false)
const rows = ref([])
const query = reactive({ keyword: '', page: 1, pageSize: 10, total: 0 })
const columns = computed(() => [
  { title: t('file.name'), dataIndex: 'originalName' },
  { title: t('file.type'), dataIndex: 'contentType', width: 180 },
  { title: t('file.size'), key: 'size', width: 110 },
  { title: t('common.createdAt'), dataIndex: 'createdAt', width: 180 },
  { title: t('common.actions'), key: 'actions', width: 120 }
])
const pagination = computed(() => ({ current: query.page, pageSize: query.pageSize, total: query.total, showSizeChanger: true }))

onMounted(load)

async function load() {
  loading.value = true
  try {
    const data = await getFiles(query)
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

async function customUpload({ file, onSuccess, onError }) {
  uploading.value = true
  const form = new FormData()
  form.append('file', file)
  try {
    await uploadFile(form)
    message.success(t('file.uploadSuccess'))
    onSuccess()
    load()
  } catch (error) {
    onError(error)
  } finally {
    uploading.value = false
  }
}

async function download(id) {
  const data = await getDownloadUrl(id)
  window.open(data.url, '_blank', 'noopener')
}

async function remove(id) {
  await deleteFile(id)
  message.success(t('file.deleted'))
  load()
}

function formatSize(size = 0) {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / 1024 / 1024).toFixed(1)} MB`
}
</script>

<style scoped>
.danger {
  color: var(--danger);
}
</style>
