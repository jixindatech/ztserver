<template>
  <el-dialog
    :title="title"
    :visible.sync="visible"
    width="500px"
    :before-close="handleClose"
  >
    <el-form
      ref="formData"
      :rules="rules"
      :model="formData"
      label-width="100px"
      label-position="right"
      style="width: 400px"
      status-icon
    >
      <el-form-item label="名称:" prop="name">
        <el-input v-model="formData.name" :disabled="typeof(formData.id) !== 'undefined' && formData.id !== 0" maxlength="254" />
      </el-form-item>
      <el-form-item prop="host">
        <span slot="label">域名:
          <el-tooltip placement="top" effect="light">
            <div slot="content">
              web服务的访问域名，支持通配符
            </div>
            <i class="el-icon-question" />
          </el-tooltip>
        </span>
        <el-input v-model="formData.host" maxlength="254" />
      </el-form-item>
      <el-form-item prop="path">
        <span slot="label">路径:
          <el-tooltip placement="top" effect="light">
            <div slot="content">
              web服务的访问路径，支持通配符
            </div>
            <i class="el-icon-question" />
          </el-tooltip>
        </span>
        <el-input v-model="formData.path" maxlength="254" />
      </el-form-item>
      <el-form-item label="请求方法" prop="method">
        <el-checkbox-group v-model="methods">
          <el-checkbox v-for="method in METHOD_OPTIONS" :key="method" :label="method">{{ method }}</el-checkbox>
        </el-checkbox-group>
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
import * as api from '@/api/resource'
import { METHOD_OPTIONS } from '@/utils/resource'
export default {
  props: {
    title: {
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
      METHOD_OPTIONS,
      methods: [],
      rules: {
        name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
        host: [{ required: true, message: '请输入请求域名', trigger: 'blur' }],
        path: [{ required: true, message: '请输入请求路径', trigger: 'blur' }]
      }
    }
  },
  watch: {
    visible(newVal, oldVal) {
      if (newVal) {
        if (this.formData.id !== undefined) {
          this.methods = this.formData.method
        }
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

    async submitData() {
      let response = null
      this.formData.method = this.methods
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
      this.methods = []
      this.$refs['formData'].resetFields()
      this.remoteClose()
    }
  }
}
</script>
