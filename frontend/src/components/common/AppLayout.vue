<template>
  <a-layout class="app-layout">
    <a-layout-sider
      v-model:collapsed="innerCollapsed"
      class="app-layout__sider"
      collapsible
      width="240"
      breakpoint="lg"
      :trigger="null"
      @breakpoint="onBreakpoint"
    >
      <router-link class="app-layout__brand" to="/dashboard" :aria-label="t('common.backDashboard')">
        <span class="app-layout__brand-mark">A</span>
        <span v-if="!innerCollapsed" class="app-layout__brand-copy">
          <strong>{{ t('app.title') }}</strong>
          <small>{{ t('app.subtitle') }}</small>
        </span>
      </router-link>

      <a-menu
        class="app-layout__menu"
        mode="inline"
        theme="light"
        :selectedKeys="[route.path]"
        :openKeys="openKeys"
        @openChange="emit('update:openKeys', $event)"
        @click="onMenuClick"
      >
        <AppMenuTree :menus="menus" />
      </a-menu>
    </a-layout-sider>

    <a-layout class="app-layout__main">
      <a-layout-header class="app-layout__header">
        <div class="app-layout__header-left">
          <a-button class="app-layout__icon-button" type="text" :aria-label="toggleLabel" @click="toggleCollapsed">
            <template #icon>
              <MenuUnfoldOutlined v-if="innerCollapsed" />
              <MenuFoldOutlined v-else />
            </template>
          </a-button>
          <router-link class="app-layout__home-link" to="/dashboard" :aria-label="t('common.backDashboard')">
            <HomeOutlined />
          </router-link>
          <span class="app-layout__page-title">{{ pageTitle }}</span>
        </div>

        <div class="app-layout__header-actions">
          <a-dropdown placement="bottomRight">
            <a-button class="app-layout__locale" type="text" :aria-label="t('common.language')">
              <template #icon>
                <GlobalOutlined />
              </template>
              <span class="app-layout__locale-label">{{ localeLabel }}</span>
            </a-button>
            <template #overlay>
              <a-menu :selectedKeys="[locale]" @click="onLocaleMenu">
                <a-menu-item key="zh">{{ t('common.languageZh') }}</a-menu-item>
                <a-menu-item key="en">{{ t('common.languageEn') }}</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
          <a-button class="app-layout__icon-button app-layout__notify" type="text" :aria-label="t('common.notifications')">
            <a-badge dot>
              <BellOutlined />
            </a-badge>
          </a-button>
          <a-button class="app-layout__icon-button" type="text" :aria-label="t('common.settings')">
            <template #icon>
              <SettingOutlined />
            </template>
          </a-button>
          <a-dropdown placement="bottomRight">
            <a-button class="app-layout__user" type="text">
              <a-avatar :size="24" :src="user?.avatarUrl">{{ userInitial }}</a-avatar>
            </a-button>
            <template #overlay>
              <a-menu @click="onUserMenu">
                <a-menu-item key="profile">
                  <UserOutlined />
                  {{ t('common.profile') }}
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout">
                  <LogoutOutlined />
                  {{ t('common.logout') }}
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <a-layout-content class="app-layout__content">
        <router-view v-slot="{ Component }">
          <transition name="route-fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
import { computed } from 'vue'
import { BellOutlined, GlobalOutlined, HomeOutlined, LogoutOutlined, MenuFoldOutlined, MenuUnfoldOutlined, SettingOutlined, UserOutlined } from '@ant-design/icons-vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import AppMenuTree from './AppMenuTree.vue'

const props = defineProps({
  collapsed: {
    type: Boolean,
    default: false
  },
  menus: {
    type: Array,
    default: () => []
  },
  openKeys: {
    type: Array,
    default: () => []
  },
  user: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:collapsed', 'update:openKeys', 'logout'])

const route = useRoute()
const router = useRouter()
const { locale, t, te } = useI18n()

const innerCollapsed = computed({
  get: () => props.collapsed,
  set: (value) => emit('update:collapsed', value)
})
const username = computed(() => props.user?.nickname || props.user?.username || t('common.unknownUser'))
const userInitial = computed(() => username.value.slice(0, 1).toUpperCase())
const pageTitle = computed(() => {
  if (route.path === '/dashboard') return t('common.home')
  if (route.path === '/profile') return t('common.profile')
  const name = route.name ? `menu.${String(route.name)}` : ''
  if (name && te(name)) return t(name)
  return route.meta.title || t('menu.Dashboard')
})
const toggleLabel = computed(() => (innerCollapsed.value ? t('common.expandMenu') : t('common.collapseMenu')))
const localeLabel = computed(() => (locale.value === 'zh' ? '中' : 'EN'))

function onMenuClick({ key }) {
  if (key && key !== route.path) router.push(key)
}

function onBreakpoint(isBroken) {
  innerCollapsed.value = isBroken
}

function toggleCollapsed() {
  innerCollapsed.value = !innerCollapsed.value
}

function onLocaleMenu({ key }) {
  if (!key || key === locale.value) return
  locale.value = key
  localStorage.setItem('locale', key)
  document.documentElement.lang = key
}

function onUserMenu({ key }) {
  if (key === 'profile') {
    router.push('/profile')
    return
  }
  if (key === 'logout') emit('logout')
}
</script>
