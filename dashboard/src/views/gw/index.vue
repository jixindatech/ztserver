<template>
  <div class="app-container">
    <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
      <el-tab-pane label="证书配置" name="ssl">
        <SSL :name="activeName" />
      </el-tab-pane>
      <el-tab-pane label="代理配置" name="upstream">
        <Upstream :name="activeName" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>
import Upstream from './components/upstream/index'
import SSL from './components/ssl/index'

export default {
  name: 'Gw',
  components: { Upstream, SSL },
  data() {
    return {
      activeName: 'ssl'
    }
  },
  methods: {
    handleClick(tab, event) {
      console.log(tab, event)
    }
  }
}
</script>
