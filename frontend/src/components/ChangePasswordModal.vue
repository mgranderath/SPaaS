<template>
  <div class="modal" id="changePasswordModal" :class="{ 'is-active': changePasswordModal }">
    <div class="modal-background" v-on:click="closeModal"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">Change Password</p>
        <button class="delete" aria-label="close" v-on:click="closeModal"></button>
      </header>
      <section class="modal-card-body">
        <div class="field">
          <label class="label">New Password:</label>
          <div class="control">
              <input class="input" type="password" name="password" v-model="password" placeholder="New password" autofocus="">
          </div>
        </div>
      </section>
      <footer class="modal-card-foot">
        <button class="button is-success" v-on:click="handleSubmit" :disabled="changePassword > 0" :class="{ 'is-loading' : changePassword }">Change Password</button>
        <button class="button" v-on:click="closeModal">Cancel</button>
      </footer>
    </div>
  <button class="modal-close is-large" aria-label="close" v-on:click="closeModal"></button>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
export default {
  name: "ChangePasswordModal",
  data() {
    return {
      password: ""
    };
  },
  computed: {
    ...mapGetters({
      changePasswordModal: "viewstate/getChangePasswordModal",
      changePassword: "authentication/changingPassword"
    })
  },
  methods: {
    closeModal: function() {
      this.$store.dispatch("viewstate/closeModal", "changePasswordModal")
      this.password = ""
    },
    handleSubmit: function(e) {
      this.$store.dispatch("authentication/changePassword", this.password)
    }
  }
};
</script>