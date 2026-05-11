<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">角色管理</h1>
        <p class="page-subtitle">配置角色、数据范围与菜单权限</p>
      </div>
      <a-button v-permission="'system:role:create'" type="primary" @click="openCreate">新增角色</a-button>
    </div>
    <div class="panel">
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="false">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'enabled' ? 'green' : 'red'">{{ record.status }}</a-tag>
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a v-permission="'system:role:update'" @click="openEdit(record)">编辑</a>
              <a-popconfirm title="确认删除该角色？" @confirm="remove(record.id)">
                <a v-permission="'system:role:delete'" class="danger">删除</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
    <a-drawer v-model:open="drawerOpen" :title="editing?.id ? '编辑角色' : '新增角色'" width="520">
      <a-form layout="vertical" :model="form" @finish="submit">
        <a-form-item label="角色编码" name="code" :rules="[{ required: true, message: '请输入角色编码' }]">
          <a-input v-model:value="form.code" />
        </a-form-item>
        <a-form-item label="角色名称" name="name" :rules="[{ required: true, message: '请输入角色名称' }]">
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="数据范围">
          <a-select v-model:value="form.dataScope">
            <a-select-option value="all">全部数据</a-select-option>
            <a-select-option value="dept">本部门</a-select-option>
            <a-select-option value="self">本人数据</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status">
            <a-select-option value="enabled">enabled</a-select-option>
            <a-select-option value="disabled">disabled</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="菜单权限">
          <a-tree v-model:checkedKeys="form.menuIds" checkable :tree-data="menuTree" :field-names="{ title: 'title', key: 'id', children: 'children' }" />
        </a-form-item>
        <a-button type="primary" html-type="submit" :loading="saving">保存</a-button>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { createRole, deleteRole, getMenuTree, getRoles, updateRole } from '../../../api/system'

const loading = ref(false)
const saving = ref(false)
const drawerOpen = ref(false)
const editing = ref(null)
const rows = ref([])
const menuTree = ref([])
const form = reactive({ code: '', name: '', sort: 0, status: 'enabled', dataScope: 'self', menuIds: [] })
const columns = [
  { title: 'ID', dataIndex: 'id', width: 72 },
  { title: '编码', dataIndex: 'code' },
  { title: '名称', dataIndex: 'name' },
  { title: '数据范围', dataIndex: 'dataScope' },
  { title: '状态', key: 'status', width: 100 },
  { title: '操作', key: 'actions', width: 140 }
]

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
    message.success('已保存')
    load()
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await deleteRole(id)
  message.success('已删除')
  load()
}
</script>

<style scoped>
.danger {
  color: var(--danger);
}
</style>
