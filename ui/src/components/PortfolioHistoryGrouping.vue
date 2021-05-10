<template>
  <div class="btn-group" role="group">
    <button type="button" class="btn btn-outline-primary btn-sm" @click="handleSelection('DAY')">Daily</button>
    <button type="button" class="btn btn-outline-primary btn-sm" @click="handleSelection('MONTH')">Monthly</button>
    <button type="button" class="btn btn-outline-primary btn-sm" @click="handleSelection('YEAR')">Yearly</button>
  </div>
</template>

<script>
import {DateTime} from "luxon";

export default {
  name: "PortfolioHistoryGrouping",
  data() {
    return {
      daily: [],
      monthly: [],
      yearly: [],
    }
  },
  props: {
    history: {
      type: Array,
      default: () => []
    }
  },
  mounted: function () {
    this.daily =  this.groupHistory('DAY', this.history)
    this.monthly =  this.groupHistory('MONTH', this.history)
    this.yearly =  this.groupHistory('YEAR', this.history)

    this.handleSelection('DAY')
  },
  methods: {
    handleSelection(period){
      switch (period){
        case 'YEAR':
          this.$emit("historyGrouped", this.yearly)
          break
        case 'MONTH':
          this.$emit("historyGrouped", this.monthly)
          break
        default:
          this.$emit("historyGrouped", this.daily)
      }
    },
    groupHistory(grouping, history) {
      let format = grouping === 'DAY' ? 'yyyy-MM-dd'
          : grouping === 'MONTH' ? 'yyyy-MM' : 'yyyy'

      //reduce by grouping items by date
      const groupBy = (arr) => arr.reduce((acc, ele) => ((acc[DateTime.fromISO(ele.date).toFormat(format)] = acc[DateTime.fromISO(ele.date).toFormat(format)] || []).push(ele), acc), {})

      //get last day in month  map entry and return result
      return Object.entries(groupBy(history)).map(([key, val]) => ({date: key, value: val[val.length - 1].value}))

    }
  }
}
</script>

<style scoped>

</style>