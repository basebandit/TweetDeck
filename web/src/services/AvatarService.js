import Api from "@/services/Api"

export default{
  uploadAvatars(token,avatarFile){
    return Api(token).post('/api/avatar/upload',avatarFile)
  },
  getAvatars(token){
    return Api(token).get('/api/avatar/')
  }
}