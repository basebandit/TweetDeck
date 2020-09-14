import Api from "@/services/Api"

export default{
  uploadAvatars(token,avatarFile){
    return Api(token).post('/api/avatar/upload',avatarFile)
  }
}