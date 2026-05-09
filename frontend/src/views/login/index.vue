<template>
  <main class="login-page">
    <section class="login-page__intro" aria-labelledby="login-title">
      <a-tag color="processing">{{ t('app.title') }}</a-tag>
      <h1 id="login-title">{{ t('page.login.title') }}</h1>
      <p>{{ t('page.login.description') }}</p>
      <div class="login-page__features">
        <span><SafetyCertificateOutlined />{{ t('page.login.featureSecurity') }}</span>
        <span><AuditOutlined />{{ t('page.login.featureAudit') }}</span>
        <span><CloudServerOutlined />{{ t('page.login.featureStorage') }}</span>
      </div>
    </section>

    <section class="login-panel">
      <div class="login-panel__header">
        <span class="login-panel__mark">A</span>
        <div>
          <strong>{{ t('app.title') }}</strong>
          <small>{{ t('app.subtitle') }}</small>
        </div>
      </div>
      <a-form layout="vertical" :model="form" @finish="submit">
        <a-form-item :label="t('page.login.username')" name="username" :rules="[{ required: true, message: t('page.login.usernameRequired') }]">
          <a-input v-model:value="form.username" autocomplete="username" />
        </a-form-item>
        <a-form-item :label="t('page.login.password')" name="password" :rules="[{ required: true, message: t('page.login.passwordRequired') }]">
          <a-input-password v-model:value="form.password" autocomplete="current-password" />
        </a-form-item>
        <a-button type="primary" html-type="submit" block :loading="loading">{{ t('page.login.submit') }}</a-button>
      </a-form>
    </section>
  </main>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { AuditOutlined, CloudServerOutlined, SafetyCertificateOutlined } from '@ant-design/icons-vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const form = reactive({ username: 'admin', password: 'Admin@123' })
const loading = ref(false)
const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

async function submit() {
  loading.value = true
  try {
    await auth.login(form)
    router.push(route.query.redirect || '/dashboard')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 420px;
  align-items: center;
  gap: var(--space-8);
  padding: clamp(24px, 6vw, 72px);
  background:
    linear-gradient(135deg, rgba(22, 119, 255, 0.08), transparent 34%),
    var(--color-bg-layout);

  &__intro {
    max-width: 640px;
  }

  h1 {
    margin: var(--space-4) 0 var(--space-3);
    color: var(--color-text-heading);
    font-size: clamp(34px, 5vw, 56px);
    line-height: 1.08;
    font-weight: 700;
  }

  p {
    margin: 0;
    color: var(--color-text-secondary);
    font-size: 16px;
    line-height: 1.7;
  }

  &__features {
    display: flex;
    gap: var(--space-3);
    flex-wrap: wrap;
    margin-top: var(--space-6);
  }

  &__features span {
    display: inline-flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-3);
    color: var(--color-text);
    background: var(--color-bg-container);
    border: 1px solid var(--color-border);
    border-radius: 999px;
  }
}

.login-panel {
  width: 100%;
  padding: var(--space-6);
  background: var(--color-bg-container);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-card);

  &__header {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    margin-bottom: var(--space-6);
  }

  &__mark {
    width: 40px;
    height: 40px;
    display: grid;
    place-items: center;
    border-radius: var(--radius-lg);
    color: #fff;
    font-weight: 700;
    background: var(--color-primary);
  }

  &__header div {
    display: grid;
    gap: 2px;
  }

  &__header strong {
    color: var(--color-text-heading);
    font-size: 16px;
  }

  &__header small {
    color: var(--color-text-secondary);
  }
}

@media (max-width: 860px) {
  .login-page {
    grid-template-columns: 1fr;
  }
}
</style>
