<template>
  <div class="dashboard-container">
    <panel-group :online-total="onlineTotal" :user-total="userTotal" :resource-total="resourceTotal" />
    <el-row :gutter="40">
      <el-col>
        <UserLine v-if="showUser" :data="statics" />
      </el-col>
    </el-row>
  </div>
</template>

<script>
import PanelGroup from './components/PanelGroup'
import UserLine from './components/UserLine.vue'
import * as user from '@/api/user'
import * as resource from '@/api/resource'

export default {
  name: 'Dashboard',
  components: { PanelGroup, UserLine },
  data() {
    return {
      onlineTotal: 0,
      userTotal: 0,
      resourceTotal: 0,

      showUser: false,
      categoryTotal: {},
      topInfo: {},
      statics: {}
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.getUserCount()
      this.getResourceCount()
      this.getUserOnlineCount()
      this.getUserGwInfo(null, 0, 0)
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
    },
    async getUserGwInfo(name, start, end) {
      const query = {}
      if (start === 0 && end === 0) {
        query.start = new Date().getTime() - 3600 * 1000 * 24 * 7
        query.end = new Date().getTime()
      } else {
        query.start = start
        query.end = end
      }
      query.user = name

      await user.getUserGwInfo(query).then((response) => {
        this.statics.info = response.data.statics
        this.statics.start = query.start
        this.statics.end = query.end

        this.showUser = true
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
