<template>
  <el-dialog
    :title="title"
    :visible.sync="visible"
    width="550px"
    :before-close="handleClose"
  >
    <el-form
      ref="formData"
      :rules="rules"
      :model="formData"
      label-width="120px"
      label-position="right"
      style="width: 500px"
      status-icon
    >
      <el-form-item prop="name" label="名称">
        <el-input v-model="formData.name" :disabled="typeof(formData.id) !== 'undefined' && formData.id !== 0" maxlength="30" />
      </el-form-item>
      <el-form-item prop="host">
        <span slot="label">域名
          <el-tooltip placement="top" effect="light">
            <div slot="content">
              web服务的域名，支持通配符
            </div>
            <i class="el-icon-question" />
          </el-tooltip>
        </span>
        <el-col :span="12">
          <el-input v-model="formData.host" />
        </el-col>
      </el-form-item>
      <el-form-item prop="path">
        <span slot="label">路径
          <el-tooltip placement="top" effect="light">
            <div slot="content">
              web服务的访问路径，支持通配符
            </div>
            <i class="el-icon-question" />
          </el-tooltip>
        </span>
        <el-input v-model="formData.path" />
      </el-form-item>
      <el-form-item label="Upstream" prop="upstreamRef">
        <el-col :span="12">
          <el-select v-model="formData.upstreamRef" placeholder="请选择">
            <el-option
              v-for="item in upstreams"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-col>
      </el-form-item>
      <el-form-item label="备注：" prop="remark">
        <el-input v-model="formData.remark" type="textarea" />
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button
        type="primary"
        size=""
        @click="submitForm('formData')"
      >确定</el-button>
      <el-button size="" @click="handleClose">取消</el-button>
    </div>
  </el-dialog>
</template>

<script>
import * as api from '@/api/router'
import { getList } from '@/api/upstream'
import { LBOPTIONS } from '@/utils/upstream'

export default {
  props: {
    title: {
      // 弹窗的标题
      type: String,
      default: ''
    },
    visible: {
      type: Boolean,
      default: false
    },
    formData: {
      type: Object,
      default: () => {}
    },
    remoteClose: {
      type: Function,
      default: () => {}
    }
  },

  data() {
    return {
      upstreams: [],
      LBOPTIONS,
      rules: {
        name: [{ required: true, message: '请输入Router名称', trigger: 'blur' }],
        host: [{ required: true, message: '请输入域名', trigger: 'change' }],
        path: [{ required: true, message: '请输入路径', trigger: 'change' }],
        upstreamRef: [{ required: true, message: '请选择Upstream', trigger: 'change' }]
      }
    }
  },
  watch: {
    visible(newVal, oldVal) {
      if (newVal) {
        if (this.formData.backend === undefined) {
          this.formData.backend = []
          const item = { ip: '', port: '', weight: '' }
          this.formData.backend.push(item)
        }

        if (this.formData.id !== undefined) {
          this.formData.upstreamRef = this.formData.upstreamRef[0].id
        }
        this.fetchUpstreams()
      }
    }
  },

  methods: {
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.submitData()
        } else {
          // console.log('error submit!!');
          return false
        }
      })
    },
    async fetchUpstreams() {
      const response = await getList(null, 0, 0)
      this.upstreams = []
      if ((response.code === 200)) {
        response.data.records.forEach((item, index) => {
          const tmp = {
            value: item.id,
            label: item.name
          }
          this.upstreams.push(tmp)
        })
      }
    },
    async submitData() {
      let response = null
      if (this.formData.id) {
        response = await api.put(this.formData.id, this.formData)
      } else {
        response = await api.add(this.formData)
      }

      if ((response.code === 200)) {
        this.$message({ message: '保存成功', type: 'success' })
        this.handleClose()
      } else {
        this.$message({ message: '保存失败', type: 'error' })
      }
    },
    handleClose() {
      this.$refs['formData'].resetFields()
      this.remoteClose()
    },
    addRule() {
      console.log('add rule')
      if (this.formData.backend.length === 5) {
        this.$message({ message: '至多包含五个条件', type: 'error' })
        return
      }
      const item = { ip: '', port: '' }
      this.formData.backend.push(item)
    },
    deleteRule(row, index) {
      if (this.formData.backend.length === 1) {
        this.$message({ message: '至少包含一个条件', type: 'error' })
        return
      }
      this.formData.backend.splice(index, 1)
    }
  }
}
</script>

<style scoped>
::v-deep .el-dialog__body{padding: 0 20px;}
::v-deep .el-table th, .el-table tr .el-form-item{margin-bottom: 0}
::v-deep .el-input--mini .el-input__inner{ border-radius: 0;}
::v-deep .cell .el-form-item__content .el-form-item__error{left: 5px; top: 55%}
</style>
