<template>
  <div class="d-flex justify-content-between">
  <h2 class="title has-text-centered" v-if="portfolio">{{ portfolio ? portfolio.name : '' }} Portfolio
    <span :class="{'blur':privateMode}">{{ formatCurrency(portfolio.total_value) }} <h5>{{ formatCurrency(portfolio.total_value_usd, 'USD') }}</h5>
      <h5>ATH: {{ formatCurrency(portfolio.all_time_high) }} <span class="negative" v-if="portfolio.all_time_high > portfolio.total_value">
                          ({{ formatPercentage((portfolio.total_value - portfolio.all_time_high) / portfolio.all_time_high) }})
</span> </h5>
    </span></h2>
  <div class="row">
    <button class="btn btn-sm fs-3" @click="handleVisibility">{{privateMode?'🐵':'🙈'}}</button>
  </div>
  </div>
</template>

<script>
import {formatCurrency,formatPercentage} from "../utils/utils";
import {state,mutations} from "../store/store";
export default {
  name: "PortfolioHeader",
  data(){
    return{
    }
  },
  methods: {
    formatCurrency,
    formatPercentage,
    handleVisibility(){
      mutations.setPrivacy(!this.privateMode, this)
    }
  },
  props: {
    portfolio: {
      type: Object,
      default: null
    },
  },
  computed: {
    privateMode: function () {
      return state.privacy
    },
  },
  mounted: async function () {
  }
}
</script>

<style scoped>

</style>