import { defineStore } from 'pinia'
import { login as loginApi, logout as logoutApi, myMenus, permissions, profile } from '@/api/auth'
import { avatarUrl } from '@/api/file'
import { addDynamicRoutes } from '@/router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    accessToken: localStorage.getItem('accessToken') || '',
    refreshToken: localStorage.getItem('refreshToken') || '',
    user: null,
    menus: [],
    permissions: [],
    routesLoaded: false
  }),
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
      const [user, perms, menus] = await Promise.all([profile(), permissions(), myMenus()])
      this.user = await this.withAvatarUrl(user)
      this.permissions = perms || []
      this.menus = menus || []
      addDynamicRoutes(this.menus)
      this.routesLoaded = true
    },
    async refreshProfile() {
      const user = await profile()
      this.user = await this.withAvatarUrl(user)
    },
    async withAvatarUrl(user) {
      if (!user?.avatarId) return user
      try {
        const data = await avatarUrl(user.avatarId)
        return { ...user, avatarUrl: data.url }
      } catch {
        return user
      }
    },
    async logout() {
      try {
        await logoutApi()
      } finally {
        this.reset()
      }
    },
    reset() {
      this.accessToken = ''
      this.refreshToken = ''
      this.user = null
      this.menus = []
      this.permissions = []
      this.routesLoaded = false
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
    },
    can(permission) {
      return this.permissions.includes(permission)
    }
  }
})
