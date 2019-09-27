import Vue from 'vue';
import Router from 'vue-router';

import HomePage from '../home/HomePage'
import LoginPage from '../login/LoginPage'
import SettingsPage from "../settings/SettingsPage";
import App from '../home/App';

Vue.use(Router);

export const router = new Router({
  routes: [
    { path: '/', component: HomePage },
    { path: '/login', component: LoginPage },
    { path: '/settings', component: SettingsPage },
    { path: '/app/:name', component: App},

    // otherwise redirect to home
    { path: '*', redirect: '/' }
  ]
});

router.beforeEach((to, from, next) => {
  // redirect to login page if not logged in and trying to access a restricted page
  const publicPages = ['/login'];
  const authRequired = !publicPages.includes(to.path);
  const loggedIn = localStorage.getItem('user');

  if (authRequired && !loggedIn) {
    return next('/login');
  }

  next();
});