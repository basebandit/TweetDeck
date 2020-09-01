<template>
  <v-row justify="center">
    <v-dialog max-width="600px">
      <template v-slot:activator="{on,attrs}">
        <v-btn color="success" v-on="on" v-bind="attrs">Add new avatar</v-btn>
      </template>
      <v-card>
        <v-card-title>
          <span class="headline">Add a New Avatar</span>
        </v-card-title>
        <v-card-text>
          <v-form class="px-3">
            <v-text-field v-model="handle" label="Twitter handle" prepend-icon="mdi-account-circle"></v-text-field>
            <v-textarea v-model="bio" label="Twitter bio" prepend-icon="mdi-clipboard-account"></v-textarea>
            <v-btn class="success mx-0 mt-3 ml-8" @click="submit">Add Avatar</v-btn>
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
export default {
  name: "Popup",
  data() {
    return {
      handle: "",
      bio: "",
    };
  },
  methods: {
    submit() {
      console.log(this.handle, this.bio);
    },
  },
};
</script>