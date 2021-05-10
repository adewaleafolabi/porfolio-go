<template>
  <div class="d-flex justify-content-between">
  <h2 class="title has-text-centered" v-if="portfolio">{{ portfolio ? portfolio.name : '' }} Portfolio
    <span :class="{'blur':privateMode}">{{ formatCurrency(portfolio.total_value) }} <h5>{{ formatCurrency(portfolio.total_value_usd, 'USD') }}</h5></span></h2>
  <button class="btn btn-sm" @click="handleVisibility">{{privateMode?'🐵':'🙈'}}</button>
  </div>
</template>

<script>
import {formatCurrency} from "../utils/utils";
export default {
  name: "PortfolioHeader",
  data(){
    return{
      privateMode:false,
    }
  },
  methods: {
    formatCurrency,
    handleVisibility(){
      this.privateMode = !this.privateMode
      this.$emit("privateMode", this.privateMode)
      localStorage.setItem("privateMode", this.privateMode)
    }
  },
  props: {
    portfolio: {
      type: Object,
      default: null
    },
  },
  mounted: async function () {
    this.privateMode = localStorage.getItem("privateMode") === 'true'
  }
}
</script>

<style scoped>

</style>