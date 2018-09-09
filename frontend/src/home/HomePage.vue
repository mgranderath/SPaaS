<template>
    <section class="section">
        <p v-model="apps" v-show="apps.length == 0">COll</p>
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
};
</script>