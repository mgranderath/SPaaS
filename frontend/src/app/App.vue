<template>
  <div class="full-height">
    <div v-if="alert.message" class="is-danger notification">{{alert.message}}</div>
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: "app",
  computed: {
    alert() {
      return this.$store.state.alert;
    }
  },
  watch: {
    $route(to, from) {
      // clear alert on location change
      this.$store.dispatch("alert/clear");
    }
  }
};
</script>

<style>
.full-height {
  height: 100%;
}

body, html {
  height: 100%;
}
</style>
