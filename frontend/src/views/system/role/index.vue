<template>
  <main class="page">
    <PageHeader :title="t('page.role.title')" :description="t('page.role.description')">
      <template #actions>
        <a-button type="primary" v-permission="'system:role:create'" @click="openCreate">
          {{ t('page.role.create') }}
        </a-button>
      </template>
    </PageHeader>

    <DataTable row-key="id" :title="t('page.role.tableTitle')" :columns="columns" :data-source="rows" :pagination="false" :loading="loading">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'dataScope'">
          <StatusTag :status="record.dataScope" />
        </template>
        <template v-if="column.key === 'status'">
          <StatusTag :status="record.status" />
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" v-permission="'system:role:update'" @click="openEdit(record)">
              {{ t('common.edit') }}
            </a-button>
            <ConfirmModal :title="t('page.role.deleteConfirm')" @confirm="remove(record.id)">
              <a-button size="small" danger v-permission="'system:role:delete'">{{ t('common.delete') }}</a-button>
            </ConfirmModal>
          </a-space>
        </template>
      </template>
    </DataTable>

    <a-modal v-model:open="modalOpen" :title="editingId ? t('page.role.edit') : t('page.role.create')" width="720px" @ok="save">
      <a-form layout="vertical" :model="form">
        <a-form-item :label="t('field.roleCode')" required>
          <a-input v-model:value="form.code" />
        </a-form-item>
        <a-form-item :label="t('field.roleName')" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item :label="t('field.dataScope')">
          <a-select v-model:value="form.dataScope" :options="dataScopeOptions" />
        </a-form-item>
        <a-form-item :label="t('field.status')">
          <a-select v-model:value="form.status" :options="statusOptions" />
        </a-form-item>
        <a-form-item :label="t('field.sort')">
          <a-input-number v-model:value="form.sort" :min="0" />
        </a-form-item>
        <a-form-item :label="t('field.menuPermission')">
          <a-tree
            v-model:checkedKeys="form.menuIds"
            checkable
            :tree-data="menuTree"
            :field-names="{ title: 'title', key: 'id', children: 'children' }"
          />
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
import { menuApi, roleApi } from '@/api/system'

const { t } = useI18n()
const loading = ref(false)
const rows = ref([])
const menuTree = ref([])
const modalOpen = ref(false)
const editingId = ref(null)
const form = reactive({ code: '', name: '', sort: 0, status: 'enabled', dataScope: 'self', menuIds: [] })
const statusOptions = computed(() => [
  { label: t('status.enabled'), value: 'enabled' },
  { label: t('status.disabled'), value: 'disabled' }
])
const dataScopeOptions = computed(() => [
  { label: t('status.all'), value: 'all' },
  { label: t('status.dept'), value: 'dept' },
  { label: t('status.self'), value: 'self' }
])
const columns = computed(() => [
  { title: t('field.roleCode'), dataIndex: 'code' },
  { title: t('field.roleName'), dataIndex: 'name' },
  { title: t('field.dataScope'), key: 'dataScope', width: 140 },
  { title: t('field.status'), key: 'status', width: 120 },
  { title: t('field.action'), key: 'actions', width: 180, fixed: 'right' }
])

async function load() {
  loading.value = true
  try {
    rows.value = await roleApi.list()
  } finally {
    loading.value = false
  }
}

async function loadMenus() {
  menuTree.value = await menuApi.tree()
}

function resetForm() {
  Object.assign(form, { code: '', name: '', sort: 0, status: 'enabled', dataScope: 'self', menuIds: [] })
}

function openCreate() {
  editingId.value = null
  resetForm()
  modalOpen.value = true
}

function openEdit(record) {
  editingId.value = record.id
  Object.assign(form, {
    code: record.code,
    name: record.name,
    sort: record.sort,
    status: record.status,
    dataScope: record.dataScope,
    menuIds: (record.menus || []).map((menu) => menu.id)
  })
  modalOpen.value = true
}

async function save() {
  if (editingId.value) {
    await roleApi.update(editingId.value, form)
  } else {
    await roleApi.create(form)
  }
  message.success(t('message.saveSuccess'))
  modalOpen.value = false
  load()
}

async function remove(id) {
  await roleApi.remove(id)
  message.success(t('message.deleteSuccess'))
  load()
}

onMounted(() => {
  load()
  loadMenus()
})
</script>
