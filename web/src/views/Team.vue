<template>
  <div class="team">
    <h1>Team</h1>
    <v-container class="my-5">
      <v-card class="pa-3">
        <v-card-title class="my-3">Team Members</v-card-title>
        <v-card-text>
          <v-row class="mx-0">
            <v-btn color="success">
              <v-icon small left>add</v-icon>Add new member
            </v-btn>
            <v-spacer></v-spacer>
            <v-text-field
              v-model="search"
              append-icon="mdi-magnify"
              label="Search"
              single-line
              hide-details
            ></v-text-field>
          </v-row>
        </v-card-text>
        <v-data-table :headers="headers" :items="team" :search="search" @click:row="handleClick"></v-data-table>
      </v-card>
    </v-container>
  </div>
</template>
<script>
import { mapGetters } from "vuex";
export default {
  data() {
    return {
      search: "",
      headers: [
        {
          text: "Firstname",
          sortable: false,
          align: "start",
          value: "firstname",
        },
        { text: "Lastname", sortable: false, value: "lastname" },
        { text: "Avatars", value: "avatars" },
        { text: "CreatedAt", value: "createdAt" },
        { text: "Followers", value: "followers" },
        { text: "Following", value: "following" },
        { text: "Tweets", value: "tweets" },
      ],
    };
  },
  mounted() {
    this.$store.dispatch("people/getPeople", { token: this.token });
  },
  computed: {
    ...mapGetters("people", ["team"]),
    token() {
      return window.localStorage.getItem("user");
    },
  },
  methods: {
    handleClick(item) {
      /**eslint-disable */
      console.log(item.id);

      this.$router.push({ path: `/team/member/${item.id}` }); // -> /user/123
    },
  },
};
</script>
