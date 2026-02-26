<script setup lang="ts">
import { computed } from 'vue'
import { useReportStore } from '../../api/reportStore'
import MonacoDiffEditor from '../../components/MonacoDiffEditor.vue'

const { report } = useReportStore()

const diffItems = computed(() => report.value?.configDiff?.items || [])
const baselineYaml = computed(() => {
  if (!report.value) return ''
  const rows = diffItems.value.map((item: any) => `${item.configKey}: ${item.baselineValue ?? 'null'}`)
  return rows.join('\n')
})

const packageYaml = computed(() => {
  if (!report.value) return ''
  const rows = diffItems.value.map((item: any) => `${item.configKey}: ${item.packageValue ?? 'null'}`)
  return rows.join('\n')
})
</script>

<template>
  <div>
    <h2 class="page-title">配置比对</h2>
    <p class="page-tip">上传包配置与目标环境基线做结构化对比，并识别生产禁止项。</p>

    <div v-if="!report" class="panel">暂无数据，请先在上传中心执行检查。</div>

    <template v-else>
      <div class="panel" style="margin-bottom: 14px">
        <el-table :data="diffItems" stripe>
          <el-table-column prop="configKey" label="配置项" min-width="280" />
          <el-table-column prop="packageValue" label="包内值" min-width="200" />
          <el-table-column prop="baselineValue" label="基线值" min-width="200" />
          <el-table-column prop="diffType" label="状态" width="110" />
          <el-table-column prop="riskLevel" label="风险" width="90" />
          <el-table-column prop="reason" label="说明" min-width="250" />
        </el-table>
      </div>

      <div class="panel">
        <h3 style="margin-top: 0">YAML 差异视图</h3>
        <MonacoDiffEditor :original="baselineYaml" :modified="packageYaml" language="yaml" height="440px" />
      </div>
    </template>
  </div>
</template>
