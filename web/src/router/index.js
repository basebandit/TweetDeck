import Vue from 'vue'
import VueRouter from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Team from '@/views/Team.vue'
import Member from '@/views/Member.vue'

Vue.use(VueRouter)

  const routes = [
  {
    path: '/',
    name: 'dashboard',
    component: Dashboard
  },
  {
    path: '/avatars',
    name: 'Avatars',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/Avatars.vue')
  },
  {
    path:'/team',
    name:'Team',
    component:Team
  },
  { 
    path: '/team/member/:id',
  name:'Member',
   component: Member, 
   props: true 
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
