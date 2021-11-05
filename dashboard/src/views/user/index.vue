<template>
  <div class="app-container">
    <!-- 条件查询 -->
    <el-form :inline="true" :model="query" size="mini">
      <el-form-item
        label="用户名:"
      >
        <el-input v-model.trim="query.username" />
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
      <el-table-column align="center" prop="name" label="用户名" />
      <el-table-column align="center" prop="phone" label="手机号" />
      <el-table-column align="center" prop="email" label="邮箱" />
      <el-table-column align="center" prop="status" label="帐号锁定">
        <!-- (1 未锁定，0已锁定) -->
        <template slot-scope="scope">
          <el-tag v-if="scope.row.status === 0" type="danger">锁定</el-tag>
          <el-tag v-if="scope.row.status === 1" type="success">正常</el-tag>
        </template>
      </el-table-column>
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
          <el-button
            type="primary"
            size="mini"
            @click="handleResource(scope.row.id)"
          >配置资源</el-button>
          <el-button
            type="primary"
            size="mini"
            @click="handleMail(scope.row.id, scope.row)"
          >初始化</el-button>
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
    <!--
    <el-dialog title="设置角色" :visible.sync="role.visible" width="65%">
       roleIds是当前用户所拥有的角色id, saveUserRole事件是子组件进行触发提交选择的角色id
      <Role :role-ids="role.roleIds" @saveUserRole="saveUserRole" />
    </el-dialog>
    -->
    <el-dialog title="配置资源" :visible.sync="resource.visible" width="65%">
      <Resource :ids="resource.ids" @saveUserResource="saveUserResource" />
    </el-dialog>
  </div>
</template>

<script>
import * as api from '@/api/user'
import Edit from './edit'
import Resource from '../resource'

export default {
  name: 'User', // 和对应路由表中配置的name值一致
  components: { Edit, Resource /* Role*/ },

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

      role: {
        visible: false,
        // 传递到子组件中时,至少会传递一个空数据[], 子组件判断是否有roleIds值
        roleIds: [], // 封装当前用户所拥有的 角色id

        userId: null // 点击哪个用户，就是哪个用户id,当保存用户角色时，需要使用
      },
      resource: {
        visible: false,
        ids: [],
        userId: null
      }
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
        this.edit.formData = response.data.user
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

    handleResource(id) {
      this.resource.userId = id
      api.getUserResource(id).then((response) => {
        this.resource.ids = response.data.records
        this.resource.visible = true
      })
    },

    saveUserResource(ids) {
      const data = { ids: ids }
      api.saveUserResource(this.resource.userId, data).then((response) => {
        if (response.code === 200) {
          this.$message({ message: '配置资源成功', type: 'success' })
          this.resource.visible = false
          this.resource.ids = []
          this.resource.userId = null
        } else {
          this.$message({ message: '配置资源失败', type: 'error' })
        }
      })
    },

    handleMail(id, data) {
      console.log('handMail')
      api.sendMail(id, data).then((response) => {
        if (response.code === 200) {
          this.$message({ message: '已经发送邮件', type: 'success' })
        } else {
          this.$message({ message: '发送邮件失败', type: 'error' })
        }
      })
    }
  }
}
</script>
