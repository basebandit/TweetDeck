<template>
  <v-row justify="center">
    <v-dialog
      v-model="showDialog"
      fullscreen
      hide-overlay
      transition="dialog-bottom-transition"
    >
      <v-card>
        <v-toolbar dark color="primary">
          <v-btn icon dark @click="closeDialog">
            <v-icon>mdi-close</v-icon>
          </v-btn>
          <v-toolbar-title>Report</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-toolbar-items>
            <v-btn dark text @click="exportPDF"> Save To PDF </v-btn>
          </v-toolbar-items>
        </v-toolbar>
        <daily-pdf
          :date="dailyReportStats.date"
          :activeAccts="dailyReportStats.activeAccts"
          :totalLikes="dailyReportStats.totalLikes"
          :totalFollowing="dailyReportStats.totalFollowing"
          :totalFollowers="dailyReportStats.totalFollowers"
          :totalTweets="dailyReportStats.totalTweets"
          :topFiveAvatarsByTweets="dailyReportStats.topFiveAvatarsByTweets"
          :topFiveAvatarsByFollowing="
            dailyReportStats.topFiveAvatarsByFollowing
          "
          :topFiveAvatarsByFollowers="
            dailyReportStats.topFiveAvatarsByFollowers
          "
          :topFiveAvatarsByLikes="dailyReportStats.topFiveAvatarsByLikes"
          :tops="dailyReportStats.tops"
        />
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
import DailyPdf from "./DailyPdf";
import html2pdf from "html2pdf.js";
import print from "vue-print-nb";
import { mapGetters } from "vuex";
import printJS from "print-js";
export default {
  name: "ReportDialog",
  components: { DailyPdf },
  directives: {
    print,
  },
  // props: {
  //   date: {
  //     type: String,
  //   },
  //   activeAccts: {
  //     type: Number,
  //   },
  //   newAccts: {
  //     type: Number,
  //   },
  //   totalLikes: {
  //     type: Number,
  //   },
  //   totalFollowers: {
  //     type: Number,
  //   },
  //   totalFollowing: {
  //     type: Number,
  //   },
  //   totalTweets: {
  //     type: Number,
  //   },
  // },
  computed: {
    ...mapGetters("report", ["showDialog", "dailyReportStats"]),
  },
  methods: {
    exportPDF() {
      printJS({
        printable: "document",
        showModal: true,
        scanStyles: true,
        modalMessage: "Generating report...",
        type: "html",
        // targetStyles: ["*"],
      });
    },
    exportToPDF() {
      let options_ = {
        margin: 10,
        filename: "report.pdf",
        image: { type: "jpeg", quality: 0.98 },
        html2canvas: {
          scale: 2,
          logging: true,
          dpi: 192,
          letterRendering: false,
        },
        jsPDF: { unit: "mm", format: "a4", orientation: "portrait" },
      };
      // let options = {
      //   margin: 1,
      //   filename: "document.pdf",
      //   image: { type: "jpeg", quality: 0.98 },
      //   html2canvas: { dpi: 192, letterRendering: true },
      //   jsPDF: { unit: "in", format: "", orientation: "landscape" },
      // };
      let element = document.getElementById("document");

      // Create instance of html2pdf class
      let exporter = new html2pdf(element, options_);
      this.closeDialog();
      exporter.then((pdf) => {
        /**eslint-disable */
        console.log(pdf);
      });

      // Download the PDF or...
      // exporter.getPdf(true).then((pdf) => {
      //   // console.log('pdf file downloaded');
      //   pdf.save();
      //
      // });
    },
    closeDialog() {
      this.$store.dispatch("report/hideDialog");
    },
  },
};
</script>