<template>
  <div class="dashboard-container">
    <panel-group :online-total="onlineTotal" :user-total="userTotal" :resource-total="resourceTotal" />
    <el-row :gutter="40">
      <el-col :xs="24" :sm="24" :lg="12">
        <el-card>
          <bar-chart v-if="flag" />
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="12">
        <el-card>
          <pie-chart v-if="flag" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import PanelGroup from './components/PanelGroup'
import PieChart from './components/PieChart'
import BarChart from './components/BarChart'
import * as user from '@/api/user'
import * as resource from '@/api/resource'

export default {
  name: 'MyDashboard',
  components: { PanelGroup, PieChart, BarChart },
  data() {
    return {
      onlineTotal: 0,
      userTotal: 0,
      resourceTotal: 0,

      flag: false,
      categoryTotal: {},
      topInfo: {}
    }
  },
  created() {
    this.fetchData()
    this.flag = true
  },
  methods: {
    fetchData() {
      this.getUserCount()
      this.getResourceCount()
      this.getUserOnlineCount()
    },
    async getUserOnlineCount() {
      await user.getOnlineCount().then((response) => {
        this.onlineTotal = response.data.count
      })
    },
    async getUserCount() {
      await user.getCount().then((response) => {
        this.userTotal = response.data.count
      })
    },
    async getResourceCount() {
      await resource.getCount().then((response) => {
        this.resourceTotal = response.data.count
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
