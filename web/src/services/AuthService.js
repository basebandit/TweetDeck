import Api from "@/services/Api"

export default{
  token(user){
    return Api().get("/api/token",{auth:{
      username:user.email,
      password:user.password
    }})
  }
}
