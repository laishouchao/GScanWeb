import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/RegisterView.vue')
  },
  {
    path: '/tasks',
    name: 'tasks',
    component: () => import('../views/TaskListView.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/task/:id',
    name: 'taskDetail',
    component: () => import('../views/TaskDetailView.vue'),
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/users',
    name: 'users',
    component: () => import('../views/UserListView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const requiresAdmin = to.matched.some(record => record.meta.requiresAdmin)
  const user = localStorage.getItem('user')
  
  if (requiresAuth && !user) {
    next({ name: 'login' })
  } else if (requiresAdmin && user) {
    const userObj = JSON.parse(user)
    if (userObj.role !== 'admin') {
      next({ name: 'home' })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
