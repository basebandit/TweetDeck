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
        <v-data-table :headers="headers" :items="team" :search="search"></v-data-table>
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
        { text: "Member", sortable: false, align: "start", value: "name" },
        { text: "Avatars", value: "avatars" },
        { text: "CreatedAt", value: "createdAt" },
        { text: "Followers", value: "followers" },
        { text: "Following", value: "following" },
        { text: "Likes", value: "likes" },
        { text: "Tweets", value: "tweets" },
      ],
      team: [
        {
          name: "Evanson Mwangi",
          avatars: 120,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          name: "Marcus Mwangi",
          avatars: 320,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          name: "Mercy Orangi",
          avatars: 620,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          name: "Millicent Achieng",
          avatars: 120,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          name: "Edward Kitili",
          avatars: 80,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          name: "Changaresi Mugamura",
          avatars: 60,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          name: "Everlyne Waithera",
          avatars: 140,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          name: "Wilberforce Juma",
          avatars: 89,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          name: "Yvette Anyango",
          avatars: 96,
          createdAt: "May 26, 2020 10:09am",
          followers: 780,
          following: 300,
          likes: 600,
          tweets: 500,
        },
        {
          name: "Linete Wavinya",
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
};
</script>
