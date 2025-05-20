import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    {
        path: '/',
        name: 'Layout',
        component: () => import('../views/layout/index.vue'),
        redirect: '/home',
        children: [
            {
                path: 'home',
                name: 'Home',
                component: () => import('../views/home/index.vue'),
                meta: { title: '首页', keepAlive: true }
            },
            {
                path: 'category',
                name: 'Category',
                component: () => import('../views/category/index.vue'),
                meta: { title: '商品分类', keepAlive: true }
            },
            {
                path: 'product/:id',
                name: 'ProductDetail',
                component: () => import('../views/product/detail.vue'),
                meta: { title: '商品详情', keepAlive: false }
            },
            {
                path: 'cart',
                name: 'Cart',
                component: () => import('../views/cart/index.vue'),
                meta: { title: '购物车', keepAlive: false }
            },
            {
                path: 'user',
                name: 'User',
                component: () => import('../views/user/index.vue'),
                meta: { title: '个人中心', keepAlive: true },
                children: [
                    {
                        path: 'profile',
                        name: 'UserProfile',
                        component: () => import('../views/user/profile.vue'),
                        meta: { title: '个人资料', keepAlive: false }
                    },
                    {
                        path: 'orders',
                        name: 'UserOrders',
                        component: () => import('../views/user/orders.vue'),
                        meta: { title: '我的订单', keepAlive: false }
                    },
                    {
                        path: 'order/:id',
                        name: 'OrderDetail',
                        component: () => import('../views/user/order-detail.vue'),
                        meta: { title: '订单详情', keepAlive: false }
                    },
                    {
                        path: 'address',
                        name: 'UserAddress',
                        component: () => import('../views/user/address.vue'),
                        meta: { title: '收货地址', keepAlive: false }
                    },
                    {
                        path: 'favorites',
                        name: 'UserFavorites',
                        component: () => import('../views/user/favorites.vue'),
                        meta: { title: '我的收藏', keepAlive: false }
                    }
                ]
            }
        ]
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/login/index.vue'),
        meta: { title: '登录', keepAlive: false }
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('../views/register/index.vue'),
        meta: { title: '注册', keepAlive: false }
    },
    {
        path: '/checkout',
        name: 'Checkout',
        component: () => import('../views/checkout/index.vue'),
        meta: { title: '订单结算', keepAlive: false }
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('../views/error/404.vue')
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 全局导航守卫
router.beforeEach((to, from, next) => {
    // 设置页面标题
    document.title = to.meta.title ? `${to.meta.title} - 商城` : '商城'

    // 需要登录的页面
    const requiresAuth = ['Cart', 'Checkout', 'User', 'UserProfile', 'UserOrders', 'UserAddress', 'UserFavorites', 'OrderDetail']

    // 判断是否需要登录
    if (requiresAuth.includes(to.name)) {
        const token = localStorage.getItem('token')

        if (!token) {
            next({
                path: '/login',
                query: { redirect: to.fullPath }
            })
            return
        }
    }

    next()
})

export default router
