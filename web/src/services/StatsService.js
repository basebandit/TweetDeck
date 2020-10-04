import Api from "@/services/Api"

export default{
  getTotals(token){
 return Api(token).get('/api/totals')
  },
  getTops(token){
    return Api(token).get('/api/totals/top')
  },
  // getTopByTweets(token){
  //   return Api(token).get('/api/top/tweets')
  // },
  // getTopByFollowers(token){
  //   return Api(token).get('/api/top/followers')
  // }
}