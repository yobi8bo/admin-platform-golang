<template>
  <div class="page">
    <div class="page-header">
      <div>
        <h1 class="page-title">个人中心</h1>
        <p class="page-subtitle">维护基础资料与登录密码</p>
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
              <div class="avatar-title">头像</div>
              <div class="muted">支持 JPG、PNG、WebP 等图片格式</div>
              <a-upload
                accept="image/*"
                :show-upload-list="false"
                :custom-request="uploadProfileAvatar"
              >
                <a-button :loading="uploadingAvatar">上传新头像</a-button>
              </a-upload>
            </div>
          </div>
          <a-form layout="vertical" :model="profileForm" @finish="saveProfile">
            <a-form-item label="昵称" name="nickname" :rules="[{ required: true, message: '请输入昵称' }]">
              <a-input v-model:value="profileForm.nickname" />
            </a-form-item>
            <a-form-item label="邮箱">
              <a-input v-model:value="profileForm.email" />
            </a-form-item>
            <a-form-item label="手机号">
              <a-input v-model:value="profileForm.mobile" />
            </a-form-item>
            <a-button type="primary" html-type="submit" :loading="savingProfile">保存资料</a-button>
          </a-form>
        </div>
      </a-col>
      <a-col :xs="24" :lg="10">
        <div class="panel">
          <a-form layout="vertical" :model="passwordForm" @finish="savePassword">
            <a-form-item label="原密码" name="oldPassword" :rules="[{ required: true, message: '请输入原密码' }]">
              <a-input-password v-model:value="passwordForm.oldPassword" />
            </a-form-item>
            <a-form-item label="新密码" name="newPassword" :rules="[{ required: true, min: 6, message: '至少 6 位' }]">
              <a-input-password v-model:value="passwordForm.newPassword" />
            </a-form-item>
            <a-button html-type="submit" :loading="savingPassword">修改密码</a-button>
          </a-form>
        </div>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { updatePassword, updateProfile } from '../../api/auth'
import { getAvatarUrl, uploadAvatar } from '../../api/file'
import { useAuthStore } from '../../stores/auth'

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
    message.success('资料已更新')
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
    message.success('密码已更新')
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
    message.success('头像已更新')
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
