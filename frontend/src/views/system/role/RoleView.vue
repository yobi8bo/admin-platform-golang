<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('system.role.title') }}</h1>
        <p class="page-subtitle">{{ t('system.role.desc') }}</p>
      </div>
      <a-button v-permission="'system:role:create'" type="primary" @click="openCreate">{{ t('system.role.create') }}</a-button>
    </div>
    <div class="panel">
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="false">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'enabled' ? 'green' : 'red'">{{ t(`common.${record.status}`) }}</a-tag>
          </template>
          <template v-if="column.key === 'dataScope'">
            {{ t(`system.role.dataScopes.${record.dataScope}`) }}
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a v-permission="'system:role:update'" @click="openEdit(record)">{{ t('common.edit') }}</a>
              <a-popconfirm :title="t('system.role.deleteConfirm')" @confirm="remove(record.id)">
                <a v-permission="'system:role:delete'" class="danger">{{ t('common.delete') }}</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
    <a-drawer v-model:open="drawerOpen" :title="editing?.id ? t('system.role.editTitle') : t('system.role.createTitle')" width="520">
      <a-form layout="vertical" :model="form" @finish="submit">
        <a-form-item :label="t('system.role.roleCode')" name="code" :rules="[{ required: true, message: t('system.role.codeRequired') }]">
          <a-input v-model:value="form.code" />
        </a-form-item>
        <a-form-item :label="t('system.role.roleName')" name="name" :rules="[{ required: true, message: t('system.role.nameRequired') }]">
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item :label="t('system.role.dataScope')">
          <a-select v-model:value="form.dataScope">
            <a-select-option value="all">{{ t('system.role.dataScopes.all') }}</a-select-option>
            <a-select-option value="dept">{{ t('system.role.dataScopes.dept') }}</a-select-option>
            <a-select-option value="self">{{ t('system.role.dataScopes.self') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="t('common.status')">
          <a-select v-model:value="form.status">
            <a-select-option value="enabled">{{ t('common.enabled') }}</a-select-option>
            <a-select-option value="disabled">{{ t('common.disabled') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="t('system.role.menuPermission')">
          <a-tree v-model:checkedKeys="form.menuIds" checkable :tree-data="menuTree" :field-names="{ title: 'title', key: 'id', children: 'children' }" />
        </a-form-item>
        <a-button type="primary" html-type="submit" :loading="saving">{{ t('common.save') }}</a-button>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import { createRole, deleteRole, getMenuTree, getRoles, updateRole } from '../../../api/system'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const drawerOpen = ref(false)
const editing = ref(null)
const rows = ref([])
const menuTree = ref([])
const form = reactive({ code: '', name: '', sort: 0, status: 'enabled', dataScope: 'self', menuIds: [] })
const columns = computed(() => [
  { title: t('common.id'), dataIndex: 'id', width: 72 },
  { title: t('system.role.code'), dataIndex: 'code' },
  { title: t('system.role.name'), dataIndex: 'name' },
  { title: t('system.role.dataScope'), key: 'dataScope' },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('common.actions'), key: 'actions', width: 140 }
])

onMounted(async () => {
  menuTree.value = await getMenuTree()
  load()
})

async function load() {
  loading.value = true
  try {
    rows.value = await getRoles()
  } finally {
    loading.value = false
  }
}

function resetForm() {
  Object.assign(form, { code: '', name: '', sort: 0, status: 'enabled', dataScope: 'self', menuIds: [] })
}

function openCreate() {
  editing.value = null
  resetForm()
  drawerOpen.value = true
}

function flattenMenuIds(items) {
  return items.flatMap((item) => [item.id, ...flattenMenuIds(item.children || [])])
}

function openEdit(record) {
  editing.value = record
  Object.assign(form, { ...record, menuIds: flattenMenuIds(record.menus || []) })
  drawerOpen.value = true
}

async function submit() {
  saving.value = true
  try {
    if (editing.value?.id) await updateRole(editing.value.id, form)
    else await createRole(form)
    drawerOpen.value = false
    message.success(t('common.saved'))
    load()
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await deleteRole(id)
  message.success(t('common.deleted'))
  load()
}
</script>

<style scoped>
.danger {
  color: var(--danger);
}
</style>
