import Vue from "vue";
import Snotify, { SnotifyPosition } from "vue-snotify";
import './index.css';

import { store } from "./_store";
import { router } from "./_helpers";
import App from "./app/App";

const options = {
  toast: {
    position: SnotifyPosition.rightTop
  }
};

Vue.use(Snotify, options);

new Vue({
  el: "#app",
  router,
  store,
  render: h => h(App)
});