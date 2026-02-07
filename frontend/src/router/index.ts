import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

/**
 * 路由配置
 * 定义了应用的导航结构
 */
const routes: RouteRecordRaw[] = [
  {
    // 扫描上传页面
    path: '/scan-upload',
    name: 'scan-upload',
    component: () => import('../views/ScanUpload.vue'),
    meta: {
      title: '扫描上传',
      icon: 'scan'
    }
  },
  {
    // URL上传页面
    path: '/url-upload',
    name: 'url-upload',
    component: () => import('../views/UrlUpload.vue'),
    meta: {
      title: 'URL上传',
      icon: 'link'
    }
  },
  {
    // 设置页面
    path: '/settings',
    name: 'settings',
    component: () => import('../views/Settings.vue'),
    meta: {
      title: '设置',
      icon: 'settings'
    }
  },
  {
    // 日志页面
    path: '/logs',
    name: 'logs',
    component: () => import('../views/Log.vue'),
    meta: {
      title: '日志',
      icon: 'logs'
    }
  },
  {
    // 关于页面
    path: '/about',
    name: 'about',
    component: () => import('../views/About.vue'),
    meta: {
      title: '关于',
      icon: 'about'
    }
  },
  {
    // 默认重定向到扫描上传页面
    path: '/',
    redirect: '/scan-upload'
  },
  {
    // 404 页面 - 匹配所有未定义的路由
    path: '/:pathMatch(.*)*',
    redirect: '/scan-upload'
  }
]

// 创建路由实例
const router = createRouter({
  // 使用 HTML5 历史模式
  history: createWebHistory(),
  routes
})

/**
 * 全局前置守卫
 * 每次路由切换前执行
 */
router.beforeEach((to, from, next) => {
  // 更新页面标题
  document.title = `${to.meta.title || 'ZPIC'} - 图床客户端`
  next()
})

export default router
