<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { approvalAction } from '../../api/client'
import { useReportStore } from '../../api/reportStore'

const { report } = useReportStore()
const approver = ref('')
const comment = ref('')
const latestStatus = ref('-')

async function act(action: string) {
  if (!report.value) {
    ElMessage.warning('请先生成检查报告')
    return
  }
  if (!approver.value) {
    ElMessage.warning('请填写审批人')
    return
  }
  const requestNo = report.value.parse.requestNo
  const res: any = await approvalAction(requestNo, {
    approver: approver.value,
    action,
    comment: comment.value,
  })
  if (!res.success) {
    ElMessage.error(res.message || '审批失败')
    return
  }
  latestStatus.value = res.data
  ElMessage.success(`审批动作已提交: ${action}`)
}
</script>

<template>
  <div>
    <h2 class="page-title">上线审批流</h2>
    <p class="page-tip">支持开发负责人、DBA、运维负责人逐级审批，P0 自动驳回。</p>

    <div v-if="!report" class="panel">暂无数据，请先在上传中心执行检查。</div>

    <div v-else class="panel">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="申请单号">{{ report.parse.requestNo }}</el-descriptions-item>
        <el-descriptions-item label="目标环境">{{ report.parse.targetEnv }}</el-descriptions-item>
        <el-descriptions-item label="P0问题数">{{ report.sqlAudit.summary.P0 + report.configDiff.summary.P0 }}</el-descriptions-item>
        <el-descriptions-item label="当前审批状态">{{ latestStatus }}</el-descriptions-item>
      </el-descriptions>

      <el-form label-width="80px" style="margin-top: 14px">
        <el-form-item label="审批人">
          <el-input v-model="approver" placeholder="如: zhangsan" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="comment" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>

      <el-space>
        <el-button type="primary" @click="act('APPROVE')">通过</el-button>
        <el-button type="danger" @click="act('REJECT')">驳回</el-button>
        <el-button type="success" @click="act('DEPLOY')">标记上线完成</el-button>
      </el-space>

      <el-divider />
      <h3>手动操作清单（来自 config.txt）</h3>
      <el-table :data="report.manualChecklist" stripe>
        <el-table-column prop="operationDesc" label="操作项" min-width="420" />
        <el-table-column prop="operationType" label="类型" width="140" />
        <el-table-column prop="execStatus" label="执行状态" width="130" />
      </el-table>
    </div>
  </div>
</template>
