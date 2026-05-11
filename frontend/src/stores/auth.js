import { defineStore } from 'pinia'
import { getMyMenuTree } from '../api/system'
import { getPermissions, getProfile, login as loginApi, logout as logoutApi } from '../api/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    accessToken: localStorage.getItem('accessToken') || '',
    refreshToken: localStorage.getItem('refreshToken') || '',
    profile: null,
    permissions: [],
    menus: []
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.accessToken),
    menuPermissions: (state) => flattenMenuPermissions(state.menus),
    hasPermission: (state) => (permission) => {
      if (!permission) return true
      return state.permissions.includes(permission) || flattenMenuPermissions(state.menus).includes(permission)
    }
  },
  actions: {
    async login(payload) {
      const data = await loginApi(payload)
      this.accessToken = data.accessToken
      this.refreshToken = data.refreshToken
      localStorage.setItem('accessToken', data.accessToken)
      localStorage.setItem('refreshToken', data.refreshToken)
      await this.bootstrap()
    },
    async bootstrap() {
      if (!this.accessToken) return
      const [profile, permissions, menus] = await Promise.all([getProfile(), getPermissions(), getMyMenuTree()])
      this.profile = profile
      this.permissions = permissions || []
      this.menus = menus || []
    },
    async logout() {
      try {
        if (this.accessToken) await logoutApi()
      } finally {
        this.clearSession()
      }
    },
    clearSession() {
      this.accessToken = ''
      this.refreshToken = ''
      this.profile = null
      this.permissions = []
      this.menus = []
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
    }
  }
})

function flattenMenuPermissions(menus = []) {
  return menus.flatMap((menu) => [
    menu.permission,
    ...flattenMenuPermissions(menu.children || [])
  ]).filter(Boolean)
}
