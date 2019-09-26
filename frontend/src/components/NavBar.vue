<template>
  <nav class="flex items-center justify-between flex-wrap p-2 bg-white text-blue-600">
    <div>
      <a href="/" class="flex items-center flex-shrink-0 mr-6">
        <svg class="fill-current h-12 w-12 mr-2" viewBox="0 0 640 480" xmlns="http://www.w3.org/2000/svg">
          <path fill="none" d="M-1-1h642v482H-1z"/>
          <g>
            <path
              d="M97.5954 122.4996c72.2238 0 102.0411 1.9482 181.493 122.3009 0 0 42.8417 58.1988 94.836 55.2413 0 0 35.648 7.8631 67.3897-67.0712 31.7762-74.9696 57.7676-116.3977 138.6513-105.5536 0 0-37.5495 33.5297-42.8531 47.836-5.2922 14.3061-91.2964 120.8338-138.3077 128.2275 6.896 1.9716 16.541 4.4244 32.8644 1.467 0 0 16.8503 1.0092 19.5079 9.377-16.5983-1.467-71.4792 15.3038-83.0259 24.188-11.5352 8.9075-51.135 32.6143-82.8654 8.4264 6.7813-1.021 44.789-8.1918 62.3724-25.5493-29.5768.352-78.6615-9.518-115.283-44.5264-36.553-35.0202-125.1574-140.5738-172.3406-148.4604l37.561-5.9032"/>
          </g>
        </svg>
        <span class="font-bold text-xl tracking-tight">SPaaS</span>
      </a>
    </div>
    <div class="text-sm flex-grow flex justify-center" v-show="enableSearch">
      <input type="search" class="bg-purple-white shadow rounded border-0 p-3 w-1/2" placeholder="Search by name...">
    </div>
    <div class="flex items-center w-auto">
      <div>
        <router-link to="/settings" @click="open = !open"
                     class="inline-block text-xl px-4 py-2 leading-none text-blue-500 hover:text-black"><i
          class="fa fa-cog" aria-hidden="true"></i></router-link>
        <button @click="logout"
                class="text-xl text-center px-2 py-3 text-grey-darkest bg-white hover:text-black hover:bg-blue"><i
          class="fa fa-sign-out" aria-hidden="true"></i></button>
      </div>
    </div>
  </nav>
</template>

<script>
  export default {
    name: "NavBar",
    props: {
      username: String,
      enableSearch: Boolean
    },
    data() {
      return {};
    },
    methods: {
      logout: function () {
        this.$store.dispatch("authentication/logout");
      },
      toggleNav: function () {
        this.navOpen = !this.navOpen;
      }
    }
  }
</script>

<style scoped>

</style>