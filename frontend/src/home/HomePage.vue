<template>
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
    </section>
</template>

<script>
import { mapGetters } from "vuex";
import CreateAppModal from "../components/CreateAppModal"
import AppListPanel from "../components/AppListPanel"
import AppDetailFragment from "../components/AppDetailFragment"
export default {
  components: {
    CreateAppModal,
    AppListPanel,
    AppDetailFragment
  },
  computed: {
    ...mapGetters({
      apps: "api/getApps",
    })
  },
  methods: {
    openModal: function() {
      this.$store.dispatch("viewstate/openModal", "createModal");
    },
  }
};
</script>

<style scoped>
.columns, .section, .level {
  height: 100%;
}

.text-top {
  margin-top: 1rem;
}
</style>
