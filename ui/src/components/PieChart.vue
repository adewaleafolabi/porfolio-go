<script>


import {Doughnut} from "vue-chartjs";
import {formatCurrency} from "../utils/utils";

export default {
  extends: Doughnut,
  props: {
    chartData: {
      type: Object,
      default: null
    },
    options: {
      type: Object,
      default: null
    }
  },
  mounted() {
    let tooltips = {
      tooltips: {
        callbacks: {
          label: function (tooltipItem, data) {
            //get the concerned dataset
            const dataset = data.datasets[tooltipItem.datasetIndex];
            //calculate the total of this data set
            // eslint-disable-next-line
            const total = dataset.data.reduce(function (previousValue, currentValue, currentIndex, array) {
              return previousValue + currentValue;
            });
            //get the current items value
            const currentValue = dataset.data[tooltipItem.index];
            //calculate the percentage based on the total and current item, also this does a rough rounding to give a whole number
            const percentage = Math.floor(((currentValue / total) * 100) + 0.5);

            return `${data.labels[tooltipItem.index]} ${formatCurrency(currentValue)} (${percentage}%)`;
          }
        }
      },
    }
    this.renderChart(this.chartData, {...this.options, ...tooltips})
  }
}
</script>
