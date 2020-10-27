<template>
  <div class="home">
    <v-toolbar
      class="pt-3"
      extended
      flat
      style="position: -webkit-sticky; position: sticky; top: 4rem; z-index: 1"
    >
      <v-row no-gutters class="mt-3 mb-0">
        <v-col align="left">
          <span class="headline grey--text">Dashboard</span>
        </v-col>

        <v-col align="right">
          <v-menu
            v-model="menu"
            :close-on-content-click="false"
            :nudge-width="200"
            offset-x
          >
            <template v-slot:activator="{ on, attrs }">
              <v-btn color="grey" icon v-bind="attrs" v-on="on">
                <v-icon>mdi-dots-vertical</v-icon>
              </v-btn>
            </template>

            <v-list>
              <v-list-item @click.stop="showDailyReport">
                <v-list-item-icon>
                  <v-icon>mdi-file-pdf</v-icon>
                </v-list-item-icon>
                <v-list-item-title>Daily Report</v-list-item-title>
              </v-list-item>
              <v-list-item @click.stop="weeklyReportDialog = true">
                <v-list-item-icon>
                  <v-icon>mdi-file-pdf</v-icon>
                </v-list-item-icon>
                <v-list-item-title>Weekly Report</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
          <v-dialog v-model="weeklyReportDialog" persistent max-width="900">
            <v-card>
              <v-card-title class="headline">
                Choose the date range for weekly report
              </v-card-title>
              <v-card-text>
                <v-row>
                  <v-col cols="12" sm="6">
                    <v-date-picker
                      v-model="dates"
                      :min="minDate"
                      :max="maxDate"
                      range
                    ></v-date-picker>
                  </v-col>
                  <v-col cols="12" sm="6">
                    <v-text-field
                      v-model="dateRangeText"
                      label="Date range"
                      prepend-icon="mdi-calendar"
                      readonly
                    ></v-text-field>
                    <span v-if="dateRangeError != ''">{{
                      dateRangeError
                    }}</span>
                  </v-col>
                </v-row>
              </v-card-text>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="green darken-1"
                  text
                  @click="weeklyReportDialog = false"
                >
                  Cancel
                </v-btn>
                <v-btn
                  color="green darken-1"
                  text
                  @click="prepareWeeklyReport()"
                >
                  Ok
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
          <!-- <v-btn
            rounded
            color="primary"
            dark
            class="text-center"
            @click.stop="showDailyReport"
          >
            <v-icon left>mdi-file-pdf</v-icon> Print Report
          </v-btn> -->
        </v-col>
      </v-row>
    </v-toolbar>
    <v-container class="my-3">
      <div ref="contentHTML" id="contentHTML">
        <v-row class="my-5">
          <v-col justify="space-around">
            <v-card class="pa-2">
              <v-card-title class="subheading grey--text">
                Accounts Overview
              </v-card-title>
              <v-row justify="space-around">
                <v-btn text color="success">
                  <v-icon dark left>mdi-account-check</v-icon>
                  <span class="mx-2">{{ accounts.active }}</span
                  >Active Accounts
                </v-btn>
                <span class="mx-1">|</span>
                <v-btn text color="error">
                  <v-icon>mdi-account-alert</v-icon>
                  <span class="mx-2">{{ accounts.suspended }}</span> Suspended
                  Accounts
                </v-btn>
                <span class="mx-1">|</span>
                <v-btn text color="primary">
                  <v-icon>mdi-account-multiple-check</v-icon>
                  <span class="mx-2">{{ accounts.assigned }}</span> Assigned
                  Active Accounts
                </v-btn>
                <span class="mx-1">|</span>
                <v-btn text disabled>
                  <v-icon>mdi-account-multiple-remove</v-icon>
                  <span class="mx-2">{{ accounts.unassigned }}</span> Unassigned
                  Active Accounts
                </v-btn>
              </v-row>
            </v-card>
          </v-col>
        </v-row>

        <v-row class="mb-3">
          <v-col v-for="(stat, index) in stats" :key="index">
            <v-card :class="`stat ${stat.name}`">
              <v-card-title class="headline">
                <span>{{ stat.total }}</span>
                <v-spacer></v-spacer>
                <span>
                  <v-icon right color="primary">{{ stat.icon }}</v-icon>
                </span>
              </v-card-title>
              <v-card-subtitle>{{ stat.name }}</v-card-subtitle>
            </v-card>
          </v-col>
        </v-row>

        <!-- <v-row class="mb-3">
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
      </v-row> -->

        <v-row class="mb-3">
          <v-col cols="12" lg="8" md="8">
            <v-card flat class="mx-auto" v-if="tops.tweets">
              <v-card-title class="align-start"
                >Top 5 Avatars By Tweets</v-card-title
              >
              <div
                v-for="(avatar, index) in tops.tweets"
                :key="avatar.username"
                :class="`avatar-${index + 1}`"
              >
                <v-row class="px-3">
                  <v-col>
                    <div class="caption grey--text">Person</div>
                    <div>{{ avatar.person || "Unassigned" }}</div>
                  </v-col>
                  <v-col>
                    <div class="caption grey--text">Avatar</div>
                    <div>@{{ avatar.username }}</div>
                  </v-col>
                  <v-col>
                    <div class="caption grey--text">Number of tweets</div>
                    <div>{{ avatar.tweets }}</div>
                  </v-col>
                </v-row>
                <v-divider></v-divider>
              </div>
            </v-card>

            <v-card flat disabled class="mx-auto" v-else>
              <v-card-title class="align-start"
                >Top 5 Avatars By Tweets</v-card-title
              >

              <v-card-text>
                <v-icon large left>mdi-cloud-alert</v-icon>
                <span class="grey--text">Oops not fetched yet!</span>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="4" lg="4">
            <v-card
              class="mx-auto pa-3"
              v-if="Object.keys(topFiveAvatarsByTweets).length > 0"
            >
              <v-card-title class="headline"
                >Top 5 Avatars By Tweets</v-card-title
              >
              <v-card-text>
                <pie-chart
                  :round="2"
                  suffix="%"
                  :data="topFiveAvatarsByTweets"
                />
              </v-card-text>
            </v-card>

            <v-card flat disabled class="mx-auto" v-else>
              <v-card-title class="align-start"
                >Top 5 Avatars By Tweets</v-card-title
              >

              <v-card-text>
                <v-icon large left>mdi-cloud-alert</v-icon>
                <span class="grey--text">Oops not fetched yet!</span>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <v-row class="mb-3">
          <v-col cols="12" lg="8" md="8">
            <v-card flat class="mx-auto" v-if="tops.following">
              <v-card-title class="align-start"
                >Top 5 Avatars By Following</v-card-title
              >
              <div
                v-for="(avatar, index) in tops.following"
                :key="avatar.username"
                :class="`avatar-${index + 1}`"
              >
                <v-row class="px-3">
                  <v-col>
                    <div class="caption grey--text">Person</div>
                    <div>{{ avatar.person || "Unassigned" }}</div>
                  </v-col>
                  <v-col>
                    <div class="caption grey--text">Avatar</div>
                    <div>@{{ avatar.username }}</div>
                  </v-col>
                  <v-col>
                    <div class="caption grey--text">Number of followings</div>
                    <div>{{ avatar.following }}</div>
                  </v-col>
                </v-row>
                <v-divider></v-divider>
              </div>
            </v-card>

            <v-card flat disabled class="mx-auto" v-else>
              <v-card-title class="align-start"
                >Top 5 Avatars By Following</v-card-title
              >

              <v-card-text>
                <v-icon large left>mdi-cloud-alert</v-icon>
                <span class="grey--text">Oops not fetched yet!</span>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="4" lg="4">
            <v-card
              class="mx-auto pa-3"
              v-if="Object.keys(topFiveAvatarsByFollowing).length > 0"
            >
              <v-card-title class="headline"
                >Top 5 Avatars By Following</v-card-title
              >
              <v-card-text>
                <pie-chart
                  :round="2"
                  suffix="%"
                  :data="topFiveAvatarsByFollowing"
                />
              </v-card-text>
            </v-card>

            <v-card flat disabled class="mx-auto" v-else>
              <v-card-title class="align-start"
                >Top 5 Avatars By Following</v-card-title
              >

              <v-card-text>
                <v-icon large left>mdi-cloud-alert</v-icon>
                <span class="grey--text">Oops not fetched yet!</span>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- <v-row class="mb-3">
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
      </v-row> -->
        <v-row class="mb-3">
          <v-col cols="12" lg="8" md="8">
            <v-card flat class="mx-auto" v-if="tops.followers">
              <v-card-title class="align-start"
                >Top 5 Avatars By Followers</v-card-title
              >
              <div
                v-for="(avatar, index) in tops.followers"
                :key="avatar.username"
                :class="`avatar-${index + 1}`"
              >
                <v-row class="px-3">
                  <v-col>
                    <div class="caption grey--text">Person</div>
                    <div>{{ avatar.person || "Unassigned" }}</div>
                  </v-col>
                  <v-col>
                    <div class="caption grey--text">Avatar</div>
                    <div>@{{ avatar.username }}</div>
                  </v-col>
                  <v-col>
                    <div class="caption grey--text">Number of followers</div>
                    <div>{{ avatar.followers }}</div>
                  </v-col>
                </v-row>
                <v-divider></v-divider>
              </div>
            </v-card>

            <v-card flat disabled class="mx-auto" v-else>
              <v-card-title class="align-start"
                >Top 5 Avatars By Followers</v-card-title
              >

              <v-card-text>
                <v-icon large left>mdi-cloud-alert</v-icon>
                <span class="grey--text">Oops not fetched yet!</span>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="4" lg="4">
            <v-card
              class="pa-3"
              v-if="Object.keys(topFiveAvatarsByFollowers).length > 0"
            >
              <v-card-title class="headline"
                >Top 5 Avatars By Followers</v-card-title
              >
              <v-card-text>
                <pie-chart
                  :round="2"
                  suffix="%"
                  :data="topFiveAvatarsByFollowers"
                />
              </v-card-text>
            </v-card>

            <v-card flat disabled class="mx-auto" v-else>
              <v-card-title class="align-start"
                >Top 5 Avatars By Followers</v-card-title
              >

              <v-card-text>
                <v-icon large left>mdi-cloud-alert</v-icon>
                <span class="grey--text">Oops not fetched yet!</span>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <v-row class="mb-3">
          <v-col cols="12" lg="8" md="8">
            <v-card flat class="mx-auto" v-if="tops.likes">
              <v-card-title class="align-start"
                >Top 5 Avatars By Likes</v-card-title
              >
              <div
                v-for="(avatar, index) in tops.likes"
                :key="avatar.username"
                :class="`avatar-${index + 1}`"
              >
                <v-row class="px-3">
                  <v-col>
                    <div class="caption grey--text">Person</div>
                    <div>{{ avatar.person || "Unassigned" }}</div>
                  </v-col>
                  <v-col>
                    <div class="caption grey--text">Avatar</div>
                    <div>@{{ avatar.username }}</div>
                  </v-col>
                  <v-col>
                    <div class="caption grey--text">Number of likes</div>
                    <div>{{ avatar.likes }}</div>
                  </v-col>
                </v-row>
                <v-divider></v-divider>
              </div>
            </v-card>

            <v-card flat disabled class="mx-auto" v-else>
              <v-card-title class="align-start"
                >Top 5 Avatars By Likes</v-card-title
              >

              <v-card-text>
                <v-icon large left>mdi-cloud-alert</v-icon>
                <span class="grey--text">Oops not fetched yet!</span>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="4" lg="4">
            <v-card
              class="pa-3"
              v-if="Object.keys(topFiveAvatarsByLikes).length > 0"
            >
              <v-card-title class="headline"
                >Top 5 Avatars By Likes</v-card-title
              >
              <v-card-text>
                <pie-chart
                  :round="2"
                  suffix="%"
                  :data="topFiveAvatarsByLikes"
                />
              </v-card-text>
            </v-card>

            <v-card flat disabled class="mx-auto" v-else>
              <v-card-title class="align-start"
                >Top 5 Avatars By Likes</v-card-title
              >

              <v-card-text>
                <v-icon large left>mdi-cloud-alert</v-icon>
                <span class="grey--text">Oops not fetched yet!</span>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <v-row>
          <v-col>
            <v-card class="pa-5">
              <v-card-title>
                Suspended Accounts
                <v-spacer></v-spacer>
                <v-text-field
                  v-model="search"
                  append-icon="mdi-magnify"
                  label="Search"
                  single-line
                  hide-details
                ></v-text-field>
              </v-card-title>
              <v-data-table
                :headers="headers"
                :items="suspended"
                :search="search"
              ></v-data-table>
            </v-card>
          </v-col>
        </v-row>
      </div>
    </v-container>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
export default {
  name: "Home",
  data() {
    return {
      menu: false,
      search: "",
      dates: [],
      today: new Date().toString(),
      weeklyReportDialog: false,
      dateRangeError: "",
      headers: [
        {
          text: "Avatar",
          align: "start",
          sortable: true,
          value: "avatar",
        },
        { text: "Person", value: "person" },
      ],
    };
  },
  mounted() {
    this.$store.dispatch("people/getPeople", { token: this.token });
    this.$store.dispatch("stats/getTotals", { token: this.token });
    this.$store.dispatch("stats/getTops", { token: this.token });
    this.$store.dispatch("avatars/getAvatars", { token: this.token });
    this.$store.dispatch("avatars/getSuspendedAvatars", { token: this.token });
    this.$store.dispatch("report/getDateRange", { token: this.token });
  },
  computed: {
    ...mapGetters("people", ["team"]),
    ...mapGetters("avatars", ["avatars", "suspendedAvatars"]),
    ...mapGetters("stats", ["totals", "tops"]),
    ...mapGetters("report", ["minDate", "maxDate"]),
    dateRangeText() {
      return this.dates.join(" to ");
    },

    token() {
      return window.localStorage.getItem("user");
    },
    accounts() {
      // let suspended = 0;
      let active = 0;
      let assigned = 0;
      let unassigned = 0;
      if (this.avatars.length > 0) {
        this.avatars.forEach((avatar) => {
          if (Object.keys(avatar).length > 6) {
            active++;
            //we also calculate assigned and unassigned here. Since you can only assign an account if
            //it is not already suspended.
            if (avatar.assigned === 1) {
              assigned++;
            } else if (avatar.assigned === 0) {
              unassigned++;
            }
          }
        });
      }

      return {
        suspended: this.suspendedAvatars.length,
        active: active,
        assigned: assigned,
        unassigned: unassigned,
      };
    },
    suspended() {
      const accts = [];
      if (this.suspendedAvatars.length > 0) {
        this.suspendedAvatars.forEach((avatar) => {
          accts.push({
            avatar: `@${avatar.username}`,
            person: avatar.person || "Unassigned",
          });
        });
      }
      return accts;
    },
    stats() {
      return [
        {
          total: this.totals.avatars + this.suspendedAvatars.length,
          icon: "mdi-account-supervisor",
          name: "avatars",
        },
        { total: this.totals.tweets, icon: "mdi-twitter", name: "tweets" },
        {
          total: this.totals.followers,
          icon: "mdi-account-arrow-left",
          name: "followers",
        },
        {
          total: this.totals.following,
          icon: "mdi-account-arrow-right",
          name: "following",
        },
        { total: this.totals.likes, icon: "mdi-thumb-up", name: "likes" },
      ];
    },
    topFiveAvatarsByFollowers() {
      const obj = {};
      if (this.tops.followers) {
        this.tops.followers.forEach((avatar) => {
          let p = (avatar.followers / this.totals.followers) * 100;
          let key = `${avatar.username} (${p.toFixed()}%)`;
          obj[key] = p;
        });
      }
      return obj;
    },
    topFiveAvatarsByLikes() {
      const obj = {};
      if (this.tops.likes) {
        this.tops.likes.forEach((avatar) => {
          let p = (avatar.likes / this.totals.likes) * 100;
          let key = `${avatar.username} (${p.toFixed()}%)`;
          obj[key] = p;
        });
      }
      return obj;
    },
    topFiveAvatarsByFollowing() {
      const obj = {};
      if (this.tops.following) {
        this.tops.following.forEach((avatar) => {
          let p = (avatar.following / this.totals.following) * 100;
          let key = `${avatar.username} (${p.toFixed()}%)`;
          obj[key] = p;
        });
      }
      return obj;
    },
    topFiveAvatarsByTweets() {
      const obj = {};
      if (this.tops.tweets) {
        this.tops.tweets.forEach((avatar) => {
          let p = (avatar.tweets / this.totals.tweets) * 100;
          let key = `${avatar.username} (${p.toFixed()}%)`;
          obj[key] = p;
        });
      }
      return obj;
    },
  },
  methods: {
    prepareWeeklyReport() {
      //Lets send the date range for the weekly report here
      /**eslint-disable */
      console.log(this.dates);
      if (this.dates.length > 0) {
        if (new Date(this.dates[0]) < new Date(this.dates[1])) {
          this.weeklyReportDialog = false;
          //Set the start and end dates for display in the report pdf
          this.$store.dispatch("report/weeklyReportDateRange", {
            startDate: this.dates[0],
            endDate: this.dates[1],
          });

          this.$store.dispatch("stats/getWeeklyStats", {
            token: this.token,
            start: this.dates[0],
            end: this.dates[1],
            router: this.$router,
          });
        } else {
          /**eslint-disable */
          console.error("Start Date is greater than End Date");
          this.dateRangeError = "Start Date cannot be  greater than End Date";
        }
      }
    },
    showDailyReport() {
      let dailyStats = {
        date: this.today,
        activeAccts: this.accounts.active,
        totalLikes: this.totals.likes,
        totalFollowers: this.totals.followers,
        totalFollowing: this.totals.following,
        totalTweets: this.totals.tweets,
        newAccounts: this.totals.newAccounts,
        highestLikes: this.tops.likes[0],
        totalSuspendedAccounts: this.suspendedAvatars.length,
        topFiveAvatarsByFollowers: this.topFiveAvatarsByFollowers,
        topFiveAvatarsByFollowing: this.topFiveAvatarsByFollowing,
        topFiveAvatarsByLikes: this.topFiveAvatarsByLikes,
        topFiveAvatarsByTweets: this.topFiveAvatarsByTweets,
        tops: this.tops,
      };
      this.$store.dispatch("report/showDialog", {
        dailyStats,
        router: this.$router,
      });
    },
  },
};
</script>
<style scoped>
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