<template>
  <div class="row">
    <div class="col-12">
      <div class=" mr-3">
        <line-chart v-if="chartData" :chart-data="chartData" :options="chartOptions" :privacy-mode="privacyMode"
                    :styles="lineChartStyle"></line-chart>
      </div>
    </div>
  </div>
</template>

<script>
import LineChart from "./LineChart";
import { formatDate} from "../utils/utils";
export default {
  name: "PortfolioGrowth",
  components: {LineChart},
  props: {
    history: {
      type: Array,
      default: ()=>[]
    },
    privacyMode: {
      type: Boolean,
      default: ()=>false
    },
    width:{
      type: Number
    }
  },
  data() {
    return {
      chartOptions: {
        scaleShowLabels: this.privacyMode,
        legend: {
          display: true
        },
        scales: {
          yAxes: [
            {
              ticks: {
                beginAtZero: true,
                display: true
              }
            }]
        },
        height: 500,
        responsive: true,
        maintainAspectRatio: false
      },
      lineChartStyle: {
        height: 500
      },
      chartData: null,
    }
  },
  mounted: function() {
    if(this.width){
      this.chartOptions.width = this.width
      this.lineChartStyle.width = this.width
    }
    this.chartData = {
      labels: this.history.map(i => formatDate(i.date)),
      datasets: [
        {
          backgroundColor: ['#466eb1'],
          fill: true,
          data: this.history.map(i => i.value)
        }
      ]
    }
  }
}
</script>

<style scoped>

</style>