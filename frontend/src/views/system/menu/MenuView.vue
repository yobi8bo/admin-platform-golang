<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">菜单管理</h1>
        <p class="page-subtitle">维护路由菜单、按钮权限与排序</p>
      </div>
      <a-button type="primary" @click="openCreate">新增菜单</a-button>
    </div>
    <div class="panel">
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="false" default-expand-all-rows>
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'type'"><a-tag>{{ record.type }}</a-tag></template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a @click="openEdit(record)">编辑</a>
              <a-popconfirm title="确认删除该菜单？" @confirm="remove(record.id)"><a class="danger">删除</a></a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
    <a-drawer v-model:open="drawerOpen" :title="editing?.id ? '编辑菜单' : '新增菜单'" width="520">
      <a-form layout="vertical" :model="form" @finish="submit">
        <a-form-item label="上级菜单"><a-tree-select v-model:value="form.parentId" allow-clear :tree-data="parentTree" :field-names="{ label: 'title', value: 'id', children: 'children' }" /></a-form-item>
        <a-form-item label="名称" name="name" :rules="[{ required: true, message: '请输入名称' }]"><a-input v-model:value="form.name" /></a-form-item>
        <a-form-item label="标题" name="title" :rules="[{ required: true, message: '请输入标题' }]"><a-input v-model:value="form.title" /></a-form-item>
        <a-form-item label="类型"><a-segmented v-model:value="form.type" :options="['catalog', 'menu', 'button']" /></a-form-item>
        <a-form-item label="路径"><a-input v-model:value="form.path" /></a-form-item>
        <a-form-item label="组件"><a-input v-model:value="form.component" /></a-form-item>
        <a-form-item label="图标"><a-input v-model:value="form.icon" /></a-form-item>
        <a-form-item label="权限标识"><a-input v-model:value="form.permission" /></a-form-item>
        <a-form-item label="排序"><a-input-number v-model:value="form.sort" :min="0" /></a-form-item>
        <a-form-item><a-checkbox v-model:checked="form.hidden">隐藏菜单</a-checkbox></a-form-item>
        <a-button type="primary" html-type="submit" :loading="saving">保存</a-button>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { createMenu, deleteMenu, getMenuTree, updateMenu } from '../../../api/system'

const loading = ref(false)
const saving = ref(false)
const drawerOpen = ref(false)
const editing = ref(null)
const rows = ref([])
const form = reactive({ parentId: 0, name: '', title: '', type: 'menu', path: '', component: '', icon: '', permission: '', hidden: false, sort: 0 })
const columns = [
  { title: '标题', dataIndex: 'title' },
  { title: '名称', dataIndex: 'name' },
  { title: '类型', key: 'type', width: 110 },
  { title: '路径', dataIndex: 'path' },
  { title: '权限标识', dataIndex: 'permission' },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '操作', key: 'actions', width: 140 }
]
const parentTree = computed(() => [{ id: 0, title: '根目录', children: rows.value }])

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
    message.success('已保存')
    load()
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await deleteMenu(id)
  message.success('已删除')
  load()
}
</script>

<style scoped>.danger { color: var(--danger); }</style>
