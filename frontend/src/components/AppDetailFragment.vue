<template>
<div class="column is-fullheight" :class="{ 'is-vertical-center' : (appSelected == '') }">
  <div class="is-fullheight" v-show="appSelected != ''">
    <nav class="level">
      <div class="level-item has-text-centered">
        <h1 class="title is-h1">{{ appSelected.toUpperCase() }}</h1>
      </div> 
    </nav>
    <div class="tabs is-centered">
      <ul>
        <li :class="{ 'is-active' : tabSelected == 0 }" v-on:click="selectTab(0)"><a>Home</a></li>
        <li :class="{ 'is-active' : tabSelected == 1 }" v-on:click="selectTab(1)"><a>Logs</a></li>
        <li :class="{ 'is-active' : tabSelected == 2 }" v-on:click="selectTab(2)"><a>Settings</a></li>
      </ul>
    </div>
    <div v-show="tabSelected == 0" style="height: calc(100% - 77px)">
      <div v-show="notDeployed" class="level is-vcentered is-fullheight">
        <div class="level-item has-text-centered">
          <div>
            <h3 class="title is-3">Deploy App to see more!</h3>
            <a class="button is-link is-large text-top" v-on:click="deployApp" :class="{ 'is-loading' : deployState }">Deploy</a>
          </div>
        </div>
      </div>
      <div class="tile is-ancestor is-vertical" v-show="!notDeployed">
        <div class="tile">
          <div class="tile has-text-centered is-parent">
            <div class="tile is-child level box">
              <div>
                <p class="heading">Deployed</p>
                <p class="title is-6">{{ new Date(appState.Created).toTimeString() }}</p>
              </div>
            </div>
          </div>
          <div class="tile has-text-centered is-parent">
            <div class="tile is-child level box">
              <div>
                <p class="heading">Container Name</p>
                <p class="title is-6">{{ appState.Name.substring(1) }}</p>
              </div>
            </div>
          </div>
          <div class="tile has-text-centered is-parent">
            <div class="tile is-child level box">
              <div>
                <p class="heading">State</p>
                <p class="title is-6">{{ appState.State.Status }}</p>
              </div>
            </div>
          </div>
        </div>
        <div class="tile">
          <div class="tile is-parent level">
            <div class="level-item has-text-centered">
              <div>
                <p class="heading">Deploy</p>
                <a class="button is-link is-large" v-on:click="deployApp" :class="{ 'is-loading' : deployState }">Deploy</a>
              </div>
            </div>
            <div class="level-item has-text-centered" v-show="appState.State.Running">
              <div>
                <p class="heading">Stop</p>
                <a class="button is-danger is-large" v-on:click="stopApp" :class="{ 'is-loading' : stopState }">Stop</a>
              </div>
            </div>
            <div class="level-item has-text-centered" v-show="!appState.State.Running">
              <div>
                <p class="heading">Start</p>
                <a class="button is-success is-large" v-on:click="startApp" :class="{ 'is-loading' : startState }">Start</a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-show="tabSelected == 1">
      <table class="table is-striped is-hoverable is-fullwidth">
        <thead>
          <tr>
            <th>Time</th>
            <th>Log</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in logs">
            <td>{{ item.date }}</td>
            <td>{{ item.message }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-show="tabSelected == 2">
      <table class="table is-fullwidth is-narrow">
        <tbody>
          <tr>
            <td class="settings-td has-text-centered">Delete: </td>
            <td class="settings-td">
              <a class="button is-danger" :class="{ 'is-loading' : deleteInProgress }" v-on:click="deleteApp">Delete</a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
  <div v-show="appSelected == ''" class="level is-vcentered">
    <div class="level-item has-text-centered">
      <div>
        <p class="title is-3">Select a App</p>
        <p><i class="fas fa-rocket fa-2x" aria-hidden="true"></i></p>
      </div>
    </div>
  </div>
</div>
</template>

<script>
import { mapGetters } from "vuex";
import { apiService, alertService } from "../_services";
import dayjs from "dayjs"

export default {
  name: "AppDetailFragment",
  data() {
    return {
      tabSelected: 0,
      logs: [],
      deleteInProgress: false,
    };
  },
  computed: {
    ...mapGetters({
      appSelected: "viewstate/getAppSelected",
      appState: "api/INSPECT_APP_STATE",
      notDeployed: "api/INSPECT_APP_NOT_DEPLOYED",
      deployState: "api/DEPLOY_APP_STATE",
      stopState: "api/STOP_APP_STATE",
      startState: "api/START_APP_STATE"
    }),
    user() {
      return this.$store.state.authentication.user;
    }
  },
  methods: {
    selectTab: function(tab) {
      this.tabSelected = tab;
      if (tab == 1) {
        this.logs = []
        this.logsApp();
      }
    },
    deployApp: function() {
      this.$store.dispatch("api/deployApp", this.appSelected);
    },
    stopApp: function() {
      this.$store.dispatch("api/stopApp", this.appSelected);
    },
    startApp: function() {
      this.$store.dispatch("api/startApp", this.appSelected);
    },
    deleteApp: function() {
      this.deleteInProgress = true
      apiService.deleteApp(this.appSelected)
        .then( result => {
          this.$store.dispatch("viewstate/selectApp", '', { root: true });
          this.$store.dispatch("api/getAll", { root: true });
          alertService.success("Deleting App", this.appSelected);
          this.deleteInProgress = false
        })
    },
    logsApp: function() {
      apiService.logs(this.appSelected).then(reader => {
        var decoder = new TextDecoder();
        function search(ref) {
          return reader.read().then(function(result) {
            var decoded = decoder.decode(result.value || new Uint8Array(), {
              stream: !result.done
            });
            const responseObjects = decoded.split("\n");
            responseObjects
              .filter(value => {
                return value != "";
              })
              .forEach(value => {
                const Error = JSON.parse(value).type == "error";
                if (Error) {
                  dispatch("alert/error", JSON.parse(value).message, {
                    root: true
                  });
                }
                const message = JSON.parse(value)["message"]
                const test = { 
                  date: dayjs(message.substr(8, message.indexOf(" ") - 8)).format("YYYY-MM-DD HH:mm:ss"),
                  message: message.substr(message.indexOf(" "))
                }
                ref.logs.push(test);
              });

            return search();
          });
        }

        return search(this);
      });
    }
  }
};
</script>

<style scoped>
@import url(./AppDetailFragment.css);

.level {
  width: 100%;
  min-width: 100%;
}

.settings-td {
  vertical-align: center;
}

.is-fullheight {
  height: 100%;
}

.text-top {
  margin-top: 1rem;
}
</style>

