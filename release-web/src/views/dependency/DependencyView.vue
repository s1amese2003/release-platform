<script setup lang="ts">
import { computed } from 'vue'
import { useReportStore } from '../../api/reportStore'

const { report } = useReportStore()
const changes = computed(() => report.value?.dependency?.changes || [])
</script>

<template>
  <div>
    <h2 class="page-title">依赖安全扫描</h2>
    <p class="page-tip">扫描 BOOT-INF/lib 依赖并提示高风险组件与版本变更。</p>

    <div v-if="!report" class="panel">暂无数据，请先在上传中心执行检查。</div>

    <div v-else class="panel">
      <el-descriptions :column="3" border>
        <el-descriptions-item label="依赖总数">{{ report.dependency.totalDependencies }}</el-descriptions-item>
        <el-descriptions-item label="高风险数量">{{ report.dependency.highRiskDependencies }}</el-descriptions-item>
        <el-descriptions-item label="变更项数量">{{ changes.length }}</el-descriptions-item>
      </el-descriptions>

      <el-table :data="changes" style="margin-top: 12px" stripe>
        <el-table-column label="坐标" min-width="240">
          <template #default="scope">{{ scope.row.groupId }}:{{ scope.row.artifactId }}</template>
        </el-table-column>
        <el-table-column prop="oldVersion" label="旧版本" width="120" />
        <el-table-column prop="newVersion" label="新版本" width="120" />
        <el-table-column prop="changeType" label="变更类型" width="110" />
        <el-table-column prop="riskLevel" label="风险" width="90" />
        <el-table-column label="CVE" min-width="240">
          <template #default="scope">{{ (scope.row.cveIds || []).join(', ') || '-' }}</template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>
