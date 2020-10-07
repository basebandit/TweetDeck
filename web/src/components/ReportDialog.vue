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
            <v-btn dark text @click="closeDialog"> Save </v-btn>
          </v-toolbar-items>
        </v-toolbar>
        <daily-pdf />
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
import DailyPdf from "./DailyPdf";
import { mapGetters } from "vuex";
export default {
  name: "ReportDialog",
  components: { DailyPdf },
  computed: {
    ...mapGetters("report", ["showDialog"]),
  },
  methods: {
    closeDialog() {
      this.$store.dispatch("report/hideDialog");
    },
  },
};
</script>