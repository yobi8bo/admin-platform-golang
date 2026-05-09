<template>
  <template v-for="item in menus" :key="item.path || item.id">
    <a-sub-menu v-if="item.children?.length" :key="item.path || item.id">
      <template #icon>
        <component :is="iconMap[item.icon] || AppstoreOutlined" />
      </template>
      <template #title>{{ menuTitle(item) }}</template>
      <AppMenuTree :menus="item.children" />
    </a-sub-menu>
    <a-menu-item v-else :key="item.path">
      <template #icon>
        <component :is="iconMap[item.icon] || AppstoreOutlined" />
      </template>
      {{ menuTitle(item) }}
    </a-menu-item>
  </template>
</template>

<script setup>
import {
  ApartmentOutlined,
  AppstoreOutlined,
  AuditOutlined,
  DashboardOutlined,
  FileOutlined,
  LoginOutlined,
  MenuOutlined,
  ProfileOutlined,
  SettingOutlined,
  TeamOutlined,
  UserOutlined
} from '@ant-design/icons-vue'
import { useI18n } from 'vue-i18n'

defineProps({
  menus: {
    type: Array,
    default: () => []
  }
})

const { t, te } = useI18n()

const iconMap = {
  ApartmentOutlined,
  AppstoreOutlined,
  AuditOutlined,
  DashboardOutlined,
  FileOutlined,
  LoginOutlined,
  MenuOutlined,
  ProfileOutlined,
  SettingOutlined,
  TeamOutlined,
  UserOutlined
}

function menuTitle(item) {
  const key = `menu.${item.name || ''}`
  return te(key) ? t(key) : item.title
}
</script>
