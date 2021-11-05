<template>
  <div class="app-container" style="width: 600px" align="center">
    <el-form
      ref="formData"
      :rules="rules"
      :model="formData"
      label-width="80px"
      label-position="left"
      status-icon
    >
      <el-form-item label="服务器" prop="server">
        <el-input v-model="formData.server" maxlength="254" />
      </el-form-item>
      <el-form-item label="端口" prop="port">
        <el-input v-model.number="formData.port" maxlength="254" />
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="formData.email" maxlength="254" />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input v-model="formData.password" maxlength="254" show-password />
      </el-form-item>
      <el-form-item label="备注" prop="remark">
        <el-input v-model="formData.remark" type="textarea" />
      </el-form-item>
    </el-form>
    <el-button
      type="primary"
      size=""
      @click="submitForm('formData')"
    >确定</el-button>
  </div>
</template>

<script>
import * as api from '@/api/email'
export default {
  name: 'Email',
  props: {
    name: {
      type: String,
      default: ''
    }
  },

  data() {
    return {
      formData: {},
      rules: {
        server: [{ required: true, message: '请输入服务名称', trigger: 'blur' }],
        port: [
          { required: true, message: '端口不能为空' },
          { type: 'number', message: '端口必须为数字值' }],
        email: [
          { required: true, message: '请输入邮箱地址', trigger: 'blur' },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }],
        password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
      }
    }
  },
  watch: {
    name(newVal, oldVal) {
      this.fetchData()
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      api.get().then((response) => {
        console.log(response)
        if (response.data.email.id !== 0) {
          console.log(response.data.email)
          this.formData = response.data.email
        }
      })
    },
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
      if (this.formData.id !== 0) {
        console.log('put')
        response = await api.put(this.formData.id, this.formData)
      } else {
        console.log('add')
        response = await api.add(this.formData)
      }

      if ((response.code === 200)) {
        this.$message({ message: '保存成功', type: 'success' })
      } else {
        this.$message({ message: '保存失败', type: 'error' })
      }
    }
  }
}
</script>
