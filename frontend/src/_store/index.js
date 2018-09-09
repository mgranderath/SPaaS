import Vue from 'vue';
import Vuex from 'vuex';

import { alert } from './alert.module';
import { authentication } from './authentication.module';
import { api } from './api.module'
import { viewstate } from './viewstate.module'

Vue.use(Vuex);

export const store = new Vuex.Store({
    modules: {
        alert,
        api,
        viewstate,
        authentication
    }
});
