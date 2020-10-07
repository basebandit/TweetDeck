<template>
  <div>
    <v-container>
      <v-row justify="center">
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
          <div>TOTAL TWITTER ACCOUNTS</div>
          <div>{{ activeAccts }}</div>
        </v-col>
        <v-col class="text-center">
          <div>TOTAL NEW ACCOUNTS</div>
          <div>{{ newAccts }}</div>
        </v-col>
        <v-col class="text-right">
          <div>TIMEFRAME (PERIOD)</div>
          <div>{{ date | formatDate }}</div>
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
  },
};
</script> 