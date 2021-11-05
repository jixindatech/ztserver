<template>
  <!-- 具备一个宽高的dom容器 -->
  <div ref="main" :class="className" :style="{height: height, width: width}" />
</template>

<script>
import * as echarts from 'echarts'
require('echarts/theme/macarons')
import resize from './mixins/resize'
// import api from '@/api/event'

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
      default: '400px'
    }
  },
  data() {
    return {
      chart: null, // 接收echarts实例的属性
      legendData: [],
      seriesData: []
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initPieChart()
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
    async initPieChart() {
      /* await api.getEventsStatisticsInfo().then((response) => {
        this.legendData = response.data.types
        this.seriesData = response.data.values
      })
      */
      this.chart = echarts.init(this.$refs.main, 'macarons')
      this.chart.setOption({
        title: { // 标题
          text: '触发事件汇总',
          left: 'center' // 居中
        },
        tooltip: { // 鼠标放上去的提示框格式
          trigger: 'item',
          formatter: '{a} <br/>{b} : {c} ({d}%)'
        },
        legend: { // 左上角
          orient: 'vertical',
          left: 'left',
          data: this.legendData
        },
        series: [ // 序列，展示的具体数据
          {
            name: '统计内容',
            type: 'pie', // 饼状图
            radius: '55 %', // 圆大小
            center: ['50%', '50%'], // 饼图位置【左，上】
            data: this.seriesData,
            emphasis: {
              itemStyle: {
                shadowBlur: 10, // 图形阴影的模糊大小
                shadowOffsetX: 0, // 阴影水平方向偏移距离
                shadowColor: 'rgba(0, 0, 0, 0.5)' // 阴影颜色
              }
            }
          }
        ]
      })
    }
  }
}
</script>
