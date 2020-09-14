import Vue from 'vue'
import App from './App.vue'
import Nprogress from "nprogress"
import "nprogress/nprogress.css"
import router from './router'
import Chartkick from 'vue-chartkick'
import Chart from 'chart.js'
import vuetify from './plugins/vuetify';


Vue.config.productionTip = false

Vue.use(Chartkick.use(Chart))


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

new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
