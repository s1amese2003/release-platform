<script setup lang="ts">
import { computed } from 'vue'
import { useReportStore } from '../../api/reportStore'
import SqlHighlight from '../../components/SqlHighlight.vue'

const { report } = useReportStore()

const sqlIssues = computed(() => report.value?.sqlAudit?.results || [])
const summary = computed(() => report.value?.sqlAudit?.summary || { P0: 0, P1: 0, P2: 0, SAFE: 0 })
</script>

<template>
  <div>
    <h2 class="page-title">SQL 审计报告</h2>
    <p class="page-tip">基于静态规则检测 DDL/DML 风险，P0 自动阻断上线流程。</p>

    <div v-if="!report" class="panel">暂无数据，请先在上传中心执行检查。</div>

    <template v-else>
      <div class="panel" style="margin-bottom: 14px">
        <el-space>
          <el-tag type="danger">P0: {{ summary.P0 }}</el-tag>
          <el-tag type="warning">P1: {{ summary.P1 }}</el-tag>
          <el-tag>P2: {{ summary.P2 }}</el-tag>
          <el-tag type="success">SAFE: {{ summary.SAFE }}</el-tag>
        </el-space>

        <el-table :data="sqlIssues" style="margin-top: 14px" stripe>
          <el-table-column prop="line" label="行号" width="80" />
          <el-table-column prop="level" label="等级" width="90" />
          <el-table-column prop="reason" label="风险原因" min-width="260" />
          <el-table-column prop="suggestion" label="建议" min-width="280" />
        </el-table>
      </div>

      <div class="panel">
        <h3 style="margin-top: 0">完整 SQL 预览</h3>
        <SqlHighlight :sql="report.parse.sqlContent || ''" :issues="sqlIssues" />
      </div>
    </template>
  </div>
</template>
