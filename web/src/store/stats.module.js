import StatsService from "@/services/StatsService"

export default{
  namespaced:true,
  state:{
    totals:{},
    fetchTotals:false,
  },
  getters:{
    totals:(state) => state.totals,
    fetching:(state) => state.fetchTotals,
  },
  actions:{
    getTotals({commit},payload){
      const {token}=payload
      commit("totalsFetchStatus")
      StatsService.getTotals(token).then(response=>{
        if (response.status === 200){
          commit("totalsFetchSuccess",{
            totals:response.data,
          })
        }
      }).catch(err => {
        commit("totalsFetchFailure",{
          message:err.response.data.error
        })
      })
    }
  },
  mutations: {
    totalsFetchStatus(state){
      state.fetchTotals = true
    },
    totalsFetchSuccess(state,payload){
      const {totals} = payload
      state.totals = totals
      state.fetchTotals = false
  },
  totalsFetchFailure(state,payload){
    const {message } = payload
    state.fetchTotals = false
    /**eslint-disable */
    console.error("TOTALS_FETCH_FAILURE",message)
  },
},

}