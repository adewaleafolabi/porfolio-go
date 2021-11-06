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
            Road to <span :class="{'blur':privacy}">{{
              formatCurrency(this.portfolio.goal, this.portfolio.base_currency, true)
            }}</span>
          </div>
          <div class="card-body">
            <pie-chart :chart-data="roadToTarget" :options="chartOptions" :privacy-mode="privacy" v-if="roadToTarget"
                       :key="hashCode(`${privacy}${currentPortfolioID}`)"></pie-chart>
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
                       :key="hashCode(`${privacy}${currentPortfolioID}`)"></pie-chart>
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
                              :key="hashCode(`${privacy}${currentPortfolioID}`)"></portfolio-growth>
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
                <th class="text-end">Amount</th>
                <th class="text-end">Total</th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="(item,index) in portfolio.items" :key="index">
                <td>{{ item.label || item.symbol }} <span class="text-sm fw-bold"
                                                          v-if="item.unit_price">({{
                    formatCurrency(item.unit_price)
                  }})</span>
                </td>
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
                <td :class="{ 'blur':privacy, negative: item.total_value < 0 }" class="text-end">{{
                    formatCurrency(item.total_value, portfolio.base_currency, true)
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
import {formatCurrency, formatDate, formatDecimal, formatNumber, formatPercentage, hashCode} from "../utils/utils";
import PieChart from "../components/PieChart";
import {interpolateRdYlBu} from "d3-scale-chromatic";
import {interpolateColors} from "../utils/colors";
import PortfolioHeader from "../components/PortfolioHeader";
import PortfolioGrowth from "../components/PortfolioGrowth";
import PortfolioHistoryGrouping from "../components/PortfolioHistoryGrouping";
import LoadingIndicator from "../components/LoadingIndicator";
import {mutations, state} from "../store/store";

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
      portfolios: {},
      growthComponentKey: 0,
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
    hashCode,
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
    },
    setupPortfolio: async function () {
      try {
        let portfolios = await this.getPortfolios()
        if (portfolios) {
          for (let index = 0; index < portfolios.length; ++index) {
            portfolios[index].items = portfolios[index].items.filter((i) => i.quantity != 0);
          }

          mutations.setPortfolioKeyValue(portfolios.map(k => {
            this.$set(this.portfolios, k.id, k)
            return {id: k.id, name: k.name}
          }))
        }
        if (this.currentPortfolioID) {
          if (!this.portfolios[this.currentPortfolioID]) {
            //if currentPortfolioID is not found on the server set the first item as the
            //current id
            mutations.setPortfolioID(portfolios[0].id)
          }
        } else {
          mutations.setPortfolioID(portfolios[0].id)
        }
      } catch (e) {
        console.log(e)
      }
    }
  },
  watch: {
    // eslint-disable-next-line no-unused-vars
    currentPortfolioID: function (old, newId) {
      // this.setupPortfolio()
    }
  },
  computed: {
    privacy: function () {
      return state.privacy
    },
    currentPortfolioID: function () {
      return state.portfolioID
    },
    loading: function () {
      return this.portfolio === undefined
    },
    portfolio: function () {
      return this.portfolios[this.currentPortfolioID]
    },
    colors: function () {
      if (!this.portfolio) {
        return []
      }
      let length = (this.portfolio.items || []).length
      if (length === 1) {
        length += 1
      }
      return this.generateChartColors(length)
    },
    distribution: function () {
      if (!this.portfolio) {
        return {}
      }
      let assetDistribution = this.portfolio.items.reduce(function (r, o) {
        (r[o.asset_type]) ? r[o.asset_type] += o.total_value : r[o.asset_type] = o.total_value;
        return r;
      }, {});

      return {
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
    },
    roadToTarget: function () {
      if (!this.portfolio) {
        return {}
      }
      return {
        labels: ["Complete", "Remaining"],
        datasets: [{
          backgroundColor: this.colors,
          hoverBackgroundColor: this.colors,
          data: [this.portfolio.total_value, this.portfolio.goal - this.portfolio.total_value]
        }]
      }
    },
    lineChartData: function () {
      if (!this.portfolio) {
        return {}
      }
      return {
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
    },
    chartData: function () {
      if (!this.portfolio) {
        return {}
      }
      return {
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
    }
  },
  mounted: async function () {
    await this.setupPortfolio()
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

em {
  font-weight: 700;
}

.text-sm {
  font-size: x-small;
}
</style>
