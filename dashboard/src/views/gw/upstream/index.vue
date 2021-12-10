<template>
  <div class="app-container">
    <el-form :inline="true" :model="query" size="mini">
      <el-form-item
        label="名称:"
      >
        <el-input v-model.trim="query.server" />
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
        <el-button
          icon="el-icon-circle-plus-outline"
          type="primary"
          @click="openAdd"
        >新增</el-button>
      </el-form-item>
    </el-form>
    <el-table
      :data="list"
      stripe
      border
      style="width: 100%"
    >
      <el-table-column align="center" type="index" label="序号" width="60" />
      <el-table-column align="center" prop="name" label="Upstream名称" />
      <el-table-column align="center" prop="lb" label="负载均衡方式">
        <template slot-scope="scope">
          <template v-if="scope.row.key.length > 0">
            {{ scope.row.lb + ":" + scope.row.key }}
          </template>
          <template v-else>
            {{ scope.row.lb }}
          </template>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="backend" :render-header="renderLBHeader">
        <template slot-scope="scope">
          <el-input v-for="(item, index) in scope.row.backend" :key="index" :value="item.ip + ':' + item.port + ':' + item.weight" size="mini" />
        </template>
      </el-table-column>
      <el-table-column align="center" prop="remark" label="备注" />
      <el-table-column align="center" label="操作" width="330">
        <template slot-scope="scope">
          <el-button
            type="success"
            size="mini"
            @click="handleEdit(scope.row.id)"
          >编辑</el-button>
          <el-button
            type="danger"
            size="mini"
            @click="handleDelete(scope.row.id)"
          >删除</el-button>
        </template>
      </el-table-column>
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
    <edit
      :title="edit.title"
      :form-data="edit.formData"
      :visible="edit.visible"
      :remote-close="remoteClose"
    />

  </div>
</template>

<script>
import * as api from '@/api/upstream'
import Edit from './edit'

export default {
  name: 'Upstream',
  components: { Edit },
  props: {
    name: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      list: [],
      page: {
        current: 1,
        size: 20,
        total: 0
      },
      query: {},
      edit: {
        title: '',
        visible: false,
        formData: {}
      }
    }
  },
  created() {
    this.fetchData()
  },

  methods: {
    renderLBHeader() {
      return (
        <div style='display : inline; '>
          <el-tooltip class='tooltip' effect='light' placement='top'>
            <i class='el-icon-question'>后端</i>
            <div slot='content' >IP 端口及权重按顺序以 ':' 分割</div>
          </el-tooltip>
        </div>
      )
    },
    async fetchData() {
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
    },

    reload() {
      this.query = {}
      this.fetchData()
    },

    handleDelete(id) {
      this.$confirm('确认删除这条记录吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(() => {
          api.deleteById(id).then((response) => {
            this.$message({
              type: response.code === 200 ? 'success' : 'error',
              message: '删除成功!'
            })
            this.fetchData()
          })
        })
        .catch(() => {
        })
    },

    openAdd() {
      this.edit.title = '新增'
      this.edit.visible = true
    },

    handleEdit(id) {
      api.get(id).then((response) => {
        this.edit.formData = response.data.record
        this.edit.title = '编辑'
        this.edit.visible = true
      })
    },
    remoteClose() {
      this.edit.formData = {}
      this.edit.visible = false
      this.fetchData()
    }
  }
}
</script>
