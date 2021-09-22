<template>
  <loading-indicator v-if="!portfolioData"></loading-indicator>
  <div class="container" v-else>
    <div class="row">
      <div class="col-12">
        <portfolio-header :portfolio="portfolioData" @privateMode="handlePrivacy" ></portfolio-header>
      </div>
    </div>

    <div class="row mt-2">
      <div class="col-12">
        <div class="card tile">
          <div class="card-header">
            Portfolio Growth
          </div>
          <div class="card-body">
            <portfolio-growth :history="portfolioData.history"
                              :privacy-mode="privacy"
                              v-if="portfolioData" :width="900"
                              :key="privacy"></portfolio-growth>
          </div>
        </div>
      </div>
    </div>

    <div class="row">

      <div class="d-flex justify-content-between">
        <portfolio-history-grouping :history="portfolioData.history"
                                    @historyGrouped="handleGrouping"></portfolio-history-grouping>

        <button :class="{'is-loading':loading}" class="btn btn-sm btn-outline-warning" @click="logPortfolioValue">Update
          History
        </button>
      </div>
    </div>

    <div class="row mt-4">
      <div class="col rounded">
        <div class="card shadow-sm rounded">
          <div class="table-responsive rounded">
            <table class="table table-striped table-hover table-bordered rounded" >
              <thead class="table-dark">
              <tr>
                <th>Date</th>
                <th class="text-end">Value</th>
                <th class="text-end">% Change</th>
              </tr>
              </thead>
              <tbody>
              <tr v-for="(item,index) in portfolioData.history" :key="index">
                <td>{{ (item.date) }}</td>
                <td :class="{'blur':privacy}" class="text-end">{{ formatCurrency(item.value, portfolioData.base_currency, true) }}</td>
                <td v-if="index ===0" class="text-end">--</td>
                <td v-else :class="{'negative':item.value < portfolioData.history[index-1].value}" class="text-end">
                  {{ formatPercentage((item.value - portfolioData.history[index - 1].value) / portfolioData.history[index - 1].value) }}
                </td>
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
import {createPortfoliosDailyValue, getPortfolio} from "../service/portfolio_service";
import {formatCurrency, formatDate, formatDecimal, formatPercentage} from "../utils/utils";
import {DateTime} from "luxon";
import PortfolioHeader from "../components/PortfolioHeader";
import PortfolioGrowth from "../components/PortfolioGrowth";
import PortfolioHistoryGrouping from "../components/PortfolioHistoryGrouping";
import LoadingIndicator from "../components/LoadingIndicator";
import {state} from "../store/store";

export default {

  name: 'PortfolioHistory',
  components: {
    LoadingIndicator,
    PortfolioHistoryGrouping,
    PortfolioGrowth,
    PortfolioHeader,
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
    createPortfoliosDailyValue,
    formatDate,
    formatDecimal,
    formatCurrency,
    formatPercentage,
    handleGrouping(data) {
      this.portfolioData.history = data
      this.growthComponentKey++
    },
    handlePrivacy(privacy) {
      this.privacy = privacy
      this.growthComponentKey++
    },
    async logPortfolioValue() {
      //let today = DateTime.now().toFormat('yyyy-MM-dd')
      if (this.daily.length > 0) {
        let last = this.daily[this.daily.length - 1].date
        if (last === DateTime.now().toFormat('yyyy-MM-dd')) {
          alert("value for day is already logged")
          return
        }
      }

      try {
        this.loading = true
        await createPortfoliosDailyValue()
        alert("operation submitted reload after some minutes")
      } catch (e) {
        console.log(e)
        alert("error processing logValue. check console for details")
        this.loading = false
        return
      }
      this.loading = false
    }
  },
  data() {
    return {
      growthComponentKey: 0,
      loading: false,
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
        width: '100%',
        responsive: false,
        maintainAspectRatio: false
      },
      lineChartStyle: {
        height: 700
      }
    }
  },
  mounted: async function () {
    this.portfolioData = this.portfolio
    if (!this.portfolioData) {
      this.portfolioData = await this.getPortfolio(this.id)
    }

    this.portfolioData.history.sort(function (a, b) {
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
    privacy: function (){
      return state.privacy
    },
    chartData: function () {
      return {
        labels: this.portfolioData.history.map(i => i.date),
        datasets: [
          {
            label: this.portfolioData.name,
            backgroundColor: ['#466eb1'],
            fill: true,
            data: this.portfolioData.history.map(i => i.value)
          }
        ]
      }
    }
  }
}
</script>
