<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('system.menu.title') }}</h1>
        <p class="page-subtitle">{{ t('system.menu.desc') }}</p>
      </div>
      <a-button type="primary" @click="openCreate">{{ t('system.menu.create') }}</a-button>
    </div>
    <div class="panel">
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="false" default-expand-all-rows>
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'type'"><a-tag>{{ t(`system.menu.types.${record.type}`) }}</a-tag></template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a @click="openEdit(record)">{{ t('common.edit') }}</a>
              <a-popconfirm :title="t('system.menu.deleteConfirm')" @confirm="remove(record.id)"><a class="danger">{{ t('common.delete') }}</a></a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
    <a-drawer v-model:open="drawerOpen" :title="editing?.id ? t('system.menu.editTitle') : t('system.menu.createTitle')" width="520">
      <a-form layout="vertical" :model="form" @finish="submit">
        <a-form-item :label="t('system.menu.parent')"><a-tree-select v-model:value="form.parentId" allow-clear :tree-data="parentTree" :field-names="{ label: 'title', value: 'id', children: 'children' }" /></a-form-item>
        <a-form-item :label="t('system.menu.name')" name="name" :rules="[{ required: true, message: t('system.menu.nameRequired') }]"><a-input v-model:value="form.name" /></a-form-item>
        <a-form-item :label="t('system.menu.displayTitle')" name="title" :rules="[{ required: true, message: t('system.menu.titleRequired') }]"><a-input v-model:value="form.title" /></a-form-item>
        <a-form-item :label="t('system.menu.type')"><a-segmented v-model:value="form.type" :options="menuTypeOptions" /></a-form-item>
        <a-form-item :label="t('system.menu.path')"><a-input v-model:value="form.path" /></a-form-item>
        <a-form-item :label="t('system.menu.component')"><a-input v-model:value="form.component" /></a-form-item>
        <a-form-item :label="t('system.menu.icon')"><a-input v-model:value="form.icon" /></a-form-item>
        <a-form-item :label="t('system.menu.permission')"><a-input v-model:value="form.permission" /></a-form-item>
        <a-form-item :label="t('system.menu.sort')"><a-input-number v-model:value="form.sort" :min="0" /></a-form-item>
        <a-form-item><a-checkbox v-model:checked="form.hidden">{{ t('system.menu.hidden') }}</a-checkbox></a-form-item>
        <a-button type="primary" html-type="submit" :loading="saving">{{ t('common.save') }}</a-button>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import { createMenu, deleteMenu, getMenuTree, updateMenu } from '../../../api/system'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const drawerOpen = ref(false)
const editing = ref(null)
const rows = ref([])
const form = reactive({ parentId: 0, name: '', title: '', type: 'menu', path: '', component: '', icon: '', permission: '', hidden: false, sort: 0 })
const columns = computed(() => [
  { title: t('system.menu.displayTitle'), dataIndex: 'title' },
  { title: t('system.menu.name'), dataIndex: 'name' },
  { title: t('system.menu.type'), key: 'type', width: 110 },
  { title: t('system.menu.path'), dataIndex: 'path' },
  { title: t('system.menu.permission'), dataIndex: 'permission' },
  { title: t('system.menu.sort'), dataIndex: 'sort', width: 80 },
  { title: t('common.actions'), key: 'actions', width: 140 }
])
const menuTypeOptions = computed(() => ['catalog', 'menu', 'button'].map((value) => ({ label: t(`system.menu.types.${value}`), value })))
const parentTree = computed(() => [{ id: 0, title: t('system.menu.root'), children: rows.value }])

onMounted(load)

async function load() {
  loading.value = true
  try {
    rows.value = await getMenuTree()
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  Object.assign(form, { parentId: 0, name: '', title: '', type: 'menu', path: '', component: '', icon: '', permission: '', hidden: false, sort: 0 })
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
    if (editing.value?.id) await updateMenu(editing.value.id, form)
    else await createMenu(form)
    drawerOpen.value = false
    message.success(t('common.saved'))
    load()
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await deleteMenu(id)
  message.success(t('common.deleted'))
  load()
}
</script>

<style scoped>.danger { color: var(--danger); }</style>
