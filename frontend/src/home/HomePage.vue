<template>
  <div class="h-screen">
    <NavBar v-bind:username="user.username" enable-search></NavBar>

    <section class="sm:mx-8 h-auto flex flex-col justify-start">
      <div class="flex justify-end">
        <button class="bg-transparent flex items-baseline hover:bg-blue-500 text-blue-700 font-semibold hover:text-white py-2 px-4 border border-blue-500 hover:border-transparent rounded">
          New <i class="fa fa-plus ml-2" aria-hidden="true"></i>
        </button>
      </div>
      <div class="flex flex-col justify-start shadow h-full my-4">
        <div v-for="app in apps" class="flex flex-row justify-between px-8 py-4 border-b sm:px-16 lg:px-32 hover:bg-blue-100">
          <h1>{{ app.name }}</h1>
          <h1>
            <Node v-if="app.type === 'node'" class="w-6 h-6 inline-block fill-current"/>
            <Python v-else-if="app.type === 'python'" class="w-6 h-6 inline-block fill-current"/>
            <Ruby v-else-if="app.type === 'ruby'" class="w-6 h-6 inline-block fill-current"/>
          </h1>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import CreateAppModal from "../components/CreateAppModal";
import AppListPanel from "../components/AppListPanel";
import AppDetailFragment from "../components/AppDetailFragment";
import ChangePasswordModal from "../components/ChangePasswordModal";
import NavBar from "../components/NavBar";
import Node from "../assets/node.svg";
import Python from "../assets/python.svg";
import Ruby from "../assets/ruby.svg";

export default {
  components: {
    CreateAppModal,
    AppListPanel,
    AppDetailFragment,
    ChangePasswordModal,
    NavBar,
    Node, Python, Ruby
  },
  data() {
    return {
      navOpen: false,
    };
  },
  created() {
    this.$store.dispatch("api/getAll");
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
    },
    appTypeToIcon: function (appType) {
      return appTypeToIcon(appType);
    }
  }
};
</script>

<style scoped>
</style>
