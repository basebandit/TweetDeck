import Vue from "vue"
import Vuex from "vuex"
import avatars from "./avatars.module"
import people from "./people.module"

Vue.use(Vuex)

export const store = new Vuex.Store({
  modules:{
    avatars,
    people,
  },
});