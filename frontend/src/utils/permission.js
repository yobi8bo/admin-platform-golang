import { useAuthStore } from '../stores/auth'

export default {
  mounted(el, binding) {
    const auth = useAuthStore()
    const value = binding.value
    const allowed = Array.isArray(value)
      ? value.some((item) => auth.hasPermission(item))
      : auth.hasPermission(value)
    if (!allowed) el.remove()
  }
}
