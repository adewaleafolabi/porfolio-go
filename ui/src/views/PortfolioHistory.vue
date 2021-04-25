<template>
  <div class="portfolio-history">
    <div class="columns">
      <div class="column " >
        <line-chart v-if="portfolioData" :chart-data="chartData" :key="dateGrouping"  :styles="lineChartStyle"></line-chart>
      </div>
      <div class="column" ></div>
    </div>
    <div class="field has-addons">
      <p class="control">
        <button class="button" @click="dateGrouping='DAY'">
      <span class="icon is-small">
        <i class="fas fa-align-left"></i>
      </span>
          <span>Daily</span>
        </button>
      </p>
      <p class="control">
        <button class="button" @click="dateGrouping='MONTH'">
      <span class="icon is-small">
        <i class="fas fa-align-center"></i>
      </span>
          <span>Monthly</span>
        </button>
      </p>
      <p class="control">
        <button class="button" @click="dateGrouping='YEAR'">
      <span class="icon is-small">
        <i class="fas fa-align-right"></i>
      </span>
          <span>Yearly</span>
        </button>
      </p>
    </div>
    <table class="table is-striped  is-hoverable is-fullwidth" v-if="historyData">
      <thead>
      <th>Date</th>
      <th>Value</th>
      </thead>
      <tr v-for="(item,index) in historyData" :key="index">
        <td>{{ (item.date) }}</td>
        <td>{{ formatCurrency(item.value) }}</td>
      </tr>
    </table>
  </div>
</template>

<script>
import LineChart from "../components/LineChart";
import {getPortfolio} from "../service/portfolio_service";
import {formatCurrency, formatDate, formatDecimal} from "../utils/utils";
import {DateTime} from "luxon";

export default {

  name: 'PortfolioHistory',
  components: {
    LineChart
  },
  props: {
    portfolio: {
      type: Object,
      default: null
    },
    id: String,
  },
  methods: {
    getPortfolio,
    formatDate,
    formatDecimal,
    formatCurrency,
    groupHistory(grouping, history) {
      let format = grouping === 'DAY' ? 'yyyy-MM-dd'
          :grouping === 'MONTH' ? 'yyyy-MM': 'yyyy'

      //reduce by grouping items by date
      const groupBy = (arr) => arr.reduce((acc, ele)=>( (acc[DateTime.fromISO(ele.date).toFormat(format)] = acc[DateTime.fromISO(ele.date).toFormat(format)] || []).push(ele), acc),{})

      //get last day in month  map entry and return result
      return  Object.entries(groupBy(history)).map(([key, val])=> ({date:key, value: val[val.length-1].value}))

    }
  },
  data() {
    return {
      dateGrouping: 'DAY',
      portfolioData: null,
      chartOptions: {
        legend: {
          display: true
        },
        scales: {
          yAxes: [
            {
              ticks: {
                beginAtZero: true
              }
            }]
        },
        height: 300,
        responsive: false,
        maintainAspectRatio: false
      },
      lineChartStyle:{
        height:300
      }
    }
  },
  mounted: async function () {
    this.portfolioData = this.portfolio
    if (!this.portfolioData) {
      this.portfolioData = await this.getPortfolio(this.id)
    }

    this.portfolioData.history.sort(function(a,b){
      let nameA = DateTime.fromISO(a.date)
      let nameB = DateTime.fromISO(b.date)
      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }
      return 0;
    })
  },
  computed: {
    daily: function(){
      return this.groupHistory('DAY',this.portfolioData.history)
    },
    monthly: function(){
      return this.groupHistory('MONTH',this.portfolioData.history)
    },
    yearly: function(){
      return this.groupHistory('YEAR',this.portfolioData.history)
    },
    historyData: function (){
      if(!this.portfolioData){
        return null
      }
      return this.dateGrouping === 'DAY' ? this.daily
          :this.dateGrouping === 'MONTH' ? this.monthly: this.yearly
    }, chartData: function () {
      return {
        labels: this.historyData.map(i => i.date),
        datasets: [
          {
            label: this.portfolioData.name,
            backgroundColor: ['#466eb1'],
            fill: true,
            data: this.historyData.map(i => i.value)
          }
        ]
      }
    }
  }
}
</script>
