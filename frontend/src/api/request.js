import axios from 'axios'
import { message } from 'ant-design-vue'
import i18n from '@/locales'

const { t } = i18n.global

const service = axios.create({
  baseURL: '/api',
  timeout: 15000
})

let refreshing = null

service.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

service.interceptors.response.use(
  (response) => {
    const body = response.data
    if (body?.code === 0) {
      return body.data
    }
    message.error(body?.message || t('message.requestFailed'))
    return Promise.reject(body)
  },
  async (error) => {
    const status = error.response?.status
    if (status === 401 && !error.config.__retried) {
      const refreshToken = localStorage.getItem('refreshToken')
      if (refreshToken) {
        error.config.__retried = true
        refreshing = refreshing || axios.post('/api/auth/refresh', { refreshToken }).finally(() => {
          refreshing = null
        })
        try {
          const res = await refreshing
          const data = res.data.data
          localStorage.setItem('accessToken', data.accessToken)
          localStorage.setItem('refreshToken', data.refreshToken)
          error.config.headers.Authorization = `Bearer ${data.accessToken}`
          return service(error.config)
        } catch {
          clearAuth()
        }
      }
    }
    if (status === 403) {
      message.error(t('message.forbidden'))
    } else if (status === 404) {
      message.error(t('message.notFound'))
    } else {
      message.error(error.response?.data?.message || t('message.networkError'))
    }
    return Promise.reject(error)
  }
)

function clearAuth() {
  localStorage.removeItem('accessToken')
  localStorage.removeItem('refreshToken')
  window.location.href = '/login'
}

export default service
