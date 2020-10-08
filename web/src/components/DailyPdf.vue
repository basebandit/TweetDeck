<template>
  <div id="document">
    <v-container>
      <v-row justify="center" class="font-weight-bold">
        <v-col align="left"> Cyclones </v-col>
        <v-spacer></v-spacer>
        <v-col align="right">
          <div>Summary of Twitter Activity</div>
          <div>{{ date | formatDate }}</div>
        </v-col>
      </v-row>
      <v-divider></v-divider>
      <v-row align-content="space-between">
        <v-col class="text-left">
          <div class="text-decoration-underline">TOTAL TWITTER ACCOUNTS</div>
          <div class="font-weight-bold">{{ activeAccts.toLocaleString() }}</div>
        </v-col>
        <v-col class="text-center">
          <div class="text-decoration-underline">TOTAL NEW ACCOUNTS</div>
          <div class="font-weight-bold">{{ newAccts.toLocaleString() }}</div>
        </v-col>
        <v-col class="text-right">
          <div class="text-decoration-underline">TIMEFRAME (PERIOD)</div>
          <div class="font-weight-bold">{{ date | formatDate }}</div>
        </v-col>
      </v-row>
      <v-row>
        <v-col class="text-left">
          <div class="text-decoration-underline">TOTAL DAY'S FOLLOWERS</div>
          <div class="font-weight-bold">
            {{ totalFollowers.toLocaleString() }}
          </div>
        </v-col>
        <v-col class="text-center">
          <div class="text-decoration-underline">TOTAL DAY'S FOLLOWING</div>
          <div class="font-weight-bold">
            {{ totalFollowing.toLocaleString() }}
          </div>
        </v-col>
        <v-col class="text-right">
          <div class="text-decoration-underline">TOTAL DAY'S TWEETS</div>
          <div class="font-weight-bold">{{ totalTweets.toLocaleString() }}</div>
        </v-col>
      </v-row>
      <v-row>
        <v-col class="text-center">
          <div class="text-decoration-underline">TOTAL DAY'S LIKES</div>
          <div class="font-weight-bold">{{ totalLikes.toLocaleString() }}</div>
        </v-col>
      </v-row>
      <v-row>
        <v-col class="text-center">
          <div class="text-decoration-underline">
            AVATAR WITH THE HIGHEST GAINED FOLLOWERS
          </div>
          <div class="font-weight-bold">
            {{ highestGainedFollowers.avatar }}
            {{ highestGainedFollowers.count.toLocaleString() }} -
            {{ highestGainedFollowers.person }}
          </div>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col align="left">
          <div class="text-decoration-underline">
            AVATAR WITH THE HIGHEST LIKES
          </div>
          <div class="font-weight-bold">
            {{ highestLikes.avatar }}
            {{ highestLikes.count.toLocaleString() }} -
            {{ highestLikes.person }}
          </div>
        </v-col>
        <v-col align="right">
          <div class="text-decoration-underline">TOTAL UNUSED ACCOUNTS</div>
          <div class="font-weight-bold">
            {{ totalUnusedAccounts.toLocaleString() }}
          </div>
        </v-col>
      </v-row>
      <v-row class="my-3">
        <v-col>
          <v-card class="pa-3" v-if="tops.followers">
            <v-card-title class="align-start"
              >TOP 5 AVATARS (BY TWEETS)</v-card-title
            >

            <template>
              <v-simple-table>
                <template v-slot:default>
                  <thead>
                    <tr>
                      <th class="text-left">REAL NAME</th>
                      <th class="text-left">AVATAR NAME</th>
                      <th class="text-left">NUMBER OF TWEETS</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="item in tops.tweets" :key="item.username">
                      <td>{{ item.person || "unassigned" }}</td>
                      <td>{{ item.username }}</td>
                      <td>{{ item.tweets.toLocaleString() }}</td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
            </template>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <v-card
            class="pa-3"
            v-if="Object.keys(topFiveAvatarsByTweets).length > 0"
          >
            <v-card-title class="headline"
              >TOP 5 AVATARS BY TWEETS</v-card-title
            >
            <v-card-text>
              <pie-chart :round="2" suffix="%" :data="topFiveAvatarsByTweets" />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>
<script>
export default {
  name: "DailyPDF",
  data() {
    return {};
  },
  props: {
    date: {
      type: String,
      default: new Date().toString(),
    },
    activeAccts: {
      type: Number,
      default: 0,
    },
    newAccts: {
      type: Number,
      default: 0,
    },
    totalLikes: {
      type: Number,
      default: 0,
    },
    totalFollowers: {
      type: Number,
      default: 0,
    },
    totalFollowing: {
      type: Number,
      default: 0,
    },
    totalTweets: {
      type: Number,
      default: 0,
    },
    highestGainedFollowers: {
      type: Object,
      default: function () {
        return {
          avatar: "John Doe",
          person: "John Doe",
          count: 0,
        };
      },
    },
    highestLikes: {
      type: Object,
      default: function () {
        return {
          avatar: "John Doe",
          person: "John Doe",
          count: 0,
        };
      },
    },
    totalUnusedAccounts: {
      type: Number,
      default: 0,
    },
    tops: {
      type: Object,
      default: function () {
        return {
          tweets: [
            {
              username: "KairituK",
              person: "Unassigned",
              tweets: 4395,
            },
            {
              username: "ni_witu",
              person: "Unassigned",
              tweets: 3404,
            },
            {
              username: "h_nyinyi",
              tweets: 3158,
              person: "Unassigned",
            },
            {
              username: "MwasikaRomelu",
              person: "Unassigned",
              tweets: 3058,
            },
            {
              username: "NaisulaB",
              person: "Unassigned",
              tweets: 2983,
            },
          ],
          followers: [
            {
              username: "JanetKariuki16",
              person: "Unassigned",
              followers: 2739,
            },
            {
              username: "KabogoSharon",
              person: "Unassigned",
              followers: 1947,
            },
            {
              username: "CalebMutonyi1",
              person: "Unassigned",
              followers: 1610,
            },
            {
              username: "CMoyaka",
              person: "Unassigned",
              followers: 1488,
            },
            {
              username: "Lydia_Korir_C",
              person: "Unassigned",
              followers: 1309,
            },
          ],
        };
      },
    },
    topFiveAvatarsByTweets: {
      type: Object,
      default: function () {
        return {
          JohnDoe: 0.0,
          JohnDoey: 0.1,
          JohnDoek: 0.2,
        };
      },
    },
    topFiveAvatarsByFollowers: {
      type: Object,
      default: function () {
        return {
          JohnDoe: 0.0,
          JohnDoey: 0.1,
          JohnDoek: 0.2,
        };
      },
    },
    suspendedAccounts: {
      type: Object,
      default: function () {
        return {};
      },
    },
  },
};
</script> 