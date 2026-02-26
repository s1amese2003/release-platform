<script setup lang="ts">
import { computed } from 'vue'
import { useReportStore } from '../../api/reportStore'

const { report } = useReportStore()

const cards = computed(() => {
  if (!report.value) {
    return [
      { title: 'SQL语句数', value: 0 },
      { title: 'P0风险', value: 0 },
      { title: '配置差异项', value: 0 },
      { title: '依赖高风险', value: 0 },
    ]
  }

  return [
    { title: 'SQL语句数', value: report.value.sqlAudit.totalStatements },
    { title: 'P0风险', value: report.value.sqlAudit.summary.P0 + report.value.configDiff.summary.P0 },
    { title: '配置差异项', value: report.value.configDiff.items.length },
    { title: '依赖高风险', value: report.value.dependency.highRiskDependencies },
  ]
})
</script>

<template>
  <div>
    <h2 class="page-title">统计仪表盘</h2>
    <p class="page-tip">展示最近一次检查结果的核心指标，可扩展为多申请单趋势看板。</p>

    <div class="cards">
      <div v-for="card in cards" :key="card.title" class="panel card">
        <div class="label">{{ card.title }}</div>
        <div class="value">{{ card.value }}</div>
      </div>
    </div>

    <div class="panel" style="margin-top: 14px">
      <h3 style="margin-top: 0">流程节点覆盖</h3>
      <el-steps :active="report ? 4 : 0" finish-status="success" align-center>
        <el-step title="提交上线申请" />
        <el-step title="自动检查" />
        <el-step title="人工审批" />
        <el-step title="上线执行" />
        <el-step title="完成并快照" />
      </el-steps>
    </div>
  </div>
</template>

<style scoped>
.cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
}

.card .label {
  color: var(--text-sub);
}

.card .value {
  margin-top: 10px;
  font-size: 30px;
  font-weight: 700;
}
</style>
