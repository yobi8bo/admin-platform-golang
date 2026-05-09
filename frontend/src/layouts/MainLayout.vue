<template>
  <AppLayout
    v-model:collapsed="collapsed"
    v-model:openKeys="openKeys"
    :menus="visibleMenus"
    :user="auth.user"
    @logout="logout"
  />
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/common/AppLayout.vue'
import { useAuthStore } from '@/stores/auth'

const collapsed = ref(false)
const openKeys = ref(['/system', '/audit'])
const router = useRouter()
const auth = useAuthStore()

const visibleMenus = computed(() => filterMenus(auth.menus))

function filterMenus(menus) {
  return (menus || [])
    .filter((item) => !item.hidden)
    .map((item) => ({
      ...item,
      children: filterMenus(item.children || [])
    }))
    .filter((item) => item.type !== 'button')
}

async function logout() {
  await auth.logout()
  router.push('/login')
}
</script>
