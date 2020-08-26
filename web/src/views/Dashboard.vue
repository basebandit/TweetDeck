<template>
  <div class="home">
    <h1 class="subheading grey--text">Dashboard</h1>

    <v-container class="my-5">
      <v-row class="mb-3">
        <v-col v-for="(stat,index) in stats" :key="index">
          <v-card :class="`stat ${stat.name}`">
            <v-card-title class="headline">
              <span>{{stat.total}}</span>
              <v-spacer></v-spacer>
              <span>
                <v-icon right color="primary">{{stat.icon}}</v-icon>
              </span>
            </v-card-title>
            <v-card-subtitle>{{stat.name}}</v-card-subtitle>
          </v-card>
        </v-col>
      </v-row>
      <v-row class="mb-3">
        <v-col cols="12">
          <v-card>
            <v-card-title class="headline">Total Daily Tweets</v-card-title>
            <v-card-text>
              <line-chart
                thousands=","
                label="Total Tweets"
                ytitle="No of tweets"
                xtitle="Days"
                :data="totalDailyTweets"
                :legend="true"
              ></line-chart>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row class="mb-3">
        <v-col cols="12" lg="8" md="8">
          <v-card flat class="mx-auto">
            <v-card-title class="align-start">Top 5 Weekly Avatars By Tweets</v-card-title>
            <div
              v-for="(avatar,index) in topWeeklyAvatarTweets"
              :key="avatar.handle"
              :class="`avatar-${index+1}`"
            >
              <v-row class="px-3">
                <v-col>
                  <div class="caption grey--text">Real Name</div>
                  <div>{{avatar.person}}</div>
                </v-col>
                <v-col>
                  <div class="caption grey--text">Avatar</div>
                  <div>{{avatar.handle}}</div>
                </v-col>
                <v-col>
                  <div class="caption grey--text">Number of tweets</div>
                  <div>{{avatar.tweets}}</div>
                </v-col>
              </v-row>
              <v-divider></v-divider>
            </div>
          </v-card>
        </v-col>
        <v-col cols="12" md="4" lg="4">
          <v-card class="mx-auto pa-3">
            <v-card-title class="headline">Top 5 Weekly Avatars By Tweets</v-card-title>
            <v-card-text>
              <pie-chart :round="2" suffix="%" :data="topWeeklyAvatarsByTweets" />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row class="mb-3">
        <v-col cols="12">
          <v-card>
            <v-card-title class="headline">Total Daily Followers</v-card-title>
            <v-card-text>
              <line-chart
                thousands=","
                label="Total Followers"
                ytitle="No of followers"
                xtitle="Days"
                :data="totalDailyFollowers"
                :legend="true"
              ></line-chart>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row class="mb-3">
        <v-col cols="12" lg="8" md="8">
          <v-card flat class="mx-auto">
            <v-card-title class="align-start">Top 5 Weekly Avatars By Followers</v-card-title>
            <div
              v-for="(avatar,index) in topWeeklyAvatarFollowers"
              :key="avatar.handle"
              :class="`avatar-${index+1}`"
            >
              <v-row class="px-3">
                <v-col>
                  <div class="caption grey--text">Real Name</div>
                  <div>{{avatar.person}}</div>
                </v-col>
                <v-col>
                  <div class="caption grey--text">Avatar</div>
                  <div>{{avatar.handle}}</div>
                </v-col>
                <v-col>
                  <div class="caption grey--text">Number of followers</div>
                  <div>{{avatar.followers}}</div>
                </v-col>
              </v-row>
              <v-divider></v-divider>
            </div>
          </v-card>
        </v-col>
        <v-col cols="12" md="4" lg="4">
          <v-card class="pa-3">
            <v-card-title class="headline">Top 5 Weekly Avatars By Followers</v-card-title>
            <v-card-text>
              <pie-chart :round="2" suffix="%" :data="topWeeklyAvatarsByFollowers" />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
// @ is an alias to /src

export default {
  name: "Home",
  data() {
    return {
      stats: [
        { total: 560, icon: "mdi-account-supervisor", name: "avatars" },
        { total: 289, icon: "mdi-twitter", name: "tweets" },
        { total: 389, icon: "mdi-account-arrow-left", name: "followers" },
        { total: 459, icon: "mdi-account-arrow-right", name: "following" },
        { total: 897, icon: "mdi-thumb-up", name: "likes" },
      ],
      totalDailyTweets: {
        Mon: 789,
        Tue: 876,
        Wed: 980,
        Thur: 680,
        Fri: 920,
        Sat: 689,
        Sun: 920,
      },
      topWeeklyAvatarTweets: [
        { handle: "Mwasya Kyalo", tweets: 269, person: "Francis Ngaruiya" },
        { handle: "Harris Kimani", tweets: 249, person: "Stephen Ndungu" },
        { handle: "Salma kambo", tweets: 152, person: "Stephen Ndungu" },
        {
          handle: "Kairitu karuiri",
          tweets: 141,
          person: "Josephine Waithiru",
        },
        { handle: "Kimani Otieno", tweets: 107, person: "Doughlas Marigiri" },
      ],
      totalDailyFollowers: {
        Mon: 480,
        Tue: 500,
        Wed: 480,
        Thur: 550,
        Fri: 650,
        Sat: 750,
        Sun: 800,
      },
      topWeeklyAvatarFollowers: [
        {
          handle: "Maina Kamanda Elder",
          followers: 400,
          person: "Jane Wanjiku",
        },
        {
          handle: "Cate Mukami",
          followers: 365,
          person: "Eric Kinyua",
        },
        {
          handle: "Mwasya kyalo",
          followers: 323,
          person: "Francis Ngaruiya",
        },
        {
          handle: "Harris Kimani",
          followers: 235,
          person: "Stephen Ndungu",
        },
        {
          handle: "Salma kambo",
          followers: 189,
          person: "Stephen Ndungu",
        },
      ],
    };
  },
  computed: {
    topWeeklyAvatarsByTweets() {
      let total = 0;
      this.topWeeklyAvatarTweets.forEach((avatar) => {
        total += avatar.tweets;
      });

      const obj = {};
      this.topWeeklyAvatarTweets.forEach(
        (avatar) => (obj[avatar.handle] = (avatar.tweets / total) * 100)
      );

      return obj;
    },
    topWeeklyAvatarsByFollowers() {
      let total = 0;
      this.topWeeklyAvatarFollowers.forEach((avatar) => {
        total += avatar.followers;
      });

      const obj = {};
      this.topWeeklyAvatarFollowers.forEach(
        (avatar) => (obj[avatar.handle] = (avatar.followers / total) * 100)
      );

      return obj;
    },
  },
};
</script>
<style>
.stat.avatars {
  border-left: 4px solid #3cd1c2;
}
.stat.tweets {
  border-left: 4px solid #ffaa2c;
}
.stat.followers {
  border-left: 4px solid #f83e70;
}
.stat.following {
  border-left: 4px solid #3cd1c2;
}
.stat.likes {
  border-left: 4px solid #ffaa2c;
}

.avatar-1 {
  border-left: 4px solid #3cd1c2;
}
.avatar-2 {
  border-left: 4px solid #ffaa2c;
}
.avatar-3 {
  border-left: 4px solid #f83e70;
}
.avatar-4 {
  border-left: 4px solid #3cd1c2;
}
.avatar-5 {
  border-left: 4px solid #ffaa2c;
}
</style>