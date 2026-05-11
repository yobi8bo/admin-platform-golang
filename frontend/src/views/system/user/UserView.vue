<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('system.user.title') }}</h1>
        <p class="page-subtitle">{{ t('system.user.desc') }}</p>
      </div>
      <a-button v-permission="'system:user:create'" type="primary" @click="openCreate">{{ t('system.user.create') }}</a-button>
    </div>
    <div class="panel">
      <div class="toolbar">
        <a-input-search v-model:value="query.keyword" :placeholder="t('system.user.keyword')" allow-clear @search="load" />
      </div>
      <a-table row-key="id" :columns="columns" :data-source="rows" :loading="loading" :pagination="pagination" @change="onPageChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'roles'">
            <a-tag v-for="role in record.roles" :key="role.id">{{ role.name }}</a-tag>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'enabled' ? 'green' : 'red'">{{ t(`common.${record.status}`) }}</a-tag>
          </template>
          <template v-if="column.key === 'actions'">
            <a-space>
              <a v-permission="'system:user:update'" @click="openEdit(record)">{{ t('common.edit') }}</a>
              <a-popconfirm :title="t('system.user.deleteConfirm')" @confirm="remove(record.id)">
                <a v-permission="'system:user:delete'" class="danger">{{ t('common.delete') }}</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
    <a-drawer v-model:open="drawerOpen" :title="editing?.id ? t('system.user.editTitle') : t('system.user.createTitle')" width="480">
      <a-form layout="vertical" :model="form" @finish="submit">
        <a-form-item :label="t('system.user.username')" name="username" :rules="[{ required: !editing?.id, message: t('system.user.usernameRequired') }]">
          <a-input v-model:value="form.username" :disabled="Boolean(editing?.id)" />
        </a-form-item>
        <a-form-item :label="t('system.user.nickname')" name="nickname" :rules="[{ required: true, message: t('system.user.nicknameRequired') }]">
          <a-input v-model:value="form.nickname" />
        </a-form-item>
        <a-form-item v-if="!editing?.id" :label="t('system.user.password')" name="password" :rules="[{ required: true, message: t('system.user.passwordRequired') }]">
          <a-input-password v-model:value="form.password" />
        </a-form-item>
        <a-form-item :label="t('system.user.email')"><a-input v-model:value="form.email" /></a-form-item>
        <a-form-item :label="t('system.user.mobile')"><a-input v-model:value="form.mobile" /></a-form-item>
        <a-form-item :label="t('common.status')">
          <a-select v-model:value="form.status">
            <a-select-option value="enabled">{{ t('common.enabled') }}</a-select-option>
            <a-select-option value="disabled">{{ t('common.disabled') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="t('system.user.roles')">
          <a-select v-model:value="form.roleIds" mode="multiple" :options="roleOptions" />
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
import { createUser, deleteUser, getRoles, getUsers, updateUser } from '../../../api/system'

const { t } = useI18n()
const loading = ref(false)
const saving = ref(false)
const drawerOpen = ref(false)
const editing = ref(null)
const rows = ref([])
const roles = ref([])
const query = reactive({ keyword: '', page: 1, pageSize: 10, total: 0 })
const form = reactive({ username: '', nickname: '', password: '', email: '', mobile: '', status: 'enabled', roleIds: [] })
const columns = computed(() => [
  { title: t('common.id'), dataIndex: 'id', width: 72 },
  { title: t('system.user.username'), dataIndex: 'username' },
  { title: t('system.user.nickname'), dataIndex: 'nickname' },
  { title: t('system.user.email'), dataIndex: 'email' },
  { title: t('system.user.roles'), key: 'roles' },
  { title: t('common.status'), key: 'status', width: 100 },
  { title: t('common.actions'), key: 'actions', width: 140 }
])
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
    if (editing.value?.id) await updateUser(editing.value.id, userUpdatePayload())
    else await createUser(userCreatePayload())
    drawerOpen.value = false
    message.success(t('common.saved'))
    load()
  } finally {
    saving.value = false
  }
}

async function remove(id) {
  await deleteUser(id)
  message.success(t('common.deleted'))
  load()
}

function userCreatePayload() {
  return {
    username: form.username,
    nickname: form.nickname,
    password: form.password,
    email: form.email,
    mobile: form.mobile,
    roleIds: form.roleIds
  }
}

function userUpdatePayload() {
  return {
    nickname: form.nickname,
    email: form.email,
    mobile: form.mobile,
    status: form.status,
    roleIds: form.roleIds
  }
}
</script>

<style scoped>
.danger {
  color: var(--danger);
}
</style>
