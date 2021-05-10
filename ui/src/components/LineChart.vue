<script>


import {Line} from "vue-chartjs";
import {formatCurrency, formatPercentage} from "../utils/utils";

export default {
  extends: Line,
  props: {
    chartData: {
      type: Object,
      default: null
    },
    privacyMode: {
      type: Boolean,
      default: ()=>false
    },
    options: {
      type: Object,
      default: null
    },
    styles: {
      type: Object,
      default: null
    }
  },
  mounted () {
    this.options.legend.display = false
    let privacy = this.privacyMode
    let tooltips = {
      tooltips: {
        callbacks: {
          label: function (tooltipItem, data) {
            //get the concerned dataset
            const dataset = data.datasets[tooltipItem.datasetIndex];
            const currentValue = dataset.data[tooltipItem.index];

            //calculate the total of this data set
            // eslint-disable-next-line
            let change = 0
            if(tooltipItem.index > 0){
              const previousValue = dataset.data[tooltipItem.index-1];
              change = (currentValue - previousValue)/ previousValue
            }
            return `${privacy? '':formatCurrency(currentValue)} (${change< 0 ? '⬇️':'⬆️'} ${formatPercentage(change)})`;
          }
        }
      },
    }
    if(this.privacyMode){
      this.options.scaleShowLabels = false
      this.options.scales.yAxes[0].ticks.display = false
    }
    this.renderChart(this.chartData, {...this.options,...tooltips}, this.styles)
  }
}
</script>
