<template>
  <main class="page">
    <PageHeader :title="t('page.menu.title')" :description="t('page.menu.description')">
      <template #actions>
        <a-button type="primary" @click="openCreate()">{{ t('page.menu.create') }}</a-button>
      </template>
    </PageHeader>

    <DataTable
      row-key="id"
      :title="t('page.menu.tableTitle')"
      :columns="columns"
      :data-source="rows"
      :pagination="false"
      :loading="loading"
      :children-column-name="'children'"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'type'">
          <StatusTag :status="record.type" />
        </template>
        <template v-if="column.key === 'hidden'">
          <StatusTag :status="String(record.hidden)" />
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" @click="openCreate(record)">{{ t('common.addChild') }}</a-button>
            <a-button size="small" @click="openEdit(record)">{{ t('common.edit') }}</a-button>
            <ConfirmModal :title="t('page.menu.deleteConfirm')" @confirm="remove(record.id)">
              <a-button size="small" danger>{{ t('common.delete') }}</a-button>
            </ConfirmModal>
          </a-space>
        </template>
      </template>
    </DataTable>

    <a-modal v-model:open="modalOpen" :title="editingId ? t('page.menu.edit') : t('page.menu.create')" @ok="save">
      <a-form layout="vertical" :model="form">
        <a-form-item :label="t('field.parentMenu')">
          <a-tree-select
            v-model:value="form.parentId"
            allow-clear
            tree-default-expand-all
            :tree-data="parentOptions"
            :field-names="{ label: 'title', value: 'id', children: 'children' }"
          />
        </a-form-item>
        <a-form-item :label="t('field.type')" required>
          <a-radio-group v-model:value="form.type">
            <a-radio-button value="catalog">{{ t('status.catalog') }}</a-radio-button>
            <a-radio-button value="menu">{{ t('status.menu') }}</a-radio-button>
            <a-radio-button value="button">{{ t('status.button') }}</a-radio-button>
          </a-radio-group>
        </a-form-item>
        <a-form-item :label="t('field.name')" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item :label="t('field.title')" required>
          <a-input v-model:value="form.title" />
        </a-form-item>
        <a-form-item :label="t('field.route')">
          <a-input v-model:value="form.path" />
        </a-form-item>
        <a-form-item :label="t('field.component')">
          <a-input v-model:value="form.component" />
        </a-form-item>
        <a-form-item :label="t('field.icon')">
          <a-input v-model:value="form.icon" />
        </a-form-item>
        <a-form-item :label="t('field.permission')">
          <a-input v-model:value="form.permission" />
        </a-form-item>
        <a-form-item :label="t('field.sort')">
          <a-input-number v-model:value="form.sort" :min="0" />
        </a-form-item>
        <a-form-item :label="t('field.hidden')">
          <a-switch v-model:checked="form.hidden" />
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
import { menuApi } from '@/api/system'

const { t } = useI18n()
const loading = ref(false)
const rows = ref([])
const modalOpen = ref(false)
const editingId = ref(null)
const form = reactive({
  parentId: 0,
  name: '',
  title: '',
  type: 'menu',
  path: '',
  component: '',
  icon: '',
  permission: '',
  hidden: false,
  sort: 0
})
const columns = computed(() => [
  { title: t('field.title'), dataIndex: 'title', width: 180 },
  { title: t('field.type'), key: 'type', width: 110 },
  { title: t('field.route'), dataIndex: 'path', width: 180 },
  { title: t('field.component'), dataIndex: 'component', width: 220 },
  { title: t('field.permission'), dataIndex: 'permission', width: 180 },
  { title: t('field.hidden'), key: 'hidden', width: 100 },
  { title: t('field.sort'), dataIndex: 'sort', width: 100 },
  { title: t('field.action'), key: 'actions', width: 260, fixed: 'right' }
])
const parentOptions = computed(() => [{ id: 0, title: t('page.menu.title'), children: rows.value }])

async function load() {
  loading.value = true
  try {
    rows.value = await menuApi.tree()
  } finally {
    loading.value = false
  }
}

function resetForm(parent) {
  Object.assign(form, {
    parentId: parent?.id || 0,
    name: '',
    title: '',
    type: 'menu',
    path: '',
    component: '',
    icon: '',
    permission: '',
    hidden: false,
    sort: 0
  })
}

function openCreate(parent) {
  editingId.value = null
  resetForm(parent)
  modalOpen.value = true
}

function openEdit(record) {
  editingId.value = record.id
  Object.assign(form, {
    parentId: record.parentId,
    name: record.name,
    title: record.title,
    type: record.type,
    path: record.path,
    component: record.component,
    icon: record.icon,
    permission: record.permission,
    hidden: record.hidden,
    sort: record.sort
  })
  modalOpen.value = true
}

async function save() {
  if (editingId.value) {
    await menuApi.update(editingId.value, form)
  } else {
    await menuApi.create(form)
  }
  message.success(t('message.saveSuccess'))
  modalOpen.value = false
  load()
}

async function remove(id) {
  await menuApi.remove(id)
  message.success(t('message.deleteSuccess'))
  load()
}

onMounted(load)
</script>
