import { createRouter, createWebHistory, createWebHashHistory } from 'vue-router'
import { getToken } from '../utils/auth'
import HomeView from '../view/HomeView.vue'
import AuthView from '../view/AuthView.vue'

const routes = [
    {
        path: '/',
        name: 'home',
        component: HomeView,
        meta: { requiresAuth: true }
    },
    {
        path: '/auth',
        name: 'auth',
        component: AuthView
    },
    {
        path: '/suppliers',
        name: 'suppliers',
        component: () => import('../view/SuppliersView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/documents',
        name: 'documents',
        component: () => import('../view/DocumentsView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/movements',
        name: 'movements',
        component: () => import('../view/MovementsView.vue'),
        meta: { requiresAuth: true }
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

router.beforeEach((to, from, next) => {
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
    const hasToken = getToken()

    if (requiresAuth && !hasToken) {
        next('/auth')
    } else if (to.path === '/auth' && hasToken) {
        next('/')
    } else {
        next()
    }
})

export default router 