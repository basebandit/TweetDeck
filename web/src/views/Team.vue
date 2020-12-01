<template>
  <div class="team">
    <v-container class="my-5">
      <v-row>
        <v-col>
          <v-breadcrumbs :items="breadcrumbs">
            <template v-slot:item="{ item }">
              <v-breadcrumbs-item :href="item.href" :disabled="item.disabled">
                {{ item.text.toUpperCase() }}
              </v-breadcrumbs-item>
            </template>
            <template v-slot:divider>
              <v-icon>mdi-chevron-right</v-icon>
            </template>
          </v-breadcrumbs>
        </v-col>
      </v-row>
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
                  <v-btn color="blue darken-1" text @click="dialog = false"
                    >Cancel</v-btn
                  >
                  <v-btn color="blue darken-1" text @click="addMember"
                    >Save</v-btn
                  >
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
        <v-data-table
          :headers="headers"
          :items="team"
          :search="search"
          @click:row="handleClick"
        >
          <template v-slot:item.firstname="props">
            <v-edit-dialog
              :return-value.sync="props.item.firstname"
              @save="saveFirstname(props.item.id)"
              @cancel="cancel"
              @open="open"
              @close="close"
            >
              {{ props.item.firstname }}
              <template v-slot:input>
                <v-text-field
                  v-model="props.item.firstname"
                  @change="editFirstname"
                  :rules="[max25chars]"
                  label="Edit"
                  single-line
                  counter
                ></v-text-field>
              </template>
            </v-edit-dialog>
          </template>
          <template v-slot:item.lastname="props">
            <v-edit-dialog
              :return-value.sync="props.item.lastname"
              @save="saveLastname(props.item.id)"
              @cancel="cancel"
              @open="open"
              @close="close"
            >
              {{ props.item.lastname }}
              <template v-slot:input>
                <v-text-field
                  v-model="props.item.lastname"
                  @change="editLastname"
                  :rules="[max25chars]"
                  label="Edit"
                  single-line
                  counter
                ></v-text-field>
              </template>
            </v-edit-dialog>
          </template>
          <template v-slot:item.avatars="{ item }">
            <v-chip :color="getColor(item.avatars)" dark>{{
              item.avatars
            }}</v-chip>
          </template>
          <template v-slot:item.createdAt="{ item }">{{
            item.createdAt | formatDate
          }}</template>
        </v-data-table>
        <v-snackbar v-model="snack" :timeout="3000" :color="snackColor">
          {{ snackText }}
          <template v-slot:action="{ attrs }">
            <v-btn v-bind="attrs" text @click="snack = false">Close</v-btn>
          </template>
        </v-snackbar>
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
      firstname: "",
      lastname: "",
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
      breadcrumbs: [
        {
          text: "Dashboard",
          disabled: false,
          href: "/dashboard",
        },
        {
          text: "Team",
          disabled: true,
          href: "/team",
        },
      ],
      dialog: false,
      valid: true,
      snack: false,
      snackColor: "",
      snackText: "",
      max25chars: (v) => v.length <= 25 || "Input too long!",
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
    editFirstname(v) {
      this.firstname = v;
    },
    editLastname(v) {
      this.lastname = v;
    },
    getColor(avatars) {
      if (avatars > 10) return "red";
      else if (avatars > 20) return "orange";
      else return "green";
    },
    handleClick(item) {
      /**eslint-disable */
      console.log(item.id);

      this.$router.push({ path: `/team/member/${item.id}` }); // -> /user/123
    },
    addMember() {
      if (this.$refs.form.validate()) {
        let person = {
          firstname: this.newMember.firstname,
          lastname: this.newMember.lastname,
          email: this.newMember.email,
        };
        this.$store.dispatch("people/addPerson", {
          token: this.token,
          person,
          router: this.$router,
          refs: this.$refs,
        });
      }
    },
    saveFirstname(id) {
      /**eslint-disable */
      console.log("Firstname", this.firstname);
      this.snack = true;
      this.snackColor = "success";
      this.snackText = "Data saved";
      this.$store.dispatch("people/updateFirstname", {
        token: this.token,
        id,
        firstname: this.firstname,
        router: this.$router,
      });
    },
    saveLastname(id) {
      /**eslint-disable */
      console.log("Lastname", this.lastname);
      this.snack = true;
      this.snackColor = "success";
      this.snackText = "Data saved";
      this.$store.dispatch("people/updateLastname", {
        token: this.token,
        id,
        lastname: this.lastname,
        router: this.$router,
      });
    },
    cancel() {
      this.snack = true;
      this.snackColor = "error";
      this.snackText = "Canceled";
    },
    open() {
      this.snack = true;
      this.snackColor = "info";
      this.snackText = "Dialog opened";
    },
    close() {
      console.log("Dialog closed");
    },
  },
};
</script>
