import { apiService } from "../_services";

const initialState = {
  apps: [],
  messages: [],
  createAppStatus: -1,
  deployAppStatus: -1,
  inspectAppState: {
    Name: "",
    Created: Date.now(),
    State: {
      Status: "",
      Running: false
    }
  },
  inspectAppNotDeployed: true,
};

export const requestStatus = {
  "": -1,
  "pending": 0,
  "success": 1
}

export const api = {
  namespaced: true,
  state: initialState,
  actions: {
    getAll({ dispatch, commit }) {
      apiService
        .getAll()
        .then(data => {
          commit("getAllSuccess", data);
        })
        .catch(error => {
          dispatch("alert/error", error, { root: true });
        });
    },
    resetCreateApp({ commit }) {
      commit("CREATE_APP_RESET")
    },
    clearMessages({ commit }) {
      commit("clearMessage")
    },
    createApp({ dispatch, commit }, name) {
      commit("clearMessage")
      commit("CREATE_APP_PENDING")
      apiService.createApp(name).then(data => {
        new ReadableStream({
          start(controller) {
            function push() {
              return data
                .read()
                .then(({ done, value }) => {
                  if (done) {
                    commit("CREATE_APP_SUCCESS")
                    dispatch("getAll")
                    controller.close()
                    return;
                  }
                  const responseObjects = new TextDecoder("utf-8").decode(value).split("\n");
                  responseObjects
                    .filter(value => {
                      return value != "";
                    })
                    .forEach(value => {
                      const Error = JSON.parse(value).type == "error"
                      if (Error) {
                        dispatch("alert/error", JSON.parse(value).message, { root: true })
                      }
                      commit(
                        "appendMessage",
                        JSON.parse(value)
                      );
                    });
                })
                .then(push)
                .catch( error => {})
            }
            push();
          }
        });
      })
      .catch( error => {
        dispatch("alert/error", error);
      })
    },
    inspectApp({ commit, dispatch }, name) {
      apiService.inspectApp(name)
        .then( text => {
          commit("INSPECT_APP_STATE", JSON.parse(text))
        })
    },
    deployApp({ commit, dispatch }, name) {
      commit("DEPLOY_APP_PENDING")
      apiService.deployApp(name)
        .then( text => {
          commit("DEPLOY_APP_SUCCESS")
        })
    }
  },
  mutations: {
    getAllSuccess(state, list) {
      state.apps = list;
    },
    appendMessage(state, item) {
      state.messages.push(item);
    },
    clearMessage(state) {
      state.messages = [];
    },
    CREATE_APP_PENDING(state) {
      state.createAppStatus = requestStatus["pending"];
    },
    CREATE_APP_SUCCESS(state) {
      state.createAppStatus = requestStatus["success"];
    },
    CREATE_APP_RESET(state) {
      state.createAppStatus = requestStatus[""];
    },
    DEPLOY_APP_PENDING(state) {
      state.deployAppStatus = requestStatus["pending"];
    },
    DEPLOY_APP_SUCCESS(state) {
      state.deployAppStatus = requestStatus["success"]
    },
    DEPLOY_APP_RESET(state) {
      state.deployAppStatus = requestStatus[""]
    },
    INSPECT_APP_STATE(state, newState) {
      if (newState["message"]) {
        state.inspectAppNotDeployed = newState["message"].includes("No such container");
      } else {
        state.inspectAppNotDeployed = false;
        state.inspectAppState = newState;
      }
    }
  },
  getters: {
    getApps: state => state.apps,
    getMessages: state => state.messages,
    CREATE_APP: state => state.createAppStatus,
    INSPECT_APP_STATE: state => state.inspectAppState,
    INSPECT_APP_NOT_DEPLOYED: state => state.inspectAppNotDeployed,
    DEPLOY_APP_STATE: state => state.deployAppStatus == requestStatus["pending"]
  }
};
