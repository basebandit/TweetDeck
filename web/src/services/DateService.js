import Api from "@/services/Api"

export default {
 getMinDate(token){
  return Api(token).get('/api/totals/weekly/mindate')
 }
}