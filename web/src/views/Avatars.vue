<template>
  <div class="avatars">
    <v-container class="my-5">
      <v-data-iterator
        :items="avatars"
        item-key="username"
        :items-per-page.sync="avatarsPerPage"
        :page="page"
        :search="search"
        :sort-desc="sortDesc"
        hide-default-footer
        v-model="selected"
      >
        <template v-slot:header>
          <v-toolbar
            class="pa-5"
            extended
            flat
            style="
              position: -webkit-sticky;
              position: sticky;
              top: 4rem;
              z-index: 1;
            "
          >
            <v-text-field
              v-model="search"
              clearable
              hide-details
              prepend-icon="search"
              label="Search"
            ></v-text-field>
            <template v-if="$vuetify.breakpoint.mdAndUp">
              <v-spacer></v-spacer>

              <v-spacer></v-spacer>

              <v-dialog
                v-model="assignDialog"
                persistent
                v-if="showAssign"
                max-width="600px"
              >
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                    class="ma-2"
                    small
                    color="success"
                    v-on="on"
                    v-bind="attrs"
                    slot="activator"
                  >
                    <v-icon small left>mdi-account-multiple-plus</v-icon>Assign
                  </v-btn>
                </template>

                <v-card>
                  <v-card-title>
                    <span class="headline">Assign Selected Avatar(s)</span>
                  </v-card-title>
                  <v-card-text>
                    <v-container>
                      <v-row>
                        <v-col cols="12">
                          <v-autocomplete
                            :items="members"
                            v-model="assignee"
                            label="Member to assign"
                          ></v-autocomplete>
                        </v-col>
                      </v-row>
                    </v-container>
                    <small>*indicates required field</small>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn
                      color="blue darken-1"
                      text
                      @click="assignDialog = false"
                      >Close</v-btn
                    >
                    <v-btn
                      color="blue darken-1"
                      text
                      :loading="assigning"
                      @click="assignAvatar"
                      >Assign</v-btn
                    >
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <v-dialog
                v-model="reassignDialog"
                persistent
                v-if="showReassign"
                max-width="600px"
              >
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                    class="ma-2"
                    small
                    color="success"
                    v-on="on"
                    v-bind="attrs"
                    slot="activator"
                    :loading="assigning"
                  >
                    <v-icon small left>mdi-account-multiple-plus</v-icon
                    >Reassign
                  </v-btn>
                </template>

                <v-card>
                  <v-card-title>
                    <span class="headline">Reassign Selected Avatar(s)</span>
                  </v-card-title>
                  <v-card-text>
                    <v-container>
                      <v-row>
                        <v-col cols="12">
                          <v-autocomplete
                            :items="members"
                            v-model="reassignee"
                            label="Member to reassign"
                          ></v-autocomplete>
                        </v-col>
                      </v-row>
                    </v-container>
                    <small>*indicates required field</small>
                  </v-card-text>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn
                      color="blue darken-1"
                      text
                      @click="reassignDialog = false"
                      >Close</v-btn
                    >
                    <v-btn color="blue darken-1" text @click="reassignAvatar"
                      >Reassign</v-btn
                    >
                  </v-card-actions>
                </v-card>
              </v-dialog>

              <div class="h6 text--primary">
                <v-tooltip top>
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      text
                      class="ma-2"
                      v-bind="attrs"
                      v-on="on"
                      color="success"
                    >
                      {{ stats.active }}
                      <v-icon dark right>mdi-account-check</v-icon>
                    </v-btn>
                  </template>
                  <span>{{ stats.active }} Active Accounts</span>
                </v-tooltip>
                <span class="mx-1">|</span>
                <v-tooltip top>
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      text
                      class="ma-2"
                      v-bind="attrs"
                      v-on="on"
                      color="indigo"
                    >
                      {{ stats.assigned }}
                      <v-icon dark right>mdi-account-multiple-check</v-icon>
                    </v-btn>
                  </template>
                  <span>{{ stats.assigned }} Assigned Active Accounts</span>
                </v-tooltip>
                <span class="mx-1">|</span>
                <v-tooltip top>
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      text
                      class="ma-2"
                      v-bind="attrs"
                      v-on="on"
                      color="grey"
                    >
                      {{ stats.unassigned }}
                      <v-icon dark right>mdi-account-multiple-remove</v-icon>
                    </v-btn>
                  </template>
                  <span>{{ stats.unassigned }} Unassigned Active Accounts</span>
                </v-tooltip>
                <span class="mx-1">|</span>

                <v-tooltip top>
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn
                      text
                      class="ma-2"
                      v-bind="attrs"
                      v-on="on"
                      color="error"
                    >
                      {{ stats.suspended }}
                      <v-icon dark right>mdi-account-alert</v-icon>
                    </v-btn>
                  </template>
                  <span>{{ stats.suspended }} Suspended Accounts</span>
                </v-tooltip>
              </div>

              <v-tooltip top>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                    small
                    text
                    color="grey"
                    v-on="on"
                    v-bind="attrs"
                    @click="sortAscendingBy('tweets')"
                    slot="activator"
                  >
                    <v-icon small left>arrow_upward</v-icon>
                  </v-btn>
                </template>
                <span>Sort by ascending tweets</span>
              </v-tooltip>
              <v-tooltip top>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                    small
                    text
                    color="grey"
                    v-on="on"
                    v-bind="attrs"
                    @click="sortDescendingBy('tweets')"
                    slot="activator"
                  >
                    <v-icon small left>arrow_downward</v-icon>
                  </v-btn>
                </template>
                <span>Sort by descending tweets</span>
              </v-tooltip>
            </template>
          </v-toolbar>
        </template>

        <template v-slot:default="props">
          <v-container>
            <v-row>
              <v-col
                cols="12"
                sm="6"
                md="4"
                lg="3"
                v-for="item in props.items"
                :key="item.username"
              >
                <v-card
                  :color="props.isSelected(item) ? 'secondary' : ''"
                  class="ma-4"
                  align="center"
                  :disabled="isAssigning"
                  @click.stop="reassign(props, item)"
                  v-if="item.assigned === 1"
                >
                  <v-responsive class="pt-4">
                    <v-avatar size="100" class="grey lighten-2">
                      <img :src="item.profileImageURL" />
                    </v-avatar>
                  </v-responsive>

                  <v-card-subtitle>
                    <div class="grey--text">@{{ item.username }}</div>
                    <div class="grey--text text--darken-4">
                      {{ item.bio }}
                    </div>
                  </v-card-subtitle>
                  <v-card-text>
                    <div class="grey--text text--darken-4">
                      {{ item.following }}
                      <span class="grey--text text-caption">following</span>
                      {{ item.followers }}
                      <span class="grey--text text-caption">followers</span>
                    </div>
                    <div>
                      <v-chip class="ma-2" color="indigo" text-color="white">
                        <v-icon left>mdi-account-multiple-check</v-icon>assigned
                      </v-chip>
                    </div>
                  </v-card-text>
                  <v-card-actions>
                    <v-row class="ma-3 text-sm grey--text">
                      <div>
                        <v-icon small left>mdi-twitter</v-icon>
                        <span class="grey--text text--darken-4">{{
                          item.tweets
                        }}</span>
                        tweets
                      </div>
                      <v-spacer></v-spacer>
                      <div class="text-subtitle-1">
                        <v-icon small left color="secondary"
                          >mdi-calendar-month-outline</v-icon
                        >
                        <span class="grey--text text--darken-4">Joined</span>
                        {{ item.joinDate | formatDate }}
                      </div>
                    </v-row>
                  </v-card-actions>
                </v-card>

                <v-card
                  :color="props.isSelected(item) ? 'secondary' : ''"
                  class="ma-4"
                  align="center"
                  :disabled="reassigning"
                  @click.stop="assign(props, item)"
                  v-else-if="item.assigned === 0"
                >
                  <v-responsive class="pt-4">
                    <v-avatar size="100" class="grey lighten-2">
                      <img :src="item.profileImageURL" />
                    </v-avatar>
                  </v-responsive>

                  <v-card-subtitle>
                    <div class="grey--text">@{{ item.username }}</div>
                    <div class="grey--text text--darken-4">
                      {{ item.bio }}
                    </div>
                  </v-card-subtitle>
                  <v-card-text>
                    <div class="grey--text text--darken-4">
                      {{ item.following }}
                      <span class="grey--text text-caption">following</span>
                      {{ item.followers }}
                      <span class="grey--text text-caption">followers</span>
                    </div>

                    <div>
                      <v-chip class="ma-2" text-color="grey">
                        <v-icon left>mdi-account-multiple-remove</v-icon
                        >unasssigned
                      </v-chip>
                    </div>
                  </v-card-text>
                  <v-card-actions>
                    <v-row class="ma-3 text-sm grey--text">
                      <div>
                        <v-icon small left>mdi-twitter</v-icon>
                        <span class="grey--text text--darken-4">{{
                          item.tweets
                        }}</span>
                        tweets
                      </div>

                      <v-spacer></v-spacer>
                      <div class="text-subtitle-1">
                        <v-icon small left color="secondary"
                          >mdi-calendar-month-outline</v-icon
                        >
                        <span class="grey--text text--darken-4">Joined</span>
                        {{ item.joinDate | formatDate }}
                      </div>
                    </v-row>
                  </v-card-actions>
                </v-card>
              </v-col>
            </v-row>
          </v-container>
        </template>

        <template v-slot:footer>
          <v-row class="mt-2" align="center" justify="center">
            <span class="grey--text">Avatars per page</span>
            <v-menu offset-y>
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  dark
                  text
                  color="primary"
                  class="ml-2"
                  v-bind="attrs"
                  v-on="on"
                >
                  {{ avatarsPerPage }}
                  <v-icon>mdi-chevron-down</v-icon>
                </v-btn>
              </template>
              <v-list>
                <v-list-item
                  v-for="(number, index) in avatarsPerPageArray"
                  :key="index"
                  @click="updateAvatarsPerPage(number)"
                >
                  <v-list-item-title>{{ number }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>

            <v-spacer></v-spacer>

            <span class="mr-4 grey--text"
              >Page {{ page }} of {{ numberOfPages }}</span
            >
            <v-btn
              fab
              dark
              color="blue darken-3"
              class="mr-1"
              @click="formerPage"
            >
              <v-icon>mdi-chevron-left</v-icon>
            </v-btn>
            <v-btn
              fab
              dark
              color="blue darken-3"
              class="ml-1"
              @click="nextPage"
            >
              <v-icon>mdi-chevron-right</v-icon>
            </v-btn>
          </v-row>
        </template>
      </v-data-iterator>
    </v-container>
  </div>
</template>
<script>
import { mapGetters } from "vuex";
export default {
  data() {
    return {
      avatarsPerPage: 16,
      avatarsPerPageArray: [16, 32, 64],
      search: "",
      filter: {},
      sortDesc: false,
      page: 1,
      count: 0,
      // sortByHandle: "handle",
      selected: [],
      // chosen: "",
      // selectedAvatarsIDs: [],
      // selectedItems: [],
      assignDialog: false,
      reassignDialog: false,
      assignCard: false,
      reassignCard: false,
      assignee: "",
      assigneeID: [],
      reassignee: "",
      reassigneeID: [],
    };
  },
  mounted() {
    this.$store.dispatch("avatars/getAvatars", { token: this.token });
    this.$store.dispatch("people/getPeople", { token: this.token });
  },
  computed: {
    ...mapGetters("avatars", ["avatars", "fetching", "assigning"]),
    ...mapGetters("people", ["team"]),
    isAssigning() {
      if (this.selected.length > 0 && this.assignCard) {
        return true;
      }
      return false;
    },
    reassigning() {
      if (this.selected.length > 0 && this.reassignCard) {
        return true;
      }
      return false;
    },
    token() {
      return window.localStorage.getItem("user");
    },

    numberOfPages() {
      return Math.ceil(this.avatars.length / this.avatarsPerPage);
    },
    showAssign() {
      return this.selected.length > 0 && !this.reassigning;
    },
    showReassign() {
      return this.selected.length > 0 && this.reassigning;
    },
    selectedAvatarsIDs() {
      let selectedIDs = [];
      if (this.selected.length > 0) {
        this.selected.forEach((avatar) => {
          /**eslint-disable */
          console.log("AVATAR_ID", avatar.id);
          selectedIDs.push(avatar.id);
        });
      }
      return selectedIDs;
    },
    members() {
      const names = [];
      this.team.forEach((member) => {
        names.push(member.firstname + " " + member.lastname);
      });

      return names;
    },
    stats() {
      let suspended = 0;
      let active = 0;
      let assigned = 0;
      let unassigned = 0;

      this.avatars.forEach((avatar) => {
        if (Object.keys(avatar).length <= 6) {
          suspended++;
        } else {
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

      return {
        avatars: this.avatars.length,
        suspended: suspended,
        active: active,
        assigned: assigned,
        unassigned: unassigned,
      };
    },
  },
  methods: {
    handlePaginate(v) {
      /**eslint-disable */
      console.log("PAGINATION_EVENT", v);
    },
    itemSelected(v) {
      /**eslint-disable */
      console.log("ITEM_SELECTED", v);
    },
    nextPage() {
      if (this.page + 1 <= this.numberOfPages) this.page += 1;
    },
    formerPage() {
      if (this.page - 1 >= 1) this.page -= 1;
    },
    updateAvatarsPerPage(number) {
      this.avatarsPerPage = number;
    },
    sortAscendingBy(prop) {
      this.avatars = this.avatars.sort(function (a, b) {
        var pa = prop in a ? a[prop] : 0; //check if a[prop] exists if not return 0
        var pb = prop in b ? b[prop] : 0;
        return pa < pb ? -1 : 1;
      });
    },
    sortDescendingBy(prop) {
      this.avatars = this.avatars.sort(function (a, b) {
        var pa = prop in a ? a[prop] : 0; //check if a[prop] exists if not return 0
        var pb = prop in b ? b[prop] : 0;
        return pa > pb ? -1 : 1;
      });
    },
    assign(props, item) {
      props.select(item, !props.isSelected(item));
      if (props.isSelected(item) && this.selected.length > 0) {
        this.assignCard = true;
      } else {
        this.assignCard = false;
      }
    },
    reassign(props, item) {
      props.select(item, !props.isSelected(item));
      if (props.isSelected(item) && this.selected.length > 0) {
        this.reassignCard = true;
      } else {
        this.reassignCard = false;
      }
    },
    // once(fn, context) {
    //   let result;

    //   return function () {
    //     if (fn) {
    //       result = fn.apply(context || this, arguments);
    //       fn = null;
    //     }

    //     return result;
    //   };
    // },
    assignAvatar() {
      this.team.forEach((member) => {
        let name = member.firstname + " " + member.lastname;
        if (name === this.assignee) {
          this.assigneeID = member.id;
        }
      });
      this.assignDialog = false;
      this.isAssigning = false;
      let payload = {
        userID: this.assigneeID,
        avatars: this.selectedAvatarsIDs,
      };
      /**eslint-disable */
      // console.log("PAYLOAD", payload);
      this.$store.dispatch("avatars/assignAvatars", {
        token: this.token,
        assign: payload,
        router: this.$router,
      });
    },
    reassignAvatar() {
      this.team.forEach((member) => {
        let name = member.firstname + " " + member.lastname;
        if (name === this.reassignee) {
          this.reassigneeID = member.id;
        }
      });
      this.reassignDialog = false;
      this.reassigning = false;
      let payload = {
        userID: this.reassigneeID,
        avatars: this.selectedAvatarsIDs,
      };
      this.$store.dispatch("avatars/assignAvatars", {
        token: this.token,
        assign: payload,
        router: this.$router,
      });
      /**eslint-disable */
      // console.log(payload);
    },
  },
};
</script>
<style scoped>
</style>