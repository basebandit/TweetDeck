import Api from "@/services/Api"

export default {
  getAssignedPeople(token) {
    return Api(token).get('/api/people/assigned')
  },
  getUnassignedPeople(token) {
    return Api(token).get('/api/people/unassigned')
  },
  addPerson(token, person) {
    return Api(token).post('/api/people/new', person)
  },
  updateFirstname(token,{id,firstname}){
    return Api(token).put(`/api/people/${id}/edit`,{firstname:firstname})
  },
  updateLastname(token,{id,lastname}){
    return Api(token).put(`/api/people/${id}/edit`,{lastname:lastname})
  }
}