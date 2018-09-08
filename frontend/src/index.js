import Vue from 'vue';

import { store } from './_store';
import { router } from './_helpers';
import App from './app/App';

new Vue({
    el: '#app',
    router,
    store,
    render: h => h(App)
});