export default{
  namespaced: true,
  state:{
  dialog: false
  },
  getters:{
    showDialog:(state) => state.dialog,
  },
  actions:{
    showDialog({commit}){
      commit("mutateDialog",{show:true})
    },
    hideDialog({commit}){
      commit("mutateDialog",{show:false})
    }
  },
  mutations:{
    mutateDialog(state,payload){
      const {show} = payload
      state.dialog = show
    } 
  }
};