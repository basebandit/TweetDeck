import Api from "@/services/Api"

export default{
  getTotals(token){
 return Api(token).get('/api/totals')
  },
}