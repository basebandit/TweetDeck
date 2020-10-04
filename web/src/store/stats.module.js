import StatsService from "@/services/StatsService"

export default{
  namespaced:true,
  state:{
    totals:{},
    tops:{},
    fetchTotals:false,
    fetchTops: false,
  },
  getters:{
    totals:(state) => state.totals,
    tops:(state) => state.tops,
    fetching:(state) => state.fetchTotals,
    fetchingTops:(state) => state.fetchTops,
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
    },
    getTops({commit},payload){
      const {token} = payload
      commit("topsFetchStatus")
      StatsService.getTops(token).then(response => {
        if (response.status=== 200){
          commit("topsFetchSuccess",{
            tops: response.data,
          })
        }
      }).catch(err => {
        commit("topsFetchFailure",{
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
  topsFetchStatus(state){
    state.fetchTops = true
  },
  topsFetchSuccess(state,payload){
    const {tops} = payload
    state.tops = tops
    state.fetchTops = false
},
topsFetchFailure(state,payload){
  const {message } = payload
  state.fetchTops = false
  /**eslint-disable */
  console.error("TOPS_FETCH_FAILURE",message)
},
},

}