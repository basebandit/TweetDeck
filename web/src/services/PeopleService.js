import Api from "@/services/Api"

export default{
  getPeople(token){
    return Api(token).get('/api/persons')
  }
}