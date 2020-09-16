import PeopleService from "@/services/PeopleService"

export default {
  namespaced: true,
  state:{
    people:[],
    fetchPeople: false,
  },

  getters:{
    team: (state)=> state.people,
    fetching: (state) => state.fetchPeople
  },

  actions:{
    getPeople({commit},payload){
      const {token} = payload
      commit("updatePeopleFetchStatus")
      PeopleService.getPeople(token).then(response => {
        setTimeout(()=>{
          if (response.status === 200){
            commit("peopleFetchSuccess",{
              people: response.data,
            })
          }
        },500)
      }).catch(err => {
        commit("peopleFetchFailure",{
          message: err.response.data.error
        })
      })
    }
  },
  mutations:{
    updatePeopleFetchStatus(state){
      state.fetchPeople = true;
    },
    peopleFetchSuccess(state,payload){
      const {people} = payload
      state.fetchPeople = false
      state.people = people
    },
    peopleFetchFailure(state,payload){
      const {message} = payload
      state.fetchPeople = false
      /**eslint-disable */
      console.error("peopleFetchFailure",message)
    }
  }
}