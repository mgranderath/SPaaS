<template>
  <div class="column is-one-third">
    <nav class="panel">
      <p class="panel-heading">
        Apps
      </p>
      <div class="panel-block">
        <button class="button is-link is-outlined is-fullwidth" v-on:click="openModal">
          <span class="icon is-small">
              <i class="fas fa-plus"></i>
          </span>
          <span>Add</span>
        </button>
      </div>
      <a class="panel-block" v-model="apps" v-for="item in apps" v-on:click="selectApp(item)" :class="{ 'is-active' : appSelected == item }">
        <span class="panel-icon">
        <i class="fas fa-rocket" aria-hidden="true"></i>
        </span>
        <span>{{ item }}</span>
      </a>
    </nav>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
export default {
  name: "AppListPanel",
  created() {
    this.$store.dispatch("api/getAll")
  },
  computed: {
    ...mapGetters({
      apps: "api/getApps",
      appSelected: "viewstate/getAppSelected"
    })
  },
  methods: {
    openModal: function() {
      this.$store.dispatch("viewstate/openModal", "createModal");
    },
    selectApp: function(name) {
      this.$store.dispatch("viewstate/selectApp", name)
    }
  }
};
</script>
