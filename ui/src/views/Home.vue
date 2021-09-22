<template>
  <loading-indicator v-if="loading"></loading-indicator>
  <div class="container" v-else>

    <div class="row">
      <div class="col-12">
        <portfolio-header :portfolio="portfolio" @privateMode="handlePrivateMode"></portfolio-header>
      </div>
    </div>
    <div class="row row-cols-1 row-cols-md-3 g-4 mt-2">
      <div class="col">
        <div class="card tile">
          <div class="card-header" v-if="portfolio">
            Road to <span :class="{'blur':privacy}">{{ formatCurrency(this.portfolio.goal, this.portfolio.base_currency, true) }}</span>
          </div>
          <div class="card-body">
            <pie-chart :chart-data="roadToTarget" :options="chartOptions" :privacy-mode="privacy" v-if="roadToTarget"
                       :key="privacy"></pie-chart>
          </div>
        </div>
      </div>
      <div class="col">
        <div class="card tile">
          <div class="card-header">
            Asset Distribution
          </div>
          <div class="card-body">
            <pie-chart :chart-data="distribution" :options="chartOptions" :privacy-mode="privacy" v-if="distribution"
                       :key="privacy"></pie-chart>
          </div>
        </div>
      </div>
      <div class="col">
        <div class="card tile">
          <div class="card-header">
            Portfolio Growth
            <router-link class="btn btn-primary btn-sm float-end"
                         :to="{name:'PortfolioHistory', params:{id:portfolio.id, portfolio:portfolio}}">View History
            </router-link>
          </div>
          <div class="card-body">
            <portfolio-growth :history="portfolio.history" :privacy-mode="privacy"
                              :key="growthComponentKey"></portfolio-growth>
            <portfolio-history-grouping :history="portfolio.history" v-if="portfolio"
                                        @historyGrouped="handleGrouping"></portfolio-history-grouping>
          </div>
        </div>
      </div>
    </div>

    <div class="row mt-4">
      <div class="col rounded">
        <div class="card shadow-sm rounded">
          <div class="table-responsive rounded">
            <table class="table table-striped table-hover table-bordered rounded" v-if="portfolio">
              <thead class="table-dark">
              <tr>
                <th>Asset</th>
                <th>Allocation</th>
                <th  class="text-end">Amount</th>
                <th  class="text-end">Total</th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="(item,index) in portfolio.items" :key="index">
                <td>{{ item.label || item.symbol }} <span class="text-sm fw-bold" v-if="item.unit_price">({{formatCurrency(item.unit_price)}})</span> </td>
                <td>
                  <div class="progress position-relative">
                    <div class="progress-bar progress-bar-striped" role="progressbar"
                         :style="`width: ${formatPercentage(item.total_value / portfolio.total_value)}`"
                         aria-valuenow="(Math.abs(item.total_value) / portfolio.total_value)*100"
                         aria-valuemin="0" aria-valuemax="100"></div>
                    <small class="justify-content-center d-flex position-absolute w-100">
                      {{ formatPercentage(item.total_value / portfolio.total_value) }}
                    </small>
                  </div>
                </td>
                <td :class="{'blur':privacy,}" class="text-end">{{ formatNumber(item.quantity, 3) }}</td>
                <td :class="{ 'blur':privacy, negative: item.total_value < 0 }"  class="text-end">{{
                    formatCurrency(item.total_value,portfolio.base_currency, true)
                  }}
                </td>
              </tr>
              <tr>
                <td colspan="3">Total</td>
                <td class="has-text-weight-bold text-end" :class="{'blur':privacy}"><em>{{
                    formatCurrency(portfolio.total_value, portfolio.base_currency, true)
                  }}</em></td>
              </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {getPortfolio, getPortfolios} from "../service/portfolio_service";
import {formatCurrency, formatDate, formatDecimal, formatNumber, formatPercentage} from "../utils/utils";
import PieChart from "../components/PieChart";
import {interpolateRdYlBu} from "d3-scale-chromatic";
import {interpolateColors} from "../utils/colors";
import PortfolioHeader from "../components/PortfolioHeader";
import PortfolioGrowth from "../components/PortfolioGrowth";
import PortfolioHistoryGrouping from "../components/PortfolioHistoryGrouping";
import LoadingIndicator from "../components/LoadingIndicator";
import {state,mutations} from "../store/store";

export default {

  name: 'Home',
  components: {
    LoadingIndicator,
    PortfolioHistoryGrouping,
    PortfolioGrowth,
    PortfolioHeader,
    PieChart
  },
  data() {
    return {
      portfolios: [],
      portfolio: null,
      chartData: null,
      growthComponentKey: 0,
      colors: [],
      distribution: null,
      roadToTarget: null,
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
    formatNumber,
    formatDecimal,
    formatPercentage,
    handleGrouping(data) {
      this.portfolio.history = data
      this.growthComponentKey++
    },
    handlePrivateMode(visible) {
      mutations.setPrivacy(!visible)
      this.growthComponentKey++
    },
    generateChartColors: function (dataLength) {
      return interpolateColors(dataLength, interpolateRdYlBu, {
        colorStart: 0.75,
        colorEnd: 2,
        useEndAsStart: false,
      });
    }
  },
  computed: {
    privacy: function () {
      return state.privacy
    },
    loading: function () {
      return this.portfolio === null
    }
  },
  mounted: async function () {
    try {
      this.portfolios = await this.getPortfolios()
      this.portfolio = await this.getPortfolio(this.portfolios[0].id)
      this.portfolio.items = this.portfolio.items.filter((i)=>i.quantity > 0)
      this.colors = this.generateChartColors((this.portfolio.items || []).length)
      // this.privacy = localStorage.getItem("privateMode") === 'true'

      let assetDistribution = this.portfolio.items.reduce(function (r, o) {
        (r[o.asset_type]) ? r[o.asset_type] += o.total_value : r[o.asset_type] = o.total_value;
        return r;
      }, {});

      this.distribution = {
        labels: Object.keys(assetDistribution),
        datasets: [
          {
            label: this.portfolio.name,
            backgroundColor: this.colors,
            hoverBackgroundColor: this.colors,
            data: Object.values(assetDistribution)
          }
        ]
      }
      this.roadToTarget = {
        labels: ["Complete", "Remaining"],
        datasets: [{
          backgroundColor: this.colors,
          hoverBackgroundColor: this.colors,
          data: [this.portfolio.total_value, this.portfolio.goal - this.portfolio.total_value]
        }]
      }
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

.tile {
  min-height: 525px;
  max-height: 525px;
}

.rounded {
  border-radius: 0.5rem !important;
}

.blur {
  filter: blur(0.5rem)
}

em{
  font-weight: 700;
}
.text-sm{
  font-size: x-small;
}
</style>
