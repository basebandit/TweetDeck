<template>
  <v-app>
    <v-main>
      <v-row justify="center">
        <v-card class="ma-3">
          <weekly-pdf
            :startDate="weeklyStartDate"
            :endDate="weeklyEndDate"
            :newAccts="weeklyStats.newAccts"
            :activeAccts="weeklyStats.activeAccts"
            :totalWeeklyFollowers="weeklyStats.followers"
            :totalWeeklyLikes="weeklyStats.likes"
            :totalWeeklyFollowing="weeklyStats.following"
            :totalWeeklyTweets="weeklyStats.tweets"
            :highestGainedFollowers="weeklyStats.highestGainedByFollowers"
            :highestWeeklyLikes="weeklyStats.highestGainedByLikes"
            :totalSuspendedAccounts="weeklyStats.totalSuspendedAccounts"
            :tops="weeklyStats.tops"
            :topFiveWeeklyAvatarsByFollowers="topFiveAvatarsByFollowers"
            :topFiveWeeklyAvatarsByFollowing="topFiveAvatarsByFollowing"
            :topFiveWeeklyAvatarsByLikes="topFiveAvatarsByLikes"
            :topFiveWeeklyAvatarsByTweets="topFiveAvatarsByTweets"
            :suspendedAccounts="weeklyStats.suspendedAccounts"
            :totalUnusedAccounts="weeklyStats.totalUnusedAccounts"
            :unusedAccounts="weeklyStats.unusedAccounts"
            :inactiveAccts="weeklyStats.inactiveAccts"
          />
        </v-card>
      </v-row>
    </v-main>
  </v-app>
</template>
<script>
import WeeklyPdf from "@/components/WeeklyPdf";
import print from "vue-print-nb";
import { mapGetters } from "vuex";
export default {
  name: "ReportDialog",
  components: { WeeklyPdf },
  directives: {
    print,
  },
  computed: {
    ...mapGetters("report", [
      "showWeeklyDialog",
      "weeklyStartDate",
      "weeklyEndDate",
    ]),
    ...mapGetters("stats", ["weeklyStats"]),
    topFiveAvatarsByFollowers() {
      const obj = {};
      if (this.weeklyStats.tops.followers) {
        this.weeklyStats.tops.followers.forEach((avatar) => {
          let p = (avatar.followers / this.weeklyStats.followers) * 100;
          let key = `${avatar.username} (${p.toFixed()}%)`;
          obj[key] = p;
        });
      }
      return obj;
    },
    topFiveAvatarsByLikes() {
      const obj = {};
      if (this.weeklyStats.tops.likes) {
        this.weeklyStats.tops.likes.forEach((avatar) => {
          let p = (avatar.likes / this.weeklyStats.likes) * 100;
          let key = `${avatar.username} (${p.toFixed()}%)`;
          obj[key] = p;
        });
      }
      return obj;
    },
    topFiveAvatarsByFollowing() {
      const obj = {};
      if (this.weeklyStats.tops.following) {
        this.weeklyStats.tops.following.forEach((avatar) => {
          let p = (avatar.following / this.weeklyStats.following) * 100;
          let key = `${avatar.username} (${p.toFixed()}%)`;
          obj[key] = p;
        });
      }
      return obj;
    },
    topFiveAvatarsByTweets() {
      const obj = {};
      if (this.weeklyStats.tops.tweets) {
        this.weeklyStats.tops.tweets.forEach((avatar) => {
          let p = (avatar.tweets / this.weeklyStats.tweets) * 100;
          let key = `${avatar.username} (${p.toFixed()}%)`;
          obj[key] = p;
        });
      }
      return obj;
    },
  },
};
</script>