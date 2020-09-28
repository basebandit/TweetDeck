import Vue from "vue"
import Vuex from "vuex"
import avatars from "./avatars.module"
import people from "./people.module"
import stats from "./stats.module"

Vue.use(Vuex)

export const store = new Vuex.Store({
  modules:{
    avatars,
    people,
    stats,
  },
});