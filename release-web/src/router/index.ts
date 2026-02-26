import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  { path: '/', redirect: '/upload' },
  { path: '/upload', component: () => import('../views/upload/UploadView.vue') },
  { path: '/sql-audit', component: () => import('../views/sql-audit/SqlAuditView.vue') },
  { path: '/config-diff', component: () => import('../views/config-diff/ConfigDiffView.vue') },
  { path: '/dependency', component: () => import('../views/dependency/DependencyView.vue') },
  { path: '/approval', component: () => import('../views/approval/ApprovalView.vue') },
  { path: '/baseline', component: () => import('../views/baseline/BaselineView.vue') },
  { path: '/dashboard', component: () => import('../views/dashboard/DashboardView.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
