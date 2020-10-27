import Api from "@/services/Api"

export default {
 getDateRange(token){
  return Api(token).get('/api/totals/weekly/mindate')
 }
}