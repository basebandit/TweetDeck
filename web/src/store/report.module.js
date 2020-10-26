import DateService from "../services/DateService";

export default{
  namespaced: true,
  state:{
  dialog: false,
  dailyReportStats:{},
  minDate:'',
  },
  getters:{
    showDialog:(state) => state.dialog,
    dailyReportStats:(state)=> state.dailyReportStats,
    minDate:(state) => state.minDate,
  },
  actions:{
    showDialog({commit},payload){
      const {dailyStats,router} = payload
      commit("openDialog",{dailyStats,router})
    },
    hideDialog({commit}){
      commit("closeDialog")
    },
    getMinDate({commit},payload){
      const {token} = payload
      DateService.getMinDate(token).then(response =>{
        if (response.status === 200){
          commit("minDateSuccess",{
            minDate:response.data,
          })
        }
      }).catch(err =>{
        commit("minDateFailure",{
          message:err.response.data.error
        })
      })
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
    },
    minDateSuccess(state,payload){
   const {minDate} = payload
   state.minDate = minDate.startDate
    },
    minDateFailure(payload){
const {message} = payload
/**eslint-disable */
console.log("MIN_DATE_FAILURE",message)
    }
  }
};