import PeopleService from "@/services/PeopleService"

export default {
  namespaced: true,
  state: {
    assignedPeople: [],
    unassignedPeople: [],
    fetchAssignedPeople: false,
    fetchUnassignedPeople: false,
    addPerson: false,
  },

  getters: {
    assignedTeam: (state) => state.assignedPeople,
    unassignedTeam: (state) => state.unassignedPeople,
    fetchingAssigned: (state) => state.fetchAssignedPeople,
    fetchingUnassigned: (state) => state.fetchUnassignedPeople,
    creating: (state) => state.addPerson
  },

  actions: {
    getAssignedPeople({ commit }, payload) {
      const { token } = payload
      commit("asPeopleFetchStatus")
      PeopleService.getAssignedPeople(token).then(response => {
        setTimeout(() => {
          if (response.status === 200) {
            commit("asPeopleFetchSuccess", {
              people: response.data,
            })
          }
        }, 500)
      }).catch(err => {
        commit("asPeopleFetchFailure", {
          message: err.response.data.error
        })
      })
    },
    getUnassignedPeople({ commit }, payload) {
      const { token } = payload
      commit("usPeopleFetchStatus")
      PeopleService.getUnassignedPeople(token).then(response => {
        setTimeout(() => {
          if (response.status === 200) {
            commit("usPeopleFetchSuccess", {
              people: response.data,
            })
          }
        }, 500)
      }).catch(err => {
        commit("usPeopleFetchFailure", {
          message: err.response.data.error
        })
      })
    },
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
    }
  },
  mutations: {
    usPeopleFetchStatus(state) {
      state.fetchUnassignedPeople = true;
    },
    usPeopleFetchSuccess(state, payload) {
      const { people } = payload
      state.fetchUnassignedPeople = false
      state.unassignedPeople = people
    },
    usPeopleFetchFailure(state, payload) {
      const { message } = payload
      state.fetchUnasignedPeople = false
      /**eslint-disable */
      console.error("usPeopleFetchFailure", message)
    },
    asPeopleFetchStatus(state) {
      state.fetchAssignedPeople = true;
    },
    asPeopleFetchSuccess(state, payload) {
      const { people } = payload
      state.fetchAssignedPeople = false
      state.assignedPeople = people
    },
    asPeopleFetchFailure(state, payload) {
      const { message } = payload
      state.fetchAssignedPeople = false
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
    }
  }
}