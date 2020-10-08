<template>
  <div class="home">
    <report-dialog> </report-dialog>
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
          <v-btn
            rounded
            color="primary"
            dark
            class="text-center"
            @click.stop="openReportDialog"
          >
            <v-icon left>mdi-file-pdf</v-icon> Print Report
          </v-btn>
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
import { jsPDF } from "jspdf";
import ReportDialog from "@/components/ReportDialog";
// import html2canvas from "html2canvas";
import { mapGetters } from "vuex";
export default {
  name: "Home",
  components: { ReportDialog },
  data() {
    return {
      search: "",
      today: new Date().toString(),
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
  },
  computed: {
    ...mapGetters("people", ["team"]),
    ...mapGetters("avatars", ["avatars", "suspendedAvatars"]),
    ...mapGetters("stats", ["totals", "tops"]),
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
        { total: this.totals.followers, icon: "mdi-thumb-up", name: "likes" },
      ];
    },
    topFiveAvatarsByFollowers() {
      const obj = {};
      if (this.tops.followers) {
        this.tops.followers.forEach(
          (avatar) =>
            (obj[avatar.username] =
              (avatar.followers / this.totals.followers) * 100)
        );
      }
      return obj;
    },
    topFiveAvatarsByLikes() {
      const obj = {};
      if (this.tops.likes) {
        this.tops.likes.forEach(
          (avatar) =>
            (obj[avatar.username] = (avatar.likes / this.totals.likes) * 100)
        );
      }
      return obj;
    },
    topFiveAvatarsByFollowing() {
      const obj = {};
      if (this.tops.following) {
        this.tops.following.forEach(
          (avatar) =>
            (obj[avatar.username] =
              (avatar.following / this.totals.following) * 100)
        );
      }
      return obj;
    },
    topFiveAvatarsByTweets() {
      const obj = {};
      if (this.tops.tweets) {
        this.tops.tweets.forEach(
          (avatar) =>
            (obj[avatar.username] = (avatar.tweets / this.totals.tweets) * 100)
        );
      }
      return obj;
    },
  },
  methods: {
    openReportDialog() {
      let dailyStats = {
        date: this.today,
        activeAccts: this.accounts.active,
        totalLikes: this.totals.likes,
        totalFollowers: this.totals.followers,
        totalFollowing: this.totals.following,
        totalTweets: this.totals.tweets,
        topFiveAvatarsByFollowers: this.topFiveAvatarsByFollowers,
        topFiveAvatarsByFollowing: this.topFiveAvatarsByFollowing,
        topFiveAvatarsByLikes: this.topFiveAvatarsByLikes,
        topFiveAvatarsByTweets: this.topFiveAvatarsByTweets,
        tops: this.tops,
      };
      this.$store.dispatch("report/showDialog", { dailyStats });
    },
    dailyReport() {
      let doc = new jsPDF("p", "pt", "A4");
      let margins = {
        top: 80,
        bottom: 60,
        left: 40,
        width: 822,
      };

      doc.html(document.getElementById("contentHTML"), {
        callback: function (pdf) {
          pdf.save("a4.pdf");
        },
        x: margins.left,
        y: margins.top,
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