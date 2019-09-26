<template>
  <div class="h-screen flex bg-grey-darker flex-col">
    <nav class="flex flex-row items-center justify-between flex-wrap p-2">
      <div class="flex items-center flex-shrink-0 text-blue-600 mr-6">
        <svg class="fill-current h-12 w-12 mr-2" viewBox="0 0 640 480" xmlns="http://www.w3.org/2000/svg">
          <path fill="none" d="M-1-1h642v482H-1z"/>
          <g>
            <path
              d="M97.5954 122.4996c72.2238 0 102.0411 1.9482 181.493 122.3009 0 0 42.8417 58.1988 94.836 55.2413 0 0 35.648 7.8631 67.3897-67.0712 31.7762-74.9696 57.7676-116.3977 138.6513-105.5536 0 0-37.5495 33.5297-42.8531 47.836-5.2922 14.3061-91.2964 120.8338-138.3077 128.2275 6.896 1.9716 16.541 4.4244 32.8644 1.467 0 0 16.8503 1.0092 19.5079 9.377-16.5983-1.467-71.4792 15.3038-83.0259 24.188-11.5352 8.9075-51.135 32.6143-82.8654 8.4264 6.7813-1.021 44.789-8.1918 62.3724-25.5493-29.5768.352-78.6615-9.518-115.283-44.5264-36.553-35.0202-125.1574-140.5738-172.3406-148.4604l37.561-5.9032"/>
          </g>
        </svg>
        <span class="font-bold text-xl tracking-tight">SPaaS</span>
      </div>
      <div class="flex-grow flex items-center w-auto">
        <div class="flex-grow"></div>
        <div>
          <a href="https://github.com/mgranderath/SPaaS" class="px-4 py-2 text-xl"><i class="fa fa-github"
                                                                                                   aria-hidden="true"></i></a>
        </div>
      </div>
    </nav>
    <div class="max-w-sm m-auto">
      <form class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4" @submit.prevent="handleSubmit">
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="username">
            Username
          </label>
          <input v-model="username"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                 id="username" type="text" placeholder="Your Username" :class="{ 'border-red-500' : submitted && password === '' }">
        </div>
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="password">
            Password
          </label>
          <input v-model="password" v-on:keyup.enter="handleSubmit"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                 id="password" type="password" placeholder="Your Password" :class="{ 'border-red-500' : submitted && password === '' }">
        </div>
        <div class="w-full ld-over-full" :class="{ 'running' : loggingIn }">
          <button v-on:click="handleSubmit"
                  class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 w-full rounded focus:outline-none focus:shadow-outline"
                  type="button">
            Sign In
          </button>
          <div class="ld ld-ring ld-spin"></div>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
  @import url("index.css");
</style>


<script>
  import {mapGetters} from "vuex";

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
      handleSubmit() {
        this.submitted = true;
        const {username, password} = this;
        const {dispatch} = this.$store;
        if (username && password) {
          dispatch("authentication/login", {username, password});
        }
      }
    }
  };
</script>