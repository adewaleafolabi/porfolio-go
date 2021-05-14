<template>
  <div id="app">
    <nav :class="`navbar navbar-expand-lg navbar-${theme} bg-${theme}`">
      <div class="container-fluid">
        <router-link class="navbar-brand" to="/">
          <img src="./assets/logo.svg" alt="" width="30" height="24" class="d-inline-block align-text-top">
          Portfolio
        </router-link>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNavAltMarkup"
                aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
          <div class="navbar-nav">
            <router-link class="nav-link" aria-current="page" to="/">Home</router-link>
          </div>
        </div>
        <theme-switcher @theme="handleThemeChange"></theme-switcher>
      </div>
    </nav>
    <div class="container">
      <router-view/>
    </div>
  </div>
</template>
<script>
import ThemeSwitcher from "./components/ThemeSwitcher";

export default {
  components: {ThemeSwitcher},
  data() {
    return {
      theme: ''
    }
  },
  methods: {
    handleThemeChange(theme) {
      this.theme = theme;
    }
  },
  mounted: function () {
    this.theme = localStorage.getItem("theme")
  }
}
</script>