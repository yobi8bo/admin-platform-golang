import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    collapsed: false,
    locale: localStorage.getItem('locale') || 'zh-CN'
  }),
  actions: {
    toggleCollapsed() {
      this.collapsed = !this.collapsed
    },
    setLocale(locale) {
      this.locale = locale
      localStorage.setItem('locale', locale)
    }
  }
})
