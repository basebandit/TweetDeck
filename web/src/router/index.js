import Vue from 'vue'
import VueRouter from 'vue-router'
import Dashboard from '@/views/Dashboard.vue'
import Team from '@/views/Team.vue'
import Member from '@/views/Member.vue'
import Auth from '@/components/Auth.vue'
import Layout from '@/components/Layout.vue'
import Login from '@/views/Login.vue'
import Report from '@/components/Report.vue'
import DailyReport from '@/views/DailyReport.vue'

Vue.use(VueRouter)

  const routes = [
    {
      path: "/report",
      name:"Report",
      component: Report,
      children:[
        {
          path:"daily",
          name:"DailyReport",
          component: DailyReport
        }
      ]
    },
    {
      path: "/auth",
      name: "Auth",
      component: Auth,
      children: [
        {
          path: "login",
          name: "Login",
          component: Login,
        },
      ],
    },
    {
      path:"/",
      component:Layout,
      children:[
        {
          path:"/",
          redirect:"/auth/login",
        },
        {
          path: '/dashboard',
          name: 'Dashboard',
          component: Dashboard,
          meta: {
            requiresAuth: true,
          },
        },
        {
          path: '/avatars',
          name: 'Avatars',
          // route level code-splitting
          // this generates a separate chunk (about.[hash].js) for this route
          // which is lazy-loaded when the route is visited.
          component: () => import(/* webpackChunkName: "about" */ '../views/Avatars.vue'),
          meta: {
            requiresAuth: true,
          },
        },
        {
          path:'/team',
          name:'Team',
          component:Team,
          meta: {
            requiresAuth: true,
          },
        },
        { 
          path: '/team/member/:id',
          name:'Member',
          component: Member, 
          props: true ,
          meta: {
            requiresAuth: true,
          },
        },
      ]
    }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
