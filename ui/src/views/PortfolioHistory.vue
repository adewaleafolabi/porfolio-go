<template>
  <div class="portfolio-history">
    <line-chart></line-chart>
    <table class="table is-striped  is-hoverable is-fullwidth" v-if="portfolio">
      <tr v-for="(item,index) in portfolio.history" :key="index">
        <td>{{item.date}}</td>
        <td>{{item.value}}</td>
      </tr>
    </table>
  </div>
</template>

<script>
import LineChart from "../components/LineChart";

export default {

  name: 'PortfolioHistory',
  components: {
    LineChart
  },
  props:{
    portfolio: {
      type: Object,
      default: null
    }
  }
}
</script>
