<script setup lang="ts">
import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { analyzeRelease } from '../../api/client'
import { useReportStore } from '../../api/reportStore'

const appName = ref('dis-modules-itsm')
const targetEnv = ref('prod')
const submitter = ref('release-user')
const file = ref<File | null>(null)
const loading = ref(false)

const { report, setReport, clearReport } = useReportStore()

const selectedFileText = computed(() => {
  if (!file.value) return '未选择文件'
  const mb = (file.value.size / 1024 / 1024).toFixed(2)
  return `${file.value.name} (${mb} MB)`
})

function onFileChange(uploadFile: any) {
  file.value = uploadFile.raw || null
}

async function handleAnalyze() {
  if (!file.value) {
    ElMessage.warning('请先上传 JAR/ZIP 文件')
    return
  }

  const form = new FormData()
  form.append('file', file.value)
  form.append('appName', appName.value)
  form.append('targetEnv', targetEnv.value)
  form.append('submitter', submitter.value)

  loading.value = true
  try {
    const res: any = await analyzeRelease(form)
    if (!res?.success) {
      ElMessage.error(res?.message || '检查失败')
      return
    }
    setReport(res.data)
    ElMessage.success(`检查完成，申请单号 ${res.data.parse.requestNo}`)
  } catch (error: any) {
    const message =
      error?.response?.data?.message ||
      error?.message ||
      '上传或检查失败，请查看后端日志'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <h2 class="page-title">上传与自动检查</h2>
    <p class="page-tip">上传 JAR/ZIP 后自动执行包解析、SQL 审计、配置比对和依赖扫描。</p>

    <div class="panel">
      <el-form label-width="92px" inline>
        <el-form-item label="应用名称">
          <el-input v-model="appName" style="width: 220px" />
        </el-form-item>
        <el-form-item label="目标环境">
          <el-select v-model="targetEnv" style="width: 150px">
            <el-option label="dev" value="dev" />
            <el-option label="sit" value="sit" />
            <el-option label="uat" value="uat" />
            <el-option label="prod" value="prod" />
          </el-select>
        </el-form-item>
        <el-form-item label="提交人">
          <el-input v-model="submitter" style="width: 180px" />
        </el-form-item>
      </el-form>

      <el-upload drag :auto-upload="false" :show-file-list="true" :limit="1" @change="onFileChange">
        <div class="el-upload__text">拖拽上传 JAR/ZIP 文件，或点击选择</div>
      </el-upload>

      <div style="margin-top: 10px; color: #9cb0ce">当前文件：{{ selectedFileText }}</div>

      <div style="margin-top: 14px; display: flex; gap: 10px">
        <el-button type="primary" :loading="loading" @click="handleAnalyze">
          {{ loading ? '检查中，请稍候...' : '开始检查' }}
        </el-button>
        <el-button @click="clearReport">清空结果</el-button>
      </div>
    </div>

    <div v-if="report" class="panel" style="margin-top: 14px">
      <h3>解析结果</h3>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="申请单号">{{ report.parse.requestNo }}</el-descriptions-item>
        <el-descriptions-item label="应用版本">{{ report.parse.appVersion }}</el-descriptions-item>
        <el-descriptions-item label="构建时间">{{ report.parse.buildTime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="升级版本">{{ report.parse.latestUpgradeVersion || '-' }}</el-descriptions-item>
        <el-descriptions-item label="SQL语句数">{{ report.sqlAudit.totalStatements }}</el-descriptions-item>
        <el-descriptions-item label="依赖数量">{{ report.dependency.totalDependencies }}</el-descriptions-item>
        <el-descriptions-item label="是否自动驳回">
          <el-tag :type="report.rejectedByP0 ? 'danger' : 'success'">{{ report.rejectedByP0 ? '是' : '否' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="包MD5">{{ report.parse.packageMd5 }}</el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>
