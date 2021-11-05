<template>
  <div class="app-container">
    <!-- 条件查询 -->
    <el-form :inline="true" :model="query" size="mini">
      <el-form-item
        label="用户邮箱:"
      >
        <el-input v-model.trim="query.user" />
      </el-form-item>
      <el-form-item>
        <el-button
          icon="el-icon-search"
          type="primary"
          @click="queryData"
        >查询</el-button>
        <el-button
          icon="el-icon-refresh"
          @click="reload"
        >重置</el-button>
      </el-form-item>
      <el-form-item>
        <el-date-picker
          ref="picker"
          v-model="queryTime"
          type="datetimerange"
          :picker-options="pickerOptions"
          range-separator="-"
          start-placeholder=""
          end-placeholder=""
          value-format="timestamp"
          align="right"
        />
      </el-form-item>
    </el-form>
    <el-table
      :data="list"
      stripe
      border
      style="width: 100%"
    >
      <el-table-column align="center" type="index" label="序号" width="60" />
      <el-table-column align="center" prop="_source.user" label="用户" />
      <el-table-column align="center" prop="_source.client" label="IP" />
      <el-table-column align="center" prop="_source.time" label="时间">
        <template slot-scope="scope">
          {{ new Date(scope.row._source.time * 1000).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column align="center" prop="_source.dev" label="设备">
        <template slot-scope="scope">
          <el-tooltip class="item" effect="light" :content="scope.row._source.dev" placement="top-start">
            <span>{{ scope.row._source.dev.substring(0,32) }}...</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="_source.gw" label="网关">
        <template slot-scope="scope">
          <el-input v-for="(item, index) in scope.row._source.gw" :key="index" :value="item" size="mini" />
        </template>
      </el-table-column>
      <el-table-column align="center" prop="_source.event" label="事件" />
    </el-table>
    <el-pagination
      :current-page="page.current"
      :page-sizes="[10, 20, 50]"
      :page-size="page.size"
      layout="total, sizes, prev, pager, next, jumper"
      :total="page.total"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script>
import * as api from '@/api/ws'
const ONE_DAY = new Date().getTime() - 3600 * 1000 * 24 * 1
export default {
  name: 'WsEvent',
  data() {
    return {
      list: [],
      page: {
        current: 1,
        size: 20,
        total: 0
      },
      query: {},
      ONE_DAY,
      queryTime: [],
      pickerOptions: {
        shortcuts: [{
          text: '最近30分钟',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 1800 * 1000)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近一小时',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近24小时',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24)
            picker.$emit('pick', [start, end])
          }
        },
        {
          text: '最近一周',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.$emit('pick', [start, end])
          }
        },
        {
          text: '最近一个月',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
            picker.$emit('pick', [start, end])
          }
        }]
      }
    }
  },

  created() {
    this.queryTime[0] = ONE_DAY
    this.queryTime[1] = new Date().getTime()

    this.fetchData()
  },

  methods: {
    async fetchData() {
      this.query.start = this.queryTime[0]
      this.query.end = this.queryTime[1]

      const { data } = await api.getList(
        this.query,
        this.page.current,
        this.page.size
      )

      this.list = data.records
      this.page.total = data.total
    },

    handleSizeChange(val) {
      this.page.size = val
      this.fetchData()
    },

    handleCurrentChange(val) {
      this.page.current = val
      this.fetchData()
    },

    queryData() {
      this.page.current = 1
      this.fetchData()

      console.log(this.queryTime)
    },

    reload() {
      this.$set(this.queryTime, 0, ONE_DAY)
      this.$set(this.queryTime, 1, new Date().getTime())

      this.query = {}
      this.fetchData()
    }
  }
}
</script>

<style>
  .el-tooltip__popper {
    max-width: 400px;
    line-height: 180%;
  }
</style>
