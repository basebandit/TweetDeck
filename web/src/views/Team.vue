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
      team: [
        {
          id: 1,
          firstname: "Evanson ",
          lastname: "Mwangi",
          avatars: 120,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          id: 2,
          firstname: "Marcus",
          lastname: "Mwangi",
          avatars: 320,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          id: 3,
          firstname: "Mercy ",
          lastname: "Orangi",
          avatars: 620,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          id: 4,
          firstname: "Millicent",
          lastname: "Achieng",
          avatars: 120,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          id: 5,
          firstname: "Edward",
          lastname: "Kitili",
          avatars: 80,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          id: 6,
          firstname: "Changaresi",
          lastname: "Mugamura",
          avatars: 60,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          id: 7,
          firstname: "Everlyne ",
          lastname: "Waithera",
          avatars: 140,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          id: 8,
          firstname: "Wilberforce",
          lastname: "Juma",
          avatars: 89,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          id: 9,
          firstname: "Yvette ",
          lastname: "Anyango",
          avatars: 96,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          id: 10,
          firstname: "Linete ",
          lastname: "Wavinya",
          avatars: 80,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
      ],
    };
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
