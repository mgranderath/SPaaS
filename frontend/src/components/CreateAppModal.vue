<template>
  <div class="modal" id="createAppModal" :class="{ 'is-active': createModal }">
    <div class="modal-background" v-on:click="closeModal"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">Create App</p>
        <button class="delete" aria-label="close" v-on:click="closeModal"></button>
      </header>
      <section class="modal-card-body">
        <div class="field">
          <label class="label">Name</label>
          <div class="control">
              <input class="input" type="text" name="name" v-model="name" placeholder="App name" autofocus="">
          </div>
        </div>
        <div class="box" v-show="messages.length > 0">
          <ul>
            <li v-for="item in messages" v-model="messages">
              <span v-if="item.type == 'info'">INFO: </span>
              <span v-else-if="item.type == 'success'">SUCCESS: </span>
              <span v-else>ERROR: </span> {{ item.message }}
            </li>
          </ul>
        </div>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-success" v-on:click="handleSubmit">Create App</button>
        <button class="button" v-on:click="closeModal">Cancel</button>
      </footer>
    </div>
  <button class="modal-close is-large" aria-label="close" v-on:click="closeModal"></button>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
export default {
  name: "CreateAppModal",
  data() {
    return {
      name: ""
    };
  },
  computed: {
    ...mapGetters({
      messages: "api/getMessages",
      createModal: "viewstate/getCreateModal"
    })
  },
  methods: {
    closeModal: function() {
      this.$store.dispatch("viewstate/closeModal", "createModal")
    },
    handleSubmit: function(e) {
      this.$store.dispatch("api/createApp", this.name)
    }
  }
};
</script>