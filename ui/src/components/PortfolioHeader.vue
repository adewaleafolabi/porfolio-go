<template>
  <h2 class="title has-text-centered" v-if="portfolio">{{ portfolio ? portfolio.name : '' }} Portfolio
    {{ formatCurrency(portfolio.total_value) }} <h5>{{ formatCurrency(portfolio.total_value_usd, 'USD') }}</h5></h2>
</template>

<script>
import {formatCurrency} from "../utils/utils";
export default {
  name: "PortfolioHeader",
  methods: {
    formatCurrency,
  },
  props: {
    portfolio: {
      type: Object,
      default: null
    },
  }
}
</script>

<style scoped>

</style>