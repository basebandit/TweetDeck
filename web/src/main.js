import Vue from 'vue'
import App from './App.vue'
import Nprogress from "nprogress"
import "nprogress/nprogress.css"
import router from './router'
import Chartkick from 'vue-chartkick'
import Chart from 'chart.js'
import { store } from "./store";
import vuetify from './plugins/vuetify';





router.beforeEach((to,from,next) => {
  Nprogress.configure({showSpinner:true})
  Nprogress.start();
  if (to.matched.some(record => record.meta.requiresAuth)){
    //this route requires auth,check if client is logged in if not redirect to login page
    if(!window.localStorage.getItem('user')){
      next("/auth/login")
    }else{
      next()
    }
  }else{
    next()
  }
})

router.afterEach(()=>{
  Nprogress.done();
})

Vue.config.productionTip = false

Vue.use(Chartkick.use(Chart))


Vue.filter("formatDate",function(dateString){
  if (dateString){
  return new Date(dateString).toLocaleDateString(undefined,{ weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' })
  }
  return
})

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
