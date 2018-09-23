const initialState = {
  "createModal": false,
  "changePasswordModal": false,
  "appSelected": ""
};

export const viewstate = {
  namespaced: true,
  state: initialState,
  actions: {
    closeModal({ commit }, message) {
      commit("closeModal", message);
    },
    openModal({ commit }, message) {
      commit("openModal", message);
    },
    selectApp({ commit, dispatch }, app) {
      if (app != "") {
        dispatch("api/inspectApp", app, { root: true })
      }
      commit("selectApp", app)
    }
  },
  mutations: {
    closeModal(state, message) {
      state[message] = false;
    },
    openModal(state, message) {
      state[message] = true;
    },
    selectApp(state, app) {
      state["appSelected"] = app
    }
  },
  getters: {
    getCreateModal: state => state["createModal"],
    getChangePasswordModal: state => state["changePasswordModal"],
    getAppSelected: state => state["appSelected"]
  }
};
