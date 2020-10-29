<template>
  <v-container>
    <v-row justify="center" class="font-weight-bold">
      <v-col align="left"> Cyclones </v-col>
      <v-spacer></v-spacer>
      <v-col align="right">
        <div>Summary of Twitter Activity</div>
        <div>{{ startDate | formatDate }} - {{ endDate | formatDate }}</div>
      </v-col>
    </v-row>
    <v-divider></v-divider>

    <v-row align-content="space-between">
      <v-col class="text-left">
        <div class="text-decoration-underline">TOTAL ACTIVE ACCOUNTS</div>
        <div class="font-weight-bold">{{ activeAccts.toLocaleString() }}</div>
      </v-col>
      <v-col class="text-center">
        <div class="text-decoration-underline">TOTAL NEW ACCOUNTS</div>
        <div class="font-weight-bold">{{ newAccts.toLocaleString() }}</div>
      </v-col>
      <v-col class="text-right">
        <div class="text-decoration-underline">TIMEFRAME (PERIOD)</div>
        <div class="font-weight-bold">
          {{ startDate | formatDate }} - {{ endDate | formatDate }}
        </div>
      </v-col>
    </v-row>
    <v-row>
      <v-col class="text-left">
        <div class="text-decoration-underline">TOTAL WEEK'S FOLLOWERS</div>
        <div class="font-weight-bold">
          {{ totalWeeklyFollowers.toLocaleString() }}
        </div>
      </v-col>
      <v-col class="text-center">
        <div class="text-decoration-underline">TOTAL WEEK'S FOLLOWING</div>
        <div class="font-weight-bold">
          {{ totalWeeklyFollowing.toLocaleString() }}
        </div>
      </v-col>
      <v-col class="text-right">
        <div class="text-decoration-underline">TOTAL WEEK'S TWEETS</div>
        <div class="font-weight-bold">
          {{ totalWeeklyTweets.toLocaleString() }}
        </div>
      </v-col>
    </v-row>
    <v-row>
      <v-col class="text-center">
        <div class="text-decoration-underline">TOTAL WEEK'S LIKES</div>
        <div class="font-weight-bold">
          {{ totalWeeklyLikes.toLocaleString() }}
        </div>
      </v-col>
    </v-row>
    <v-row>
      <v-col class="text-center">
        <div class="text-decoration-underline">
          AVATAR WITH THE HIGHEST GAINED FOLLOWERS
        </div>
        <div class="font-weight-bold">
          @{{ highestGainedFollowers.username }} ({{
            highestGainedFollowers.followers.toLocaleString()
          }}) -
          {{ highestGainedFollowers.person || "unassigned" }}
        </div>
      </v-col>
    </v-row>

    <v-row justify="center">
      <v-col align="left">
        <div class="text-decoration-underline">
          AVATAR WITH THE HIGHEST LIKES
        </div>
        <div class="font-weight-bold">
          @{{ highestWeeklyLikes.username }} ({{
            highestWeeklyLikes.likes.toLocaleString()
          }}) -
          {{ highestWeeklyLikes.person || "unassigned" }}
        </div>
      </v-col>
      <v-col align="right">
        <div class="text-decoration-underline">TOTAL SUSPENDED ACCOUNTS</div>
        <div class="font-weight-bold">
          {{ totalSuspendedAccounts.toLocaleString() }}
        </div>
      </v-col>
    </v-row>

    <v-row>
      <v-col class="text-center">
        <div class="text-decoration-underline">TOTAL UNUSED ACCOUNTS</div>
        <div class="font-weight-bold">
          {{ totalUnusedAccounts.toLocaleString() }}
        </div>
      </v-col>
    </v-row>

    <v-row class="my-3">
      <v-col>
        <v-card flat class="pa-3" v-if="tops.tweets">
          <v-card-title class="align-start"
            >TOP 5 WEEKLY AVATARS (BY TWEETS)</v-card-title
          >

          <template>
            <v-simple-table>
              <template v-slot:default>
                <thead>
                  <tr>
                    <th class="text-center">REAL NAME</th>
                    <th class="text-center">AVATAR NAME</th>
                    <th class="text-center">NUMBER OF TWEETS</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in tops.tweets" :key="item.username">
                    <td class="text-center">
                      {{ item.person || "unassigned" }}
                    </td>
                    <td class="text-center">{{ item.username }}</td>
                    <td class="text-center">
                      {{ item.tweets.toLocaleString() }}
                    </td>
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
          flat
          class="pa-3"
          v-if="
            topFiveWeeklyAvatarsByTweets &&
            Object.keys(topFiveWeeklyAvatarsByTweets).length > 0
          "
        >
          <v-card-title class="headline"
            >TOP 5 WEEKLY AVATARS (BY TWEETS)</v-card-title
          >
          <v-card-text>
            <pie-chart
              :round="2"
              suffix="%"
              :data="topFiveWeeklyAvatarsByTweets"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row class="my-3">
      <v-col>
        <v-card flat class="pa-3" v-if="tops.followers">
          <v-card-title class="align-start"
            >TOP 5 WEEKLY AVATARS (BY FOLLOWERS)</v-card-title
          >

          <template>
            <v-simple-table>
              <template v-slot:default>
                <thead>
                  <tr>
                    <th class="text-center">REAL NAME</th>
                    <th class="text-center">AVATAR NAME</th>
                    <th class="text-center">NUMBER OF FOLLOWERS</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in tops.followers" :key="item.username">
                    <td class="text-center">
                      {{ item.person || "unassigned" }}
                    </td>
                    <td class="text-center">{{ item.username }}</td>
                    <td class="text-center">
                      {{ item.followers.toLocaleString() }}
                    </td>
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
          flat
          class="pa-3"
          v-if="
            topFiveWeeklyAvatarsByFollowers &&
            Object.keys(topFiveWeeklyAvatarsByFollowers).length > 0
          "
        >
          <v-card-title class="headline"
            >TOP 5 WEEKLY AVATARS (BY FOLLOWERS)
          </v-card-title>
          <v-card-text>
            <pie-chart
              :round="2"
              suffix="%"
              :data="topFiveWeeklyAvatarsByFollowers"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row class="my-3">
      <v-col>
        <v-card
          flat
          class="pa-3"
          v-if="suspendedAccounts && suspendedAccounts.length > 0"
        >
          <v-card-title class="align-start">SUSPENDED ACCOUNTS</v-card-title>

          <template>
            <v-simple-table>
              <template v-slot:default>
                <thead>
                  <tr>
                    <th class="text-center">#</th>
                    <th class="text-center">REAL NAME</th>
                    <th class="text-center">AVATAR NAME</th>
                    <!-- <th class="text-center">NUMBER OF FOLLOWERS</th> -->
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(item, index) in suspendedAccounts"
                    :key="item.username"
                  >
                    <td class="text-center">{{ index + 1 }}</td>
                    <td class="text-center">
                      {{ item.person || "unassigned" }}
                    </td>
                    <td class="text-center">@{{ item.username }}</td>
                    <!-- <td class="text-center">
                      {{ item.followers.toLocaleString() }}
                    </td> -->
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </template>
        </v-card>
      </v-col>
    </v-row>

    <v-row class="my-3">
      <v-col>
        <v-card flat class="pa-3" v-if="inactiveAccts.length > 0">
          <v-card-title class="align-start"
            >INACTIVE ACCOUNTS (BY WEEKLY TWEETS)</v-card-title
          >

          <template>
            <v-simple-table>
              <template v-slot:default>
                <thead>
                  <tr>
                    <th class="text-center">#</th>
                    <th class="text-center">REAL NAME</th>
                    <th class="text-center">AVATAR NAME</th>
                    <!-- <th class="text-center">NUMBER OF FOLLOWERS</th> -->
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(item, index) in inactiveAccts"
                    :key="item.username"
                  >
                    <td class="text-center">{{ index + 1 }}</td>
                    <td class="text-center">
                      {{ item.person || "unassigned" }}
                    </td>
                    <td class="text-center">@{{ item.username }}</td>
                    <!-- <td class="text-center">
                      {{ item.followers.toLocaleString() }}
                    </td> -->
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </template>
        </v-card>
      </v-col>
    </v-row>

    <v-row class="my-3">
      <v-col>
        <v-card flat class="pa-3" v-if="unusedAccounts.length > 0">
          <v-card-title class="align-start"
            >UNUSED ACCOUNTS (BY WEEKLY TWEETS)</v-card-title
          >

          <template>
            <v-simple-table>
              <template v-slot:default>
                <thead>
                  <tr>
                    <th class="text-center">#</th>
                    <th class="text-center">REAL NAME</th>
                    <th class="text-center">AVATAR NAME</th>
                    <!-- <th class="text-center">NUMBER OF FOLLOWERS</th> -->
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(item, index) in unusedAccounts"
                    :key="item.username"
                  >
                    <td class="text-center">{{ index + 1 }}</td>
                    <td class="text-center">
                      {{ item.person || "unassigned" }}
                    </td>
                    <td class="text-center">@{{ item.username }}</td>
                    <!-- <td class="text-center">
                      {{ item.followers.toLocaleString() }}
                    </td> -->
                  </tr>
                </tbody>
              </template>
            </v-simple-table>
          </template>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
<script>
export default {
  name: "WeeklyPDF",
  data() {
    return {};
  },
  props: {
    startDate: {
      type: String,
      default: new Date("2020-10-21").toString(),
    },
    endDate: {
      type: String,
      default: new Date().toString(),
    },
    activeAccts: {
      type: Number,
      default: 1,
    },
    newAccts: {
      type: Number,
      default: 1,
    },
    totalWeeklyLikes: {
      type: Number,
      default: 1,
    },
    totalWeeklyFollowers: {
      type: Number,
      default: 1,
    },
    totalWeeklyFollowing: {
      type: Number,
      default: 1,
    },
    totalWeeklyTweets: {
      type: Number,
      default: 1,
    },
    highestGainedFollowers: {
      type: Object,
      default: function () {
        return {
          username: "John Doe",
          person: "John Doe",
          followers: 1,
        };
      },
    },
    highestWeeklyLikes: {
      type: Object,
      default: function () {
        return {
          username: "John Doe",
          person: "John Doe",
          likes: 0,
        };
      },
    },
    totalSuspendedAccounts: {
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
    topFiveWeeklyAvatarsByTweets: {
      type: Object,
      default: function () {
        return {
          JohnDoe: 0.0,
          JohnDoey: 0.1,
          JohnDoek: 0.2,
        };
      },
    },
    topFiveWeeklyAvatarsByFollowers: {
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
      type: Array,
      default: function () {
        return [
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
        ];
      },
    },
    unusedAccounts: {
      type: Array,
      default: function () {
        return [
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
        ];
      },
    },
    totalUnusedAccounts: {
      type: Number,
      default: 0,
    },
    inactiveAccts: {
      type: Array,
      default: function () {
        return [
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
        ];
      },
    },
  },
};
</script> 