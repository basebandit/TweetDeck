<template>
  <nav>
    <v-app-bar app flat>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title class="text-upercase grey--text">
        <span class="font-weight-light">Avatar</span>
        <span>lysis</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-menu offset-y>
        <template v-slot:activator="{on,attrs}">
          <v-btn v-bind="attrs" v-on="on" text color="grey">
            <v-icon left>expand_more</v-icon>menu
          </v-btn>
        </template>
        <v-list>
          <v-list-item v-for="link in links" :key="link.text" :to="link.route">
            <v-list-item-title>{{link.text}}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-btn text color="grey" @click="signout">
        <span>Sign Out</span>
        <v-icon right>exit_to_app</v-icon>
      </v-btn>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" app class="primary">
      <v-row class="mt-5">
        <v-col align="center">
          <v-avatar size="100">
            <img src="/avatar-1.png" alt />
          </v-avatar>
          <p class="white--text subheading mt-1">The Basebandit</p>
          <v-col class="mt-4 mb-3">
            <Popup />
          </v-col>
        </v-col>
      </v-row>

      <v-list>
        <v-list-item v-for="link in links" :key="link.text" :to="link.route">
          <v-list-item-icon>
            <v-icon class="white--text">{{link.icon}}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title class="white--text">{{link.text}}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
  </nav>
</template>
<script>
import Popup from "./Popup";
export default {
  name: "Navbar",
  components: { Popup },
  data() {
    return {
      drawer: true,
      links: [
        { icon: "dashboard", text: "Dashboard", route: "/" },
        { icon: "mdi-account-supervisor", text: "Avatars", route: "/avatars" },
        { icon: "mdi-account-group", text: "Team", route: "/team" },
      ],
    };
  },
  methods: {
    signout() {
      window.localStorage.removeItem("user");
      this.$router.push({ name: "Login" });
    },
  },
};
</script>