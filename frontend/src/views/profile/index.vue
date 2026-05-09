<template>
  <main class="page profile-page">
    <section class="profile-shell">
      <aside class="profile-nav">
        <button
          v-for="item in navItems"
          :key="item.key"
          type="button"
          :class="{ 'is-active': activeTab === item.key }"
          @click="activeTab = item.key"
        >
          <span>{{ item.label }}</span>
          <a-badge v-if="item.count" :count="item.count" />
        </button>
      </aside>

      <div class="profile-content">
        <section v-if="activeTab === 'profile'" class="profile-section">
          <h2>{{ t('page.profile.title') }}</h2>
          <div class="profile-section__body">
            <a-form class="profile-form" layout="vertical" :model="form" @finish="save">
              <a-form-item :label="t('field.username')">
                <a-input :value="auth.user?.username" disabled />
              </a-form-item>
              <a-form-item :label="t('field.nickname')" name="nickname" :rules="[{ required: true, message: t('page.profile.nicknameRequired') }]">
                <a-input v-model:value="form.nickname" />
              </a-form-item>
              <a-form-item :label="t('field.email')">
                <a-input v-model:value="form.email" />
              </a-form-item>
              <a-form-item :label="t('field.mobile')">
                <a-input v-model:value="form.mobile" />
              </a-form-item>

              <div class="profile-form__actions">
                <a-button html-type="submit" type="primary" :loading="saving">{{ t('common.save') }}</a-button>
              </div>
            </a-form>

            <div class="avatar-panel">
              <h3>{{ t('page.profile.avatar') }}</h3>
              <a-upload accept="image/*" :show-upload-list="false" :custom-request="uploadAvatar">
                <button class="avatar-panel__upload" type="button" :disabled="uploading">
                  <a-avatar :size="150" :src="avatarPreview">{{ userInitial }}</a-avatar>
                  <span>{{ uploading ? t('page.profile.uploading') : t('page.profile.changeAvatar') }}</span>
                </button>
              </a-upload>
            </div>
          </div>
        </section>

        <section v-else-if="activeTab === 'password'" class="profile-section profile-section--narrow">
          <h2>{{ t('page.profile.changePassword') }}</h2>
          <a-form class="profile-form profile-form--password" layout="vertical" :model="passwordForm" @finish="savePassword">
            <a-form-item :label="t('page.profile.oldPassword')" name="oldPassword" :rules="[{ required: true, message: t('page.profile.oldPasswordRequired') }]">
              <a-input-password v-model:value="passwordForm.oldPassword" autocomplete="current-password" />
            </a-form-item>
            <a-form-item :label="t('page.profile.newPassword')" name="newPassword" :rules="[{ required: true, min: 6, message: t('page.profile.newPasswordRequired') }]">
              <a-input-password v-model:value="passwordForm.newPassword" autocomplete="new-password" />
            </a-form-item>
            <a-form-item :label="t('page.profile.confirmPassword')" name="confirmPassword" :rules="confirmPasswordRules">
              <a-input-password v-model:value="passwordForm.confirmPassword" autocomplete="new-password" />
            </a-form-item>

            <div class="profile-form__actions">
              <a-button html-type="submit" type="primary" :loading="passwordSaving">{{ t('common.save') }}</a-button>
            </div>
          </a-form>
        </section>

        <section v-else class="profile-section profile-section--empty">
          <h2>{{ currentNavLabel }}</h2>
        </section>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import { updatePassword, updateProfile } from '@/api/auth'
import { uploadAvatar as uploadAvatarFile } from '@/api/file'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const auth = useAuthStore()
const saving = ref(false)
const passwordSaving = ref(false)
const uploading = ref(false)
const activeTab = ref('profile')
const form = reactive({
  nickname: '',
  email: '',
  mobile: '',
  avatarId: null
})
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const avatarPreview = computed(() => auth.user?.avatarUrl || '')
const userInitial = computed(() => (auth.user?.nickname || auth.user?.username || t('common.unknownUser')).slice(0, 1).toUpperCase())
const navItems = computed(() => [
  { key: 'profile', label: t('page.profile.title') },
  { key: 'password', label: t('page.profile.changePassword') }
])
const currentNavLabel = computed(() => navItems.value.find((item) => item.key === activeTab.value)?.label || '')
const confirmPasswordRules = computed(() => [
  { required: true, message: t('page.profile.confirmPasswordRequired') },
  {
    validator: async (_rule, value) => {
      if (value && value !== passwordForm.newPassword) {
        return Promise.reject(t('page.profile.passwordMismatch'))
      }
      return Promise.resolve()
    }
  }
])

onMounted(() => {
  fillForm()
})

function fillForm() {
  form.nickname = auth.user?.nickname || ''
  form.email = auth.user?.email || ''
  form.mobile = auth.user?.mobile || ''
  form.avatarId = auth.user?.avatarId || null
}

async function save() {
  saving.value = true
  try {
    await updateProfile({
      nickname: form.nickname,
      email: form.email,
      mobile: form.mobile,
      avatarId: form.avatarId
    })
    await auth.refreshProfile()
    fillForm()
    message.success(t('message.saveSuccess'))
  } finally {
    saving.value = false
  }
}

async function uploadAvatar({ file, onSuccess, onError }) {
  uploading.value = true
  try {
    const record = await uploadAvatarFile(file)
    form.avatarId = record.id
    await save()
    onSuccess?.(record)
  } catch (error) {
    onError?.(error)
  } finally {
    uploading.value = false
  }
}

async function savePassword() {
  passwordSaving.value = true
  try {
    await updatePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    })
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    message.success(t('page.profile.passwordUpdated'))
  } finally {
    passwordSaving.value = false
  }
}
</script>

<style scoped lang="scss">
.profile-shell {
  display: grid;
  grid-template-columns: 200px minmax(0, 1fr);
  gap: 14px;
  align-items: stretch;
}

.profile-nav,
.profile-content {
  background: var(--color-bg-container);
  border: 1px solid var(--color-border);
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.03);
}

.profile-nav {
  display: grid;
  align-content: start;
  gap: 8px;
  padding: 8px;

  button {
    min-height: 42px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    padding: 0 18px;
    border: 0;
    border-radius: 8px;
    background: transparent;
    color: var(--color-text-heading);
    font-size: 14px;
    text-align: left;
    cursor: pointer;
    transition: color 0.16s ease, background-color 0.16s ease;
  }

  button:hover,
  button:focus-visible,
  button.is-active {
    color: var(--color-primary);
    background: #e6f4ff;
    outline: none;
  }
}

.profile-content {
  min-height: 360px;
  padding: var(--space-5);
}

.profile-section {
  h2 {
    margin: 0 0 var(--space-4);
    color: var(--color-text-heading);
    font-size: 22px;
    font-weight: 650;
  }

  &__body {
    display: grid;
    grid-template-columns: minmax(320px, 440px) minmax(220px, 1fr);
    gap: var(--space-8);
    align-items: start;
  }

  &--narrow {
    max-width: 420px;
  }

  &--empty {
    color: var(--color-text-secondary);
  }
}

.avatar-panel {
  h3 {
    margin: 0 0 var(--space-5);
    color: var(--color-text-heading);
    font-size: 15px;
    font-weight: 600;
  }

  &__upload {
    display: grid;
    justify-items: center;
    gap: var(--space-4);
    padding: 0;
    border: 0;
    background: transparent;
    color: var(--color-primary);
    cursor: pointer;
  }

  &__upload:disabled {
    cursor: wait;
    opacity: 0.72;
  }

  :deep(.ant-avatar) {
    border: 1px dashed #d9d9d9;
    background: #fafafa;
  }
}

.profile-form {
  :deep(.ant-form-item-label > label) {
    color: var(--color-text-heading);
    font-size: 14px;
  }

  :deep(.ant-input),
  :deep(.ant-input-affix-wrapper) {
    font-size: 14px;
  }

  :deep(.ant-input) {
    min-height: 32px;
  }

  :deep(.ant-input-affix-wrapper) {
    min-height: 32px;
    padding-top: 4px;
    padding-bottom: 4px;
  }

  &__actions {
    display: flex;
    justify-content: flex-start;
    margin-top: var(--space-2);
  }

  &--password {
    max-width: 360px;

    :deep(.ant-input) {
      min-height: 30px;
    }

    :deep(.ant-input-affix-wrapper) {
      min-height: 30px;
      padding-top: 3px;
      padding-bottom: 3px;
    }

    :deep(.ant-form-item) {
      margin-bottom: var(--space-3);
    }

    :deep(.ant-form-item-label) {
      padding-bottom: 4px;
    }

    .profile-form__actions {
      margin-top: var(--space-1);
    }
  }
}

@media (max-width: 820px) {
  .profile-shell,
  .profile-section__body {
    grid-template-columns: 1fr;
  }

  .profile-nav {
    grid-auto-flow: column;
    overflow-x: auto;
  }

  .profile-nav button {
    min-width: 130px;
    padding: 0 18px;
  }
}
</style>
