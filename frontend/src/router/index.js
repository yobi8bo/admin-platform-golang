import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const ROOT_ROUTE_NAME = 'Root'

const viewModules = import.meta.glob('../views/**/*.vue')

const routeComponentAliases = {
  'dashboard/index': '../views/dashboard/DashboardView.vue',
  'system/user/index': '../views/system/user/UserView.vue',
  'system/role/index': '../views/system/role/RoleView.vue',
  'system/menu/index': '../views/system/menu/MenuView.vue',
  'system/dept/index': '../views/system/dept/DeptView.vue',
  'file/index': '../views/file/FileView.vue',
  'audit/login-log/index': '../views/audit/LoginLogView.vue',
  'audit/operation-log/index': '../views/audit/OperationLogView.vue'
}

export const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/LoginView.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    name: ROOT_ROUTE_NAME,
    component: () => import('../layouts/AppLayout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'profile', name: 'Profile', component: () => import('../views/profile/ProfileView.vue'), meta: { appKey: 'profile' } }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

let dynamicRouteSignature = ''
let removeDynamicRoutes = []

export function resetDynamicRoutes() {
  removeDynamicRoutes.forEach((removeRoute) => removeRoute())
  removeDynamicRoutes = []
  dynamicRouteSignature = ''
}

export function setupDynamicRoutes(menus = []) {
  const signature = JSON.stringify(flattenRouteMenus(menus).map((menu) => ({
    name: menu.name,
    path: menu.path,
    component: menu.component,
    permission: menu.permission
  })))

  if (signature === dynamicRouteSignature) return false

  resetDynamicRoutes()
  dynamicRouteSignature = signature

  flattenRouteMenus(menus).forEach((menu) => {
    const component = resolveMenuComponent(menu.component)
    if (!component) return

    removeDynamicRoutes.push(router.addRoute(ROOT_ROUTE_NAME, {
      path: normalizeChildPath(menu.path),
      name: menu.name,
      component,
      meta: {
        menuName: menu.name,
        permission: menu.permission,
        title: menu.title
      }
    }))
  })

  return true
}

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.meta.public) return true
  if (!auth.isLoggedIn) return { name: 'Login', query: { redirect: to.fullPath } }
  if (!auth.profile) await auth.bootstrap()
  const dynamicRoutesChanged = setupDynamicRoutes(auth.menus)
  if (dynamicRoutesChanged && !to.matched.length) return to.fullPath
  if (!to.matched.length) return '/dashboard'
  if (to.meta.permission && !auth.hasPermission(to.meta.permission)) return '/dashboard'
  return true
})

function flattenRouteMenus(menus = []) {
  return menus.flatMap((menu) => [
    ...(isRouteMenu(menu) ? [menu] : []),
    ...flattenRouteMenus(menu.children || [])
  ])
}

function isRouteMenu(menu) {
  return menu.type === 'menu' && Boolean(menu.path) && Boolean(menu.component)
}

function normalizeChildPath(path) {
  return path.replace(/^\/+/, '')
}

function resolveMenuComponent(component) {
  const normalized = component.replace(/^\/+/, '').replace(/\.vue$/, '')
  const candidates = [
    routeComponentAliases[normalized],
    `../views/${normalized}.vue`,
    `../views/${normalized}/index.vue`
  ].filter(Boolean)

  return viewModules[candidates.find((path) => viewModules[path])]
}

export default router
