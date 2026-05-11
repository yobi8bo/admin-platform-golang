<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">用户管理</h1>
        <p class="page-subtitle">维护登录账号、联系方式与角色授权</p>
      </div>
      <a-button v-permission="'system:user:create'" type="primary" @click="openCreate">新增用户</a-button>
    </div>
    <div class="panel">
      <div class="toolbar">
        <a-input-search v-model:value="query.keyword" placeholder="用户名 / 昵称" allow-clear @search="load" />
      </div>
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="pagination" @change="onPageChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'roles'">
            <a-tag v-for="role in record.roles" :key="role.id">{{ role.name }}</a-tag>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'enabled' ? 'green' : 'red'">{{ record.status }}</a-tag>
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a v-permission="'system:user:update'" @click="openEdit(record)">编辑</a>
              <a-popconfirm title="确认删除该用户？" @confirm="remove(record.id)">
                <a v-permission="'system:user:delete'" class="danger">删除</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
    <a-drawer v-model:open="drawerOpen" :title="editing?.id ? '编辑用户' : '新增用户'" width="480">
      <a-form layout="vertical" :model="form" @finish="submit">
        <a-form-item label="用户名" name="username" :rules="[{ required: !editing?.id, message: '请输入用户名' }]">
          <a-input v-model:value="form.username" :disabled="Boolean(editing?.id)" />
        </a-form-item>
        <a-form-item label="昵称" name="nickname" :rules="[{ required: true, message: '请输入昵称' }]">
          <a-input v-model:value="form.nickname" />
        </a-form-item>
        <a-form-item v-if="!editing?.id" label="初始密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model:value="form.password" />
        </a-form-item>
        <a-form-item label="邮箱"><a-input v-model:value="form.email" /></a-form-item>
        <a-form-item label="手机号"><a-input v-model:value="form.mobile" /></a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="form.status">
            <a-select-option value="enabled">enabled</a-select-option>
            <a-select-option value="disabled">disabled</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="角色">
          <a-select v-model:value="form.roleIds" mode="multiple" :options="roleOptions" />
        </a-form-item>
        <a-button type="primary" html-type="submit" :loading="saving">保存</a-button>
      </a-form>
    </a-drawer>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { createUser, deleteUser, getRoles, getUsers, updateUser } from '../../../api/system'

const loading = ref(false)
const saving = ref(false)
const drawerOpen = ref(false)
const editing = ref(null)
const rows = ref([])
const roles = ref([])
const query = reactive({ keyword: '', page: 1, pageSize: 10, total: 0 })
const form = reactive({ username: '', nickname: '', password: '', email: '', mobile: '', status: 'enabled', roleIds: [] })
const columns = [
  { title: 'ID', dataIndex: 'id', width: 72 },
  { title: '用户名', dataIndex: 'username' },
  { title: '昵称', dataIndex: 'nickname' },
  { title: '邮箱', dataIndex: 'email' },
  { title: '角色', key: 'roles' },
  { title: '状态', key: 'status', width: 100 },
  { title: '操作', key: 'actions', width: 140 }
]
const pagination = computed(() => ({ current: query.page, pageSize: query.pageSize, total: query.total, showSizeChanger: true }))
const roleOptions = computed(() => roles.value.map((item) => ({ label: item.name, value: item.id })))

onMounted(async () => {
  roles.value = await getRoles()
  load()
})

async function load() {
  loading.value = true
  try {
    const data = await getUsers(query)
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

function resetForm() {
  Object.assign(form, { username: '', nickname: '', password: '', email: '', mobile: '', status: 'enabled', roleIds: [] })
}

function openCreate() {
  editing.value = null
  resetForm()
  drawerOpen.value = true
}

function openEdit(record) {
  editing.value = record
  Object.assign(form, { ...record, password: '', roleIds: (record.roles || []).map((role) => role.id) })
  drawerOpen.value = true
}

async function submit() {
  saving.value = true
  try {
    if (editing.value?.id) await updateUser(editing.value.id, form)
    else await createUser(form)
    drawerOpen.value = false
    message.success('已保存')
    load()
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await deleteUser(id)
  message.success('已删除')
  load()
}
</script>

<style scoped>
.danger {
  color: var(--danger);
}
</style>
