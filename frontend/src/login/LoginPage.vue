<template>
    <section class="hero is-success is-fullheight">
        <div class="hero-body">
            <div class="container has-text-centered">
                <div class="column is-4 is-offset-4">
                    <h1 class="title has-text-black">SPaaS Dashboard</h1>
                    <p class="subtitle has-text-grey">Please login to proceed.</p>
                    <div class="box">
                        <form @submit.prevent="handleSubmit">
                            <div class="field">
                                <div class="control">
                                    <input class="input is-large" v-model="username" type="text" name="username" placeholder="Username" autofocus="" :class="{ 'is-danger': submitted && !username }">
                                </div>
                            </div>
                            <div class="field">
                                <div class="control">
                                    <input class="input is-large" v-model="password" name="password" type="password" placeholder="Your Password" :class="{ 'is-danger': submitted && !password }">
                                </div>
                            </div>
                            <button class="button is-block is-info is-large is-fullwidth" :class="{ 'is-loading' : loggingIn }">Login</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
        <vue-snotify></vue-snotify>
    </section>
</template>

<style scoped>
@import url("index.css");
</style>


<script>
import { mapGetters } from "vuex";
export default {
  data() {
    return {
      username: "",
      password: "",
      submitted: false
    };
  },
  computed: {
    ...mapGetters({
      loggingIn: "authentication/loggingIn"
    }),
    loggingIn() {
      return this.$store.state.authentication.status.loggingIn;
    }
  },
  created() {
    // reset login status
    this.$store.dispatch("authentication/logout");
  },
  methods: {
    handleSubmit(e) {
      this.submitted = true;
      const { username, password } = this;
      const { dispatch } = this.$store;
      if (username && password) {
        dispatch("authentication/login", { username, password });
      }
    }
  }
};
</script>