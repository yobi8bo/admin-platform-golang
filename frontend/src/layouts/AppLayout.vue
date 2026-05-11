<template>
  <a-layout class="shell">
    <a-layout-sider v-model:collapsed="app.collapsed" collapsible :trigger="null" class="shell-sider">
      <div class="brand" :class="{ compact: app.collapsed }">
        <div class="brand-mark">A</div>
        <div v-if="!app.collapsed">
          <strong>{{ t('app.name') }}</strong>
          <span>{{ t('app.console') }}</span>
        </div>
      </div>
      <a-menu
        theme="dark"
        mode="inline"
        :items="menuItems"
        :selected-keys="[route.path]"
        :open-keys="openKeys"
        @click="handleMenuClick"
        @openChange="openKeys = $event"
      />
    </a-layout-sider>
    <a-layout>
      <a-layout-header class="shell-header">
        <div class="header-left">
          <a-button type="text" :icon="h(MenuFoldOutlined)" @click="app.toggleCollapsed" />
          <a-breadcrumb :items="breadcrumbs" />
        </div>
        <div class="header-actions">
          <a-select :value="app.locale" class="locale-select" @change="changeLocale">
            <a-select-option value="zh-CN">中文</a-select-option>
            <a-select-option value="en-US">EN</a-select-option>
          </a-select>
          <a-dropdown>
            <button class="profile-button">
              <a-avatar :size="30" :src="avatarUrl">{{ auth.profile?.nickname?.slice(0, 1) || 'A' }}</a-avatar>
              <span>{{ auth.profile?.nickname || auth.profile?.username }}</span>
            </button>
            <template #overlay>
              <a-menu @click="handleUserAction">
                <a-menu-item key="profile">
                  <UserOutlined />
                  {{ t('app.profile') }}
                </a-menu-item>
                <a-menu-item key="logout">
                  <LogoutOutlined />
                  {{ t('app.logout') }}
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>
      <a-layout-content class="shell-content">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { computed, h, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  ApartmentOutlined,
  AuditOutlined,
  DashboardOutlined,
  FileOutlined,
  FolderOutlined,
  LoginOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuOutlined,
  ProfileOutlined,
  SettingOutlined,
  TeamOutlined,
  UserOutlined
} from '@ant-design/icons-vue'
import { useAppStore } from '../stores/app'
import { useAuthStore } from '../stores/auth'
import { getAvatarUrl } from '../api/file'

const app = useAppStore()
const auth = useAuthStore()
const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()
const openKeys = ref([])
const avatarUrl = ref('')

const icons = {
  ApartmentOutlined,
  AuditOutlined,
  DashboardOutlined,
  FileOutlined,
  FolderOutlined,
  LoginOutlined,
  MenuOutlined,
  ProfileOutlined,
  SettingOutlined,
  TeamOutlined,
  UserOutlined
}

const toMenuItem = (menu) => {
  const children = (menu.children || []).filter((item) => item.type !== 'button' && !item.hidden).map(toMenuItem)
  const Icon = icons[menu.icon] || FolderOutlined
  return {
    key: menu.path || menu.name,
    icon: () => h(Icon),
    label: t(`menu.${menu.name}`),
    children: children.length ? children : undefined
  }
}

const menuItems = computed(() => auth.menus.filter((item) => !item.hidden).map(toMenuItem))

const breadcrumbs = computed(() => {
  if (route.meta.menuName) return [{ title: t(`menu.${route.meta.menuName}`) }]
  if (route.meta.appKey) return [{ title: t(`app.${route.meta.appKey}`) }]
  return [{ title: t('app.dashboard') }]
})

watch(
  () => route.path,
  () => {
  const parent = auth.menus.find((item) => (item.children || []).some((child) => child.path === route.path))
  openKeys.value = parent ? [parent.path || parent.name] : []
  },
  { immediate: true }
)

watch(
  () => auth.profile?.avatarId,
  async (avatarId) => {
    if (!avatarId) {
      avatarUrl.value = ''
      return
    }
    const data = await getAvatarUrl(avatarId)
    avatarUrl.value = data.url
  },
  { immediate: true }
)

function handleMenuClick({ key }) {
  if (key.startsWith('/')) router.push(key)
}

function changeLocale(value) {
  app.setLocale(value)
  locale.value = value
}

async function handleUserAction({ key }) {
  if (key === 'profile') router.push('/profile')
  if (key === 'logout') {
    await auth.logout()
    router.replace('/login')
  }
}
</script>

<style scoped lang="scss">
.shell {
  min-height: 100vh;
}

.shell-sider {
  background: oklch(0.22 0.025 210);
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 64px;
  padding: 0 18px;
  color: #fff;

  &.compact {
    justify-content: center;
    padding: 0;
  }

  strong,
  span {
    display: block;
    white-space: nowrap;
  }

  strong {
    font-size: 15px;
  }

  span {
    color: rgba(255, 255, 255, 0.58);
    font-size: 12px;
  }
}

.brand-mark {
  display: grid;
  flex: 0 0 34px;
  width: 34px;
  height: 34px;
  place-items: center;
  border-radius: 8px;
  background: var(--primary);
  font-weight: 800;
}

.shell-header {
  position: sticky;
  top: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 18px;
  border-bottom: 1px solid var(--border);
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
}

.header-left,
.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.locale-select {
  width: 94px;
}

.profile-button {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 38px;
  border: 0;
  background: transparent;
  cursor: pointer;
}

.shell-content {
  padding: 20px;
}

@media (max-width: 720px) {
  .shell-sider {
    position: fixed;
    z-index: 20;
    min-height: 100vh;
  }

  .shell-content {
    padding: 14px;
  }

  .profile-button span,
  .locale-select {
    display: none;
  }
}
</style>
