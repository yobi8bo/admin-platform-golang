<template>
  <main class="page">
    <PageHeader :title="t('page.dept.title')" :description="t('page.dept.description')">
      <template #actions>
        <a-button type="primary" @click="openCreate()">{{ t('page.dept.create') }}</a-button>
      </template>
    </PageHeader>

    <DataTable
      row-key="id"
      :title="t('page.dept.tableTitle')"
      :columns="columns"
      :data-source="rows"
      :pagination="false"
      :loading="loading"
      :children-column-name="'children'"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <StatusTag :status="record.status" />
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="openCreate(record)">{{ t('common.addChild') }}</a-button>
            <a-button size="small" @click="openEdit(record)">{{ t('common.edit') }}</a-button>
            <ConfirmModal :title="t('page.dept.deleteConfirm')" @confirm="remove(record.id)">
              <a-button size="small" danger>{{ t('common.delete') }}</a-button>
            </ConfirmModal>
          </a-space>
        </template>
      </template>
    </DataTable>

    <a-modal v-model:open="modalOpen" :title="editingId ? t('page.dept.edit') : t('page.dept.create')" @ok="save">
      <a-form layout="vertical" :model="form">
        <a-form-item :label="t('field.parentDept')">
          <a-tree-select
            v-model:value="form.parentId"
            allow-clear
            tree-default-expand-all
            :tree-data="parentOptions"
            :field-names="{ label: 'name', value: 'id', children: 'children' }"
          />
        </a-form-item>
        <a-form-item :label="t('field.name')" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item :label="t('field.status')">
          <a-select v-model:value="form.status" :options="statusOptions" />
        </a-form-item>
        <a-form-item :label="t('field.sort')">
          <a-input-number v-model:value="form.sort" :min="0" />
        </a-form-item>
      </a-form>
    </a-modal>
  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import ConfirmModal from '@/components/common/ConfirmModal.vue'
import DataTable from '@/components/common/DataTable.vue'
import PageHeader from '@/components/common/PageHeader.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { deptApi } from '@/api/system'

const { t } = useI18n()
const loading = ref(false)
const rows = ref([])
const modalOpen = ref(false)
const editingId = ref(null)
const form = reactive({ parentId: 0, name: '', status: 'enabled', sort: 0 })
const statusOptions = computed(() => [
  { label: t('status.enabled'), value: 'enabled' },
  { label: t('status.disabled'), value: 'disabled' }
])
const columns = computed(() => [
  { title: t('field.name'), dataIndex: 'name' },
  { title: t('field.status'), key: 'status', width: 120 },
  { title: t('field.sort'), dataIndex: 'sort', width: 100 },
  { title: t('field.action'), key: 'actions', width: 240, fixed: 'right' }
])
const parentOptions = computed(() => [{ id: 0, name: t('page.dept.title'), children: rows.value }])

async function load() {
  loading.value = true
  try {
    rows.value = await deptApi.tree()
  } finally {
    loading.value = false
  }
}

function resetForm(parent) {
  Object.assign(form, { parentId: parent?.id || 0, name: '', status: 'enabled', sort: 0 })
}

function openCreate(parent) {
  editingId.value = null
  resetForm(parent)
  modalOpen.value = true
}

function openEdit(record) {
  editingId.value = record.id
  Object.assign(form, { parentId: record.parentId, name: record.name, status: record.status, sort: record.sort })
  modalOpen.value = true
}

async function save() {
  if (editingId.value) {
    await deptApi.update(editingId.value, form)
  } else {
    await deptApi.create(form)
  }
  message.success(t('message.saveSuccess'))
  modalOpen.value = false
  load()
}

async function remove(id) {
  await deptApi.remove(id)
  message.success(t('message.deleteSuccess'))
  load()
}

onMounted(load)
</script>
