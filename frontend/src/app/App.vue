<template>
  <div class="full-height">
    <router-view></router-view>
  </div>
</template>

<script>
import { alertService } from "../_services";

export default {
  name: "app",
  created() {
    alertService.init(this);
  }
};
</script>

<style>
@import "~vue-snotify/styles/material.css";

.full-height {
  height: 100%;
}

body,
html {
  height: 100%;
}
</style>
