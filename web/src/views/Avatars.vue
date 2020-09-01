<template>
  <div class="avatars">
    <h1>Avatars</h1>
    <v-container class="my-5">
      <v-data-iterator
        :items:="avatars"
        :items-per-page.sync="avatarsPerPage"
        :page="page"
        :search="search"
        :sort-desc="sortDesc"
      >
        <template v-slot:header>
          <v-toolbar flat class="mb-1">
            <v-text-field
              v-model="search"
              clearable
              flat
              solo-inverted
              hide-details
              prepend-inner-icon="search"
              label="Search"
            ></v-text-field>
            <template v-if="$vuetify.breakpoint.mdAndUp">
              <v-spacer></v-spacer>
              <v-btn text depressed large color="grey" @click="sortAscendingBy('likes')">
                <v-icon small left>arrow_upward</v-icon>
              </v-btn>
              <v-btn text depressed large color="grey" @click="sortDescendingBy('likes')">
                <v-icon small left>arrow_downward</v-icon>
              </v-btn>
            </template>
          </v-toolbar>
        </template>

        <template v-slot:default="props">
          <v-row class="mb-3">
            <v-col
              cols="12"
              sm="6"
              lg="3"
              md="4"
              v-for="avatar in props.avatars"
              :key="avatar.handle"
            >
              <v-card align="center" class="ma-4">
                <v-responsive class="pt-4">
                  <v-avatar size="100" class="grey lighten-2">
                    <img :src="avatar.profile" />
                  </v-avatar>
                </v-responsive>
                <v-card-subtitle>
                  <div class="grey--text">@{{avatar.handle}}</div>
                  <div class="grey--text text--darken-4">{{avatar.bio}}</div>
                </v-card-subtitle>
                <v-card-text>
                  <div class="grey--text text--darken-4">
                    {{avatar.following}}
                    <span class="grey--text text-caption">following</span>
                    {{avatar.followers}}
                    <span
                      class="grey--text text-caption"
                    >followers</span>
                  </div>
                  <div v-if="avatar.assigned">
                    <v-chip class="ma-2" color="indigo" text-color="white">
                      <v-icon left>mdi-account-multiple-check</v-icon>assigned
                    </v-chip>
                  </div>
                  <div v-else>
                    <v-chip class="ma-2" text-color="grey">
                      <v-icon left>mdi-account-multiple-remove</v-icon>unasssigned
                    </v-chip>
                  </div>
                </v-card-text>
                <v-card-actions>
                  <v-row class="ma-3 text-sm grey--text">
                    <div>
                      <v-icon small left>mdi-thumb-up</v-icon>
                      <span class="grey--text text--darken-4">{{avatar.likes}}</span> likes
                    </div>

                    <v-spacer></v-spacer>
                    <div class="text-subtitle-1">
                      <v-icon small left color="secondary">mdi-calendar-month-outline</v-icon>
                      <span class="grey--text text--darken-4">Joined</span>
                      {{avatar.joined}}
                    </div>
                  </v-row>
                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
        </template>
        <template v-slot:footer>
          <v-row class="mt-2" align="center" justify="center">
            <span class="grey--text">Avatars per page</span>
            <v-menu offset-y>
              <template v-slot:activator="{ on, attrs }">
                <v-btn dark text color="primary" class="ml-2" v-bind="attrs" v-on="on">
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

            <span class="mr-4 grey--text">Page {{ page }} of {{ numberOfPages }}</span>
            <v-btn fab dark color="blue darken-3" class="mr-1" @click="formerPage">
              <v-icon>mdi-chevron-left</v-icon>
            </v-btn>
            <v-btn fab dark color="blue darken-3" class="ml-1" @click="nextPage">
              <v-icon>mdi-chevron-right</v-icon>
            </v-btn>
          </v-row>
        </template>
      </v-data-iterator>
    </v-container>
  </div>
</template>
<script>
export default {
  data() {
    return {
      avatarsPerPage: 16,
      avatarsPerPageArray: [16, 32, 64],
      search: "",
      filter: {},
      sortDesc: false,
      page: 1,
      sortByHandle: "handle",
      avatars: [
        {
          handle: "mwasyakyalo",
          bio: "Real distance trader | A woman of her word",
          following: 323,
          followers: 230,
          likes: 710,
          joined: "April 2010",
          profile: "/avatar-3.png",
          assigned: true,
        },
        {
          handle: "catemukami",
          bio: "Staunch Kikuyu | Lover of fine things",
          following: 324,
          followers: 710,
          likes: 300,
          joined: "April 2010",
          profile: "/avatar-4.png",
          assigned: true,
        },
        {
          handle: "mainakamandaelder",
          bio: "For the people by the people",
          following: 340,
          followers: 349,
          likes: 490,
          joined: "April 2010",
          profile: "/avatar-5.png",
          assigned: true,
        },
        {
          handle: "salmakambo",
          bio: "Maji yakimwagika hayazoleki | Leave the past to be ",
          following: 710,
          followers: 800,
          likes: 1020,
          joined: "April 2010",
          profile: "/avatar-1.png",
          assigned: false,
        },
        {
          handle: "the_basebandit",
          bio:
            "Distributed Systems Enthusiast | Teacher | Mentor | Father to many",
          following: 540,
          followers: 620,
          likes: 230,
          joined: "April 2010",
          profile: "/avatar-2.png",
          assigned: true,
        },
        {
          handle: "mamakhe",
          bio: "Who jah bless no woman curse",
          following: 410,
          followers: 301,
          likes: 370,
          joined: "April 2010",
          profile: "/avatar-3.png",
          assigned: false,
        },
        {
          handle: "jaduong_jeuri",
          bio: "Live and let live | Fisher of women | Opinions are my own",
          following: 210,
          followers: 310,
          likes: 200,
          joined: "April 2010",
          profile: "/avatar-4.png",
          assigned: false,
        },
        {
          handle: "mwanziakaluki",
          bio: "Real distance trader | Connecting people",
          following: 323,
          followers: 230,
          likes: 710,
          joined: "April 2010",
          profile: "/avatar-5.png",
          assigned: false,
        },
        {
          handle: "catemukamizo",
          bio: "Staunch Kikuyu | Lover of fine things",
          following: 324,
          followers: 710,
          likes: 300,
          joined: "April 2010",
          profile: "/avatar-1.png",
          assigned: false,
        },
        {
          handle: "mainakamandz",
          bio: "For the people by the people",
          following: 340,
          followers: 349,
          likes: 490,
          joined: "April 2010",
          profile: "/avatar-3.png",
          assigned: true,
        },
        {
          handle: "welmakambui",
          bio: "Maji yakimwagika hayazoleki | Leave the past to be ",
          following: 710,
          followers: 800,
          likes: 1020,
          joined: "April 2010",
          profile: "/avatar-2.png",
          assigned: true,
        },
        {
          handle: "mr.Parish",
          bio: "Substance over hype | Good over evil | Go getter",
          following: 540,
          followers: 620,
          likes: 230,
          joined: "April 2010",
          profile: "/avatar-3.png",
          assigned: true,
        },
        {
          handle: "mamakhellen",
          bio: "Who jah bless no woman curse",
          following: 410,
          followers: 301,
          likes: 370,
          joined: "April 2010",
          profile: "/avatar-4.png",
          assigned: true,
        },
        {
          handle: "jaduong_ndogo",
          bio: "Live and let live",
          following: 210,
          followers: 310,
          likes: 200,
          joined: "April 2010",
          profile: "/avatar-5.png",
          assigned: true,
        },
      ],
    };
  },
  computed: {
    numberOfPages() {
      return Math.ceil(this.avatars.length / this.avatarsPerPage);
    },
    filteredKeys() {
      return this.keys.filter((key) => key !== `Name`);
    },
  },
  methods: {
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
      this.avatars.sort((a, b) => (a[prop] < b[prop] ? -1 : 1));
    },
    sortDescendingBy(prop) {
      this.avatars.sort((a, b) => (a[prop] > b[prop] ? -1 : 1));
    },
  },
};
</script>