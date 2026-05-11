<template>
  <main class="login-page">
    <section class="login-visual">
      <div>
        <p>ADMIN PLATFORM</p>
        <h1>{{ t('login.heading') }}</h1>
        <span>{{ t('login.tagline') }}</span>
      </div>
    </section>
    <section class="login-panel">
      <div class="login-card">
        <h2>{{ t('login.title') }}</h2>
        <p class="muted">{{ t('login.desc') }}</p>
        <a-form ref="formRef" layout="vertical" :model="form" @finish="submit">
          <a-form-item name="username" :label="t('login.username')" :rules="[{ required: true, message: t('login.usernameRequired') }]">
            <a-input v-model:value="form.username" size="large" autocomplete="username" />
          </a-form-item>
          <a-form-item name="password" :label="t('login.password')" :rules="[{ required: true, message: t('login.passwordRequired') }]">
            <a-input-password v-model:value="form.password" size="large" autocomplete="current-password" />
          </a-form-item>
          <a-button type="primary" size="large" block :loading="loading" @click="submit">{{ t('login.submit') }}</a-button>
        </a-form>
        <div class="demo-account">{{ t('login.demo') }}</div>
      </div>
    </section>
  </main>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const { t } = useI18n()
const auth = useAuthStore()
const router = useRouter()
const route = useRoute()
const formRef = ref()
const loading = ref(false)
const form = reactive({ username: 'admin', password: 'Admin@123' })

async function submit() {
  if (loading.value) return
  await formRef.value?.validate()
  loading.value = true
  try {
    await auth.login({ username: form.username, password: form.password })
    router.replace(route.query.redirect || '/dashboard')
  } catch (error) {
    message.error(error.response?.data?.message || error.message || t('login.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-page {
  display: grid;
  grid-template-columns: minmax(420px, 0.9fr) minmax(360px, 1fr);
  min-height: 100vh;
}

.login-visual {
  display: flex;
  align-items: flex-end;
  padding: 56px;
  color: #fff;
  background:
    linear-gradient(145deg, rgba(10, 31, 38, 0.82), rgba(13, 73, 62, 0.78)),
    url("https://images.unsplash.com/photo-1497366754035-f200968a6e72?auto=format&fit=crop&w=1400&q=80") center/cover;

  p {
    margin: 0 0 14px;
    letter-spacing: 0;
    opacity: 0.72;
  }

  h1 {
    margin: 0;
    font-size: 42px;
    line-height: 1.15;
  }

  span {
    display: block;
    margin-top: 16px;
    max-width: 420px;
    color: rgba(255, 255, 255, 0.76);
    font-size: 16px;
  }
}

.login-panel {
  display: grid;
  place-items: center;
  padding: 32px;
  background: var(--bg);
}

.login-card {
  width: min(420px, 100%);
  padding: 30px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--surface);
  box-shadow: var(--shadow);

  h2 {
    margin: 0;
    font-size: 28px;
  }
}

.demo-account {
  margin-top: 18px;
  padding: 10px 12px;
  border-radius: 6px;
  color: var(--text-muted);
  background: var(--surface-soft);
  font-size: 13px;
}

@media (max-width: 820px) {
  .login-page {
    grid-template-columns: 1fr;
  }

  .login-visual {
    min-height: 260px;
    padding: 32px;
  }
}
</style>
