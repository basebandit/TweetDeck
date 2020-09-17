import Api from "@/services/Api"

export default {
  getPeople(token) {
    return Api(token).get('/api/people')
  },
  addPerson(token, person) {
    return Api(token).post('/api/people/new', person)
  }
}