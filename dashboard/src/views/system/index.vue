<template>
  <div class="app-container">
    <el-tabs v-model="activeName" type="border-card" @tab-click="handleClick">
      <el-tab-pane label="邮箱配置" name="email">
        <Email :name="activeName" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>
import Email from './components/email'
export default {
  components: { Email },
  data() {
    return {
      activeName: 'email'
    }
  },
  methods: {
    handleClick(tab, event) {
      // console.log(tab, event)
    }
  }
}
</script>
