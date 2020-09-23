import PeopleService from "@/services/PeopleService"

export default {
  namespaced: true,
  state: {
    people: [],
    // unassignedPeople: [],
    fetchPeople: false,
    // fetchUnassignedPeople: false,
    updatingFirstname:false,
    updatingLastname: false,
    addPerson: false,
  },

  getters: {
    team: (state) => state.people,
    // unassignedTeam: (state) => state.unassignedPeople,
    fetchingPeople: (state) => state.fetchAssignedPeople,
    // fetchingUnassigned: (state) => state.fetchUnassignedPeople,
    creating: (state) => state.addPerson,
    updatingFirstname:(state) => state.updatingFirstname,
    updatingLastname:(state) => state.updatingLastname,
  },

  actions: {
    getPeople({ commit }, payload) {
      const { token } = payload
      commit("peopleFetchStatus")
      PeopleService.getPeople(token).then(response => {
        setTimeout(() => {
          if (response.status === 200) {
            commit("peopleFetchSuccess", {
              people: response.data,
            })
          }
        }, 500)
      }).catch(err => {
        commit("peopleFetchFailure", {
          message: err.response.data.error
        })
      })
    },
    // getUnassignedPeople({ commit }, payload) {
    //   const { token } = payload
    //   commit("usPeopleFetchStatus")
    //   PeopleService.getUnassignedPeople(token).then(response => {
    //     setTimeout(() => {
    //       if (response.status === 200) {
    //         commit("usPeopleFetchSuccess", {
    //           people: response.data,
    //         })
    //       }
    //     }, 500)
    //   }).catch(err => {
    //     commit("usPeopleFetchFailure", {
    //       message: err.response.data.error
    //     })
    //   })
    // },
    addPerson({ commit }, payload) {
      const { token, person, router } = payload
      commit("addPersonStatus")
      PeopleService.addPerson(token, person).then(response => {
        if (response.status === 201) {
          commit("addPersonSuccess", { message: "added new member successfully" })
          setTimeout(() => {
            router.go()//refresh page to see added member
          })
        }
      }).catch(err => {
        commit("addPersonFailure", { message: err.response.data.error })
      })
    },
    updateFirstname({commit},payload){
      const {token,firstname,id,router} = payload
      commit("updateFirstnameStatus")
      PeopleService.updateFirstname(token,{id,firstname}).then(response =>{
        if (response.status === 200){
          commit("updateFirstnameSuccess",{message:"successfully changed firstname"})
          setTimeout(()=>{
            router.go()
          },2000)
        }
      }).catch(err =>{
        commit("updateFirstnameFailure",{message: err.response.data.error})
      })
    },
    updateLastname({commit},payload){
      const {token,lastname,id,router} = payload
      commit("updateLastnameStatus")
      PeopleService.updateLastname(token,{id,lastname}).then(response =>{
        if (response.status === 200){
          commit("updateLastnameSuccess",{message:"successfully changed lastname"})
          setTimeout(()=>{
            router.go()
          },2000)
        }
      }).catch(err =>{
        commit("updateLastnameFailure",{message: err.response.data.error})
      })
    }
  },
  mutations: {
    // usPeopleFetchStatus(state) {
    //   state.fetchUnassignedPeople = true;
    // },
    // usPeopleFetchSuccess(state, payload) {
    //   const { people } = payload
    //   state.fetchUnassignedPeople = false
    //   state.unassignedPeople = people
    // },
    // usPeopleFetchFailure(state, payload) {
    //   const { message } = payload
    //   state.fetchUnasignedPeople = false
    //   /**eslint-disable */
    //   console.error("usPeopleFetchFailure", message)
    // },
    peopleFetchStatus(state) {
      state.fetchPeople = true;
    },
    peopleFetchSuccess(state, payload) {
      const { people } = payload
      state.fetchPeople = false
      state.people = people
    },
    peopleFetchFailure(state, payload) {
      const { message } = payload
      state.fetchPeople = false
      /**eslint-disable */
      console.error("asPeopleFetchFailure", message)
    },
    addPersonStatus(state) {
      state.addPerson = true
    },
    addPersonSuccess(state, payload) {
      const { message } = payload
      state.addPerson = false
      /**eslint-disable */
      console.log(message)
    },
    addPersonFailure(state, payload) {
      const { message } = payload
      state.addPerson = false
      /**eslint-disable */
      console.error(message)
    },
    updateFirstnameStatus(state){
      state.updatingFirstnameStatus = true
    },
    updateFirstnameSuccess(state,payload){
      const {message} = payload
      state.updatingFirstnameStatus = false
      /**eslint-disable */
      console.log(message)
    },
    updateFirstnameFailure(state,payload){
      const {message} = payload
      state.updatingFirstnameStatus = false
      /**eslint-disable */
      console.log(message)
    },
    updateLastnameSuccess(state,payload){
      const {message} = payload
      state.updatingFirstnameStatus = false
      /**eslint-disable */
      console.log(message)
    },
    updateLastnameFailure(state,payload){
      const {message} = payload
      state.updatingFirstnameStatus = false
      /**eslint-disable */
      console.log(message)
    }
  }
}