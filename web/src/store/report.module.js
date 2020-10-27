import DateService from "../services/DateService";

export default{
  namespaced: true,
  state:{
  dialog: false,
  dailyReportStats:{},
  minDate:'',
  maxDate:'',
  startDate:'',
  endDate:''
  },
  getters:{
    showDialog:(state) => state.dialog,
    dailyReportStats:(state)=> state.dailyReportStats,
    minDate:(state) => state.minDate,
    maxDate:(state) => state.maxDate,
    weeklyStartDate:(state) => state.startDate,
    weeklyEndDate:(state)=> state.endDate
  },
  actions:{
    showDialog({commit},payload){
      const {dailyStats,router} = payload
      commit("openDialog",{dailyStats,router})
    },
    hideDialog({commit}){
      commit("closeDialog")
    },
    weeklyReportDateRange({commit},payload){
       const {startDate,endDate} = payload
       commit("setWeeklyReportDate",{startDate,endDate})
    },
    getDateRange({commit},payload){
      const {token} = payload
      DateService.getDateRange(token).then(response =>{
        if (response.status === 200){
          commit("dateRangeSuccess",{
            range:response.data,
          })
        }
      }).catch(err =>{
        commit("dateRangeFailure",{
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
    setWeeklyReportDate(state,payload){
      const{startDate,endDate} = payload
      state.startDate= startDate
      state.endDate=endDate
    },
    dateRangeSuccess(state,payload){
   const {range} = payload
   state.minDate = range.startDate
   state.maxDate = range.endDate
    },
    dateRangeFailure(payload){
const {message} = payload
/**eslint-disable */
console.log("MIN_DATE_FAILURE",message)
    }
  }
};