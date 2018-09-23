import { userService, alertService } from "../_services";
import { router } from "../_helpers";

const user = JSON.parse(localStorage.getItem("user"));
const initialState = user
  ? { status: { loggedIn: true }, user: user }
  : { status: {}, user: null };

export const authentication = {
  namespaced: true,
  state: initialState,
  actions: {
    login({ dispatch, commit }, { username, password }) {
      commit("loginRequest", { username });

      userService
        .login(username, password)
        .then(user => {
          commit("loginSuccess", user);
          router.push("/");
        })
        .catch(error => {
          alertService.error("Login Error", error);
          commit("loginFailure", error);
        });
    },
    logout({ commit }) {
      userService.logout();
      commit("logout");
      router.push("/login");
    },
    changePassword({ dispatch, commit }, newPassword) {
      commit("changePasswordRequest");

      userService.changePassword(newPassword)
      .then(response => {
        commit("changePasswordSuccess")
        dispatch("viewstate/closeModal", "changePasswordModal", { root: true });
      })
    }
  },
  mutations: {
    loginRequest(state, user) {
      state.status = { loggingIn: true };
      state.user = user;
    },
    loginSuccess(state, user) {
      state.status = { loggedIn: true };
      state.user = user;
    },
    loginFailure(state) {
      state.status = {};
      state.user = null;
    },
    logout(state) {
      state.status = {};
      state.user = null;
    },
    changePasswordRequest(state) {
      state.status = { ...state.status, changingPassword: true };
    },
    changePasswordSuccess(state) {
      state.status = { ...state.status, changingPassword: false };
    }
  },
  getters: {
    getUser: state => state.user,
    changingPassword: state => state.status.changingPassword,
    loggingIn: state => state.status.loggingIn
  }
};
