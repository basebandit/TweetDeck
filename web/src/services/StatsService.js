import Api from "@/services/Api"

export default{
  getTotals(token){
 return Api(token).get('/api/totals')
  },
  getTops(token){
    return Api(token).get('/api/totals/top')
  },
  getWeeklyStats(token,start,end){
    return Api(token).get(`/api/totals/weekly?start=${start}&end=${end}`)
  }
  // getTopByTweets(token){
  //   return Api(token).get('/api/top/tweets')
  // },
  // getTopByFollowers(token){
  //   return Api(token).get('/api/top/followers')
  // }
}