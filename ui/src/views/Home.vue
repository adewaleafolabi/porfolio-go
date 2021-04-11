<template>
  <div class="home">
    <h2 class="title has-text-centered" v-if="portfolio">{{ portfolio ? portfolio.name : '' }} Portfolio
      {{ formatCurrency(portfolio.total_value) }}</h2>

    <div class="columns pt-5 pb-5">
      <div class="column">
        <pie-chart :chart-data="chartData" :options="chartOptions" v-if="chartData"></pie-chart>
      </div>
      <div class="column">
        <line-chart :chart-data="lineChartData" :options="chartOptions" v-if="lineChartData"></line-chart>
      </div>

    </div>
    <table class="table is-striped  is-hoverable is-fullwidth" v-if="portfolio">
      <tr>
        <th>Asset</th>
        <th>Allocation</th>
        <th>Amount</th>
        <th>Total</th>
      </tr>
      <tr v-for="(item,index) in portfolio.items" :key="index">
        <td>{{ item.label || item.symbol }}</td>
        <td>
          <div class="progress-wrapper">
            <progress class="progress is-info " :value="(Math.abs(item.total_value) / portfolio.total_value)*100"
                      max="100">{{ formatPercentage(item.total_value / portfolio.total_value) }}
            </progress>
            <p class="progress-value ">{{ formatPercentage(item.total_value / portfolio.total_value) }}</p>
          </div>
        </td>
        <td>{{ formatDecimal(item.quantity,3) }}</td>
        <td :class="{ negative: item.total_value < 0 }">{{ formatCurrency(item.total_value) }}</td>
      </tr>
      <tr>
        <td colspan="3">Total</td>
        <td class="has-text-weight-bold"><em>{{ formatCurrency(portfolio.total_value) }}</em></td>
      </tr>
    </table>
  </div>
</template>

<script>
import {getPortfolio, getPortfolios} from "../service/portfolio_service";
import {formatCurrency, formatDate, formatDecimal, formatPercentage} from "../utils/utils";
import PieChart from "../components/PieChart";
import {interpolateRdYlBu} from "d3-scale-chromatic";
import {interpolateColors} from "../utils/colors";
import LineChart from "../components/LineChart";

export default {

  name: 'Home',
  components: {
    LineChart,
    PieChart
  },
  data() {
    return {
      portfolios: [],
      portfolio: null,
      chartData: null,
      colors: [],
      lineChartData: null,
      chartOptions: {
        legend: {
          display: true
        },
        responsive: true,
        maintainAspectRatio: false
      }
    }
  },
  methods: {
    getPortfolio,
    getPortfolios,
    formatCurrency,
    formatDecimal,
    formatPercentage,
    generateChartColors: function (dataLength) {
      return interpolateColors(dataLength, interpolateRdYlBu, {
        colorStart: 0.1,
        colorEnd: 1,
        useEndAsStart: false,
      });
    }
  },
  mounted: async function () {
    try {
      this.portfolios = await this.getPortfolios()
      this.portfolio = await this.getPortfolio(this.portfolios[0].id)
      this.colors = this.generateChartColors((this.portfolio.items || []).length)
      this.chartData = {
        labels: this.portfolio.items.map(i => i.label || i.symbol),
        datasets: [
          {
            label: this.portfolio.name,
            backgroundColor: this.colors,
            hoverBackgroundColor: this.colors,
            data: this.portfolio.items.map(i => i.total_value)
          }
        ]
      }
      this.lineChartData = {
        labels: this.portfolio.history.map(i => formatDate(i.date)),
        datasets: [
          {
            label: this.portfolio.name,
            backgroundColor: ['#466eb1'],
            fill: true,
            data: this.portfolio.history.map(i => i.value)
          }
        ]
      }
    } catch (e) {
      console.log(e)
    }
  },
}
</script>

<style>

.negative {
  color: red
}

.progress-wrapper {
  position: relative;
}

.progress-value {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  font-size: calc(1rem / 1.5);
  line-height: 1rem;
  font-weight: bold;
}

.progress.is-small + .progress-value {
  font-size: calc(0.75rem / 1.5);
  line-height: 0.75rem;
}

.progress.is-medium + .progress-value {
  font-size: calc(1.25rem / 1.5);
  line-height: 1.25rem;
}

.progress.is-large + .progress-value {
  font-size: calc(1.5rem / 1.5);
  line-height: 1.5rem;
}
</style>
