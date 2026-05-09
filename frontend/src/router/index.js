import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layouts/MainLayout.vue'
import { useAuthStore } from '@/stores/auth'

const viewModules = import.meta.glob('../views/**/*.vue')

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'Login', component: () => import('@/views/login/index.vue') },
    {
      path: '/',
      name: 'Root',
      component: Layout,
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'Dashboard',
          component: () => import('@/views/dashboard/index.vue'),
          meta: { title: 'Dashboard' }
        },
        {
          path: 'profile',
          name: 'Profile',
          component: () => import('@/views/profile/index.vue'),
          meta: { title: 'Profile' }
        }
      ]
    },
    { path: '/:pathMatch(.*)*', redirect: '/dashboard' }
  ]
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.path === '/login') {
    return true
  }
  if (!auth.accessToken) {
    return `/login?redirect=${encodeURIComponent(to.fullPath)}`
  }
  if (!auth.routesLoaded) {
    try {
      await auth.bootstrap()
      return to.fullPath
    } catch {
      auth.reset()
      return `/login?redirect=${encodeURIComponent(to.fullPath)}`
    }
  }
  return true
})

export function addDynamicRoutes(menus) {
  const routes = flattenMenus(menus)
  routes.forEach((menu) => {
    const component = viewModules[`../views/${menu.component}.vue`]
    if (!component) {
      return
    }
    if (!router.hasRoute(menu.name)) {
      router.addRoute('Root', {
        path: menu.path,
        name: menu.name,
        component,
        meta: {
          title: menu.title,
          icon: menu.icon,
          permission: menu.permission,
          hidden: menu.hidden
        }
      })
    }
  })
}

function flattenMenus(menus) {
  const out = []
  menus.forEach((menu) => {
    if (menu.type === 'menu') {
      out.push(menu)
    }
    if (menu.children?.length) {
      out.push(...flattenMenus(menu.children))
    }
  })
  return out
}

export default router
