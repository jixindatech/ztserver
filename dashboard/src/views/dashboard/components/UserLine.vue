<template>
  <div ref="main" :class="className" :style="{height:height,width:width}" />
</template>

<script>
import * as echarts from 'echarts'
require('echarts/theme/macarons') // echarts theme
import resize from './mixins/resize'

export default {
  mixins: [resize],
  props: {
    className: {
      type: String,
      default: 'chart'
    },
    width: {
      type: String,
      default: '100%'
    },
    height: {
      type: String,
      default: '450px'
    },
    autoResize: {
      type: Boolean,
      default: true
    },
    data: {
      type: Object,
      default: function() {
        return {}
      }
    }
  },
  data() {
    return {
      chart: null,
      // userData: [130, 140, 141, 142, 145, 150, 160, 11, 87, 20],
      userData: [],
      timeData: [],
      defaultData: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    }
  },
  watch: {
    data: {
      handler: function(newValue, oldValue) {
        this.data = newValue
        this.setData()
        this.initChart()
      },
      deep: true
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.setData()
      this.initChart()
    })
  },
  beforeDestroy() {
    if (!this.chart) {
      return
    }
    this.chart.dispose()
    this.chart = null
  },
  methods: {
    setData() {
      this.userData = this.data.info ? this.data.info : this.defaultData

      const intervalTime = (this.data.end - this.data.start) / 10
      const timesSplice = []
      for (var i = 0; i < 10; i++) {
        const timeData = this.data.start + i * intervalTime
        timesSplice.push((new Date(timeData)).toLocaleString())
      }
      timesSplice.push((new Date(this.data.end)).toLocaleString())
      this.timeData = timesSplice
    },
    initChart() {
      this.chart = echarts.init(this.$refs.main, 'macarons')
      this.setOptions(this.chartData)
    },
    setOptions({ expectedData, actualData } = {}) {
      this.chart.setOption({
        title: { // 标题
          text: '在线用户统计',
          left: 'left'
        },
        xAxis: {
          name: '日期',
          data: this.timeData,
          boundaryGap: false,
          axisTick: {
            show: false
          },
          axisLabel: {
            interval: 0,
            rotate: -30
          }
        },
        grid: {
          left: 10,
          right: 10,
          bottom: 20,
          top: 30,
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross'
          },
          padding: [5, 10]
        },
        yAxis: {
          axisTick: {
            show: false
          }
        },
        legend: {
          data: ['USER']
        },
        series: [
          {
            name: 'USER',
            smooth: true,
            type: 'line',
            itemStyle: {
              normal: {
                color: '#3888fa',
                lineStyle: {
                  color: '#3888fa',
                  width: 2
                },
                areaStyle: {
                  color: '#f3f8ff'
                }
              }
            },
            data: this.userData,
            animationDuration: 2800,
            animationEasing: 'quadraticOut'
          }
        ]
      })
    }
  }
}
</script>
