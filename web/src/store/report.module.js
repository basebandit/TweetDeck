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
      const {dailyStats,router} = payload
      commit("openDialog",{dailyStats,router})
    },
    hideDialog({commit}){
      commit("closeDialog")
    }
  },
  mutations:{
    openDialog(state,payload){
      const {dailyStats,router} = payload
      router.push({name:"DailyReport"})
      state.dialog = true
      state.dailyReportStats = dailyStats
    } ,
    closeDialog(state){
      state.dialog = false
    }
  }
};