import { apiService } from "../_services";

const initialState = {
  apps: [],
  messages: [],
  status: {}
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
    createApp({ dispatch, commit }, name) {
      commit("clearMessage")
      apiService.createApp(name).then(data => {
        new ReadableStream({
          start(controller) {
            function push() {
              return data
                .read()
                .then(({ done, value }) => {
                  if (done) {
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
      state.messages = []
    }
  },
  getters: {
    getApps: state => state.apps,
    getMessages: state => state.messages
  }
};
