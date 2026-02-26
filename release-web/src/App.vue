<template>
  <el-container class="app-shell">
    <el-aside width="230px" class="side-nav">
      <div class="brand">Release Platform</div>
      <el-menu
        :default-active="$route.path"
        router
        class="menu"
        background-color="transparent"
        text-color="#d7deff"
        active-text-color="#ffffff"
      >
        <el-menu-item index="/upload">上传中心</el-menu-item>
        <el-menu-item index="/sql-audit">SQL审计</el-menu-item>
        <el-menu-item index="/config-diff">配置比对</el-menu-item>
        <el-menu-item index="/dependency">依赖扫描</el-menu-item>
        <el-menu-item index="/approval">上线审批流</el-menu-item>
        <el-menu-item index="/baseline">环境基线</el-menu-item>
        <el-menu-item index="/dashboard">统计仪表盘</el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="top-bar">
        <h1>华泰期货 ITSM 上线发布管控平台</h1>
      </el-header>
      <el-main class="page-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.app-shell {
  min-height: 100vh;
  background: radial-gradient(circle at 15% 20%, rgba(57, 89, 189, 0.2), transparent 45%),
    radial-gradient(circle at 90% 10%, rgba(41, 167, 168, 0.22), transparent 40%),
    #071322;
}

.side-nav {
  border-right: 1px solid rgba(127, 146, 182, 0.24);
  background: linear-gradient(180deg, rgba(16, 35, 64, 0.96), rgba(10, 21, 36, 0.98));
}

.brand {
  padding: 20px 18px;
  font-size: 18px;
  font-weight: 700;
  color: #eff3ff;
  letter-spacing: 0.6px;
}

.menu {
  border-right: none;
}

.top-bar {
  display: flex;
  align-items: center;
  padding: 0 24px;
  border-bottom: 1px solid rgba(177, 196, 228, 0.2);
  background: rgba(6, 19, 34, 0.42);
  backdrop-filter: blur(8px);
}

.top-bar h1 {
  font-size: 18px;
  color: #f7f9ff;
  margin: 0;
}

.page-main {
  padding: 20px;
}
</style>
