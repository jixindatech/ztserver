import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/* Router Modules */
// import componentsRouter from './modules/components'
// import chartsRouter from './modules/charts'
// import tableRouter from './modules/table'
// import nestedRouter from './modules/nested'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    noCache: true                if set true, the page will no be cached(default is false)
    affix: true                  if set true, the tag will affix in the tags-view
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index')
      }
    ]
  },
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/auth-redirect',
    component: () => import('@/views/login/auth-redirect'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error-page/401'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index'),
        name: 'dashboard',
        meta: { title: 'Dashboard', icon: 'dashboard', affix: true }
      }
    ]
  },
  {
    path: '/user',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/user/index'),
        name: 'User',
        meta: { title: '用户管理', icon: 'el-icon-user-solid', affix: true }
      }
    ]
  },
  {
    path: '/resource',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/resource/index'),
        name: 'Resource',
        meta: { title: '资源管理', icon: 'el-icon-s-grid', affix: true }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    children: [
      {
        path: 'index',
        component: () => import('@/views/system/index'),
        name: 'System',
        meta: { title: '系统管理', icon: 'el-icon-s-tools', affix: true }
      }
    ]
  },
  {
    path: '/gw',
    component: Layout,
    meta: {
      title: '网关管理',
      icon: 'el-icon-s-order'
    },
    children: [
      {
        path: 'ssl',
        component: () => import('@/views/gw/ssl'),
        name: 'SSL',
        meta: { title: '证书管理', icon: 'el-icon-setting', affix: true }
      },
      {
        path: 'router',
        component: () => import('@/views/gw/router'),
        name: 'Router',
        meta: { title: '路由管理', icon: 'el-icon-setting', affix: true }
      },
      {
        path: 'upstream',
        component: () => import('@/views/gw/upstream'),
        name: 'Upstream',
        meta: { title: 'Upstream管理', icon: 'el-icon-setting', affix: true }
      }
    ]
  },
  {
    path: '/event',
    component: Layout,
    redirect: '/event/gwevent',
    meta: {
      title: '行为审计',
      icon: 'el-icon-s-order'
    },
    children: [
      {
        path: 'gwevent',
        component: () => import('@/views/event/gw'),
        name: 'GwEvent',
        meta: { title: '网关日志', icon: 'el-icon-document', affix: true }
      },
      {
        path: 'wsevent',
        component: () => import('@/views/event/ws'),
        name: 'WsEvent',
        meta: { title: '网络日志', icon: 'el-icon-tickets', affix: true }
      }
    ]
  }
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const asyncRoutes = [
  {
    path: 'external-link',
    component: Layout,
    children: [
      {
        path: 'http://www.jixindatech.com',
        meta: { title: 'External Link', icon: 'link' }
      }
    ]
  },

  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
