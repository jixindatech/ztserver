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
      <el-form-item label="代理方式" prop="lb">
        <el-col :span="12">
          <el-select v-model="formData.lb" placeholder="请选择">
            <el-option
              v-for="item in LBOPTIONS"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-col>
      </el-form-item>
      <el-form-item>
        <span slot="label">后端配置
          <el-tooltip placement="right">
            <div slot="content">
              目前只支持IP和端口的方式
            </div>
            <i class="el-icon-question" />
          </el-tooltip>
        </span>
        <el-table
          :row-style="{height:'10px'}"
          :cell-style="{padding:'1px'}"
          style="font-size: 8px; margin-top: 0px;"
          size="mini"
          show-header
          :data="formData.backend"
        >
          <el-table-column align="center" label="服务器IP" width="150px">
            <template slot-scope="scope">
              <el-form-item :prop="'backend.' + scope.$index + '.ip'" :rules="rules.ip">
                <el-input v-model="scope.row.ip" size="mini" placeholder="请输入IP" />
              </el-form-item>
            </template>
          </el-table-column>
          <el-table-column align="center" label="服务器端口" width="100px">
            <template slot-scope="scope">
              <el-form-item :prop="'backend.' + scope.$index + '.port'" :rules="rules.port">
                <el-input v-model.number="scope.row.port" size="mini" placeholder="端口号" />
              </el-form-item>
            </template>
          </el-table-column>
          <el-table-column align="center" label="权重" width="100px">
            <template slot-scope="scope">
              <el-form-item :prop="'backend.' + scope.$index + '.weight'" :rules="rules.weight">
                <el-input v-model.number="scope.row.weight" size="mini" placeholder="权重" />
              </el-form-item>
            </template>
          </el-table-column>
          <el-table-column align="center" width="30px">
            <template slot-scope="scope">
              <el-button type="text" icon="el-icon-delete" size="medium" @click="deleteRule(scope.row, scope.$index)" />
            </template>
          </el-table-column>
        </el-table>
        <el-button type="text" icon="el-icon-plus" size="mini" style="margin-bottom: 20px;" @click="addRule()">新增后端</el-button><p style="display: inline;">最多添加5条</p>
      </el-form-item>
      <el-form-item v-if="formData.lb == 'chash' " label="Hash key：" prop="key">
        <el-col :span="12">
          <el-select v-model="formData.key" placeholder="请选择">
            <el-option
              v-for="item in KEYOPTIONS"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-col>
      </el-form-item>
      <el-form-item label="重试次数" prop="retry">
        <el-col :span="12">
          <el-input v-model.number="formData.retry" placeholder="重试次数" />
        </el-col>
      </el-form-item>
      <div class="timeout">
        <el-form-item label="连接超时配置" prop="timeoutConnect">
          <el-input v-model.number="formData.timeoutConnect" size="mini" placeholder="连接超时配置" />
        </el-form-item>
        <el-form-item label="发送超时配置" prop="timeoutSend">
          <el-input v-model.number="formData.timeoutSend" size="mini" placeholder="连接超时配置" />
        </el-form-item>
        <el-form-item label="接收超时配置" prop="timeoutReceive">
          <el-input v-model.number="formData.timeoutReceive" size="mini" placeholder="连接超时配置" />
        </el-form-item>
      </div>
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
import * as api from '@/api/upstream'
import { LBOPTIONS, KEYOPTIONS } from '@/utils/upstream'
import { validateIP, isPort } from '@/utils/validate'

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
      LBOPTIONS,
      KEYOPTIONS,
      rules: {
        name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
        lb: [{ required: true, message: '请选择代理方式', trigger: 'change' }],
        ip: [
          { required: true, message: '请输入IP' },
          { validator: validateIP, tirgger: 'change' }
        ],
        port: [
          { required: true, message: '请输入端口', trigger: 'change' },
          { validator: isPort, tirgger: 'change' }
        ],
        weight: [
          { required: true, message: '请输入权重', trigger: 'blur' }
        ],
        retry: [
          { required: true, message: '请输入重试次数', trigger: 'blur' }
        ],
        timeoutConnect: [
          { required: true, message: '请输入连接超时时间', trigger: 'blur' }
        ],
        timeoutSend: [
          { required: true, message: '请输入发送超时时间', trigger: 'blur' }
        ],
        timeoutReceive: [
          { required: true, message: '请输入接收超时时间', trigger: 'blur' }
        ]
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
      const item = { ip: '', port: '', weight: '' }
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
::v-deep .cell .el-form-item__content .el-form-item__error{ left: 5px; top: 55%}
::v-deep .timeout .el-form-item__error{ left: 5px; top: 55%}
::v-deep .timeout .el-form-item{ margin-bottom: 0px; margin-right: 38%;}

</style>
