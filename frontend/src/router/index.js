import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

export const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/LoginView.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    component: () => import('../layouts/AppLayout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('../views/dashboard/DashboardView.vue'), meta: { menuName: 'Dashboard' } },
      { path: 'system/user', name: 'User', component: () => import('../views/system/user/UserView.vue'), meta: { menuName: 'User', permission: 'system:user:list' } },
      { path: 'system/role', name: 'Role', component: () => import('../views/system/role/RoleView.vue'), meta: { menuName: 'Role', permission: 'system:role:list' } },
      { path: 'system/menu', name: 'Menu', component: () => import('../views/system/menu/MenuView.vue'), meta: { menuName: 'Menu', permission: 'system:menu:list' } },
      { path: 'system/dept', name: 'Dept', component: () => import('../views/system/dept/DeptView.vue'), meta: { menuName: 'Dept', permission: 'system:dept:list' } },
      { path: 'file', name: 'File', component: () => import('../views/file/FileView.vue'), meta: { menuName: 'File', permission: 'file:read' } },
      { path: 'audit/login-log', name: 'LoginLog', component: () => import('../views/audit/LoginLogView.vue'), meta: { menuName: 'LoginLog', permission: 'audit:login-log:list' } },
      { path: 'audit/operation-log', name: 'OperationLog', component: () => import('../views/audit/OperationLogView.vue'), meta: { menuName: 'OperationLog', permission: 'audit:operation-log:list' } },
      { path: 'profile', name: 'Profile', component: () => import('../views/profile/ProfileView.vue'), meta: { appKey: 'profile' } }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.meta.public) return true
  if (!auth.isLoggedIn) return { name: 'Login', query: { redirect: to.fullPath } }
  if (!auth.profile) await auth.bootstrap()
  if (to.meta.permission && !auth.hasPermission(to.meta.permission)) return '/dashboard'
  return true
})

export default router
