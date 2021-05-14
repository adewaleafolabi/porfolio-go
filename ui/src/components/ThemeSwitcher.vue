<template>
  <div class="form-check form-switch">
    <input class="form-check-input" type="checkbox" id="themeSwitcher" @change="handle" v-model="theme">
    <label class="form-check-label" for="themeSwitcher">🌙</label>
  </div>
</template>

<script>
import {DARK_THEME, LIGHT_THEME, setClassAndThemes} from "../service/theme";

export default {
  name: "ThemeSwitcher",
  data() {
    return {
      theme: false
    }
  },
  methods: {
    handle() {
      let theme = this.theme?'dark':'light'
      localStorage.setItem("theme", theme)
      this.$emit("theme", theme)
      let body = window.document.body
      if(body){
        if(theme === DARK_THEME){
          setClassAndThemes(DARK_THEME)
          return
        }
        setClassAndThemes(LIGHT_THEME)
      }
    }
  },
  mounted: async function () {
    this.theme = localStorage.getItem("theme") === 'dark'
  }
}
</script>

<style>

</style>