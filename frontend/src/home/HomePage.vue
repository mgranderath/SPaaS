<template>
  <div style="height: 100%;">
    <nav class="navbar is-dark" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <div class="navbar-item">
          <h2 class="title has-text-white">SPaaS</h2>
        </div>

        <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false" v-on:click="toggleNav">
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>
      <div class="navbar-menu" :class="{ 'is-active' : navOpen }">
        <div class="navbar-end">
          <div class="navbar-item has-dropdown is-hoverable">
            <a class="navbar-link">
              {{ user.username }}
            </a>
            <div class="navbar-dropdown is-right">
              <a class="navbar-item" v-on:click="openPasswordModal">
                Change Password
              </a>
              <a class="navbar-item" v-on:click="logout">
                Logout
              </a>
            </div>
          </div>
        </div>
      </div>
    </nav>
    <section class="section">
        <div v-model="apps" v-show="apps.length == 0" class="level is-vcentered">
            <div class="level-item has-text-centered">
            <div>
              <p class="title is-3">Create a App</p>
              <a class="button is-link is-large text-top" v-on:click="openModal">Create</a>
            </div>
          </div>
        </div>
        <div class="columns" v-model="apps" v-show="apps.length > 0">
            <AppListPanel/>
            <AppDetailFragment/>
        </div>
        <CreateAppModal/>
        <ChangePasswordModal/>>
    </section>
    <vue-snotify></vue-snotify>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import CreateAppModal from "../components/CreateAppModal";
import AppListPanel from "../components/AppListPanel";
import AppDetailFragment from "../components/AppDetailFragment";
import ChangePasswordModal from "../components/ChangePasswordModal";

export default {
  components: {
    CreateAppModal,
    AppListPanel,
    AppDetailFragment,
    ChangePasswordModal
  },
  data() {
    return {
      navOpen: false,
    };
  },
  computed: {
    ...mapGetters({
      apps: "api/getApps",
      user: "authentication/getUser"
    })
  },
  methods: {
    openPasswordModal: function() {
      this.$store.dispatch("viewstate/openModal", "changePasswordModal");
    },
    openModal: function() {
      this.$store.dispatch("viewstate/openModal", "createModal");
    },
    logout: function() {
      this.$store.dispatch("authentication/logout");
    },
    toggleNav: function() {
      this.navOpen = !this.navOpen;
    }
  }
};
</script>

<style scoped>
.columns,
.section,
.level {
  height: 100%;
}

.text-top {
  margin-top: 1rem;
}
</style>
