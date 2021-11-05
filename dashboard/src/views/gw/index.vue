<template>
  <div class="app-container">
    <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
      <el-tab-pane label="证书配置" name="cert">
        <Cert :name="activeName" />
      </el-tab-pane>
      <el-tab-pane label="代理配置" name="proxy">
        <Proxy :name="activeName" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>
import Proxy from './components/proxy/index'
import Cert from './components/cert/index'

export default {
  name: 'Gw',
  components: { Proxy, Cert },
  data() {
    return {
      activeName: 'cert'
    }
  },
  methods: {
    handleClick(tab, event) {
      console.log(tab, event)
    }
  }
}
</script>
