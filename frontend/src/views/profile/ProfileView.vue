<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('profile.title') }}</h1>
        <p class="page-subtitle">{{ t('profile.desc') }}</p>
      </div>
    </div>
    <a-row :gutter="[16, 16]">
      <a-col :xs="24" :lg="14">
        <div class="panel">
          <div class="avatar-section">
            <a-avatar :size="88" :src="avatarUrl">
              {{ profileForm.nickname?.slice(0, 1) || auth.profile?.username?.slice(0, 1) || 'A' }}
            </a-avatar>
            <div class="avatar-actions">
              <div class="avatar-title">{{ t('profile.avatar') }}</div>
              <div class="muted">{{ t('profile.avatarHelp') }}</div>
              <a-upload
                accept="image/*"
                :show-upload-list="false"
                :custom-request="uploadProfileAvatar"
              >
                <a-button :loading="uploadingAvatar">{{ t('profile.uploadAvatar') }}</a-button>
              </a-upload>
            </div>
          </div>
          <a-form layout="vertical" :model="profileForm" @finish="saveProfile">
            <a-form-item :label="t('profile.nickname')" name="nickname" :rules="[{ required: true, message: t('profile.nicknameRequired') }]">
              <a-input v-model:value="profileForm.nickname" />
            </a-form-item>
            <a-form-item :label="t('profile.email')">
              <a-input v-model:value="profileForm.email" />
            </a-form-item>
            <a-form-item :label="t('profile.mobile')">
              <a-input v-model:value="profileForm.mobile" />
            </a-form-item>
            <a-button type="primary" html-type="submit" :loading="savingProfile">{{ t('profile.saveProfile') }}</a-button>
          </a-form>
        </div>
      </a-col>
      <a-col :xs="24" :lg="10">
        <div class="panel">
          <a-form layout="vertical" :model="passwordForm" @finish="savePassword">
            <a-form-item :label="t('profile.oldPassword')" name="oldPassword" :rules="[{ required: true, message: t('profile.oldPasswordRequired') }]">
              <a-input-password v-model:value="passwordForm.oldPassword" />
            </a-form-item>
            <a-form-item :label="t('profile.newPassword')" name="newPassword" :rules="[{ required: true, min: 6, message: t('profile.passwordMin') }]">
              <a-input-password v-model:value="passwordForm.newPassword" />
            </a-form-item>
            <a-button html-type="submit" :loading="savingPassword">{{ t('profile.changePassword') }}</a-button>
          </a-form>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import { updatePassword, updateProfile } from '../../api/auth'
import { getAvatarUrl, uploadAvatar } from '../../api/file'
import { useAuthStore } from '../../stores/auth'

const { t } = useI18n()
const auth = useAuthStore()
const savingProfile = ref(false)
const savingPassword = ref(false)
const uploadingAvatar = ref(false)
const avatarUrl = ref('')
const profileForm = reactive({ nickname: '', email: '', mobile: '', avatarId: null })
const passwordForm = reactive({ oldPassword: '', newPassword: '' })

onMounted(async () => {
  Object.assign(profileForm, {
    nickname: auth.profile?.nickname || '',
    email: auth.profile?.email || '',
    mobile: auth.profile?.mobile || '',
    avatarId: auth.profile?.avatarId || null
  })
  await loadAvatarUrl()
})

async function saveProfile() {
  savingProfile.value = true
  try {
    await updateProfile(profileForm)
    await auth.bootstrap()
    message.success(t('profile.profileUpdated'))
  } finally {
    savingProfile.value = false
  }
}

async function savePassword() {
  savingPassword.value = true
  try {
    await updatePassword(passwordForm)
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    message.success(t('profile.passwordUpdated'))
  } finally {
    savingPassword.value = false
  }
}

async function uploadProfileAvatar({ file, onSuccess, onError }) {
  uploadingAvatar.value = true
  const form = new FormData()
  form.append('file', file)
  try {
    const avatar = await uploadAvatar(form)
    profileForm.avatarId = avatar.id
    await updateProfile(profileForm)
    await auth.bootstrap()
    await loadAvatarUrl()
    message.success(t('profile.avatarUpdated'))
    onSuccess?.(avatar)
  } catch (error) {
    onError?.(error)
  } finally {
    uploadingAvatar.value = false
  }
}

async function loadAvatarUrl() {
  if (!profileForm.avatarId) {
    avatarUrl.value = ''
    return
  }
  const data = await getAvatarUrl(profileForm.avatarId)
  avatarUrl.value = data.url
}
</script>

<style scoped lang="scss">
.avatar-section {
  display: flex;
  align-items: center;
  gap: 18px;
  margin-bottom: 22px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border);
}

.avatar-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.avatar-title {
  font-weight: 700;
}

@media (max-width: 560px) {
  .avatar-section {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
