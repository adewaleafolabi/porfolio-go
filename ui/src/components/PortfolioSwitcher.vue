<template>
  <div class="form-group pad">
    <select class="form-select" v-model="selected" @change="selectedChanged">
      <option v-for="(item,index) in portfolios" :key="index" :label="item.name" :value="item.id"></option>
    </select>
  </div>
</template>

<script>

import {mutations, state} from "../store/store";

export default {
  name: "PortfolioSwitcher",

  computed: {
    portfolios() {
      return state.portfolioKeyValue
    }, selected: {
      get: function () {
        return state.portfolioID
      }, set: function (value) {
        mutations.setPortfolioID(value)
      }
    }
  },
  mounted: function () {

  },
  methods: {
    selectedChanged() {
      mutations.setPortfolioID(this.selected)
    }
  }
}
</script>

<style scoped>
.pad {
  padding-right: 1em;
}
</style>