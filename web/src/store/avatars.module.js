import AvatarService from "@/services/AvatarService"

export default {
  namespaced: true,
  state: {
    avatars: [],
    fetchAvatars: false,
    assignAvatars: false,
    avatarsByUser: false,
    uploadAvatars:false,
    userAvatars: []
  },

  getters: {
    avatars: (state) => state.avatars,
    fetching: (state) => state.fetchAvatars,
    uploading:(state) => state.uploadAvatars,
    assigning: (state) => state.assignAvatars,
    fetchAvatarsByUser: (state) => state.avatarsByUser,
    userAvatars: (state) => state.userAvatars
  },

  actions: {
    getAvatars({ commit }, payload) {
      const { token } = payload
      commit("updateAvatarFetchStatus")
      AvatarService.getAvatars(token).then(response => {
        setTimeout(() => {
          if (response.status === 200) {
            commit("avatarFetchSuccess", {
              avatars: response.data,
            })
          }
        }, 500)
      }).catch(err => {
        commit("avatarFetchFailure", {
          message: err.response.data.error
        })
      })
    },
    upload({commit},payload) {
     
   const {token,formData} = payload
        commit("uploadAvatarStatus")
         AvatarService.uploadAvatars(token,formData).then(response=>{
          if (response.status === 201){
              commit("uploadAvatarSuccess",{message:"avatars uploaded successfully"})
          } 
        }).catch(err => {
          commit("updateAvatarFailure",{message:err.response.data.error})
        })
      },
    assignAvatars({ commit }, payload) {
      const { token, assign,router} = payload
      commit("assignAvatarStatus")
      AvatarService.assignAvatars(token, assign).then(response => {
        if (response.status === 200) {
          commit("assignAvatarSuccess", {
            message: "asigned avatars successfully"
          })
          setTimeout(()=>{
            router.go()
          },500)
        }
      }).catch(err => {
        commit("assignAvatarFailure", { message: err.response.data.error })
      })

    },
    getAvatarsByUser({ commit }, payload) {
      const { token, id } = payload
      commit("avatarsByUserStatus")
      AvatarService.getAvatarsByUser(token, id).then(response => {
        if (response.status === 200) {
          commit("avatarsByUserSuccess", { avatars: response.data })
        }
      }).catch(err => {
        commit("avatarsByUserFailure", { message: err.response.data.error })
      })
    }
  },

  mutations: {
    updateAvatarFetchStatus(state) {
      state.fetchAvatars = true;
    },
    avatarFetchSuccess(state, payload) {
      const { avatars } = payload;
      /**eslint-disable */
      console.log(avatars)
      state.fetchAvatars = false;
      state.avatars = avatars;
    },
    avatarFetchFailure(state, payload) {
      const { message } = payload;
      /**eslint-disable */
      console.error("avatarFetchFailure", message)
      state.fetchAvatars = false;
    },
    assignAvatarStatus(state) {
      state.assignAvatars = true
    },
    assignAvatarSuccess(state, payload) {
      const { message } = payload
      /**eslint-disable */
      console.log("assignAvatarsSuccess", message)
      state.assignAvatars = false;
    },
    assignAvatarFailure(state, payload) {
      const { message } = payload
      /**eslint-disable */
      console.log("assignAvatarFailure", message)
      state.assignAvatars = false;
    },
    avatarsByUserStatus(state) {
      state.avatarsByUser = true
    },
    avatarsByUserSuccess(state, payload) {
      const { avatars } = payload
      state.avatarsByUser = false
      state.userAvatars = avatars
    },
    avatarsByUserFailure(state, payload) {
      const { message } = payload
      state.avatarsByUser = false
      /**eslint-disable */
      console.log(message)
    },
    uploadAvatarStatus(state){
      state.uploadAvatars = true
    },
    uploadAvatarSuccess(state,payload){
      const {message} = payload
      state.uploadAvatars = false
      /**eslint-disable */
      console.log(message)
    },
    uploadAvatarFailure(state,payload){
      const {message} = payload
      state.uploadAvatars = false
      /**eslint-disable */
      console.error(message)
    }
  },
};