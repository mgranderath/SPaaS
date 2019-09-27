import { apiService, alertService } from "../_services";

const initialState = {
  apps: [],
  messages: [],
  createAppStatus: -1,
  deployAppStatus: -1,
  stopAppStatus: -1,
  startAppStatus: -1,
  inspectAppState: {
    Name: "",
    Created: Date.now(),
    State: {
      Status: "",
      Running: false
    }
  },
  inspectAppNotDeployed: true
};

export const requestStatus = {
  "": -1,
  pending: 0,
  success: 1
};

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
      commit("CREATE_APP_RESET");
    },
    clearMessages({ commit }) {
      commit("clearMessage");
    },
    createApp({ dispatch, commit }, name) {
      commit("clearMessage");
      commit("CREATE_APP_PENDING");
      apiService
        .createApp(name)
        .then(data => {
          new ReadableStream({
            start(controller) {
              function push() {
                return data
                  .read()
                  .then(({ done, value }) => {
                    if (done) {
                      commit("CREATE_APP_SUCCESS");
                      dispatch("getAll");
                      controller.close();
                      return;
                    }
                    const responseObjects = new TextDecoder("utf-8")
                      .decode(value)
                      .split("\n");
                    responseObjects
                      .filter(value => {
                        return value != "";
                      })
                      .forEach(value => {
                        const Error = JSON.parse(value).type == "error";
                        if (Error) {
                          alertService.error(
                            "Error",
                            JSON.parse(value).message
                          );
                          controller.close();
                        }
                        commit("appendMessage", JSON.parse(value));
                      });
                      push();
                  })
                  .catch(error => {});
              }
              push();
            }
          });
        })
        .catch(error => {
          alertService.error("Error", error);
        });
    },
    inspectApp({ commit, dispatch }, name) {
      apiService.inspectApp(name).then(text => {
        commit("INSPECT_APP_STATE", JSON.parse(text));
      });
    },
    deployApp({ commit, dispatch }, name) {
      commit("DEPLOY_APP_PENDING");
      apiService.deployApp(name).then(text => {
        commit("DEPLOY_APP_SUCCESS");
        dispatch("inspectApp", name);
      });
    },
    stopApp({ commit, dispatch }, name) {
      commit("STOP_APP_PENDING");
      apiService.stopApp(name).then(text => {
        commit("STOP_APP_SUCCESS");
        dispatch("inspectApp", name);
      });
    },
    startApp({ commit, dispatch }, name) {
      commit("START_APP_PENDING");
      apiService.startApp(name).then(text => {
        commit("START_APP_SUCCESS");
        dispatch("inspectApp", name);
      });
    }
  },
  mutations: {
    getAllSuccess(state, list) {
      console.log(list);
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
      alertService.success("Creating App", name);
    },
    CREATE_APP_RESET(state) {
      state.createAppStatus = requestStatus[""];
    },
    DEPLOY_APP_PENDING(state) {
      state.deployAppStatus = requestStatus["pending"];
    },
    DEPLOY_APP_SUCCESS(state) {
      state.deployAppStatus = requestStatus["success"];
      alertService.success("Deploying App", name);
    },
    DEPLOY_APP_RESET(state) {
      state.deployAppStatus = requestStatus[""];
    },
    STOP_APP_PENDING(state) {
      state.stopAppStatus = requestStatus["pending"];
    },
    STOP_APP_SUCCESS(state) {
      state.stopAppStatus = requestStatus["success"];
      alertService.success("Stopping App", name);
    },
    STOP_APP_RESET(state) {
      state.stopAppStatus = requestStatus[""];
    },
    START_APP_PENDING(state) {
      state.startAppStatus = requestStatus["pending"];
    },
    START_APP_SUCCESS(state) {
      state.startAppStatus = requestStatus["success"];
      alertService.success("Starting App", name);
    },
    START_APP_RESET(state) {
      state.startAppStatus = requestStatus[""];
    },
    INSPECT_APP_STATE(state, newState) {
      if (newState["message"]) {
        state.inspectAppNotDeployed = newState["message"].includes(
          "No such container"
        );
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
    DEPLOY_APP_STATE: state =>
      state.deployAppStatus == requestStatus["pending"],
    STOP_APP_STATE: state => state.stopAppStatus == requestStatus["pending"],
    START_APP_STATE: state => state.startAppStatus == requestStatus["pending"]
  }
};
