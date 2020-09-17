<template>
  <div class="team">
    <h1>Team</h1>
    <v-container class="my-5">
      <v-card class="pa-3">
        <v-card-title class="my-3">Team Members</v-card-title>
        <v-card-text>
          <v-row class="mx-0">
            <v-dialog v-model="dialog" max-width="500px">
              <template v-slot:activator="{ on, attrs }">
                <v-btn color="success" v-bind="attrs" v-on="on">
                  <v-icon small left>add</v-icon>Add new member
                </v-btn>
              </template>
              <v-card>
                <v-card-title>
                  <span class="headline">Add New Member</span>
                </v-card-title>

                <v-card-text>
                  <v-form ref="form" v-model="valid" lazy-validation>
                    <v-row>
                      <v-col cols="12" sm="8" md="6">
                        <v-text-field
                          v-model="newMember.firstname"
                          label="Firstname"
                          :rules="nameRules"
                          required
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="8" md="6">
                        <v-text-field
                          v-model="newMember.lastname"
                          label="Lastname"
                          :rules="nameRules"
                          required
                        ></v-text-field>
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12">
                        <v-text-field
                          v-model="newMember.email"
                          label="Email"
                          :rules="emailRules"
                          required
                        ></v-text-field>
                      </v-col>
                    </v-row>
                  </v-form>
                </v-card-text>

                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" text @click="dialog = false">Cancel</v-btn>
                  <v-btn color="blue darken-1" text @click="addMember">Save</v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>

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
      dialog: false,
      valid: true,
      nameRules: [
        (v) => !!v || "Name is required",
        (v) => (v && v.length <= 15) || "Name must be less than 15 characters",
      ],
      emailRules: [
        (v) => !!v || "E-mail is required",
        (v) => /.+@.+\..+/.test(v) || "E-mail must be valid",
      ],
      newMember: {
        firstname: "",
        lastname: "",
        email: "",
      },
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
    addMember() {
      this.$refs.form.validate();
      /**eslint-disable */
      console.log(
        this.newMember.firstname,
        this.newMember.lastname,
        this.newMember.Email
      );
    },
  },
};
</script>
