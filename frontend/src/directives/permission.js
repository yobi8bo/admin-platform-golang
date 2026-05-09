import { useAuthStore } from '@/stores/auth'

export default {
  mounted(el, binding) {
    const auth = useAuthStore()
    const required = Array.isArray(binding.value) ? binding.value : [binding.value]
    const allowed = required.some((item) => auth.permissions.includes(item))
    if (!allowed) {
      el.parentNode?.removeChild(el)
    }
  }
}
