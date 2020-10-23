<template>
  <v-app>
    <v-main>
      <v-row justify="center">
        <v-card class="ma-3">
          <daily-pdf
            :date="dailyReportStats.date"
            :activeAccts="dailyReportStats.activeAccts"
            :totalLikes="dailyReportStats.totalLikes"
            :totalFollowing="dailyReportStats.totalFollowing"
            :totalFollowers="dailyReportStats.totalFollowers"
            :totalTweets="dailyReportStats.totalTweets"
            :totalSuspendedAccounts="dailyReportStats.totalSuspendedAccounts"
            :topFiveAvatarsByTweets="dailyReportStats.topFiveAvatarsByTweets"
            :newAccts="dailyReportStats.newAccounts"
            :topFiveAvatarsByFollowing="
              dailyReportStats.topFiveAvatarsByFollowing
            "
            :topFiveAvatarsByFollowers="
              dailyReportStats.topFiveAvatarsByFollowers
            "
            :highestLikes="dailyReportStats.highestLikes"
            :topFiveAvatarsByLikes="dailyReportStats.topFiveAvatarsByLikes"
            :tops="dailyReportStats.tops"
          />
        </v-card>
      </v-row>
    </v-main>
  </v-app>
</template>
<script>
import DailyPdf from "@/components/DailyPdf";
import print from "vue-print-nb";
import { mapGetters } from "vuex";
export default {
  name: "ReportDialog",
  components: { DailyPdf },
  directives: {
    print,
  },
  computed: {
    ...mapGetters("report", ["showDialog", "dailyReportStats"]),
  },
};
</script>