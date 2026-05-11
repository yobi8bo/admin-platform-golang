<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">部门管理</h1>
        <p class="page-subtitle">维护组织层级与部门状态</p>
      </div>
      <a-button type="primary" @click="openCreate">新增部门</a-button>
    </div>
    <div class="panel">
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="false" default-expand-all-rows>
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'"><a-tag :color="record.status === 'enabled' ? 'green' : 'red'">{{ record.status }}</a-tag></template>
          <template v-if="column.key === 'actions'">
            <a-space><a @click="openEdit(record)">编辑</a><a-popconfirm title="确认删除该部门？" @confirm="remove(record.id)"><a class="danger">删除</a></a-popconfirm></a-space>
          </template>
        </template>
      </a-table>
    </div>
    <a-drawer v-model:open="drawerOpen" :title="editing?.id ? '编辑部门' : '新增部门'" width="420">
      <a-form layout="vertical" :model="form" @finish="submit">
        <a-form-item label="上级部门"><a-tree-select v-model:value="form.parentId" allow-clear :tree-data="parentTree" :field-names="{ label: 'name', value: 'id', children: 'children' }" /></a-form-item>
        <a-form-item label="部门名称" name="name" :rules="[{ required: true, message: '请输入部门名称' }]"><a-input v-model:value="form.name" /></a-form-item>
        <a-form-item label="排序"><a-input-number v-model:value="form.sort" :min="0" /></a-form-item>
        <a-form-item label="状态"><a-select v-model:value="form.status"><a-select-option value="enabled">enabled</a-select-option><a-select-option value="disabled">disabled</a-select-option></a-select></a-form-item>
        <a-button type="primary" html-type="submit" :loading="saving">保存</a-button>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { createDept, deleteDept, getDeptTree, updateDept } from '../../../api/system'

const loading = ref(false)
const saving = ref(false)
const drawerOpen = ref(false)
const editing = ref(null)
const rows = ref([])
const form = reactive({ parentId: 0, name: '', sort: 0, status: 'enabled' })
const columns = [
  { title: '部门名称', dataIndex: 'name' },
  { title: '排序', dataIndex: 'sort', width: 90 },
  { title: '状态', key: 'status', width: 110 },
  { title: '操作', key: 'actions', width: 140 }
]
const parentTree = computed(() => [{ id: 0, name: '根部门', children: rows.value }])

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
    message.success('已保存')
    load()
  } finally { saving.value = false }
}
async function remove(id) {
  await deleteDept(id)
  message.success('已删除')
  load()
}
</script>

<style scoped>.danger { color: var(--danger); }</style>
