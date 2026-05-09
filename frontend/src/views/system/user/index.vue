<template>
  <main class="page">
    <PageHeader :title="t('page.user.title')" :description="t('page.user.description')">
      <template #actions>
        <a-button type="primary" v-permission="'system:user:create'" @click="openCreate">
          {{ t('page.user.create') }}
        </a-button>
      </template>
    </PageHeader>

    <SearchForm :title="t('page.user.searchTitle')" :model="query">
      <a-form-item :label="t('field.username')">
        <a-input v-model:value="query.keyword" :placeholder="t('page.user.keywordPlaceholder')" allow-clear />
      </a-form-item>
      <a-form-item>
        <a-space>
          <a-button type="primary" @click="load">{{ t('common.query') }}</a-button>
          <a-button @click="resetQuery">{{ t('common.reset') }}</a-button>
        </a-space>
      </a-form-item>
    </SearchForm>

    <DataTable
      row-key="id"
      :title="t('page.user.tableTitle')"
      :columns="columns"
      :data-source="rows"
      :pagination="pagination"
      :loading="loading"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <StatusTag :status="record.status" />
        </template>
        <template v-if="column.key === 'roles'">
          <a-space :size="[4, 4]" wrap>
            <a-tag v-for="role in record.roles" :key="role.id">{{ role.name }}</a-tag>
          </a-space>
        </template>
        <template v-if="column.key === 'actions'">
          <a-space>
            <a-button size="small" v-permission="'system:user:update'" @click="openEdit(record)">
              {{ t('common.edit') }}
            </a-button>
            <ConfirmModal v-if="!isCurrentUser(record)" :title="t('page.user.deleteConfirm')" @confirm="remove(record.id)">
              <a-button size="small" danger v-permission="'system:user:delete'">{{ t('common.delete') }}</a-button>
            </ConfirmModal>
            <a-tooltip v-else :title="t('page.user.currentUserTip')">
              <a-button size="small" disabled>{{ t('common.delete') }}</a-button>
            </a-tooltip>
          </a-space>
        </template>
      </template>
    </DataTable>

    <a-modal v-model:open="modalOpen" :title="editingId ? t('page.user.edit') : t('page.user.create')" @ok="save">
      <a-form layout="vertical" :model="form">
        <a-form-item :label="t('field.username')" required>
          <a-input v-model:value="form.username" :disabled="!!editingId" />
        </a-form-item>
        <a-form-item :label="t('field.nickname')" required>
          <a-input v-model:value="form.nickname" />
        </a-form-item>
        <a-form-item v-if="!editingId" :label="t('field.password')" required>
          <a-input-password v-model:value="form.password" />
        </a-form-item>
        <a-form-item :label="t('field.email')">
          <a-input v-model:value="form.email" />
        </a-form-item>
        <a-form-item :label="t('field.mobile')">
          <a-input v-model:value="form.mobile" />
        </a-form-item>
        <a-form-item :label="t('field.status')">
          <a-select v-model:value="form.status" :options="statusOptions" />
        </a-form-item>
        <a-form-item :label="t('field.role')">
          <a-select v-model:value="form.roleIds" mode="multiple" :options="roleOptions" />
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
import SearchForm from '@/components/common/SearchForm.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { roleApi, userApi } from '@/api/system'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const auth = useAuthStore()
const loading = ref(false)
const rows = ref([])
const roles = ref([])
const query = reactive({ keyword: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })
const modalOpen = ref(false)
const editingId = ref(null)
const form = reactive({ username: '', nickname: '', password: '', email: '', mobile: '', status: 'enabled', roleIds: [] })
const statusOptions = computed(() => [
  { label: t('status.enabled'), value: 'enabled' },
  { label: t('status.disabled'), value: 'disabled' }
])
const roleOptions = computed(() => roles.value.map((role) => ({ label: role.name, value: role.id })))
const columns = computed(() => [
  { title: t('field.username'), dataIndex: 'username', width: 150 },
  { title: t('field.nickname'), dataIndex: 'nickname', width: 150 },
  { title: t('field.email'), dataIndex: 'email', width: 220 },
  { title: t('field.status'), key: 'status', width: 120 },
  { title: t('field.role'), key: 'roles', width: 220 },
  { title: t('field.action'), key: 'actions', width: 180, fixed: 'right' }
])

async function load() {
  loading.value = true
  try {
    const data = await userApi.list({ page: pagination.current, pageSize: pagination.pageSize, keyword: query.keyword })
    rows.value = data.list
    pagination.total = data.total
  } finally {
    loading.value = false
  }
}

async function loadRoles() {
  roles.value = await roleApi.list()
}

function resetQuery() {
  query.keyword = ''
  pagination.current = 1
  load()
}

function resetForm() {
  Object.assign(form, { username: '', nickname: '', password: '', email: '', mobile: '', status: 'enabled', roleIds: [] })
}

function openCreate() {
  editingId.value = null
  resetForm()
  modalOpen.value = true
}

function openEdit(record) {
  editingId.value = record.id
  Object.assign(form, {
    username: record.username,
    nickname: record.nickname,
    password: '',
    email: record.email,
    mobile: record.mobile,
    status: record.status,
    roleIds: (record.roles || []).map((role) => role.id)
  })
  modalOpen.value = true
}

function isCurrentUser(record) {
  return Number(record.id) === Number(auth.user?.id)
}

async function save() {
  if (editingId.value) {
    await userApi.update(editingId.value, form)
  } else {
    await userApi.create(form)
  }
  message.success(t('message.saveSuccess'))
  modalOpen.value = false
  load()
}

async function remove(id) {
  await userApi.remove(id)
  message.success(t('message.deleteSuccess'))
  load()
}

function onTableChange(p) {
  pagination.current = p.current
  pagination.pageSize = p.pageSize
  load()
}

onMounted(() => {
  load()
  loadRoles()
})
</script>
