<template>
  <v-row justify="center">
    <v-dialog max-width="600px" v-model="dialog">
      <Notification :message="snackbarMessage" :snackbar="snackbar" :type="snackbarType" />
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
            <v-file-input v-model="file" accept="text/csv" label="File input"></v-file-input>
            <v-btn class="secondary ma-2" @click="download">
              <v-icon>mdi-file-excel</v-icon>Download CSV Template
            </v-btn>
            <v-btn class="success ma-2" @click="upload">
              <v-icon>mdi-file-upload</v-icon>Upload
            </v-btn>
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
import AvatarService from "@/services/AvatarService";
export default {
  name: "Popup",
  data() {
    return {
      handle: "",
      file: "",
      dialog: false,
      snackbarType: "",
      snackbarMessage: "",
      snackbar: false,
      // bio: "",
    };
  },
  methods: {
    async upload() {
      let formData = new FormData();
      formData.append("avatars", this.file);
      try {
        const token = window.localStorage.getItem("user");
        const response = await AvatarService.uploadAvatars(token, formData);
        this.snackbarType = response.data.type;
        this.snackbarmessage = "Uploaded avatars successfully";
        this.snackbar = true;
        setTimeout(() => {
          this.dialog = false;
        }, 1000);
        /**eslint-disable */
        console.log(response.data.status, response.data);
      } catch (err) {
        setTimeout(() => {
          this.dialog = false;
        }, 500);
        /**eslint-disable */
        console.log(err.response.data.status, err.response.data);
      }
    },
    download() {
      let csvContent = "usernames\n";
      csvContent += "username1\n";
      let anchor = document.createElement("a");
      anchor.href =
        "data:text/csv;charset=utf-8," + encodeURIComponent(csvContent);
      anchor.target = "_blank";
      anchor.download = "avatars.csv";
      anchor.click();
    },
  },
};
</script>