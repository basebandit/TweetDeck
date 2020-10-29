import Api from "@/services/Api"

export default {
  uploadAvatars(token, avatarFile) {
    return Api(token).post('/api/avatar/upload', avatarFile)
  },
  getAvatars(token) {
    return Api(token).get('/api/avatar/')
  },
  getSuspendedAvatars(token){
    return Api(token).get('/api/avatar/suspended')
  },
  assignAvatars(token, payload) {
    return Api(token).post('/api/avatar/assign', payload)
  },
  getAvatarsByUser(token, id) {
    return Api(token).get(`/api/avatar/${id}`)
  },
  getAvatarsFromTwitter(token){
    return Api(token).get('/api/lookup/twitter')
  }
}