<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('system.dept.title') }}</h1>
        <p class="page-subtitle">{{ t('system.dept.desc') }}</p>
      </div>
      <a-button type="primary" @click="openCreate">{{ t('system.dept.create') }}</a-button>
    </div>
    <div class="panel">
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="false" default-expand-all-rows>
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'"><a-tag :color="record.status === 'enabled' ? 'green' : 'red'">{{ t(`common.${record.status}`) }}</a-tag></template>
          <template v-if="column.key === 'actions'">
            <a-space><a @click="openEdit(record)">{{ t('common.edit') }}</a><a-popconfirm :title="t('system.dept.deleteConfirm')" @confirm="remove(record.id)"><a class="danger">{{ t('common.delete') }}</a></a-popconfirm></a-space>
          </template>
        </template>
      </a-table>
    </div>
    <a-drawer v-model:open="drawerOpen" :title="editing?.id ? t('system.dept.editTitle') : t('system.dept.createTitle')" width="420">
      <a-form layout="vertical" :model="form" @finish="submit">
        <a-form-item :label="t('system.dept.parent')"><a-tree-select v-model:value="form.parentId" allow-clear :tree-data="parentTree" :field-names="{ label: 'name', value: 'id', children: 'children' }" /></a-form-item>
        <a-form-item :label="t('system.dept.name')" name="name" :rules="[{ required: true, message: t('system.dept.nameRequired') }]"><a-input v-model:value="form.name" /></a-form-item>
        <a-form-item :label="t('system.dept.sort')"><a-input-number v-model:value="form.sort" :min="0" /></a-form-item>
        <a-form-item :label="t('common.status')"><a-select v-model:value="form.status"><a-select-option value="enabled">{{ t('common.enabled') }}</a-select-option><a-select-option value="disabled">{{ t('common.disabled') }}</a-select-option></a-select></a-form-item>
        <a-button type="primary" html-type="submit" :loading="saving">{{ t('common.save') }}</a-button>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import { createDept, deleteDept, getDeptTree, updateDept } from '../../../api/system'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const drawerOpen = ref(false)
const editing = ref(null)
const rows = ref([])
const form = reactive({ parentId: 0, name: '', sort: 0, status: 'enabled' })
const columns = computed(() => [
  { title: t('system.dept.name'), dataIndex: 'name' },
  { title: t('system.dept.sort'), dataIndex: 'sort', width: 90 },
  { title: t('common.status'), key: 'status', width: 110 },
  { title: t('common.actions'), key: 'actions', width: 140 }
])
const parentTree = computed(() => [{ id: 0, name: t('system.dept.root'), children: rows.value }])

onMounted(load)

async function load() {
  loading.value = true
  try { rows.value = await getDeptTree() } finally { loading.value = false }
}
function openCreate() {
  editing.value = null
  Object.assign(form, { parentId: 0, name: '', sort: 0, status: 'enabled' })
  drawerOpen.value = true
}
function openEdit(record) {
  editing.value = record
  Object.assign(form, record)
  drawerOpen.value = true
}
async function submit() {
  saving.value = true
  try {
    if (editing.value?.id) await updateDept(editing.value.id, form)
    else await createDept(form)
    drawerOpen.value = false
    message.success(t('common.saved'))
    load()
  } finally { saving.value = false }
}
async function remove(id) {
  await deleteDept(id)
  message.success(t('common.deleted'))
  load()
}
</script>

<style scoped>.danger { color: var(--danger); }</style>
