<template>
  <v-row justify="center">
    <v-dialog max-width="600px" v-model="dialog">
      <Notification
        :message="snackbarMessage"
        :snackbar="snackbar"
        :type="snackbarType"
      />
      <template v-slot:activator="{ on, attrs }">
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
            <v-file-input
              v-model="file"
              accept="text/csv"
              label="File input"
            ></v-file-input>
            <v-btn class="secondary ma-2" @click="download">
              <v-icon>mdi-file-excel</v-icon>Download CSV Template
            </v-btn>
            <v-btn class="success ma-2" @click="upload" :loading="uploading">
              <v-icon>mdi-file-upload</v-icon>Upload
            </v-btn>
          </v-form>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script>
// import AvatarService from "@/services/AvatarService";
import Notification from "@/components/Notification";
import { mapGetters } from "vuex";
export default {
  name: "Popup",
  components: {
    Notification,
  },
  data() {
    return {
      handle: "",
      file: [],
      dialog: this.uploading,
      snackbarType: "",
      snackbarMessage: "",
      snackbar: false,
      // bio: "",
    };
  },
  computed: {
    token() {
      return window.localStorage.getItem("user");
    },
    ...mapGetters("avatars", ["uploading"]),
  },
  methods: {
    upload() {
      let formData = new FormData();
      /**eslint-disable */
      console.log("file.type");
      // if (this.file.type == "text/csv"){
      formData.append("avatars", this.file);

      this.$store.dispatch("avatars/upload", { token: this.token, formData });
      // } else {
      //   /**eslint-disable */
      //   console.error("non-csv file detected");
      //   this.snackbarType = "error";
      //   this.snackbarMessage = "Invalid file type";
      //   this.snackbar = true;
      // }
      //Reset the snackbar after 2 seconds
      setTimeout(() => {
        this.snackbarType = "";
        this.snackbarMessage = "";
        this.snackbar = false;
        this.dialog = false;
      }, 2000);
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