<template>
  <div>
    <v-snackbar :value="snackbar" :color="type" top>{{message}}</v-snackbar>
  </div>
</template>
<script>
export default {
  name: "Notification",
  props: {
    snackbar: {
      type: Boolean,
      default: false,
    },
    message: {
      type: String,
      default: "Notification message",
    },
    type: {
      type: String,
      default: "success",
    },
  },
};
</script>