export default{
  namespaced: true,
  state:{
  dialog: false,
  dailyReportStats:{},
  },
  getters:{
    showDialog:(state) => state.dialog,
    dailyReportStats:(state)=> state.dailyReportStats,
  },
  actions:{
    showDialog({commit},payload){
      const {dailyStats} = payload
      commit("openDialog",{dailyStats})
    },
    hideDialog({commit}){
      commit("closeDialog")
    }
  },
  mutations:{
    openDialog(state,payload){
      const {dailyStats} = payload
      state.dialog = true
      state.dailyReportStats = dailyStats
    } ,
    closeDialog(state){
      state.dialog = false
    }
  }
};