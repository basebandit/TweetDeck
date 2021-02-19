import AvatarService from "@/services/AvatarService"

export default {
  namespaced: true,
  state: {
    avatars: [],
    fetchAvatars: false,
    fetchSuspendedAvatars:false,
    assignAvatars: false,
    avatarsByUser: false,
    uploadAvatars:false,
    fetchFromTwitter:false,
    userAvatars: [],
    suspendedAvatars:[]
  },

  getters: {
    avatars: (state) => state.avatars,
    fetching: (state) => state.fetchAvatars,
    uploading:(state) => state.uploadAvatars,
    assigning: (state) => state.assignAvatars,
    fetchAvatarsByUser: (state) => state.avatarsByUser,
    suspendedAvatars:(state)=> state.suspendedAvatars,
    fetchingFromTwitter:(state)=> state.fetchFromTwitter,
    userAvatars: (state) => state.userAvatars
  },

  actions: {
    getAvatarsFromTwitter({commit},payload){
     const {token} = payload
     commit("fetchFromTwitter")
     AvatarService.getAvatarsFromTwitter(token).then(response => {
       if(response.status === 200){
        commit("fetchFromTwitterSuccess",{message:"fetched twitter avatars successsfully"})
       }
     }).catch(err =>{
       commit("fetchFromTwitterFailure",{message:err.response.data.error})
     })
    },
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
    getSuspendedAvatars({commit},payload){
  const {token} = payload
  commit("suspendedAvatarsFetchStatus")
  AvatarService.getSuspendedAvatars(token).then(response=>{
      if (response.status == 200){
        commit("fetchSuspendedAvatarsSuccess",{
          avatars:response.data,
        })
      }
  }).catch(err => {
    commit("fetchSuspendedAvatarsFailure",{
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
   fetchFromTwitter(state){
     state.fetchFromTwitter = true
   },
   fetchFromTwitterSuccess(state,payload){
   const {message} = payload
   state.fetchFromTwitter = false
   /**eslint-disable */
   console.log("FETCH_FROM_TWITTER",message)
   },
   fetchFromTwitterFailure(state,payload){
    const {message} = payload
    state.fetchFromTwitter = false
    /**eslint-disable */
    console.log("FETCH_FROM_TWITTER",message)
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
    },
    suspendedAvatarsFetchStatus(state){
      state.fetchSuspendedAvatars = true;
          },
    fetchSuspendedAvatarsSuccess(state,payload){
      const {avatars} = payload
  state.suspendedAvatars = avatars
  state.fetchSuspendedAvatars =false
     },
     fetchSuspendedAvatarsFailure(state,payload){
      const {message} = payload
  state.fetchSuspendedAvatars =false
   /**eslint-disable */
   console.error(message)
     },
  },
};