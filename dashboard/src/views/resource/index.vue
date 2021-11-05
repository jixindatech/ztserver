<template>
  <div class="app-container">
    <!-- 条件查询 -->
    <el-form :inline="true" :model="query" size="mini">
      <el-form-item
        label="资源名称:"
      >
        <el-input v-model.trim="query.name" />
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
          v-if="!ids"
          icon="el-icon-circle-plus-outline"
          type="primary"
          @click="openAdd"
        >新增</el-button>
        <el-button
          v-if="ids"
          icon="el-icon-circle-plus-outline"
          type="success"
          @click="saveUserResource"
        >配置</el-button>
      </el-form-item>
    </el-form>

    <el-table
      ref="dataTable"
      :data="list"
      stripe
      border
      style="width: 100%"
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column v-if="ids" align="center" reserve-selection type="selection" width="55" />
      <el-table-column align="center" type="index" label="序号" width="60" />
      <el-table-column align="center" prop="name" label="名称" />
      <el-table-column align="center" prop="server" label="服务" />
      <el-table-column v-if="!ids" align="center" label="操作" width="330">
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
import * as api from '@/api/resource'
import Edit from './edit'

export default {
  name: 'Resource',
  components: { Edit },
  props: {
    ids: {
      type: Array,
      default: function() {
        return null
      }
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
      },
      checkedResourceList: []
    }
  },

  watch: {
    ids() {
      this.query = {}
      this.queryData()
    }
  },

  created() {
    this.fetchData()
  },

  methods: {
    async fetchData() {
      const { data } = await api.getList(
        this.query,
        this.page.current,
        this.page.size
      )

      this.list = data.records
      this.page.total = data.total

      this.chekedResources()
    },
    chekedResources() {
      this.$refs.dataTable.clearSelection()
      if (this.ids) {
        this.list.forEach((item) => {
          if (this.ids.indexOf(item.id) !== -1) {
            this.$refs.dataTable.toggleRowSelection(item, true)
          }
        })
      }
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

    handleEdit(id) {
      api.get(id).then((response) => {
        this.edit.formData = response.data.resource
        this.edit.title = '编辑'
        this.edit.visible = true
      })
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

    remoteClose() {
      this.edit.formData = {}
      this.edit.visible = false
      this.fetchData()
    },

    handleSelectionChange(val) {
      this.checkedResourceList = val
    },

    saveUserResource() {
      const checkedResourceIds = []
      this.checkedResourceList.forEach((item) => {
        checkedResourceIds.push(item.id)
      })

      this.$emit('saveUserResource', checkedResourceIds)
    }
  }
}
</script>
