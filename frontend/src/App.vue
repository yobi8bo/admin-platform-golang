<template>
  <a-config-provider :theme="theme" :locale="antdLocale">
    <router-view />
  </a-config-provider>
</template>

<script setup>
import { computed } from 'vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import enUS from 'ant-design-vue/es/locale/en_US'
import { useAppStore } from './stores/app'

const appStore = useAppStore()
const antdLocale = computed(() => (appStore.locale === 'en-US' ? enUS : zhCN))

const theme = {
  token: {
    colorPrimary: '#256f5d',
    colorInfo: '#256f5d',
    borderRadius: 6,
    fontFamily: 'Inter, "PingFang SC", "Microsoft YaHei", system-ui, sans-serif'
  }
}
</script>
