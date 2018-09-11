<template>
<div class="column" :class="{ 'is-vertical-center' : (appSelected == '') }">
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
      <div class="tile is-ancestor" v-show="!notDeployed">
        <div class="tile has-text-centered is-parent">
          <div class="tile is-child level box">
            <div>
              <p class="heading">Deployed</p>
              <p class="title is-6">{{ new Date(appState.Created).toTimeString() }}</p>
            </div>
          </div>
        </div>
        <div class="tile has-text-centered is-parent">
          <div class="tile is-child level box">
            <div>
              <p class="heading">Container Name</p>
              <p class="title is-6">{{ appState.Name.substring(1) }}</p>
            </div>
          </div>
        </div>
        <div class="tile has-text-centered is-parent">
          <div class="tile is-child level box">
            <div>
              <p class="heading">State</p>
              <p class="title is-6">{{ appState.State.Status }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-show="tabSelected == 1">
      Cool1
    </div>
    <div v-show="tabSelected == 2">
      Cool2
    </div>
  </div>
  <div v-show="appSelected == ''" class="level is-vcentered">
    <div class="level-item has-text-centered">
      <div>
        <p class="title is-3">Select a App</p>
        <p><i class="fas fa-rocket fa-2x" aria-hidden="true"></i></p>
      </div>
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
      appSelected: "viewstate/getAppSelected",
      appState: "api/INSPECT_APP_STATE",
      notDeployed: "api/INSPECT_APP_NOT_DEPLOYED"
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

<style scoped>
@import url(./AppDetailFragment.css);

.level {
  width: 100%;
  min-width: 100%;
}
</style>

