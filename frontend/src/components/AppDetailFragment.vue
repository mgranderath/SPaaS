<template>
<div class="column">
  <div v-show="appSelected != ''">
    <nav class="level">
      <div class="level-item has-text-centered">
        <h1 class="title is-h1">{{ appSelected.toUpperCase() }}</h1>
      </div> 
    </nav>
    <div class="tabs is-centered">
      <ul>
        <li :class="{ 'is-active' : tabSelected == 0 }" v-on:click="selectTab(0)"><a>Home</a></li>
        <li :class="{ 'is-active' : tabSelected == 1 }" v-on:click="selectTab(1)"><a>Logs</a></li>
        <li :class="{ 'is-active' : tabSelected == 2 }" v-on:click="selectTab(2)"><a>Settings</a></li>
      </ul>
    </div>
    <div v-show="tabSelected == 0">
      Cool
    </div>
    <div v-show="tabSelected == 1">
      Cool1
    </div>
    <div v-show="tabSelected == 2">
      Cool2
    </div>
  </div>
</div>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  name: "AppDetailFragment",
  data() {
    return {
      tabSelected: 0
    };
  },
  computed: {
    ...mapGetters({
      appSelected: "viewstate/getAppSelected"
    }),
    user() {
      return this.$store.state.authentication.user;
    }
  },
  methods: {
    selectTab: function(tab) {
      this.tabSelected = tab
    }
  }
}
</script>
