<template>
  <v-row justify="center">
    <v-dialog max-width="600px">
      <template v-slot:activator="{on,attrs}">
        <v-btn color="success" v-on="on" v-bind="attrs">
          <v-icon left>mdi-plus</v-icon>Add new avatar
        </v-btn>
      </template>
      <v-card>
        <v-card-title>
          <span class="headline">Upload Avatar(s)</span>
        </v-card-title>
        <v-card-text>
          <v-form class="px-3">
            <!-- <v-text-field v-model="handle" label="Twitter handle" prepend-icon="mdi-account-circle"></v-text-field> -->
            <!-- <v-textarea v-model="bio" label="Twitter bio" prepend-icon="mdi-clipboard-account"></v-textarea> -->
            <v-file-input accept="text/csv" label="File input"></v-file-input>
            <v-btn class="secondary ma-2" @click="submit">
              <v-icon>mdi-file-excel</v-icon>Download CSV Template
            </v-btn>
            <v-btn class="success ma-2" @click="submit">
              <v-icon>mdi-file-upload</v-icon>Upload
            </v-btn>
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
      // bio: "",
    };
  },
  methods: {
    submit() {
      console.log(this.handle, this.bio);
    },
  },
};
</script>