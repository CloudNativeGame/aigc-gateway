import {createRouter, createWebHashHistory} from 'vue-router'
import WelcomeView from '../views/welcome.vue'
import DashboardView from '../views/dashboard.vue'

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            name: 'welcome',
            component: WelcomeView
        },
        {
            path: '/dashboard',
            name: 'dashboard',
            // route level code-splitting
            // this generates a separate chunk (About.[hash].js) for this route
            // which is lazy-loaded when the route is visited.
            component: DashboardView,
        }
    ]
})


router.beforeEach(async (to, from) => {
    if (!login && to.name != "welcome") {
        // redirect the user to the login page
        return {name: 'welcome'}
    }

    if (login && to.name == "welcome") {
        return {name: 'dashboard'}
    }
})

export default router
