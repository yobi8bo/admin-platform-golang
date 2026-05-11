import axios from 'axios'
import { message } from 'ant-design-vue'
import router from '../router'
import { useAuthStore } from '../stores/auth'

const service = axios.create({
  baseURL: '/api',
  timeout: 20000
})

service.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.accessToken) {
    config.headers.Authorization = `Bearer ${auth.accessToken}`
  }
  return config
})

service.interceptors.response.use(
  (response) => {
    const body = response.data
    if (body?.code === 0) return body.data
    message.error(body?.message || '请求失败')
    return Promise.reject(body)
  },
  (error) => {
    const status = error.response?.status
    const msg = error.response?.data?.message || error.message || '网络异常'
    const isLoginRequest = error.config?.url === '/auth/login'
    if (status === 401 && !isLoginRequest) {
      const auth = useAuthStore()
      auth.clearSession()
      router.replace({ name: 'Login', query: { redirect: router.currentRoute.value.fullPath } })
    }
    if (!isLoginRequest) message.error(msg)
    return Promise.reject(error)
  }
)

export default service
